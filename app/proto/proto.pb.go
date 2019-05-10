// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
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

type Response struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_2fcc84b9998d60d8, []int{0}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Response) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type ServerConfig struct {
	PublicIP             string   `protobuf:"bytes,1,opt,name=publicIP,proto3" json:"publicIP,omitempty"`
	PrivateIP            string   `protobuf:"bytes,2,opt,name=privateIP,proto3" json:"privateIP,omitempty"`
	User                 string   `protobuf:"bytes,3,opt,name=user,proto3" json:"user,omitempty"`
	Password             string   `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	ClusterName          string   `protobuf:"bytes,5,opt,name=clusterName,proto3" json:"clusterName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServerConfig) Reset()         { *m = ServerConfig{} }
func (m *ServerConfig) String() string { return proto.CompactTextString(m) }
func (*ServerConfig) ProtoMessage()    {}
func (*ServerConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_2fcc84b9998d60d8, []int{1}
}

func (m *ServerConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerConfig.Unmarshal(m, b)
}
func (m *ServerConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerConfig.Marshal(b, m, deterministic)
}
func (m *ServerConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerConfig.Merge(m, src)
}
func (m *ServerConfig) XXX_Size() int {
	return xxx_messageInfo_ServerConfig.Size(m)
}
func (m *ServerConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ServerConfig proto.InternalMessageInfo

func (m *ServerConfig) GetPublicIP() string {
	if m != nil {
		return m.PublicIP
	}
	return ""
}

func (m *ServerConfig) GetPrivateIP() string {
	if m != nil {
		return m.PrivateIP
	}
	return ""
}

func (m *ServerConfig) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *ServerConfig) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *ServerConfig) GetClusterName() string {
	if m != nil {
		return m.ClusterName
	}
	return ""
}

type ClusterConfig struct {
	ClusterName          string          `protobuf:"bytes,1,opt,name=clusterName,proto3" json:"clusterName,omitempty"`
	AdvertiseAddress     string          `protobuf:"bytes,2,opt,name=advertiseAddress,proto3" json:"advertiseAddress,omitempty"`
	PublicIP             string          `protobuf:"bytes,3,opt,name=publicIP,proto3" json:"publicIP,omitempty"`
	DownloadBinSite      string          `protobuf:"bytes,4,opt,name=downloadBinSite,proto3" json:"downloadBinSite,omitempty"`
	Masters              []*ServerConfig `protobuf:"bytes,5,rep,name=masters,proto3" json:"masters,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *ClusterConfig) Reset()         { *m = ClusterConfig{} }
func (m *ClusterConfig) String() string { return proto.CompactTextString(m) }
func (*ClusterConfig) ProtoMessage()    {}
func (*ClusterConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_2fcc84b9998d60d8, []int{2}
}

func (m *ClusterConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClusterConfig.Unmarshal(m, b)
}
func (m *ClusterConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClusterConfig.Marshal(b, m, deterministic)
}
func (m *ClusterConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClusterConfig.Merge(m, src)
}
func (m *ClusterConfig) XXX_Size() int {
	return xxx_messageInfo_ClusterConfig.Size(m)
}
func (m *ClusterConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ClusterConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ClusterConfig proto.InternalMessageInfo

func (m *ClusterConfig) GetClusterName() string {
	if m != nil {
		return m.ClusterName
	}
	return ""
}

func (m *ClusterConfig) GetAdvertiseAddress() string {
	if m != nil {
		return m.AdvertiseAddress
	}
	return ""
}

func (m *ClusterConfig) GetPublicIP() string {
	if m != nil {
		return m.PublicIP
	}
	return ""
}

func (m *ClusterConfig) GetDownloadBinSite() string {
	if m != nil {
		return m.DownloadBinSite
	}
	return ""
}

func (m *ClusterConfig) GetMasters() []*ServerConfig {
	if m != nil {
		return m.Masters
	}
	return nil
}

type InstanceNode struct {
	InstanceID           string   `protobuf:"bytes,1,opt,name=instanceID,proto3" json:"instanceID,omitempty"`
	InstanceName         string   `protobuf:"bytes,2,opt,name=instanceName,proto3" json:"instanceName,omitempty"`
	InstanceService      string   `protobuf:"bytes,3,opt,name=instanceService,proto3" json:"instanceService,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InstanceNode) Reset()         { *m = InstanceNode{} }
func (m *InstanceNode) String() string { return proto.CompactTextString(m) }
func (*InstanceNode) ProtoMessage()    {}
func (*InstanceNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_2fcc84b9998d60d8, []int{3}
}

func (m *InstanceNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstanceNode.Unmarshal(m, b)
}
func (m *InstanceNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstanceNode.Marshal(b, m, deterministic)
}
func (m *InstanceNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstanceNode.Merge(m, src)
}
func (m *InstanceNode) XXX_Size() int {
	return xxx_messageInfo_InstanceNode.Size(m)
}
func (m *InstanceNode) XXX_DiscardUnknown() {
	xxx_messageInfo_InstanceNode.DiscardUnknown(m)
}

var xxx_messageInfo_InstanceNode proto.InternalMessageInfo

func (m *InstanceNode) GetInstanceID() string {
	if m != nil {
		return m.InstanceID
	}
	return ""
}

func (m *InstanceNode) GetInstanceName() string {
	if m != nil {
		return m.InstanceName
	}
	return ""
}

func (m *InstanceNode) GetInstanceService() string {
	if m != nil {
		return m.InstanceService
	}
	return ""
}

type NodeConfs struct {
	BootstrapConf        []byte   `protobuf:"bytes,1,opt,name=bootstrapConf,proto3" json:"bootstrapConf,omitempty"`
	Certs                []*Cert  `protobuf:"bytes,2,rep,name=certs,proto3" json:"certs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeConfs) Reset()         { *m = NodeConfs{} }
func (m *NodeConfs) String() string { return proto.CompactTextString(m) }
func (*NodeConfs) ProtoMessage()    {}
func (*NodeConfs) Descriptor() ([]byte, []int) {
	return fileDescriptor_2fcc84b9998d60d8, []int{4}
}

func (m *NodeConfs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeConfs.Unmarshal(m, b)
}
func (m *NodeConfs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeConfs.Marshal(b, m, deterministic)
}
func (m *NodeConfs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeConfs.Merge(m, src)
}
func (m *NodeConfs) XXX_Size() int {
	return xxx_messageInfo_NodeConfs.Size(m)
}
func (m *NodeConfs) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeConfs.DiscardUnknown(m)
}

var xxx_messageInfo_NodeConfs proto.InternalMessageInfo

func (m *NodeConfs) GetBootstrapConf() []byte {
	if m != nil {
		return m.BootstrapConf
	}
	return nil
}

func (m *NodeConfs) GetCerts() []*Cert {
	if m != nil {
		return m.Certs
	}
	return nil
}

type Cert struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Cert                 []byte   `protobuf:"bytes,2,opt,name=cert,proto3" json:"cert,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Cert) Reset()         { *m = Cert{} }
func (m *Cert) String() string { return proto.CompactTextString(m) }
func (*Cert) ProtoMessage()    {}
func (*Cert) Descriptor() ([]byte, []int) {
	return fileDescriptor_2fcc84b9998d60d8, []int{5}
}

func (m *Cert) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Cert.Unmarshal(m, b)
}
func (m *Cert) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Cert.Marshal(b, m, deterministic)
}
func (m *Cert) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Cert.Merge(m, src)
}
func (m *Cert) XXX_Size() int {
	return xxx_messageInfo_Cert.Size(m)
}
func (m *Cert) XXX_DiscardUnknown() {
	xxx_messageInfo_Cert.DiscardUnknown(m)
}

var xxx_messageInfo_Cert proto.InternalMessageInfo

func (m *Cert) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Cert) GetCert() []byte {
	if m != nil {
		return m.Cert
	}
	return nil
}

func init() {
	proto.RegisterType((*Response)(nil), "proto.Response")
	proto.RegisterType((*ServerConfig)(nil), "proto.ServerConfig")
	proto.RegisterType((*ClusterConfig)(nil), "proto.ClusterConfig")
	proto.RegisterType((*InstanceNode)(nil), "proto.InstanceNode")
	proto.RegisterType((*NodeConfs)(nil), "proto.NodeConfs")
	proto.RegisterType((*Cert)(nil), "proto.Cert")
}

func init() { proto.RegisterFile("proto.proto", fileDescriptor_2fcc84b9998d60d8) }

var fileDescriptor_2fcc84b9998d60d8 = []byte{
	// 479 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0xc5, 0xc4, 0x49, 0x9a, 0x89, 0xab, 0x56, 0x53, 0x0e, 0x56, 0x84, 0x50, 0xb0, 0x38, 0x44,
	0x48, 0xa4, 0x10, 0x2e, 0x5c, 0x38, 0x40, 0x72, 0x31, 0x12, 0x55, 0xe5, 0xf0, 0x03, 0x1b, 0x7b,
	0xa8, 0x16, 0x25, 0xbb, 0xd6, 0xee, 0x26, 0x15, 0x12, 0x5f, 0xc2, 0x17, 0xf0, 0x3d, 0x7c, 0x51,
	0xb5, 0xeb, 0x75, 0xea, 0xb8, 0x39, 0xe4, 0x92, 0xcc, 0x7b, 0xde, 0x99, 0x37, 0x6f, 0xdf, 0xc2,
	0xb0, 0x54, 0xd2, 0xc8, 0xa9, 0xfb, 0xc5, 0xae, 0xfb, 0x4b, 0x3e, 0xc1, 0x59, 0x46, 0xba, 0x94,
	0x42, 0x13, 0x22, 0x84, 0xb9, 0x2c, 0x28, 0x0e, 0xc6, 0xc1, 0x64, 0x90, 0xb9, 0x1a, 0x63, 0xe8,
	0x6f, 0x48, 0x6b, 0x76, 0x47, 0xf1, 0x73, 0x47, 0xd7, 0x30, 0xf9, 0x1b, 0x40, 0xb4, 0x24, 0xb5,
	0x23, 0x35, 0x97, 0xe2, 0x27, 0xbf, 0xc3, 0x11, 0x9c, 0x95, 0xdb, 0xd5, 0x9a, 0xe7, 0xe9, 0xad,
	0x1f, 0xb1, 0xc7, 0xf8, 0x12, 0x06, 0xa5, 0xe2, 0x3b, 0x66, 0x28, 0xbd, 0xf5, 0x83, 0x1e, 0x09,
	0x2b, 0xbc, 0xd5, 0xa4, 0xe2, 0x4e, 0x25, 0x6c, 0x6b, 0x37, 0x8d, 0x69, 0x7d, 0x2f, 0x55, 0x11,
	0x87, 0x7e, 0x9a, 0xc7, 0x38, 0x86, 0x61, 0xbe, 0xde, 0x6a, 0x43, 0xea, 0x86, 0x6d, 0x28, 0xee,
	0xba, 0xcf, 0x4d, 0x2a, 0xf9, 0x1f, 0xc0, 0xf9, 0xbc, 0xc2, 0x7e, 0xbb, 0x56, 0x4f, 0xf0, 0xa4,
	0x07, 0xdf, 0xc2, 0x25, 0x2b, 0x76, 0xa4, 0x0c, 0xd7, 0xf4, 0xa5, 0x28, 0x14, 0x69, 0xed, 0x57,
	0x7d, 0xc2, 0x1f, 0x78, 0xed, 0xb4, 0xbc, 0x4e, 0xe0, 0xa2, 0x90, 0xf7, 0x62, 0x2d, 0x59, 0xf1,
	0x95, 0x8b, 0x25, 0x37, 0xe4, 0x0d, 0xb4, 0x69, 0x7c, 0x07, 0xfd, 0x0d, 0xb3, 0xfa, 0x3a, 0xee,
	0x8e, 0x3b, 0x93, 0xe1, 0xec, 0xaa, 0x0a, 0x67, 0xda, 0xbc, 0xd7, 0xac, 0x3e, 0x93, 0xfc, 0x81,
	0x28, 0x15, 0xda, 0x30, 0x91, 0xd3, 0x8d, 0xcd, 0xe6, 0x15, 0x00, 0xf7, 0x38, 0x5d, 0x78, 0x47,
	0x0d, 0x06, 0x13, 0x88, 0x6a, 0xe4, 0x3c, 0x57, 0x66, 0x0e, 0x38, 0xbb, 0x6c, 0x8d, 0xad, 0x28,
	0xcf, 0xc9, 0xfb, 0x69, 0xd3, 0xc9, 0x0f, 0x18, 0x58, 0x55, 0xbb, 0x94, 0xc6, 0x37, 0x70, 0xbe,
	0x92, 0xd2, 0x68, 0xa3, 0x58, 0x69, 0x19, 0xa7, 0x1e, 0x65, 0x87, 0x24, 0xbe, 0x86, 0x6e, 0x4e,
	0xca, 0xd8, 0x6b, 0xb4, 0xee, 0x86, 0xde, 0xdd, 0x9c, 0x94, 0xc9, 0xaa, 0x2f, 0xc9, 0x14, 0x42,
	0x0b, 0xed, 0x13, 0x10, 0x8f, 0xb9, 0xb8, 0xda, 0xbd, 0x47, 0x52, 0xc6, 0xed, 0x1d, 0x65, 0xae,
	0x9e, 0x6d, 0xa0, 0xef, 0x33, 0xc3, 0x6b, 0x08, 0x53, 0xc1, 0x0d, 0xbe, 0xa8, 0xc7, 0x36, 0xf3,
	0x1e, 0x5d, 0x78, 0xb6, 0x7e, 0xdd, 0xc9, 0x33, 0xfc, 0x00, 0xbd, 0x05, 0xad, 0xc9, 0xd0, 0xc9,
	0x2d, 0xb3, 0x5f, 0xd0, 0xfb, 0xee, 0x6e, 0x1f, 0xa7, 0x10, 0x7e, 0x93, 0x5c, 0xe0, 0xb1, 0x88,
	0x8e, 0x89, 0xbd, 0xdf, 0x8b, 0x9d, 0xd8, 0x31, 0xfb, 0x17, 0x40, 0xe8, 0x72, 0xbd, 0x6e, 0x49,
	0x35, 0x43, 0x1f, 0x5d, 0x7a, 0x72, 0x9f, 0xc5, 0x51, 0xad, 0x83, 0x96, 0x23, 0xdb, 0x7d, 0x86,
	0xab, 0x65, 0x29, 0x4d, 0x7d, 0x6c, 0x41, 0xda, 0x28, 0xf9, 0xfb, 0xd4, 0xf6, 0x55, 0xcf, 0x31,
	0x1f, 0x1f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x35, 0xcc, 0xee, 0x2d, 0x52, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ClusterClient is the client API for Cluster service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ClusterClient interface {
	Init(ctx context.Context, in *ClusterConfig, opts ...grpc.CallOption) (*Response, error)
	Delete(ctx context.Context, in *ClusterConfig, opts ...grpc.CallOption) (*Response, error)
}

type clusterClient struct {
	cc *grpc.ClientConn
}

func NewClusterClient(cc *grpc.ClientConn) ClusterClient {
	return &clusterClient{cc}
}

func (c *clusterClient) Init(ctx context.Context, in *ClusterConfig, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.cluster/Init", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) Delete(ctx context.Context, in *ClusterConfig, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.cluster/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterServer is the server API for Cluster service.
type ClusterServer interface {
	Init(context.Context, *ClusterConfig) (*Response, error)
	Delete(context.Context, *ClusterConfig) (*Response, error)
}

func RegisterClusterServer(s *grpc.Server, srv ClusterServer) {
	s.RegisterService(&_Cluster_serviceDesc, srv)
}

func _Cluster_Init_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClusterConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).Init(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.cluster/Init",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).Init(ctx, req.(*ClusterConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClusterConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.cluster/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).Delete(ctx, req.(*ClusterConfig))
	}
	return interceptor(ctx, in, info, handler)
}

var _Cluster_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.cluster",
	HandlerType: (*ClusterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Init",
			Handler:    _Cluster_Init_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Cluster_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto.proto",
}

// MasterClient is the client API for Master service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MasterClient interface {
	Join(ctx context.Context, in *ServerConfig, opts ...grpc.CallOption) (*Response, error)
	Delete(ctx context.Context, in *ServerConfig, opts ...grpc.CallOption) (*Response, error)
}

type masterClient struct {
	cc *grpc.ClientConn
}

func NewMasterClient(cc *grpc.ClientConn) MasterClient {
	return &masterClient{cc}
}

func (c *masterClient) Join(ctx context.Context, in *ServerConfig, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.Master/Join", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterClient) Delete(ctx context.Context, in *ServerConfig, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.Master/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MasterServer is the server API for Master service.
type MasterServer interface {
	Join(context.Context, *ServerConfig) (*Response, error)
	Delete(context.Context, *ServerConfig) (*Response, error)
}

func RegisterMasterServer(s *grpc.Server, srv MasterServer) {
	s.RegisterService(&_Master_serviceDesc, srv)
}

func _Master_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServerConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Master/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterServer).Join(ctx, req.(*ServerConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _Master_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServerConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Master/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterServer).Delete(ctx, req.(*ServerConfig))
	}
	return interceptor(ctx, in, info, handler)
}

var _Master_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Master",
	HandlerType: (*MasterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Join",
			Handler:    _Master_Join_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Master_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto.proto",
}

// NodeClient is the client API for Node service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NodeClient interface {
	Join(ctx context.Context, in *InstanceNode, opts ...grpc.CallOption) (*NodeConfs, error)
	Delete(ctx context.Context, in *InstanceNode, opts ...grpc.CallOption) (*Response, error)
	SpotInstanceDestroy(ctx context.Context, in *InstanceNode, opts ...grpc.CallOption) (*Response, error)
}

type nodeClient struct {
	cc *grpc.ClientConn
}

func NewNodeClient(cc *grpc.ClientConn) NodeClient {
	return &nodeClient{cc}
}

func (c *nodeClient) Join(ctx context.Context, in *InstanceNode, opts ...grpc.CallOption) (*NodeConfs, error) {
	out := new(NodeConfs)
	err := c.cc.Invoke(ctx, "/proto.Node/Join", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) Delete(ctx context.Context, in *InstanceNode, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.Node/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) SpotInstanceDestroy(ctx context.Context, in *InstanceNode, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.Node/SpotInstanceDestroy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeServer is the server API for Node service.
type NodeServer interface {
	Join(context.Context, *InstanceNode) (*NodeConfs, error)
	Delete(context.Context, *InstanceNode) (*Response, error)
	SpotInstanceDestroy(context.Context, *InstanceNode) (*Response, error)
}

func RegisterNodeServer(s *grpc.Server, srv NodeServer) {
	s.RegisterService(&_Node_serviceDesc, srv)
}

func _Node_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstanceNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Node/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Join(ctx, req.(*InstanceNode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstanceNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Node/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Delete(ctx, req.(*InstanceNode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_SpotInstanceDestroy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstanceNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).SpotInstanceDestroy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Node/SpotInstanceDestroy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).SpotInstanceDestroy(ctx, req.(*InstanceNode))
	}
	return interceptor(ctx, in, info, handler)
}

var _Node_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Node",
	HandlerType: (*NodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Join",
			Handler:    _Node_Join_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Node_Delete_Handler,
		},
		{
			MethodName: "SpotInstanceDestroy",
			Handler:    _Node_SpotInstanceDestroy_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto.proto",
}
