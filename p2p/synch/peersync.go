/*
 * Copyright (c) 2017-2020 The qitmeer developers
 */

package synch

import (
	"fmt"
	chainhash "github.com/Qitmeer/qitmeer/common/hash"
	"github.com/Qitmeer/qitmeer/core/blockchain"
	"github.com/Qitmeer/qitmeer/core/blockdag"
	"github.com/Qitmeer/qitmeer/core/protocol"
	"github.com/Qitmeer/qitmeer/core/types"
	"github.com/Qitmeer/qitmeer/p2p/peers"
	pb "github.com/Qitmeer/qitmeer/p2p/proto/v1"
	"github.com/btcsuite/btcd/wire"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// stallSampleInterval the interval at which we will check to see if our
	// sync has stalled.
	stallSampleInterval = 300 * time.Second
)

type PeerSync struct {
	sy *Sync

	splock   sync.RWMutex
	syncPeer *peers.Peer
	// dag sync
	dagSync *blockdag.DAGSync

	hslock sync.RWMutex

	started     int32
	shutdown    int32
	msgChan     chan interface{}
	wg          sync.WaitGroup
	quit        chan struct{}
	longSyncMod bool
}

func (ps *PeerSync) Start() error {
	// Already started?
	if atomic.AddInt32(&ps.started, 1) != 1 {
		return nil
	}

	log.Info("P2P PeerSync Start")
	ps.dagSync = blockdag.NewDAGSync(ps.sy.p2p.BlockChain().BlockDAG())
	ps.longSyncMod = false

	ps.wg.Add(1)
	go ps.handler()
	return nil
}

func (ps *PeerSync) Stop() error {
	if atomic.AddInt32(&ps.shutdown, 1) != 1 {
		log.Warn("PeerSync is already in the process of shutting down")
		return nil
	}
	log.Info("P2P PeerSync Stop")

	close(ps.quit)
	ps.wg.Wait()

	return nil
}

func (ps *PeerSync) handler() {
	stallTicker := time.NewTicker(stallSampleInterval)
	defer stallTicker.Stop()

out:
	for {
		select {
		case m := <-ps.msgChan:
			switch msg := m.(type) {
			case pauseMsg:
				// Wait until the sender unpauses the manager.
				<-msg.unpause

			case *ConnectedMsg:
				ps.processConnected(msg)

			case *DisconnectedMsg:
				ps.processDisconnected(msg)
			case *GetBlocksMsg:
				err := ps.processGetBlocks(msg.pe, msg.blocks)
				if err != nil {
					log.Warn(err.Error())
				}
			case *GetBlockDatasMsg:
				err := ps.processGetBlockDatas(msg.pe, msg.blocks)
				if err != nil {
					go ps.PeerUpdate(msg.pe, false)
				}
			case *MsgGetData:
				ps.OnGetData(msg.pe, msg.InvMsg)

			case *UpdateGraphStateMsg:
				ps.processUpdateGraphState(msg.pe)
			case *syncDAGBlocksMsg:
				err := ps.processSyncDAGBlocks(msg.pe)
				if err != nil {
					log.Warn(err.Error())
				}
			case *PeerUpdateMsg:
				ps.OnPeerUpdate(msg.pe, msg.orphan)
			case *getTxsMsg:
				err := ps.processGetTxs(msg.pe, msg.txs)
				if err != nil {
					log.Warn(err.Error())
				}
			case *SyncQNRMsg:
				err := ps.processQNR(msg)
				if err != nil {
					log.Warn(err.Error())
				}
			default:
				log.Warn(fmt.Sprintf("Invalid message type in task "+
					"handler: %T", msg))
			}

		case <-stallTicker.C:
			ps.handleStallSample()

		case <-ps.quit:
			break out
		}
	}

	// Drain any wait channels before going away so there is nothing left
	// waiting on this goroutine.
cleanup:
	for {
		select {
		case <-ps.msgChan:
		default:
			break cleanup
		}
	}

	ps.wg.Done()
	log.Trace("Peer Sync handler done")
}

func (ps *PeerSync) handleStallSample() {
	if atomic.LoadInt32(&ps.shutdown) != 0 {
		return
	}
}

func (ps *PeerSync) Pause() chan<- struct{} {
	c := make(chan struct{})
	ps.msgChan <- pauseMsg{c}
	return c
}

func (ps *PeerSync) SyncPeer() *peers.Peer {
	ps.splock.RLock()
	defer ps.splock.RUnlock()

	return ps.syncPeer
}

func (ps *PeerSync) SetSyncPeer(pe *peers.Peer) {
	ps.splock.Lock()
	defer ps.splock.Unlock()

	ps.syncPeer = pe
}

func (ps *PeerSync) OnPeerConnected(pe *peers.Peer) {

	ti := pe.Timestamp()
	if !ti.IsZero() {
		// Add the remote peer time as a sample for creating an offset against
		// the local clock to keep the network time in sync.
		ps.sy.p2p.TimeSource().AddTimeSample(pe.GetID().String(), ti)
	}

	if !ps.HasSyncPeer() {
		ps.startSync()
	}
}

func (ps *PeerSync) OnPeerDisconnected(pe *peers.Peer) {

	if ps.HasSyncPeer() {
		if ps.isSyncPeer(pe) {
			ps.updateSyncPeer(true)
		}
	}
}

func (ps *PeerSync) isSyncPeer(pe *peers.Peer) bool {
	if !ps.HasSyncPeer() {
		return false
	}
	if pe == ps.SyncPeer() || pe.GetID() == ps.SyncPeer().GetID() {
		return true
	}
	return false
}

func (ps *PeerSync) PeerUpdate(pe *peers.Peer, orphan bool) {
	// Ignore if we are shutting down.
	if atomic.LoadInt32(&ps.shutdown) != 0 {
		return
	}

	ps.msgChan <- &PeerUpdateMsg{pe: pe, orphan: orphan}
}

func (ps *PeerSync) OnPeerUpdate(pe *peers.Peer, orphan bool) {

	if ps.HasSyncPeer() {
		spgs := ps.SyncPeer().GraphState()
		if !ps.SyncPeer().IsActive() || spgs == nil {
			ps.updateSyncPeer(true)
			return
		}
		if pe != nil {
			pegs := pe.GraphState()
			if pegs != nil {
				if pegs.IsExcellent(spgs) {
					ps.updateSyncPeer(true)
					return
				}
			}

		}
		ps.IntellectSyncBlocks(orphan)
		return
	}
	ps.updateSyncPeer(false)
}

func (ps *PeerSync) HasSyncPeer() bool {
	return ps.SyncPeer() != nil
}

func (ps *PeerSync) Chain() *blockchain.BlockChain {
	return ps.sy.p2p.BlockChain()
}

// startSync will choose the best peer among the available candidate peers to
// download/sync the blockchain from.  When syncing is already running, it
// simply returns.  It also examines the candidates for any which are no longer
// candidates and removes them as needed.
func (ps *PeerSync) startSync() {
	// Return now if we're already syncing.
	if ps.HasSyncPeer() {
		return
	}
	best := ps.Chain().BestSnapshot()
	bestPeer := ps.getBestPeer()
	// Start syncing from the best peer if one was selected.
	if bestPeer != nil {
		gs := bestPeer.GraphState()

		log.Info(fmt.Sprintf("Syncing to state %s from peer %s cur graph state:%s", gs.String(), bestPeer.GetID().String(), best.GraphState.String()))

		// When the current height is less than a known checkpoint we
		// can use block headers to learn about which blocks comprise
		// the chain up to the checkpoint and perform less validation
		// for them.  This is possible since each header contains the
		// hash of the previous header and a merkle root.  Therefore if
		// we validate all of the received headers link together
		// properly and the checkpoint hashes match, we can be sure the
		// hashes for the blocks in between are accurate.  Further, once
		// the full blocks are downloaded, the merkle root is computed
		// and compared against the value in the header which proves the
		// full block hasn't been tampered with.
		//
		// Once we have passed the final checkpoint, or checkpoints are
		// disabled, use standard inv messages learn about the blocks
		// and fully validate them.  Finally, regression test mode does
		// not support the headers-first approach so do normal block
		// downloads when in regression test mode.

		ps.SetSyncPeer(bestPeer)
		ps.IntellectSyncBlocks(true)
		ps.dagSync.SetGraphState(gs)

	} else {
		log.Trace("You're already up to date, no synchronization is required.")
	}
}

// getBestPeer
func (ps *PeerSync) getBestPeer() *peers.Peer {
	best := ps.Chain().BestSnapshot()
	var bestPeer *peers.Peer
	equalPeers := []*peers.Peer{}
	for _, sp := range ps.sy.peers.ConnectedPeers() {
		// Remove sync candidate peers that are no longer candidates due
		// to passing their latest known block.  NOTE: The < is
		// intentional as opposed to <=.  While techcnically the peer
		// doesn't have a later block when it's equal, it will likely
		// have one soon so it is a reasonable choice.  It also allows
		// the case where both are at 0 such as during regression test.
		gs := sp.GraphState()
		if gs == nil {
			continue
		}
		if !gs.IsExcellent(best.GraphState) {
			continue
		}
		// the best sync candidate is the most updated peer
		if bestPeer == nil {
			bestPeer = sp
			continue
		}
		if gs.IsExcellent(bestPeer.GraphState()) {
			bestPeer = sp
			if len(equalPeers) > 0 {
				equalPeers = equalPeers[0:0]
			}
		} else if gs.IsEqual(bestPeer.GraphState()) {
			equalPeers = append(equalPeers, sp)
		}
	}
	if bestPeer == nil {
		return nil
	}
	if len(equalPeers) > 0 {
		for _, sp := range equalPeers {
			if sp.GetID().String() > bestPeer.GetID().String() {
				bestPeer = sp
			}
		}
	}
	return bestPeer
}

// IsCurrent returns true if we believe we are synced with our peers, false if we
// still have blocks to check
func (ps *PeerSync) IsCurrent() bool {
	if !ps.Chain().IsCurrent() {
		return false
	}

	return ps.IsCompleteForSyncPeer()
}

func (ps *PeerSync) IsCompleteForSyncPeer() bool {
	// if blockChain thinks we are current and we have no syncPeer it
	// is probably right.
	if !ps.HasSyncPeer() {
		return true
	}

	// No matter what chain thinks, if we are below the block we are syncing
	// to we are not current.
	gs := ps.SyncPeer().GraphState()
	if gs == nil {
		return true
	}
	if gs.IsExcellent(ps.Chain().BestSnapshot().GraphState) {
		//log.Trace("comparing the current best vs sync last",
		//	"current.best", ps.Chain().BestSnapshot().GraphState.String(), "sync.last", gs.String())
		return false
	}

	return true
}

func (ps *PeerSync) IntellectSyncBlocks(refresh bool) {
	if !ps.HasSyncPeer() {
		return
	}

	if ps.Chain().GetOrphansTotal() >= blockchain.MaxOrphanBlocks || refresh {
		ps.Chain().RefreshOrphans()
	}
	allOrphan := ps.Chain().GetRecentOrphansParents()

	if len(allOrphan) > 0 {
		go ps.GetBlocks(ps.SyncPeer(), allOrphan)
	} else {
		go ps.syncDAGBlocks(ps.SyncPeer())
	}
}

func (ps *PeerSync) updateSyncPeer(force bool) {
	log.Debug("Updating sync peer")
	if force {
		ps.SetSyncPeer(nil)
	}
	ps.startSync()
}

func (ps *PeerSync) RelayInventory(data interface{}) {
	ps.sy.Peers().ForPeers(peers.PeerConnected, func(pe *peers.Peer) {
		msg := &pb.Inventory{Invs: []*pb.InvVect{}}
		switch value := data.(type) {
		case []*types.TxDesc:
			// Don't relay the transaction to the peer when it has
			// transaction relaying disabled.
			if pe.DisableRelayTx() {
				return
			}
			for _, tx := range value {
				feeFilter := pe.FeeFilter()
				if feeFilter > 0 && tx.FeePerKB < feeFilter {
					return
				}
				// Don't relay the transaction if there is a bloom
				// filter loaded and the transaction doesn't match it.
				filter := pe.Filter()
				if filter.IsLoaded() {
					if !filter.MatchTxAndUpdate(tx.Tx) {
						return
					}
				}
				msg.Invs = append(msg.Invs, NewInvVect(InvTypeTx, tx.Tx.Hash()))
			}
		case types.BlockHeader:
			blockHash := value.BlockHash()
			msg.Invs = append(msg.Invs, NewInvVect(InvTypeBlock, &blockHash))
		}
		go ps.sy.sendInventoryRequest(ps.sy.p2p.Context(), pe, msg)
	})
}

// OnFilterAdd is invoked when a peer receives a filteradd qitmeer
// message and is used by remote peers to add data to an already loaded bloom
// filter.  The peer will be disconnected if a filter is not loaded when this
// message is received or the server is not configured to allow bloom filters.
func (ps *PeerSync) OnFilterAdd(sp *peers.Peer, msg *types.MsgFilterAdd) {
	// Disconnect and/or ban depending on the node bloom services flag and
	// negotiated protocol version.
	if !sp.EnforceNodeBloomFlag(msg.Command()) {
		return
	}
	filter := sp.Filter()
	if !filter.IsLoaded() {
		log.Debug(fmt.Sprintf("%s sent a filterclear request with no "+
			"filter loaded -- disconnecting", sp))
		ps.Disconnect(sp)
		return
	}

	filter.Add(msg.Data)
}

// OnFilterClear is invoked when a peer receives a filterclear qitmeer
// message and is used by remote peers to clear an already loaded bloom filter.
// The peer will be disconnected if a filter is not loaded when this message is
// received  or the server is not configured to allow bloom filters.
func (ps *PeerSync) OnFilterClear(sp *peers.Peer, msg *wire.MsgFilterClear) {
	// Disconnect and/or ban depending on the node bloom services flag and
	// negotiated protocol version.
	if !sp.EnforceNodeBloomFlag(msg.Command()) {
		return
	}
	filter := sp.Filter()

	if !filter.IsLoaded() {
		log.Debug(fmt.Sprintf("%s sent a filterclear request with no "+
			"filter loaded -- disconnecting", sp))
		ps.Disconnect(sp)
		return
	}

	filter.Unload()
}

// OnFilterLoad is invoked when a peer receives a filterload qitmeer
// message and it used to load a bloom filter that should be used for
// delivering merkle blocks and associated transactions that match the filter.
// The peer will be disconnected if the server is not configured to allow bloom
// filters.
func (ps *PeerSync) OnFilterLoad(sp *peers.Peer, msg *types.MsgFilterLoad) {
	// Disconnect and/or ban depending on the node bloom services flag and
	// negotiated protocol version.
	if !sp.EnforceNodeBloomFlag(msg.Command()) {
		return
	}
	filter := sp.Filter()
	sp.DisableRelayTx()

	filter.Reload(msg)
}

// OnMemPool is invoked when a peer receives a mempool qitmeer message.
// It creates and sends an inventory message with the contents of the memory
// pool up to the maximum inventory allowed per message.  When the peer has a
// bloom filter loaded, the contents are filtered accordingly.
func (ps *PeerSync) OnMemPool(sp *peers.Peer, msg *types.MsgMemPool) {
	// Only allow mempool requests if the server has bloom filtering
	// enabled.
	services := sp.Services()
	if services&protocol.Bloom != protocol.Bloom {
		log.Debug(fmt.Sprintf("%s sent a filterclear request with no "+
			"filter loaded -- disconnecting", sp))
		ps.Disconnect(sp)
		return
	}

	// Generate inventory message with the available transactions in the
	// transaction memory pool.  Limit it to the max allowed inventory
	// per message.  The NewMsgInvSizeHint function automatically limits
	// the passed hint to the maximum allowed, so it's safe to pass it
	// without double checking it here.
	txDescs := ps.sy.p2p.TxMemPool().TxDescs()
	invMsg := &pb.Inventory{Invs: []*pb.InvVect{}}
	for _, txDesc := range txDescs {
		// Either add all transactions when there is no bloom filter,
		// or only the transactions that match the filter when there is
		// one.
		filter := sp.Filter()
		if !filter.IsLoaded() || filter.MatchTxAndUpdate(txDesc.Tx) {
			invMsg.Invs = append(invMsg.Invs, NewInvVect(InvTypeTx, txDesc.Tx.Hash()))
		}
	}
	// Send the inventory message if there is anything to send.
	if len(invMsg.Invs) > 0 {
		go ps.sy.sendInventoryRequest(ps.sy.p2p.Context(), sp, invMsg)
	}
}

// pushTxMsg sends a tx message for the provided transaction hash to the
// connected peer.  An error is returned if the transaction hash is not known.
func (s *PeerSync) pushTxMsg(sp *serverPeer, hash *chainhash.Hash, doneChan chan<- struct{},
	waitChan <-chan struct{}, encoding wire.MessageEncoding) error {

	// Attempt to fetch the requested transaction from the pool.  A
	// call could be made to check for existence first, but simply trying
	// to fetch a missing transaction results in the same behavior.
	tx, err := s.sy.p2p.TxMemPool().FetchTransaction(hash)
	if err != nil {
		peerLog.Tracef("Unable to fetch tx %v from transaction "+
			"pool: %v", hash, err)

		if doneChan != nil {
			doneChan <- struct{}{}
		}
		return err
	}

	// Once we have fetched data wait for any previous operation to finish.
	if waitChan != nil {
		<-waitChan
	}

	sp.QueueMessageWithEncoding(tx.MsgTx(), doneChan, encoding)

	return nil
}

// pushBlockMsg sends a block message for the provided block hash to the
// connected peer.  An error is returned if the block hash is not known.
func (s *PeerSync) pushBlockMsg(sp *serverPeer, hash *chainhash.Hash, doneChan chan<- struct{},
	waitChan <-chan struct{}, encoding wire.MessageEncoding) error {

	// Fetch the raw block bytes from the database.
	var blockBytes []byte
	err := sp.server.db.View(func(dbTx database.Tx) error {
		var err error
		blockBytes, err = dbTx.FetchBlock(hash)
		return err
	})
	if err != nil {
		peerLog.Tracef("Unable to fetch requested block hash %v: %v",
			hash, err)

		if doneChan != nil {
			doneChan <- struct{}{}
		}
		return err
	}

	// Deserialize the block.
	var msgBlock wire.MsgBlock
	err = msgBlock.Deserialize(bytes.NewReader(blockBytes))
	if err != nil {
		peerLog.Tracef("Unable to deserialize requested block hash "+
			"%v: %v", hash, err)

		if doneChan != nil {
			doneChan <- struct{}{}
		}
		return err
	}

	// Once we have fetched data wait for any previous operation to finish.
	if waitChan != nil {
		<-waitChan
	}

	// We only send the channel for this message if we aren't sending
	// an inv straight after.
	var dc chan<- struct{}
	continueHash := sp.continueHash
	sendInv := continueHash != nil && continueHash.IsEqual(hash)
	if !sendInv {
		dc = doneChan
	}
	sp.QueueMessageWithEncoding(&msgBlock, dc, encoding)

	// When the peer requests the final block that was advertised in
	// response to a getblocks message which requested more blocks than
	// would fit into a single message, send it a new inventory message
	// to trigger it to issue another getblocks message for the next
	// batch of inventory.
	if sendInv {
		best := sp.server.chain.BestSnapshot()
		invMsg := wire.NewMsgInvSizeHint(1)
		iv := wire.NewInvVect(wire.InvTypeBlock, &best.Hash)
		invMsg.AddInvVect(iv)
		sp.QueueMessage(invMsg, doneChan)
		sp.continueHash = nil
	}
	return nil
}

// pushMerkleBlockMsg sends a merkleblock message for the provided block hash to
// the connected peer.  Since a merkle block requires the peer to have a filter
// loaded, this call will simply be ignored if there is no filter loaded.  An
// error is returned if the block hash is not known.
func (s *PeerSync) pushMerkleBlockMsg(sp *serverPeer, hash *chainhash.Hash,
	doneChan chan<- struct{}, waitChan <-chan struct{}, encoding wire.MessageEncoding) error {

	// Do not send a response if the peer doesn't have a filter loaded.
	if !sp.filter.IsLoaded() {
		if doneChan != nil {
			doneChan <- struct{}{}
		}
		return nil
	}

	// Fetch the raw block bytes from the database.
	blk, err := sp.server.chain.BlockByHash(hash)
	if err != nil {
		peerLog.Tracef("Unable to fetch requested block hash %v: %v",
			hash, err)

		if doneChan != nil {
			doneChan <- struct{}{}
		}
		return err
	}

	// Generate a merkle block by filtering the requested block according
	// to the filter for the peer.
	merkle, matchedTxIndices := bloom.NewMerkleBlock(blk, sp.filter)

	// Once we have fetched data wait for any previous operation to finish.
	if waitChan != nil {
		<-waitChan
	}

	// Send the merkleblock.  Only send the done channel with this message
	// if no transactions will be sent afterwards.
	var dc chan<- struct{}
	if len(matchedTxIndices) == 0 {
		dc = doneChan
	}
	sp.QueueMessage(merkle, dc)

	// Finally, send any matched transactions.
	blkTransactions := blk.MsgBlock().Transactions
	for i, txIndex := range matchedTxIndices {
		// Only send the done channel on the final transaction.
		var dc chan<- struct{}
		if i == len(matchedTxIndices)-1 {
			dc = doneChan
		}
		if txIndex < uint32(len(blkTransactions)) {
			sp.QueueMessageWithEncoding(blkTransactions[txIndex], dc,
				encoding)
		}
	}

	return nil
}

func NewPeerSync(sy *Sync) *PeerSync {
	peerSync := &PeerSync{
		sy:      sy,
		msgChan: make(chan interface{}),
		quit:    make(chan struct{}),
	}

	return peerSync
}
