// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: messages.proto

package qitmeer_p2p_v1

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
	github_com_prysmaticlabs_go_bitfield "github.com/prysmaticlabs/go-bitfield"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ErrorResponse struct {
	Message              []byte   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty" ssz-max:"256"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ErrorResponse) Reset()         { *m = ErrorResponse{} }
func (m *ErrorResponse) String() string { return proto.CompactTextString(m) }
func (*ErrorResponse) ProtoMessage()    {}
func (*ErrorResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{0}
}
func (m *ErrorResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ErrorResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ErrorResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ErrorResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrorResponse.Merge(m, src)
}
func (m *ErrorResponse) XXX_Size() int {
	return m.Size()
}
func (m *ErrorResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrorResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ErrorResponse proto.InternalMessageInfo

func (m *ErrorResponse) GetMessage() []byte {
	if m != nil {
		return m.Message
	}
	return nil
}

type MetaData struct {
	SeqNumber            uint64                                           `protobuf:"varint,1,opt,name=seq_number,json=seqNumber,proto3" json:"seq_number,omitempty"`
	Attnets              github_com_prysmaticlabs_go_bitfield.Bitvector64 `protobuf:"bytes,2,opt,name=attnets,proto3,casttype=github.com/prysmaticlabs/go-bitfield.Bitvector64" json:"attnets,omitempty" ssz-size:"8"`
	XXX_NoUnkeyedLiteral struct{}                                         `json:"-"`
	XXX_unrecognized     []byte                                           `json:"-"`
	XXX_sizecache        int32                                            `json:"-"`
}

func (m *MetaData) Reset()         { *m = MetaData{} }
func (m *MetaData) String() string { return proto.CompactTextString(m) }
func (*MetaData) ProtoMessage()    {}
func (*MetaData) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{1}
}
func (m *MetaData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MetaData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MetaData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MetaData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetaData.Merge(m, src)
}
func (m *MetaData) XXX_Size() int {
	return m.Size()
}
func (m *MetaData) XXX_DiscardUnknown() {
	xxx_messageInfo_MetaData.DiscardUnknown(m)
}

var xxx_messageInfo_MetaData proto.InternalMessageInfo

func (m *MetaData) GetSeqNumber() uint64 {
	if m != nil {
		return m.SeqNumber
	}
	return 0
}

func (m *MetaData) GetAttnets() github_com_prysmaticlabs_go_bitfield.Bitvector64 {
	if m != nil {
		return m.Attnets
	}
	return nil
}

func init() {
	proto.RegisterType((*ErrorResponse)(nil), "qitmeer.p2p.v1.ErrorResponse")
	proto.RegisterType((*MetaData)(nil), "qitmeer.p2p.v1.MetaData")
}

func init() { proto.RegisterFile("messages.proto", fileDescriptor_4dc296cbfe5ffcd5) }

var fileDescriptor_4dc296cbfe5ffcd5 = []byte{
	// 277 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcb, 0x4d, 0x2d, 0x2e,
	0x4e, 0x4c, 0x4f, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2b, 0xcc, 0x2c, 0xc9,
	0x4d, 0x4d, 0x2d, 0xd2, 0x2b, 0x30, 0x2a, 0xd0, 0x2b, 0x33, 0x94, 0xd2, 0x4d, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xcf, 0x4f, 0xcf, 0xd7, 0x07, 0x2b, 0x4b, 0x2a,
	0x4d, 0x03, 0xf3, 0xc0, 0x1c, 0x30, 0x0b, 0xa2, 0x5d, 0xc9, 0x86, 0x8b, 0xd7, 0xb5, 0xa8, 0x28,
	0xbf, 0x28, 0x28, 0xb5, 0xb8, 0x20, 0x3f, 0xaf, 0x38, 0x55, 0x48, 0x9b, 0x8b, 0x1d, 0x6a, 0x83,
	0x04, 0xa3, 0x02, 0xa3, 0x06, 0x8f, 0x93, 0xe0, 0xa7, 0x7b, 0xf2, 0xbc, 0xc5, 0xc5, 0x55, 0xba,
	0xb9, 0x89, 0x15, 0x56, 0x4a, 0x46, 0xa6, 0x66, 0x4a, 0x41, 0x30, 0x15, 0x4a, 0x9d, 0x8c, 0x5c,
	0x1c, 0xbe, 0xa9, 0x25, 0x89, 0x2e, 0x89, 0x25, 0x89, 0x42, 0xb2, 0x5c, 0x5c, 0xc5, 0xa9, 0x85,
	0xf1, 0x79, 0xa5, 0xb9, 0x49, 0xa9, 0x45, 0x60, 0xcd, 0x2c, 0x41, 0x9c, 0xc5, 0xa9, 0x85, 0x7e,
	0x60, 0x01, 0xa1, 0x38, 0x2e, 0xf6, 0xc4, 0x92, 0x92, 0xbc, 0xd4, 0x92, 0x62, 0x09, 0x26, 0xb0,
	0xc1, 0x2e, 0x9f, 0xee, 0xc9, 0xf3, 0x80, 0x0c, 0x2e, 0xce, 0xac, 0x4a, 0xb5, 0x52, 0xb2, 0x50,
	0xfa, 0x75, 0x4f, 0xde, 0x00, 0xc9, 0xf5, 0x05, 0x45, 0x95, 0xc5, 0xb9, 0x89, 0x25, 0x99, 0xc9,
	0x39, 0x89, 0x49, 0xc5, 0xfa, 0xe9, 0xf9, 0xba, 0x49, 0x99, 0x25, 0x69, 0x99, 0xa9, 0x39, 0x29,
	0x7a, 0x4e, 0x99, 0x25, 0x65, 0xa9, 0xc9, 0x25, 0xf9, 0x45, 0x66, 0x26, 0x41, 0x30, 0x43, 0x9d,
	0x04, 0x4e, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x19, 0x8f,
	0xe5, 0x18, 0x92, 0xd8, 0xc0, 0x5e, 0x34, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xdd, 0xac, 0xf3,
	0x16, 0x33, 0x01, 0x00, 0x00,
}

func (m *ErrorResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ErrorResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ErrorResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Message) > 0 {
		i -= len(m.Message)
		copy(dAtA[i:], m.Message)
		i = encodeVarintMessages(dAtA, i, uint64(len(m.Message)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MetaData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MetaData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MetaData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Attnets) > 0 {
		i -= len(m.Attnets)
		copy(dAtA[i:], m.Attnets)
		i = encodeVarintMessages(dAtA, i, uint64(len(m.Attnets)))
		i--
		dAtA[i] = 0x12
	}
	if m.SeqNumber != 0 {
		i = encodeVarintMessages(dAtA, i, uint64(m.SeqNumber))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintMessages(dAtA []byte, offset int, v uint64) int {
	offset -= sovMessages(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ErrorResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovMessages(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *MetaData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SeqNumber != 0 {
		n += 1 + sovMessages(uint64(m.SeqNumber))
	}
	l = len(m.Attnets)
	if l > 0 {
		n += 1 + l + sovMessages(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovMessages(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMessages(x uint64) (n int) {
	return sovMessages(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ErrorResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ErrorResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ErrorResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = append(m.Message[:0], dAtA[iNdEx:postIndex]...)
			if m.Message == nil {
				m.Message = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MetaData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MetaData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MetaData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SeqNumber", wireType)
			}
			m.SeqNumber = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SeqNumber |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Attnets", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Attnets = append(m.Attnets[:0], dAtA[iNdEx:postIndex]...)
			if m.Attnets == nil {
				m.Attnets = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipMessages(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthMessages
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMessages
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMessages
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMessages        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMessages          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMessages = fmt.Errorf("proto: unexpected end of group")
)
