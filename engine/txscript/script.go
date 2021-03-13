// Copyright (c) 2017-2018 The qitmeer developers
// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2018 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Qitmeer/qitmeer/common/hash"
	"github.com/Qitmeer/qitmeer/core/types"
	"github.com/Qitmeer/qitmeer/params"
)

// These are the constants specified for maximums in individual scripts.
const (
	MaxOpsPerScript       = 255  // Max number of non-push operations.
	MaxPubKeysPerMultiSig = 20   // Multisig can't have more sigs than this.
	MaxScriptElementSize  = 2048 // Max bytes pushable to the stack.
	// witnessV0ScriptHashLen is the length of a P2WSH script.
	witnessV0ScriptHashLen = 34

	// pubKeyHashLen is the length of a P2PKH script.
	pubKeyHashLen = 25

	// witnessV0PubKeyHashLen is the length of a P2WPKH script.
	witnessV0PubKeyHashLen = 22

	// scriptHashLen is the length of a P2SH script.
	scriptHashLen = 23

	// maxLen is the maximum script length supported by ParsePkScript.
	maxLen = witnessV0ScriptHashLen

	// minPubKeyHashSigScriptLen is the minimum length of a signature script
	// that spends a P2PKH output. The length is composed of the following:
	//   Signature length (1 byte)
	//   Signature (min 8 bytes)
	//   Signature hash type (1 byte)
	//   Public key length (1 byte)
	//   Public key (33 byte)
	minPubKeyHashSigScriptLen = 1 + 8 + 1 + 1 + 33

	// maxPubKeyHashSigScriptLen is the maximum length of a signature script
	// that spends a P2PKH output. The length is composed of the following:
	//   Signature length (1 byte)
	//   Signature (max 72 bytes)
	//   Signature hash type (1 byte)
	//   Public key length (1 byte)
	//   Public key (33 byte)
	maxPubKeyHashSigScriptLen = 1 + 72 + 1 + 1 + 33

	// compressedPubKeyLen is the length in bytes of a compressed public
	// key.
	compressedPubKeyLen = 33
)

var (
	// ErrUnsupportedScriptType is an error returned when we attempt to
	// parse/re-compute an output script into a PkScript struct.
	ErrUnsupportedScriptType = errors.New("unsupported script type")
)

// isSmallInt returns whether or not the opcode is considered a small integer,
// which is an OP_0, or OP_1 through OP_16.
func isSmallInt(op *Opcode) bool {
	if op.value == OP_0 || (op.value >= OP_1 && op.value <= OP_16) {
		return true
	}
	return false
}

// IsPayToScriptHash returns true if the script is in the standard
// pay-to-script-hash (P2SH) format, false otherwise.
func IsPayToScriptHash(script []byte) bool {
	pops, err := parseScript(script)
	if err != nil {
		return false
	}
	return isScriptHash(pops)
}

// extractPubKeyHash extracts the public key hash from the passed script if it
// is a standard pay-to-pubkey-hash script.  It will return nil otherwise.
func extractPubKeyHash(script []byte) []byte {
	// A pay-to-pubkey-hash script is of the form:
	//  OP_DUP OP_HASH160 <20-byte hash> OP_EQUALVERIFY OP_CHECKSIG
	if len(script) == 25 &&
		script[0] == OP_DUP &&
		script[1] == OP_HASH160 &&
		script[2] == OP_DATA_20 &&
		script[23] == OP_EQUALVERIFY &&
		script[24] == OP_CHECKSIG {

		return script[3:23]
	}

	return nil
}

// IsPubKeyHashScript returns whether or not the passed script is a standard
// pay-to-pubkey-hash script.
func IsPubKeyHashScript(script []byte) bool {
	return extractPubKeyHash(script) != nil
}

// IsStrictCompressedPubKeyEncoding returns whether or not the passed public
// key adheres to the strict compressed encoding requirements.
func IsStrictCompressedPubKeyEncoding(pubKey []byte) bool {
	if len(pubKey) == 33 && (pubKey[0] == 0x02 || pubKey[0] == 0x03) {
		// Compressed
		return true
	}
	return false
}

// isPushOnly returns true if the script only pushes data, false otherwise.
func isPushOnly(pops []ParsedOpcode) bool {
	// NOTE: This function does NOT verify opcodes directly since it is
	// internal and is only called with parsed opcodes for scripts that did
	// not have any parse errors.  Thus, consensus is properly maintained.

	for _, pop := range pops {
		// All opcodes up to OP_16 are data push instructions.
		// NOTE: This does consider OP_RESERVED to be a data push
		// instruction, but execution of OP_RESERVED will fail anyways
		// and matches the behavior required by consensus.
		if pop.opcode.value > OP_16 {
			return false
		}
	}
	return true
}

// IsPushOnlyScript returns whether or not the passed script only pushes data.
//
// False will be returned when the script does not parse.
func IsPushOnlyScript(script []byte) bool {
	pops, err := parseScript(script)
	if err != nil {
		return false
	}
	return isPushOnly(pops)
}

// HasP2SHScriptSigStakeOpCodes returns an error is the p2sh script has either
// stake opcodes or if the pkscript cannot be retrieved.
func HasP2SHScriptSigStakeOpCodes(version uint16, scriptSig,
	scriptPubKey []byte) error {
	class := GetScriptClass(version, scriptPubKey)
	if IsStakeOutput(scriptPubKey) {
		class, _ = GetStakeOutSubclass(scriptPubKey)
	}
	if class == ScriptHashTy {
		// Obtain the embedded pkScript from the scriptSig of the
		// current transaction. Then, ensure that it does not use
		// any stake tagging OP codes.
		pData, err := PushedData(scriptSig)
		if err != nil {
			return fmt.Errorf("error retrieving pushed data "+
				"from script: %v", err)
		}
		if len(pData) == 0 {
			return fmt.Errorf("script has no pushed data")
		}

		// The pay-to-hash-script is the final data push of the
		// signature script.
		shScript := pData[len(pData)-1]

		hasStakeOpCodes, err := ContainsStakeOpCodes(shScript)
		if err != nil {
			return fmt.Errorf("unexpected error checking pkscript "+
				"from p2sh transaction: %v", err.Error())
		}
		if hasStakeOpCodes {
			return ErrP2SHStakeOpCodes
		}
	}

	return nil
}

// parseScriptTemplate is the same as parseScript but allows the passing of the
// template list for testing purposes.  When there are parse errors, it returns
// the list of parsed opcodes up to the point of failure along with the error.
func parseScriptTemplate(script []byte, opcodes *[256]Opcode) ([]ParsedOpcode,
	error) {
	//log.Trace("Parsing script","script",fmt.Sprintf("%x",script))
	retScript := make([]ParsedOpcode, 0, len(script))
	for i := 0; i < len(script); {
		instr := script[i]
		op := &opcodes[instr]
		pop := ParsedOpcode{opcode: op}

		// Parse data out of instruction.
		switch {
		// No additional data.  Note that some of the opcodes, notably
		// OP_1NEGATE, OP_0, and OP_[1-16] represent the data
		// themselves.
		case op.length == 1:
			i++

		// Data pushes of specific lengths -- OP_DATA_[1-75].
		case op.length > 1:
			if len(script[i:]) < op.length {
				return retScript, ErrStackShortScript
			}

			// Slice out the data.
			pop.data = script[i+1 : i+op.length]
			i += op.length

		// Data pushes with parsed lengths -- OP_PUSHDATAP{1,2,4}.
		case op.length < 0:
			var l uint
			off := i + 1

			if len(script[off:]) < -op.length {
				return retScript, ErrStackShortScript
			}

			// Next -length bytes are little endian length of data.
			switch op.length {
			case -1:
				l = uint(script[off])
			case -2:
				l = ((uint(script[off+1]) << 8) |
					uint(script[off]))
			case -4:
				l = ((uint(script[off+3]) << 24) |
					(uint(script[off+2]) << 16) |
					(uint(script[off+1]) << 8) |
					uint(script[off]))
			default:
				return retScript,
					fmt.Errorf("invalid opcode length %d",
						op.length)
			}

			// Move offset to beginning of the data.
			off += -op.length

			// Disallow entries that do not fit script or were
			// sign extended.
			if int(l) > len(script[off:]) || int(l) < 0 {
				return retScript, ErrStackShortScript
			}

			pop.data = script[off : off+int(l)]
			i += 1 - op.length + int(l)
		}

		retScript = append(retScript, pop)
	}
	//log.Trace("Parsing script","result",fmt.Sprintf("%x",retScript))
	return retScript, nil
}

// parseScript preparses the script in bytes into a list of ParsedOpcodes while
// applying a number of sanity checks.
func parseScript(script []byte) ([]ParsedOpcode, error) {
	return parseScriptTemplate(script, &opcodeArray)
}

// unparseScript reversed the action of parseScript and returns the
// ParsedOpcodes as a list of bytes
func unparseScript(pops []ParsedOpcode) ([]byte, error) {
	script := make([]byte, 0, len(pops))
	for _, pop := range pops {
		b, err := pop.bytes()
		if err != nil {
			return nil, err
		}
		script = append(script, b...)
	}
	return script, nil
}

// DisasmString formats a disassembled script for one line printing.  When the
// script fails to parse, the returned string will contain the disassembled
// script up to the point the failure occurred along with the string '[error]'
// appended.  In addition, the reason the script failed to parse is returned
// if the caller wants more information about the failure.
func DisasmString(buf []byte) (string, error) {
	var disbuf bytes.Buffer
	opcodes, err := parseScript(buf)
	for _, pop := range opcodes {
		disbuf.WriteString(pop.print(true))
		disbuf.WriteByte(' ')
	}
	if disbuf.Len() > 0 {
		disbuf.Truncate(disbuf.Len() - 1)
	}
	if err != nil {
		disbuf.WriteString("[error]")
	}
	return disbuf.String(), err
}

// removeOpcode will remove any opcode matching ``opcode'' from the opcode
// stream in pkscript
func removeOpcode(pkscript []ParsedOpcode, opcode byte) []ParsedOpcode {
	retScript := make([]ParsedOpcode, 0, len(pkscript))
	for _, pop := range pkscript {
		if pop.opcode.value != opcode {
			retScript = append(retScript, pop)
		}
	}
	return retScript
}

// canonicalPush returns true if the object is either not a push instruction
// or the push instruction contained wherein is matches the canonical form
// or using the smallest instruction to do the job. False otherwise.
func canonicalPush(pop ParsedOpcode) bool {
	opcode := pop.opcode.value
	data := pop.data
	dataLen := len(pop.data)
	if opcode > OP_16 {
		return true
	}

	if opcode < OP_PUSHDATA1 && opcode > OP_0 && (dataLen == 1 && data[0] <= 16) {
		return false
	}
	if opcode == OP_PUSHDATA1 && dataLen < OP_PUSHDATA1 {
		return false
	}
	if opcode == OP_PUSHDATA2 && dataLen <= 0xff {
		return false
	}
	if opcode == OP_PUSHDATA4 && dataLen <= 0xffff {
		return false
	}
	return true
}

// removeOpcodeByData will return the script minus any opcodes that would push
// the passed data to the stack.
func removeOpcodeByData(pkscript []ParsedOpcode, data []byte) []ParsedOpcode {
	retScript := make([]ParsedOpcode, 0, len(pkscript))
	for _, pop := range pkscript {
		if !canonicalPush(pop) || !bytes.Contains(pop.data, data) {
			retScript = append(retScript, pop)
		}
	}
	return retScript

}

// asSmallInt returns the passed opcode, which must be true according to
// isSmallInt(), as an integer.
func asSmallInt(op *Opcode) int {
	if op.value == OP_0 {
		return 0
	}

	return int(op.value - (OP_1 - 1))
}

// getSigOpCount is the implementation function for counting the number of
// signature operations in the script provided by pops. If precise mode is
// requested then we attempt to count the number of operations for a multisig
// op. Otherwise we use the maximum.
func getSigOpCount(pops []ParsedOpcode, precise bool) int {
	nSigs := 0
	for i, pop := range pops {
		switch pop.opcode.value {
		case OP_CHECKSIG:
			fallthrough
		case OP_CHECKSIGVERIFY:
			fallthrough
		case OP_CHECKSIGALT:
			fallthrough
		case OP_CHECKSIGALTVERIFY:
			nSigs++
		case OP_CHECKMULTISIG:
			fallthrough
		case OP_CHECKMULTISIGVERIFY:
			// If we are being precise then look for familiar
			// patterns for multisig, for now all we recognize is
			// OP_1 - OP_16 to signify the number of pubkeys.
			// Otherwise, we use the max of 20.
			if precise && i > 0 &&
				pops[i-1].opcode.value >= OP_1 &&
				pops[i-1].opcode.value <= OP_16 {
				nSigs += asSmallInt(pops[i-1].opcode)
			} else {
				nSigs += MaxPubKeysPerMultiSig
			}
		default:
			// Not a sigop.
		}
	}

	return nSigs
}

// GetSigOpCount provides a quick count of the number of signature operations
// in a script. a CHECKSIG operations counts for 1, and a CHECK_MULTISIG for 20.
// If the script fails to parse, then the count up to the point of failure is
// returned.
func GetSigOpCount(script []byte) int {
	// Don't check error since parseScript returns the parsed-up-to-error
	// list of pops.
	pops, _ := parseScript(script)
	return getSigOpCount(pops, false)
}

// GetPreciseSigOpCount returns the number of signature operations in
// scriptPubKey.  If bip16 is true then scriptSig may be searched for the
// Pay-To-Script-Hash script in order to find the precise number of signature
// operations in the transaction.  If the script fails to parse, then the count
// up to the point of failure is returned.
func GetPreciseSigOpCount(scriptSig, scriptPubKey []byte, bip16 bool) int {
	// Don't check error since parseScript returns the parsed-up-to-error
	// list of pops.
	pops, _ := parseScript(scriptPubKey)

	// Treat non P2SH transactions as normal.
	if !(bip16 && isScriptHash(pops)) {
		return getSigOpCount(pops, true)
	}

	// The public key script is a pay-to-script-hash, so parse the signature
	// script to get the final item.  Scripts that fail to fully parse count
	// as 0 signature operations.
	sigPops, err := parseScript(scriptSig)
	if err != nil {
		return 0
	}

	// The signature script must only push data to the stack for P2SH to be
	// a valid pair, so the signature operation count is 0 when that is not
	// the case.
	if !isPushOnly(sigPops) || len(sigPops) == 0 {
		return 0
	}

	// The P2SH script is the last item the signature script pushes to the
	// stack.  When the script is empty, there are no signature operations.
	shScript := sigPops[len(sigPops)-1].data
	if len(shScript) == 0 {
		return 0
	}

	// Parse the P2SH script and don't check the error since parseScript
	// returns the parsed-up-to-error list of pops and the consensus rules
	// dictate signature operations are counted up to the first parse
	// failure.
	shPops, _ := parseScript(shScript)
	return getSigOpCount(shPops, true)
}

// IsUnspendable returns whether the passed public key script is unspendable, or
// guaranteed to fail at execution.  This allows inputs to be pruned instantly
// when entering the UTXO set.
// TODO, refactor the output spendable
func IsUnspendable(pkScript []byte) bool {
	pops, err := parseScript(pkScript)
	if err != nil {
		return true
	}

	return len(pops) > 0 && pops[0].opcode.value == OP_RETURN
}

// PkScript is a wrapper struct around a byte array, allowing it to be used
// as a map index.
type PkScript struct {
	// class is the type of the script encoded within the byte array. This
	// is used to determine the correct length of the script within the byte
	// array.
	class ScriptClass

	// script is the script contained within a byte array. If the script is
	// smaller than the length of the byte array, it will be padded with 0s
	// at the end.
	script [maxLen]byte
}

// Class returns the script type.
func (s PkScript) Class() ScriptClass {
	return s.class
}

// Script returns the script as a byte slice without any padding.
func (s PkScript) Script() []byte {
	var script []byte

	switch s.class {
	case PubKeyHashTy:
		script = make([]byte, pubKeyHashLen)
		copy(script, s.script[:pubKeyHashLen])

	case ScriptHashTy:
		script = make([]byte, scriptHashLen)
		copy(script, s.script[:scriptHashLen])

	default:
		// Unsupported script type.
		return nil
	}

	return script
}

// Address encodes the script into an address for the given chain.
func (s PkScript) Address(chainParams *params.Params) (types.Address, error) {
	_, addrs, _, err := ExtractPkScriptAddrs(s.Script(), chainParams)
	if err != nil {
		return nil, fmt.Errorf("unable to parse address: %v", err)
	}

	return addrs[0], nil
}

// String returns a hex-encoded string representation of the script.
func (s PkScript) String() string {
	str, _ := DisasmString(s.Script())
	return str
}

// ComputePkScript computes the script of an output by looking at the spending
// input's signature script or witness.
//
// NOTE: Only P2PKH, P2SH, P2WSH, and P2WPKH redeem scripts are supported.
func ComputePkScript(sigScript []byte) (PkScript, error) {
	switch {
	case len(sigScript) > 0:
		return computeNonWitnessPkScript(sigScript)
	default:
		return PkScript{}, ErrUnsupportedScriptType
	}
}

// computeNonWitnessPkScript computes the script of an output by looking at the
// spending input's signature script.
func computeNonWitnessPkScript(sigScript []byte) (PkScript, error) {
	switch {
	// Since we only support P2PKH and P2SH scripts as the only non-witness
	// script types, we should expect to see a push only script.
	case !IsPushOnlyScript(sigScript):
		return PkScript{}, ErrUnsupportedScriptType

	// If a signature script is provided with a length long enough to
	// represent a P2PKH script, then we'll attempt to parse the compressed
	// public key from it.
	case len(sigScript) >= minPubKeyHashSigScriptLen &&
		len(sigScript) <= maxPubKeyHashSigScriptLen:

		// The public key should be found as the last part of the
		// signature script. We'll attempt to parse it to ensure this is
		// a P2PKH redeem script.
		pubKey := sigScript[len(sigScript)-compressedPubKeyLen:]
		if IsStrictCompressedPubKeyEncoding(pubKey) {
			pubKeyHash := hash.Hash160(pubKey)
			script, err := payToPubKeyHashScript(pubKeyHash)
			if err != nil {
				return PkScript{}, err
			}

			pkScript := PkScript{class: PubKeyHashTy}
			copy(pkScript.script[:], script)
			return pkScript, nil
		}

		fallthrough

	// If we failed to parse a compressed public key from the script in the
	// case above, or if the script length is not that of a P2PKH one, we
	// can assume it's a P2SH signature script.
	default:
		// The redeem script will always be the last data push of the
		// signature script, so we'll parse the script into opcodes to
		// obtain it.
		parsedOpcodes, err := parseScript(sigScript)
		if err != nil {
			return PkScript{}, err
		}
		redeemScript := parsedOpcodes[len(parsedOpcodes)-1].data

		scriptHash := hash.Hash160(redeemScript)
		script, err := payToScriptHashScript(scriptHash)
		if err != nil {
			return PkScript{}, err
		}

		pkScript := PkScript{class: ScriptHashTy}
		copy(pkScript.script[:], script)
		return pkScript, nil
	}
}
