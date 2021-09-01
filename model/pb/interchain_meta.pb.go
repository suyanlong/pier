// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: interchain_meta.proto

package pb

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_meshplus_bitxhub_kit_types "github.com/meshplus/bitxhub-kit/types"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type InterchainMetaS struct {
	Counter *StringVerifiedIndexSliceMap                 `protobuf:"bytes,1,opt,name=counter,proto3" json:"counter,omitempty"`
	L2Roots []github_com_meshplus_bitxhub_kit_types.Hash `protobuf:"bytes,2,rep,name=l2Roots,proto3,customtype=github.com/meshplus/bitxhub-kit/types.Hash" json:"l2Roots,omitempty"`
}

func (m *InterchainMetaS) Reset()         { *m = InterchainMetaS{} }
func (m *InterchainMetaS) String() string { return proto.CompactTextString(m) }
func (*InterchainMetaS) ProtoMessage()    {}
func (*InterchainMetaS) Descriptor() ([]byte, []int) {
	return fileDescriptor_d2602fcc2fbd6b2a, []int{0}
}
func (m *InterchainMetaS) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InterchainMetaS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_InterchainMetaS.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *InterchainMetaS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InterchainMetaS.Merge(m, src)
}
func (m *InterchainMetaS) XXX_Size() int {
	return m.Size()
}
func (m *InterchainMetaS) XXX_DiscardUnknown() {
	xxx_messageInfo_InterchainMetaS.DiscardUnknown(m)
}

var xxx_messageInfo_InterchainMetaS proto.InternalMessageInfo

func (m *InterchainMetaS) GetCounter() *StringVerifiedIndexSliceMap {
	if m != nil {
		return m.Counter
	}
	return nil
}

type InterchainS struct {
	ID                   string           `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	InterchainCounter    *StringUint64Map `protobuf:"bytes,2,opt,name=InterchainCounter,proto3" json:"InterchainCounter,omitempty"`
	ReceiptCounter       *StringUint64Map `protobuf:"bytes,3,opt,name=ReceiptCounter,proto3" json:"ReceiptCounter,omitempty"`
	SourceReceiptCounter *StringUint64Map `protobuf:"bytes,4,opt,name=SourceReceiptCounter,proto3" json:"SourceReceiptCounter,omitempty"`
}

func (m *InterchainS) Reset()         { *m = InterchainS{} }
func (m *InterchainS) String() string { return proto.CompactTextString(m) }
func (*InterchainS) ProtoMessage()    {}
func (*InterchainS) Descriptor() ([]byte, []int) {
	return fileDescriptor_d2602fcc2fbd6b2a, []int{1}
}
func (m *InterchainS) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InterchainS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_InterchainS.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *InterchainS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InterchainS.Merge(m, src)
}
func (m *InterchainS) XXX_Size() int {
	return m.Size()
}
func (m *InterchainS) XXX_DiscardUnknown() {
	xxx_messageInfo_InterchainS.DiscardUnknown(m)
}

var xxx_messageInfo_InterchainS proto.InternalMessageInfo

func (m *InterchainS) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *InterchainS) GetInterchainCounter() *StringUint64Map {
	if m != nil {
		return m.InterchainCounter
	}
	return nil
}

func (m *InterchainS) GetReceiptCounter() *StringUint64Map {
	if m != nil {
		return m.ReceiptCounter
	}
	return nil
}

func (m *InterchainS) GetSourceReceiptCounter() *StringUint64Map {
	if m != nil {
		return m.SourceReceiptCounter
	}
	return nil
}

func init() {
	proto.RegisterType((*InterchainMetaS)(nil), "pb.InterchainMetaS")
	proto.RegisterType((*InterchainS)(nil), "pb.InterchainS")
}

func init() { proto.RegisterFile("interchain_meta.proto", fileDescriptor_d2602fcc2fbd6b2a) }

var fileDescriptor_d2602fcc2fbd6b2a = []byte{
	// 328 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0x41, 0x4b, 0x32, 0x41,
	0x1c, 0xc6, 0x9d, 0xf5, 0xe5, 0x95, 0xc6, 0x30, 0x9a, 0x0a, 0xc4, 0xc3, 0x2a, 0x9e, 0x24, 0x70,
	0x17, 0x2c, 0x82, 0xe8, 0x94, 0x09, 0xb9, 0x07, 0x2f, 0xbb, 0xd4, 0x35, 0x76, 0xc6, 0x71, 0x77,
	0x48, 0x77, 0x86, 0xdd, 0xff, 0x80, 0x7d, 0x8b, 0x4e, 0x7d, 0xa6, 0x8e, 0x1e, 0xa3, 0x43, 0x84,
	0x42, 0x9f, 0x23, 0x1c, 0xb3, 0x0d, 0xcb, 0xdb, 0x3e, 0xfb, 0x7f, 0x7e, 0x0f, 0x3f, 0x18, 0x7c,
	0x24, 0x12, 0xe0, 0x29, 0x8b, 0x43, 0x91, 0xdc, 0x4d, 0x38, 0x84, 0x8e, 0x4a, 0x25, 0x48, 0x62,
	0x29, 0x5a, 0x6b, 0x47, 0x02, 0x62, 0x4d, 0x1d, 0x26, 0x27, 0x6e, 0x24, 0x23, 0xe9, 0x9a, 0x13,
	0xd5, 0x23, 0x93, 0x4c, 0x30, 0x5f, 0x2b, 0xa4, 0x56, 0xa6, 0x61, 0x26, 0xd8, 0x2a, 0x34, 0x9f,
	0x10, 0xde, 0xf3, 0xbe, 0x97, 0x07, 0x1c, 0xc2, 0x80, 0x9c, 0xe3, 0x12, 0x93, 0x7a, 0xf9, 0xb3,
	0x8a, 0x1a, 0xa8, 0x55, 0xee, 0xd4, 0x1d, 0x45, 0x9d, 0x00, 0x52, 0x91, 0x44, 0xb7, 0x3c, 0x15,
	0x23, 0xc1, 0x87, 0x5e, 0x32, 0xe4, 0xd3, 0x60, 0x2c, 0x18, 0x1f, 0x84, 0xca, 0x5f, 0xf7, 0x49,
	0x1f, 0x97, 0xc6, 0x1d, 0x5f, 0x4a, 0xc8, 0xaa, 0x56, 0xa3, 0xd8, 0xda, 0xed, 0x3a, 0xaf, 0x6f,
	0xf5, 0xe3, 0x1f, 0x7e, 0x13, 0x9e, 0xc5, 0x6a, 0xac, 0x33, 0x97, 0x0a, 0x98, 0xc6, 0x9a, 0xb6,
	0xef, 0x05, 0xb8, 0xf0, 0xa0, 0x78, 0xe6, 0xf4, 0xc3, 0x2c, 0xf6, 0xd7, 0x78, 0xf3, 0x03, 0xe1,
	0x72, 0x2e, 0x16, 0x90, 0x0a, 0xb6, 0xbc, 0x9e, 0xf1, 0xd9, 0xf1, 0x2d, 0xaf, 0x47, 0x2e, 0xf1,
	0x7e, 0x7e, 0xbe, 0xfa, 0xd2, 0xb5, 0x8c, 0xee, 0x41, 0xae, 0x7b, 0x23, 0x12, 0x38, 0x3b, 0x5d,
	0x2a, 0xfe, 0x6e, 0x93, 0x0b, 0x5c, 0xf1, 0x39, 0xe3, 0x42, 0xc1, 0x9a, 0x2f, 0x6e, 0xe7, 0x37,
	0xaa, 0xe4, 0x1a, 0x1f, 0x06, 0x52, 0xa7, 0x8c, 0x6f, 0x4c, 0xfc, 0xdb, 0x3e, 0xf1, 0x27, 0xd0,
	0xad, 0x3e, 0xcf, 0x6d, 0x34, 0x9b, 0xdb, 0xe8, 0x7d, 0x6e, 0xa3, 0xc7, 0x85, 0x5d, 0x98, 0x2d,
	0xec, 0xc2, 0xcb, 0xc2, 0x2e, 0xd0, 0xff, 0xe6, 0x89, 0x4e, 0x3e, 0x03, 0x00, 0x00, 0xff, 0xff,
	0xe2, 0xb0, 0x3a, 0x6f, 0xfb, 0x01, 0x00, 0x00,
}

func (m *InterchainMetaS) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InterchainMetaS) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *InterchainMetaS) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.L2Roots) > 0 {
		for iNdEx := len(m.L2Roots) - 1; iNdEx >= 0; iNdEx-- {
			{
				size := m.L2Roots[iNdEx].Size()
				i -= size
				if _, err := m.L2Roots[iNdEx].MarshalTo(dAtA[i:]); err != nil {
					return 0, err
				}
				i = encodeVarintInterchainMeta(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Counter != nil {
		{
			size, err := m.Counter.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintInterchainMeta(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *InterchainS) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InterchainS) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *InterchainS) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SourceReceiptCounter != nil {
		{
			size, err := m.SourceReceiptCounter.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintInterchainMeta(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.ReceiptCounter != nil {
		{
			size, err := m.ReceiptCounter.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintInterchainMeta(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.InterchainCounter != nil {
		{
			size, err := m.InterchainCounter.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintInterchainMeta(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.ID) > 0 {
		i -= len(m.ID)
		copy(dAtA[i:], m.ID)
		i = encodeVarintInterchainMeta(dAtA, i, uint64(len(m.ID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintInterchainMeta(dAtA []byte, offset int, v uint64) int {
	offset -= sovInterchainMeta(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *InterchainMetaS) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Counter != nil {
		l = m.Counter.Size()
		n += 1 + l + sovInterchainMeta(uint64(l))
	}
	if len(m.L2Roots) > 0 {
		for _, e := range m.L2Roots {
			l = e.Size()
			n += 1 + l + sovInterchainMeta(uint64(l))
		}
	}
	return n
}

func (m *InterchainS) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ID)
	if l > 0 {
		n += 1 + l + sovInterchainMeta(uint64(l))
	}
	if m.InterchainCounter != nil {
		l = m.InterchainCounter.Size()
		n += 1 + l + sovInterchainMeta(uint64(l))
	}
	if m.ReceiptCounter != nil {
		l = m.ReceiptCounter.Size()
		n += 1 + l + sovInterchainMeta(uint64(l))
	}
	if m.SourceReceiptCounter != nil {
		l = m.SourceReceiptCounter.Size()
		n += 1 + l + sovInterchainMeta(uint64(l))
	}
	return n
}

func sovInterchainMeta(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozInterchainMeta(x uint64) (n int) {
	return sovInterchainMeta(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *InterchainMetaS) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInterchainMeta
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
			return fmt.Errorf("proto: InterchainMetaS: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InterchainMetaS: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Counter", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInterchainMeta
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthInterchainMeta
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthInterchainMeta
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Counter == nil {
				m.Counter = &StringVerifiedIndexSliceMap{}
			}
			if err := m.Counter.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field L2Roots", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInterchainMeta
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
				return ErrInvalidLengthInterchainMeta
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthInterchainMeta
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_meshplus_bitxhub_kit_types.Hash
			m.L2Roots = append(m.L2Roots, v)
			if err := m.L2Roots[len(m.L2Roots)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipInterchainMeta(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthInterchainMeta
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthInterchainMeta
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *InterchainS) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInterchainMeta
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
			return fmt.Errorf("proto: InterchainS: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InterchainS: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInterchainMeta
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInterchainMeta
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInterchainMeta
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterchainCounter", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInterchainMeta
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthInterchainMeta
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthInterchainMeta
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.InterchainCounter == nil {
				m.InterchainCounter = &StringUint64Map{}
			}
			if err := m.InterchainCounter.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReceiptCounter", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInterchainMeta
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthInterchainMeta
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthInterchainMeta
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ReceiptCounter == nil {
				m.ReceiptCounter = &StringUint64Map{}
			}
			if err := m.ReceiptCounter.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SourceReceiptCounter", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInterchainMeta
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthInterchainMeta
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthInterchainMeta
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SourceReceiptCounter == nil {
				m.SourceReceiptCounter = &StringUint64Map{}
			}
			if err := m.SourceReceiptCounter.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipInterchainMeta(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthInterchainMeta
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthInterchainMeta
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipInterchainMeta(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowInterchainMeta
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
					return 0, ErrIntOverflowInterchainMeta
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
					return 0, ErrIntOverflowInterchainMeta
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
				return 0, ErrInvalidLengthInterchainMeta
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupInterchainMeta
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthInterchainMeta
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthInterchainMeta        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowInterchainMeta          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupInterchainMeta = fmt.Errorf("proto: unexpected end of group")
)