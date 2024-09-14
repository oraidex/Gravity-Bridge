// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: bech32ibc/v1beta1/types.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	_ "github.com/golang/protobuf/ptypes/duration"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// An HrpIbcRecord maps a bech32 human-readable prefix to an IBC source
// channel
type HrpIbcRecord struct {
	// The bech32 human readable prefix that serves as the key
	Hrp string `protobuf:"bytes,1,opt,name=hrp,proto3" json:"hrp,omitempty" yaml:"hrp"`
	// the channel by which the packet will be sent
	SourceChannel     string        `protobuf:"bytes,2,opt,name=source_channel,json=sourceChannel,proto3" json:"source_channel,omitempty" yaml:"source_channel"`
	IcsToHeightOffset uint64        `protobuf:"varint,3,opt,name=ics_to_height_offset,json=icsToHeightOffset,proto3" json:"ics_to_height_offset,omitempty"`
	IcsToTimeOffset   time.Duration `protobuf:"bytes,4,opt,name=ics_to_time_offset,json=icsToTimeOffset,proto3,stdduration" json:"ics_to_time_offset,omitempty" yaml:"ics_to_time_offset"`
}

func (m *HrpIbcRecord) Reset()         { *m = HrpIbcRecord{} }
func (m *HrpIbcRecord) String() string { return proto.CompactTextString(m) }
func (*HrpIbcRecord) ProtoMessage()    {}
func (*HrpIbcRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_95408f5fe13d0f7d, []int{0}
}
func (m *HrpIbcRecord) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HrpIbcRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HrpIbcRecord.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HrpIbcRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HrpIbcRecord.Merge(m, src)
}
func (m *HrpIbcRecord) XXX_Size() int {
	return m.Size()
}
func (m *HrpIbcRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_HrpIbcRecord.DiscardUnknown(m)
}

var xxx_messageInfo_HrpIbcRecord proto.InternalMessageInfo

func (m *HrpIbcRecord) GetHrp() string {
	if m != nil {
		return m.Hrp
	}
	return ""
}

func (m *HrpIbcRecord) GetSourceChannel() string {
	if m != nil {
		return m.SourceChannel
	}
	return ""
}

func (m *HrpIbcRecord) GetIcsToHeightOffset() uint64 {
	if m != nil {
		return m.IcsToHeightOffset
	}
	return 0
}

func (m *HrpIbcRecord) GetIcsToTimeOffset() time.Duration {
	if m != nil {
		return m.IcsToTimeOffset
	}
	return 0
}

func init() {
	proto.RegisterType((*HrpIbcRecord)(nil), "bech32ibc.bech32ibc.v1beta1.HrpIbcRecord")
}

func init() { proto.RegisterFile("bech32ibc/v1beta1/types.proto", fileDescriptor_95408f5fe13d0f7d) }

var fileDescriptor_95408f5fe13d0f7d = []byte{
	// 381 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0xc1, 0x8a, 0xda, 0x40,
	0x18, 0xc7, 0x33, 0x2a, 0x85, 0xa6, 0xad, 0xa5, 0xc1, 0x82, 0xda, 0x36, 0x09, 0x39, 0x79, 0xa8,
	0x19, 0xd4, 0x9e, 0x3c, 0x15, 0xdb, 0x82, 0xbd, 0xb4, 0x10, 0x3c, 0xf5, 0x12, 0x32, 0xe3, 0x98,
	0x0c, 0x24, 0x7e, 0x21, 0x33, 0x29, 0xf5, 0x21, 0x0a, 0x3d, 0xf6, 0xd8, 0xd7, 0xe8, 0x1b, 0x78,
	0xf4, 0xb8, 0xa7, 0xec, 0xa2, 0x97, 0x65, 0x8f, 0x3e, 0xc1, 0xe2, 0x8c, 0xae, 0x2c, 0xde, 0xe6,
	0x9b, 0xdf, 0xef, 0x83, 0xef, 0xcf, 0xdf, 0x7c, 0x47, 0x18, 0x4d, 0x46, 0x43, 0x4e, 0x28, 0xfe,
	0x39, 0x20, 0x4c, 0x46, 0x03, 0x2c, 0x57, 0x39, 0x13, 0x7e, 0x5e, 0x80, 0x04, 0xeb, 0xcd, 0x03,
	0xf6, 0xcf, 0xaf, 0xa3, 0xd8, 0x6d, 0xc5, 0x10, 0x83, 0xf2, 0xf0, 0xe1, 0xa5, 0x57, 0xba, 0x76,
	0x0c, 0x10, 0xa7, 0x0c, 0xab, 0x89, 0x94, 0x0b, 0x3c, 0x2f, 0x8b, 0x48, 0x72, 0x58, 0x6a, 0xee,
	0xfd, 0xaf, 0x99, 0xcf, 0xa7, 0x45, 0xfe, 0x95, 0xd0, 0x80, 0x51, 0x28, 0xe6, 0x96, 0x6b, 0xd6,
	0x93, 0x22, 0x6f, 0x23, 0x17, 0xf5, 0x9e, 0x4e, 0x9a, 0xfb, 0xca, 0x31, 0x57, 0x51, 0x96, 0x8e,
	0xbd, 0xa4, 0xc8, 0xbd, 0xe0, 0x80, 0xac, 0x8f, 0x66, 0x53, 0x40, 0x59, 0x50, 0x16, 0xd2, 0x24,
	0x5a, 0x2e, 0x59, 0xda, 0xae, 0x29, 0xb9, 0xb3, 0xaf, 0x9c, 0xd7, 0x5a, 0x7e, 0xcc, 0xbd, 0xe0,
	0x85, 0xfe, 0xf8, 0xa4, 0x67, 0x0b, 0x9b, 0x2d, 0x4e, 0x45, 0x28, 0x21, 0x4c, 0x18, 0x8f, 0x13,
	0x19, 0xc2, 0x62, 0x21, 0x98, 0x6c, 0xd7, 0x5d, 0xd4, 0x6b, 0x04, 0xaf, 0x38, 0x15, 0x33, 0x98,
	0x2a, 0xf2, 0x5d, 0x01, 0xeb, 0x37, 0x32, 0xad, 0xe3, 0x86, 0xe4, 0x19, 0x3b, 0xf9, 0x0d, 0x17,
	0xf5, 0x9e, 0x0d, 0x3b, 0xbe, 0xce, 0xe8, 0x9f, 0x32, 0xfa, 0x9f, 0x8f, 0x19, 0x27, 0x5f, 0xd6,
	0x95, 0x63, 0xdc, 0x55, 0xce, 0xdb, 0xcb, 0xe5, 0xf7, 0x90, 0x71, 0xc9, 0xb2, 0x5c, 0xae, 0xf6,
	0x95, 0xd3, 0xd1, 0x67, 0x5f, 0x5a, 0xde, 0xdf, 0x6b, 0x07, 0x05, 0x2f, 0xd5, 0x45, 0x33, 0x9e,
	0x31, 0x7d, 0xcf, 0xb8, 0x71, 0xfb, 0xcf, 0x41, 0x93, 0x6f, 0xeb, 0xad, 0x8d, 0x36, 0x5b, 0x1b,
	0xdd, 0x6c, 0x6d, 0xf4, 0x67, 0x67, 0x1b, 0x9b, 0x9d, 0x6d, 0x5c, 0xed, 0x6c, 0xe3, 0xc7, 0x87,
	0x98, 0xcb, 0xa4, 0x24, 0x3e, 0x85, 0x0c, 0x83, 0xc8, 0x40, 0x70, 0xd1, 0x4f, 0x23, 0x22, 0xb0,
	0xae, 0xad, 0x7f, 0x28, 0xf8, 0x17, 0x3e, 0x97, 0xad, 0x4a, 0x26, 0x4f, 0x54, 0x80, 0xd1, 0x7d,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x0b, 0x66, 0xa9, 0x9d, 0x06, 0x02, 0x00, 0x00,
}

func (this *HrpIbcRecord) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HrpIbcRecord)
	if !ok {
		that2, ok := that.(HrpIbcRecord)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Hrp != that1.Hrp {
		return false
	}
	if this.SourceChannel != that1.SourceChannel {
		return false
	}
	if this.IcsToHeightOffset != that1.IcsToHeightOffset {
		return false
	}
	if this.IcsToTimeOffset != that1.IcsToTimeOffset {
		return false
	}
	return true
}
func (m *HrpIbcRecord) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HrpIbcRecord) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HrpIbcRecord) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.IcsToTimeOffset, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.IcsToTimeOffset):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintTypes(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x22
	if m.IcsToHeightOffset != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.IcsToHeightOffset))
		i--
		dAtA[i] = 0x18
	}
	if len(m.SourceChannel) > 0 {
		i -= len(m.SourceChannel)
		copy(dAtA[i:], m.SourceChannel)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.SourceChannel)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Hrp) > 0 {
		i -= len(m.Hrp)
		copy(dAtA[i:], m.Hrp)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Hrp)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *HrpIbcRecord) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Hrp)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.SourceChannel)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.IcsToHeightOffset != 0 {
		n += 1 + sovTypes(uint64(m.IcsToHeightOffset))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.IcsToTimeOffset)
	n += 1 + l + sovTypes(uint64(l))
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *HrpIbcRecord) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: HrpIbcRecord: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HrpIbcRecord: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hrp", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hrp = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SourceChannel", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SourceChannel = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IcsToHeightOffset", wireType)
			}
			m.IcsToHeightOffset = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.IcsToHeightOffset |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IcsToTimeOffset", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.IcsToTimeOffset, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
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
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)
