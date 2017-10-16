// Code generated by protoc-gen-go. DO NOT EDIT.
// source: app/pb/protocol.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	app/pb/protocol.proto

It has these top-level messages:
	DPeer
	SessionInfo
	UpsertSession
	PubSubMessage
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

type DPeer struct {
	DId uint64 `protobuf:"varint,1,opt,name=dId" json:"dId,omitempty"`
	PId string `protobuf:"bytes,2,opt,name=pId" json:"pId,omitempty"`
}

func (m *DPeer) Reset()                    { *m = DPeer{} }
func (m *DPeer) String() string            { return proto.CompactTextString(m) }
func (*DPeer) ProtoMessage()               {}
func (*DPeer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *DPeer) GetDId() uint64 {
	if m != nil {
		return m.DId
	}
	return 0
}

func (m *DPeer) GetPId() string {
	if m != nil {
		return m.PId
	}
	return ""
}

type SessionInfo struct {
	Owner   *DPeer            `protobuf:"bytes,1,opt,name=owner" json:"owner,omitempty"`
	Type    int32             `protobuf:"varint,2,opt,name=type" json:"type,omitempty"`
	Name    string            `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
	Address uint64            `protobuf:"varint,4,opt,name=address" json:"address,omitempty"`
	Port    uint32            `protobuf:"varint,5,opt,name=port" json:"port,omitempty"`
	Details map[string]string `protobuf:"bytes,6,rep,name=details" json:"details,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *SessionInfo) Reset()                    { *m = SessionInfo{} }
func (m *SessionInfo) String() string            { return proto.CompactTextString(m) }
func (*SessionInfo) ProtoMessage()               {}
func (*SessionInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SessionInfo) GetOwner() *DPeer {
	if m != nil {
		return m.Owner
	}
	return nil
}

func (m *SessionInfo) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *SessionInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *SessionInfo) GetAddress() uint64 {
	if m != nil {
		return m.Address
	}
	return 0
}

func (m *SessionInfo) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *SessionInfo) GetDetails() map[string]string {
	if m != nil {
		return m.Details
	}
	return nil
}

type UpsertSession struct {
	Info *SessionInfo `protobuf:"bytes,1,opt,name=info" json:"info,omitempty"`
}

func (m *UpsertSession) Reset()                    { *m = UpsertSession{} }
func (m *UpsertSession) String() string            { return proto.CompactTextString(m) }
func (*UpsertSession) ProtoMessage()               {}
func (*UpsertSession) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *UpsertSession) GetInfo() *SessionInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

type PubSubMessage struct {
	Version int64 `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	// Types that are valid to be assigned to Msg:
	//	*PubSubMessage_UpsertSession
	Msg isPubSubMessage_Msg `protobuf_oneof:"msg"`
}

func (m *PubSubMessage) Reset()                    { *m = PubSubMessage{} }
func (m *PubSubMessage) String() string            { return proto.CompactTextString(m) }
func (*PubSubMessage) ProtoMessage()               {}
func (*PubSubMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type isPubSubMessage_Msg interface {
	isPubSubMessage_Msg()
}

type PubSubMessage_UpsertSession struct {
	UpsertSession *UpsertSession `protobuf:"bytes,2,opt,name=upsertSession,oneof"`
}

func (*PubSubMessage_UpsertSession) isPubSubMessage_Msg() {}

func (m *PubSubMessage) GetMsg() isPubSubMessage_Msg {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (m *PubSubMessage) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *PubSubMessage) GetUpsertSession() *UpsertSession {
	if x, ok := m.GetMsg().(*PubSubMessage_UpsertSession); ok {
		return x.UpsertSession
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*PubSubMessage) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _PubSubMessage_OneofMarshaler, _PubSubMessage_OneofUnmarshaler, _PubSubMessage_OneofSizer, []interface{}{
		(*PubSubMessage_UpsertSession)(nil),
	}
}

func _PubSubMessage_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*PubSubMessage)
	// msg
	switch x := m.Msg.(type) {
	case *PubSubMessage_UpsertSession:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.UpsertSession); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("PubSubMessage.Msg has unexpected type %T", x)
	}
	return nil
}

func _PubSubMessage_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*PubSubMessage)
	switch tag {
	case 2: // msg.upsertSession
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(UpsertSession)
		err := b.DecodeMessage(msg)
		m.Msg = &PubSubMessage_UpsertSession{msg}
		return true, err
	default:
		return false, nil
	}
}

func _PubSubMessage_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*PubSubMessage)
	// msg
	switch x := m.Msg.(type) {
	case *PubSubMessage_UpsertSession:
		s := proto.Size(x.UpsertSession)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*DPeer)(nil), "pb.DPeer")
	proto.RegisterType((*SessionInfo)(nil), "pb.SessionInfo")
	proto.RegisterType((*UpsertSession)(nil), "pb.UpsertSession")
	proto.RegisterType((*PubSubMessage)(nil), "pb.PubSubMessage")
}

func init() { proto.RegisterFile("app/pb/protocol.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 325 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x51, 0x4d, 0x6b, 0xf2, 0x40,
	0x10, 0x7e, 0xf3, 0xa5, 0x38, 0x79, 0x43, 0xdb, 0xa5, 0x85, 0xa5, 0x14, 0x1a, 0xd2, 0x4b, 0xa0,
	0x10, 0xc1, 0x96, 0xd2, 0x7a, 0x14, 0x0b, 0x7a, 0x28, 0xc8, 0x4a, 0x7f, 0x40, 0xd2, 0x8c, 0x22,
	0x6a, 0x76, 0xd9, 0x4d, 0x2c, 0xfe, 0xf3, 0x1e, 0xcb, 0x8e, 0x0a, 0xb1, 0xb7, 0x67, 0x66, 0x9e,
	0x7d, 0x3e, 0x58, 0xb8, 0xc9, 0x95, 0xea, 0xab, 0xa2, 0xaf, 0xb4, 0xac, 0xe5, 0x97, 0xdc, 0x64,
	0x04, 0x98, 0xab, 0x8a, 0xe4, 0x11, 0x82, 0xf1, 0x0c, 0x51, 0xb3, 0x4b, 0xf0, 0xca, 0x69, 0xc9,
	0x9d, 0xd8, 0x49, 0x7d, 0x61, 0xa1, 0xdd, 0xa8, 0x69, 0xc9, 0xdd, 0xd8, 0x49, 0x7b, 0xc2, 0xc2,
	0xe4, 0xc7, 0x81, 0x70, 0x8e, 0xc6, 0xac, 0x64, 0x35, 0xad, 0x16, 0x92, 0xdd, 0x43, 0x20, 0xbf,
	0x2b, 0xd4, 0xf4, 0x2a, 0x1c, 0xf4, 0x32, 0x55, 0x64, 0xa4, 0x26, 0x0e, 0x7b, 0xc6, 0xc0, 0xaf,
	0xf7, 0x0a, 0x49, 0x23, 0x10, 0x84, 0xed, 0xae, 0xca, 0xb7, 0xc8, 0x3d, 0xd2, 0x25, 0xcc, 0x38,
	0x74, 0xf3, 0xb2, 0xd4, 0x68, 0x0c, 0xf7, 0x29, 0xc0, 0x69, 0xb4, 0x6c, 0x25, 0x75, 0xcd, 0x83,
	0xd8, 0x49, 0x23, 0x41, 0x98, 0xbd, 0x40, 0xb7, 0xc4, 0x3a, 0x5f, 0x6d, 0x0c, 0xef, 0xc4, 0x5e,
	0x1a, 0x0e, 0xee, 0xac, 0x71, 0x2b, 0x58, 0x36, 0x3e, 0x9c, 0xdf, 0xab, 0x5a, 0xef, 0xc5, 0x89,
	0x7c, 0x3b, 0x84, 0xff, 0xed, 0x83, 0x2d, 0xb8, 0xc6, 0x3d, 0x85, 0xef, 0x09, 0x0b, 0xd9, 0x35,
	0x04, 0xbb, 0x7c, 0xd3, 0xe0, 0xb1, 0xf4, 0x61, 0x18, 0xba, 0xaf, 0x4e, 0xf2, 0x0c, 0xd1, 0xa7,
	0x32, 0xa8, 0xeb, 0xa3, 0x0d, 0x7b, 0x00, 0x7f, 0x55, 0x2d, 0xe4, 0xb1, 0xfa, 0xc5, 0x9f, 0x04,
	0x82, 0x8e, 0xc9, 0x1a, 0xa2, 0x59, 0x53, 0xcc, 0x9b, 0xe2, 0x03, 0x8d, 0xc9, 0x97, 0x54, 0x74,
	0x87, 0xda, 0xb2, 0xe8, 0xa1, 0x27, 0x4e, 0x23, 0x7b, 0x83, 0xa8, 0x69, 0x1b, 0x50, 0x84, 0x70,
	0x70, 0x65, 0x85, 0xcf, 0x9c, 0x27, 0xff, 0xc4, 0x39, 0x73, 0x14, 0x80, 0xb7, 0x35, 0xcb, 0x91,
	0x3b, 0xf1, 0x8a, 0x0e, 0xfd, 0xec, 0xd3, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x7d, 0x26, 0xd3,
	0xb9, 0xf2, 0x01, 0x00, 0x00,
}