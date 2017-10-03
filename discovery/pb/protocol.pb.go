// Code generated by protoc-gen-go. DO NOT EDIT.
// source: discovery/pb/protocol.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	discovery/pb/protocol.proto

It has these top-level messages:
	Hearbeat
	PeerInfo
	Transfer
	Message
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Hearbeat struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *Hearbeat) Reset()                    { *m = Hearbeat{} }
func (m *Hearbeat) String() string            { return proto.CompactTextString(m) }
func (*Hearbeat) ProtoMessage()               {}
func (*Hearbeat) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Hearbeat) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type PeerInfo struct {
	Info map[string]string `protobuf:"bytes,2,rep,name=info" json:"info,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *PeerInfo) Reset()                    { *m = PeerInfo{} }
func (m *PeerInfo) String() string            { return proto.CompactTextString(m) }
func (*PeerInfo) ProtoMessage()               {}
func (*PeerInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *PeerInfo) GetInfo() map[string]string {
	if m != nil {
		return m.Info
	}
	return nil
}

type Transfer struct {
	Data string `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
}

func (m *Transfer) Reset()                    { *m = Transfer{} }
func (m *Transfer) String() string            { return proto.CompactTextString(m) }
func (*Transfer) ProtoMessage()               {}
func (*Transfer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Transfer) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

type Message struct {
	Version int64 `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	// Types that are valid to be assigned to Msg:
	//	*Message_Heartbeat
	//	*Message_PeerInfo
	//	*Message_Transfer
	Msg isMessage_Msg `protobuf_oneof:"msg"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type isMessage_Msg interface {
	isMessage_Msg()
}

type Message_Heartbeat struct {
	Heartbeat *Hearbeat `protobuf:"bytes,2,opt,name=heartbeat,oneof"`
}
type Message_PeerInfo struct {
	PeerInfo *PeerInfo `protobuf:"bytes,3,opt,name=peerInfo,oneof"`
}
type Message_Transfer struct {
	Transfer *Transfer `protobuf:"bytes,4,opt,name=transfer,oneof"`
}

func (*Message_Heartbeat) isMessage_Msg() {}
func (*Message_PeerInfo) isMessage_Msg()  {}
func (*Message_Transfer) isMessage_Msg()  {}

func (m *Message) GetMsg() isMessage_Msg {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (m *Message) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *Message) GetHeartbeat() *Hearbeat {
	if x, ok := m.GetMsg().(*Message_Heartbeat); ok {
		return x.Heartbeat
	}
	return nil
}

func (m *Message) GetPeerInfo() *PeerInfo {
	if x, ok := m.GetMsg().(*Message_PeerInfo); ok {
		return x.PeerInfo
	}
	return nil
}

func (m *Message) GetTransfer() *Transfer {
	if x, ok := m.GetMsg().(*Message_Transfer); ok {
		return x.Transfer
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Message) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Message_OneofMarshaler, _Message_OneofUnmarshaler, _Message_OneofSizer, []interface{}{
		(*Message_Heartbeat)(nil),
		(*Message_PeerInfo)(nil),
		(*Message_Transfer)(nil),
	}
}

func _Message_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Message)
	// msg
	switch x := m.Msg.(type) {
	case *Message_Heartbeat:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Heartbeat); err != nil {
			return err
		}
	case *Message_PeerInfo:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.PeerInfo); err != nil {
			return err
		}
	case *Message_Transfer:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Transfer); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Message.Msg has unexpected type %T", x)
	}
	return nil
}

func _Message_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Message)
	switch tag {
	case 2: // msg.heartbeat
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Hearbeat)
		err := b.DecodeMessage(msg)
		m.Msg = &Message_Heartbeat{msg}
		return true, err
	case 3: // msg.peerInfo
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(PeerInfo)
		err := b.DecodeMessage(msg)
		m.Msg = &Message_PeerInfo{msg}
		return true, err
	case 4: // msg.transfer
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Transfer)
		err := b.DecodeMessage(msg)
		m.Msg = &Message_Transfer{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Message_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Message)
	// msg
	switch x := m.Msg.(type) {
	case *Message_Heartbeat:
		s := proto.Size(x.Heartbeat)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Message_PeerInfo:
		s := proto.Size(x.PeerInfo)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Message_Transfer:
		s := proto.Size(x.Transfer)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*Hearbeat)(nil), "pb.Hearbeat")
	proto.RegisterType((*PeerInfo)(nil), "pb.PeerInfo")
	proto.RegisterType((*Transfer)(nil), "pb.Transfer")
	proto.RegisterType((*Message)(nil), "pb.Message")
}

func init() { proto.RegisterFile("discovery/pb/protocol.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 278 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0x51, 0x4b, 0xc3, 0x30,
	0x14, 0x85, 0xd7, 0x66, 0x73, 0xed, 0x9d, 0x0f, 0x12, 0x44, 0x82, 0x82, 0x8c, 0xe2, 0xc3, 0x18,
	0xd2, 0xc1, 0x7c, 0x50, 0x7c, 0x1c, 0x08, 0xf5, 0x41, 0x90, 0xe2, 0x1f, 0x48, 0xb7, 0xdb, 0x59,
	0xdc, 0x92, 0x90, 0xc4, 0x42, 0xff, 0x95, 0x3f, 0x51, 0x92, 0x35, 0xeb, 0x4b, 0xb9, 0xe7, 0x9e,
	0xaf, 0xc9, 0x39, 0x81, 0xbb, 0x5d, 0x63, 0xb6, 0xb2, 0x45, 0xdd, 0xad, 0x54, 0xb5, 0x52, 0x5a,
	0x5a, 0xb9, 0x95, 0x87, 0xdc, 0x0f, 0x34, 0x56, 0x55, 0xf6, 0x00, 0x49, 0x81, 0x5c, 0x57, 0xc8,
	0x2d, 0x65, 0x30, 0x3d, 0xa2, 0x31, 0x7c, 0x8f, 0x2c, 0x9a, 0x47, 0x8b, 0xb4, 0x0c, 0x32, 0x93,
	0x90, 0x7c, 0x22, 0xea, 0x77, 0x51, 0x4b, 0xba, 0x84, 0x71, 0x23, 0x6a, 0xc9, 0xe2, 0x39, 0x59,
	0xcc, 0xd6, 0x37, 0xb9, 0xaa, 0xf2, 0xe0, 0xe5, 0xee, 0xf3, 0x26, 0xac, 0xee, 0x4a, 0xcf, 0xdc,
	0x3e, 0x43, 0x7a, 0x5e, 0xd1, 0x2b, 0x20, 0x3f, 0xd8, 0xf5, 0x47, 0xbb, 0x91, 0x5e, 0xc3, 0xa4,
	0xe5, 0x87, 0x5f, 0x64, 0xb1, 0xdf, 0x9d, 0xc4, 0x6b, 0xfc, 0x12, 0x65, 0xf7, 0x90, 0x7c, 0x69,
	0x2e, 0x4c, 0x8d, 0x9a, 0x52, 0x18, 0xef, 0xb8, 0xe5, 0xfd, 0x8f, 0x7e, 0xce, 0xfe, 0x22, 0x98,
	0x7e, 0x9c, 0xc2, 0xb9, 0xd8, 0x2d, 0x6a, 0xd3, 0x48, 0xe1, 0x11, 0x52, 0x06, 0x49, 0x1f, 0x21,
	0xfd, 0x46, 0xae, 0xad, 0x6b, 0xe7, 0xef, 0x98, 0xad, 0x2f, 0x5d, 0xde, 0xd0, 0xb8, 0x18, 0x95,
	0x03, 0x40, 0x97, 0x90, 0xa8, 0xbe, 0x08, 0x23, 0x03, 0x1c, 0xca, 0x15, 0xa3, 0xf2, 0xec, 0x3b,
	0xd6, 0xf6, 0xf9, 0xd8, 0x78, 0x60, 0x43, 0x66, 0xc7, 0x06, 0x7f, 0x33, 0x01, 0x72, 0x34, 0xfb,
	0x4d, 0x5c, 0x90, 0xea, 0xc2, 0x3f, 0xfc, 0xd3, 0x7f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x0e, 0x95,
	0x17, 0xbb, 0x97, 0x01, 0x00, 0x00,
}