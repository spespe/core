// Code generated by protoc-gen-go. DO NOT EDIT.
// source: miner.proto

package sonm

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// grpccmd imports
import (
	"io"

	"github.com/spf13/cobra"
	"github.com/sshaman1101/grpccmd"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type MinerHandshakeRequest struct {
	Hub   string      `protobuf:"bytes,1,opt,name=hub" json:"hub,omitempty"`
	Tasks []*TaskInfo `protobuf:"bytes,2,rep,name=tasks" json:"tasks,omitempty"`
}

func (m *MinerHandshakeRequest) Reset()                    { *m = MinerHandshakeRequest{} }
func (m *MinerHandshakeRequest) String() string            { return proto.CompactTextString(m) }
func (*MinerHandshakeRequest) ProtoMessage()               {}
func (*MinerHandshakeRequest) Descriptor() ([]byte, []int) { return fileDescriptor9, []int{0} }

func (m *MinerHandshakeRequest) GetHub() string {
	if m != nil {
		return m.Hub
	}
	return ""
}

func (m *MinerHandshakeRequest) GetTasks() []*TaskInfo {
	if m != nil {
		return m.Tasks
	}
	return nil
}

type MinerHandshakeReply struct {
	Miner        string        `protobuf:"bytes,1,opt,name=miner" json:"miner,omitempty"`
	Capabilities *Capabilities `protobuf:"bytes,2,opt,name=capabilities" json:"capabilities,omitempty"`
	NatType      NATType       `protobuf:"varint,3,opt,name=natType,enum=sonm.NATType" json:"natType,omitempty"`
}

func (m *MinerHandshakeReply) Reset()                    { *m = MinerHandshakeReply{} }
func (m *MinerHandshakeReply) String() string            { return proto.CompactTextString(m) }
func (*MinerHandshakeReply) ProtoMessage()               {}
func (*MinerHandshakeReply) Descriptor() ([]byte, []int) { return fileDescriptor9, []int{1} }

func (m *MinerHandshakeReply) GetMiner() string {
	if m != nil {
		return m.Miner
	}
	return ""
}

func (m *MinerHandshakeReply) GetCapabilities() *Capabilities {
	if m != nil {
		return m.Capabilities
	}
	return nil
}

func (m *MinerHandshakeReply) GetNatType() NATType {
	if m != nil {
		return m.NatType
	}
	return NATType_NONE
}

type MinerStartRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// Container describes container settings.
	Container     *Container                `protobuf:"bytes,2,opt,name=container" json:"container,omitempty"`
	RestartPolicy *ContainerRestartPolicy   `protobuf:"bytes,3,opt,name=restartPolicy" json:"restartPolicy,omitempty"`
	Resources     *TaskResourceRequirements `protobuf:"bytes,4,opt,name=resources" json:"resources,omitempty"`
	// OrderId describes an unique order identifier.
	// It is here for proper resource allocation and limitation.
	OrderId string `protobuf:"bytes,5,opt,name=orderId" json:"orderId,omitempty"`
}

func (m *MinerStartRequest) Reset()                    { *m = MinerStartRequest{} }
func (m *MinerStartRequest) String() string            { return proto.CompactTextString(m) }
func (*MinerStartRequest) ProtoMessage()               {}
func (*MinerStartRequest) Descriptor() ([]byte, []int) { return fileDescriptor9, []int{2} }

func (m *MinerStartRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *MinerStartRequest) GetContainer() *Container {
	if m != nil {
		return m.Container
	}
	return nil
}

func (m *MinerStartRequest) GetRestartPolicy() *ContainerRestartPolicy {
	if m != nil {
		return m.RestartPolicy
	}
	return nil
}

func (m *MinerStartRequest) GetResources() *TaskResourceRequirements {
	if m != nil {
		return m.Resources
	}
	return nil
}

func (m *MinerStartRequest) GetOrderId() string {
	if m != nil {
		return m.OrderId
	}
	return ""
}

type SocketAddr struct {
	Addr string `protobuf:"bytes,1,opt,name=addr" json:"addr,omitempty"`
	//
	// Actually an `uint16` here. Protobuf is so clear and handy.
	Port uint32 `protobuf:"varint,2,opt,name=port" json:"port,omitempty"`
}

func (m *SocketAddr) Reset()                    { *m = SocketAddr{} }
func (m *SocketAddr) String() string            { return proto.CompactTextString(m) }
func (*SocketAddr) ProtoMessage()               {}
func (*SocketAddr) Descriptor() ([]byte, []int) { return fileDescriptor9, []int{3} }

func (m *SocketAddr) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *SocketAddr) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

type MinerStartReply struct {
	Container string `protobuf:"bytes,1,opt,name=container" json:"container,omitempty"`
	// PortMap represent port mapping between container network and host ones.
	PortMap map[string]*Endpoints `protobuf:"bytes,2,rep,name=portMap" json:"portMap,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *MinerStartReply) Reset()                    { *m = MinerStartReply{} }
func (m *MinerStartReply) String() string            { return proto.CompactTextString(m) }
func (*MinerStartReply) ProtoMessage()               {}
func (*MinerStartReply) Descriptor() ([]byte, []int) { return fileDescriptor9, []int{4} }

func (m *MinerStartReply) GetContainer() string {
	if m != nil {
		return m.Container
	}
	return ""
}

func (m *MinerStartReply) GetPortMap() map[string]*Endpoints {
	if m != nil {
		return m.PortMap
	}
	return nil
}

type TaskInfo struct {
	Request *MinerStartRequest `protobuf:"bytes,1,opt,name=request" json:"request,omitempty"`
	Reply   *MinerStartReply   `protobuf:"bytes,2,opt,name=reply" json:"reply,omitempty"`
}

func (m *TaskInfo) Reset()                    { *m = TaskInfo{} }
func (m *TaskInfo) String() string            { return proto.CompactTextString(m) }
func (*TaskInfo) ProtoMessage()               {}
func (*TaskInfo) Descriptor() ([]byte, []int) { return fileDescriptor9, []int{5} }

func (m *TaskInfo) GetRequest() *MinerStartRequest {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *TaskInfo) GetReply() *MinerStartReply {
	if m != nil {
		return m.Reply
	}
	return nil
}

type Endpoints struct {
	Endpoints []*SocketAddr `protobuf:"bytes,1,rep,name=endpoints" json:"endpoints,omitempty"`
}

func (m *Endpoints) Reset()                    { *m = Endpoints{} }
func (m *Endpoints) String() string            { return proto.CompactTextString(m) }
func (*Endpoints) ProtoMessage()               {}
func (*Endpoints) Descriptor() ([]byte, []int) { return fileDescriptor9, []int{6} }

func (m *Endpoints) GetEndpoints() []*SocketAddr {
	if m != nil {
		return m.Endpoints
	}
	return nil
}

type MinerStatusMapRequest struct {
}

func (m *MinerStatusMapRequest) Reset()                    { *m = MinerStatusMapRequest{} }
func (m *MinerStatusMapRequest) String() string            { return proto.CompactTextString(m) }
func (*MinerStatusMapRequest) ProtoMessage()               {}
func (*MinerStatusMapRequest) Descriptor() ([]byte, []int) { return fileDescriptor9, []int{7} }

type SaveRequest struct {
	ImageID string `protobuf:"bytes,1,opt,name=imageID" json:"imageID,omitempty"`
}

func (m *SaveRequest) Reset()                    { *m = SaveRequest{} }
func (m *SaveRequest) String() string            { return proto.CompactTextString(m) }
func (*SaveRequest) ProtoMessage()               {}
func (*SaveRequest) Descriptor() ([]byte, []int) { return fileDescriptor9, []int{8} }

func (m *SaveRequest) GetImageID() string {
	if m != nil {
		return m.ImageID
	}
	return ""
}

func init() {
	proto.RegisterType((*MinerHandshakeRequest)(nil), "sonm.MinerHandshakeRequest")
	proto.RegisterType((*MinerHandshakeReply)(nil), "sonm.MinerHandshakeReply")
	proto.RegisterType((*MinerStartRequest)(nil), "sonm.MinerStartRequest")
	proto.RegisterType((*SocketAddr)(nil), "sonm.SocketAddr")
	proto.RegisterType((*MinerStartReply)(nil), "sonm.MinerStartReply")
	proto.RegisterType((*TaskInfo)(nil), "sonm.TaskInfo")
	proto.RegisterType((*Endpoints)(nil), "sonm.Endpoints")
	proto.RegisterType((*MinerStatusMapRequest)(nil), "sonm.MinerStatusMapRequest")
	proto.RegisterType((*SaveRequest)(nil), "sonm.SaveRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Miner service

type MinerClient interface {
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*PingReply, error)
	Info(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*InfoReply, error)
	Handshake(ctx context.Context, in *MinerHandshakeRequest, opts ...grpc.CallOption) (*MinerHandshakeReply, error)
	Save(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (Miner_SaveClient, error)
	Load(ctx context.Context, opts ...grpc.CallOption) (Miner_LoadClient, error)
	Start(ctx context.Context, in *MinerStartRequest, opts ...grpc.CallOption) (*MinerStartReply, error)
	Stop(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Empty, error)
	TasksStatus(ctx context.Context, opts ...grpc.CallOption) (Miner_TasksStatusClient, error)
	TaskDetails(ctx context.Context, in *ID, opts ...grpc.CallOption) (*TaskStatusReply, error)
	TaskLogs(ctx context.Context, in *TaskLogsRequest, opts ...grpc.CallOption) (Miner_TaskLogsClient, error)
	DiscoverHub(ctx context.Context, in *DiscoverHubRequest, opts ...grpc.CallOption) (*Empty, error)
}

type minerClient struct {
	cc *grpc.ClientConn
}

func NewMinerClient(cc *grpc.ClientConn) MinerClient {
	return &minerClient{cc}
}

func (c *minerClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*PingReply, error) {
	out := new(PingReply)
	err := grpc.Invoke(ctx, "/sonm.Miner/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *minerClient) Info(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*InfoReply, error) {
	out := new(InfoReply)
	err := grpc.Invoke(ctx, "/sonm.Miner/Info", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *minerClient) Handshake(ctx context.Context, in *MinerHandshakeRequest, opts ...grpc.CallOption) (*MinerHandshakeReply, error) {
	out := new(MinerHandshakeReply)
	err := grpc.Invoke(ctx, "/sonm.Miner/Handshake", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *minerClient) Save(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (Miner_SaveClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Miner_serviceDesc.Streams[0], c.cc, "/sonm.Miner/Save", opts...)
	if err != nil {
		return nil, err
	}
	x := &minerSaveClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Miner_SaveClient interface {
	Recv() (*Chunk, error)
	grpc.ClientStream
}

type minerSaveClient struct {
	grpc.ClientStream
}

func (x *minerSaveClient) Recv() (*Chunk, error) {
	m := new(Chunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *minerClient) Load(ctx context.Context, opts ...grpc.CallOption) (Miner_LoadClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Miner_serviceDesc.Streams[1], c.cc, "/sonm.Miner/Load", opts...)
	if err != nil {
		return nil, err
	}
	x := &minerLoadClient{stream}
	return x, nil
}

type Miner_LoadClient interface {
	Send(*Chunk) error
	Recv() (*Progress, error)
	grpc.ClientStream
}

type minerLoadClient struct {
	grpc.ClientStream
}

func (x *minerLoadClient) Send(m *Chunk) error {
	return x.ClientStream.SendMsg(m)
}

func (x *minerLoadClient) Recv() (*Progress, error) {
	m := new(Progress)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *minerClient) Start(ctx context.Context, in *MinerStartRequest, opts ...grpc.CallOption) (*MinerStartReply, error) {
	out := new(MinerStartReply)
	err := grpc.Invoke(ctx, "/sonm.Miner/Start", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *minerClient) Stop(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/sonm.Miner/Stop", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *minerClient) TasksStatus(ctx context.Context, opts ...grpc.CallOption) (Miner_TasksStatusClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Miner_serviceDesc.Streams[2], c.cc, "/sonm.Miner/TasksStatus", opts...)
	if err != nil {
		return nil, err
	}
	x := &minerTasksStatusClient{stream}
	return x, nil
}

type Miner_TasksStatusClient interface {
	Send(*MinerStatusMapRequest) error
	Recv() (*StatusMapReply, error)
	grpc.ClientStream
}

type minerTasksStatusClient struct {
	grpc.ClientStream
}

func (x *minerTasksStatusClient) Send(m *MinerStatusMapRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *minerTasksStatusClient) Recv() (*StatusMapReply, error) {
	m := new(StatusMapReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *minerClient) TaskDetails(ctx context.Context, in *ID, opts ...grpc.CallOption) (*TaskStatusReply, error) {
	out := new(TaskStatusReply)
	err := grpc.Invoke(ctx, "/sonm.Miner/TaskDetails", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *minerClient) TaskLogs(ctx context.Context, in *TaskLogsRequest, opts ...grpc.CallOption) (Miner_TaskLogsClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Miner_serviceDesc.Streams[3], c.cc, "/sonm.Miner/TaskLogs", opts...)
	if err != nil {
		return nil, err
	}
	x := &minerTaskLogsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Miner_TaskLogsClient interface {
	Recv() (*TaskLogsChunk, error)
	grpc.ClientStream
}

type minerTaskLogsClient struct {
	grpc.ClientStream
}

func (x *minerTaskLogsClient) Recv() (*TaskLogsChunk, error) {
	m := new(TaskLogsChunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *minerClient) DiscoverHub(ctx context.Context, in *DiscoverHubRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/sonm.Miner/DiscoverHub", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Miner service

type MinerServer interface {
	Ping(context.Context, *Empty) (*PingReply, error)
	Info(context.Context, *Empty) (*InfoReply, error)
	Handshake(context.Context, *MinerHandshakeRequest) (*MinerHandshakeReply, error)
	Save(*SaveRequest, Miner_SaveServer) error
	Load(Miner_LoadServer) error
	Start(context.Context, *MinerStartRequest) (*MinerStartReply, error)
	Stop(context.Context, *ID) (*Empty, error)
	TasksStatus(Miner_TasksStatusServer) error
	TaskDetails(context.Context, *ID) (*TaskStatusReply, error)
	TaskLogs(*TaskLogsRequest, Miner_TaskLogsServer) error
	DiscoverHub(context.Context, *DiscoverHubRequest) (*Empty, error)
}

func RegisterMinerServer(s *grpc.Server, srv MinerServer) {
	s.RegisterService(&_Miner_serviceDesc, srv)
}

func _Miner_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MinerServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sonm.Miner/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MinerServer).Ping(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Miner_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MinerServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sonm.Miner/Info",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MinerServer).Info(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Miner_Handshake_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MinerHandshakeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MinerServer).Handshake(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sonm.Miner/Handshake",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MinerServer).Handshake(ctx, req.(*MinerHandshakeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Miner_Save_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SaveRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MinerServer).Save(m, &minerSaveServer{stream})
}

type Miner_SaveServer interface {
	Send(*Chunk) error
	grpc.ServerStream
}

type minerSaveServer struct {
	grpc.ServerStream
}

func (x *minerSaveServer) Send(m *Chunk) error {
	return x.ServerStream.SendMsg(m)
}

func _Miner_Load_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MinerServer).Load(&minerLoadServer{stream})
}

type Miner_LoadServer interface {
	Send(*Progress) error
	Recv() (*Chunk, error)
	grpc.ServerStream
}

type minerLoadServer struct {
	grpc.ServerStream
}

func (x *minerLoadServer) Send(m *Progress) error {
	return x.ServerStream.SendMsg(m)
}

func (x *minerLoadServer) Recv() (*Chunk, error) {
	m := new(Chunk)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Miner_Start_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MinerStartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MinerServer).Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sonm.Miner/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MinerServer).Start(ctx, req.(*MinerStartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Miner_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MinerServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sonm.Miner/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MinerServer).Stop(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Miner_TasksStatus_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MinerServer).TasksStatus(&minerTasksStatusServer{stream})
}

type Miner_TasksStatusServer interface {
	Send(*StatusMapReply) error
	Recv() (*MinerStatusMapRequest, error)
	grpc.ServerStream
}

type minerTasksStatusServer struct {
	grpc.ServerStream
}

func (x *minerTasksStatusServer) Send(m *StatusMapReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *minerTasksStatusServer) Recv() (*MinerStatusMapRequest, error) {
	m := new(MinerStatusMapRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Miner_TaskDetails_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MinerServer).TaskDetails(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sonm.Miner/TaskDetails",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MinerServer).TaskDetails(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Miner_TaskLogs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TaskLogsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MinerServer).TaskLogs(m, &minerTaskLogsServer{stream})
}

type Miner_TaskLogsServer interface {
	Send(*TaskLogsChunk) error
	grpc.ServerStream
}

type minerTaskLogsServer struct {
	grpc.ServerStream
}

func (x *minerTaskLogsServer) Send(m *TaskLogsChunk) error {
	return x.ServerStream.SendMsg(m)
}

func _Miner_DiscoverHub_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiscoverHubRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MinerServer).DiscoverHub(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sonm.Miner/DiscoverHub",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MinerServer).DiscoverHub(ctx, req.(*DiscoverHubRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Miner_serviceDesc = grpc.ServiceDesc{
	ServiceName: "sonm.Miner",
	HandlerType: (*MinerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Miner_Ping_Handler,
		},
		{
			MethodName: "Info",
			Handler:    _Miner_Info_Handler,
		},
		{
			MethodName: "Handshake",
			Handler:    _Miner_Handshake_Handler,
		},
		{
			MethodName: "Start",
			Handler:    _Miner_Start_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _Miner_Stop_Handler,
		},
		{
			MethodName: "TaskDetails",
			Handler:    _Miner_TaskDetails_Handler,
		},
		{
			MethodName: "DiscoverHub",
			Handler:    _Miner_DiscoverHub_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Save",
			Handler:       _Miner_Save_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Load",
			Handler:       _Miner_Load_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "TasksStatus",
			Handler:       _Miner_TasksStatus_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "TaskLogs",
			Handler:       _Miner_TaskLogs_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "miner.proto",
}

// Begin grpccmd
var _ = grpccmd.RunE

// Miner
var _MinerCmd = &cobra.Command{
	Use:   "miner [method]",
	Short: "Subcommand for the Miner service.",
}

var _Miner_PingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Make the Ping method call, input-type: sonm.Empty output-type: sonm.PingReply",
	RunE: grpccmd.RunE(
		"Ping",
		"sonm.Empty",
		func(c io.Closer) interface{} {
			cc := c.(*grpc.ClientConn)
			return NewMinerClient(cc)
		},
	),
}

var _Miner_PingCmd_gen = &cobra.Command{
	Use:   "ping-gen",
	Short: "Generate JSON for method call of Ping (input-type: sonm.Empty)",
	RunE:  grpccmd.TypeToJson("sonm.Empty"),
}

var _Miner_InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Make the Info method call, input-type: sonm.Empty output-type: sonm.InfoReply",
	RunE: grpccmd.RunE(
		"Info",
		"sonm.Empty",
		func(c io.Closer) interface{} {
			cc := c.(*grpc.ClientConn)
			return NewMinerClient(cc)
		},
	),
}

var _Miner_InfoCmd_gen = &cobra.Command{
	Use:   "info-gen",
	Short: "Generate JSON for method call of Info (input-type: sonm.Empty)",
	RunE:  grpccmd.TypeToJson("sonm.Empty"),
}

var _Miner_HandshakeCmd = &cobra.Command{
	Use:   "handshake",
	Short: "Make the Handshake method call, input-type: sonm.MinerHandshakeRequest output-type: sonm.MinerHandshakeReply",
	RunE: grpccmd.RunE(
		"Handshake",
		"sonm.MinerHandshakeRequest",
		func(c io.Closer) interface{} {
			cc := c.(*grpc.ClientConn)
			return NewMinerClient(cc)
		},
	),
}

var _Miner_HandshakeCmd_gen = &cobra.Command{
	Use:   "handshake-gen",
	Short: "Generate JSON for method call of Handshake (input-type: sonm.MinerHandshakeRequest)",
	RunE:  grpccmd.TypeToJson("sonm.MinerHandshakeRequest"),
}

var _Miner_SaveCmd = &cobra.Command{
	Use:   "save",
	Short: "Make the Save method call, input-type: sonm.SaveRequest output-type: sonm.Chunk",
	RunE: grpccmd.RunE(
		"Save",
		"sonm.SaveRequest",
		func(c io.Closer) interface{} {
			cc := c.(*grpc.ClientConn)
			return NewMinerClient(cc)
		},
	),
}

var _Miner_SaveCmd_gen = &cobra.Command{
	Use:   "save-gen",
	Short: "Generate JSON for method call of Save (input-type: sonm.SaveRequest)",
	RunE:  grpccmd.TypeToJson("sonm.SaveRequest"),
}

var _Miner_LoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Make the Load method call, input-type: sonm.Chunk output-type: sonm.Progress",
	RunE: grpccmd.RunE(
		"Load",
		"sonm.Chunk",
		func(c io.Closer) interface{} {
			cc := c.(*grpc.ClientConn)
			return NewMinerClient(cc)
		},
	),
}

var _Miner_LoadCmd_gen = &cobra.Command{
	Use:   "load-gen",
	Short: "Generate JSON for method call of Load (input-type: sonm.Chunk)",
	RunE:  grpccmd.TypeToJson("sonm.Chunk"),
}

var _Miner_StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Make the Start method call, input-type: sonm.MinerStartRequest output-type: sonm.MinerStartReply",
	RunE: grpccmd.RunE(
		"Start",
		"sonm.MinerStartRequest",
		func(c io.Closer) interface{} {
			cc := c.(*grpc.ClientConn)
			return NewMinerClient(cc)
		},
	),
}

var _Miner_StartCmd_gen = &cobra.Command{
	Use:   "start-gen",
	Short: "Generate JSON for method call of Start (input-type: sonm.MinerStartRequest)",
	RunE:  grpccmd.TypeToJson("sonm.MinerStartRequest"),
}

var _Miner_StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Make the Stop method call, input-type: sonm.ID output-type: sonm.Empty",
	RunE: grpccmd.RunE(
		"Stop",
		"sonm.ID",
		func(c io.Closer) interface{} {
			cc := c.(*grpc.ClientConn)
			return NewMinerClient(cc)
		},
	),
}

var _Miner_StopCmd_gen = &cobra.Command{
	Use:   "stop-gen",
	Short: "Generate JSON for method call of Stop (input-type: sonm.ID)",
	RunE:  grpccmd.TypeToJson("sonm.ID"),
}

var _Miner_TasksStatusCmd = &cobra.Command{
	Use:   "tasksStatus",
	Short: "Make the TasksStatus method call, input-type: sonm.MinerStatusMapRequest output-type: sonm.StatusMapReply",
	RunE: grpccmd.RunE(
		"TasksStatus",
		"sonm.MinerStatusMapRequest",
		func(c io.Closer) interface{} {
			cc := c.(*grpc.ClientConn)
			return NewMinerClient(cc)
		},
	),
}

var _Miner_TasksStatusCmd_gen = &cobra.Command{
	Use:   "tasksStatus-gen",
	Short: "Generate JSON for method call of TasksStatus (input-type: sonm.MinerStatusMapRequest)",
	RunE:  grpccmd.TypeToJson("sonm.MinerStatusMapRequest"),
}

var _Miner_TaskDetailsCmd = &cobra.Command{
	Use:   "taskDetails",
	Short: "Make the TaskDetails method call, input-type: sonm.ID output-type: sonm.TaskStatusReply",
	RunE: grpccmd.RunE(
		"TaskDetails",
		"sonm.ID",
		func(c io.Closer) interface{} {
			cc := c.(*grpc.ClientConn)
			return NewMinerClient(cc)
		},
	),
}

var _Miner_TaskDetailsCmd_gen = &cobra.Command{
	Use:   "taskDetails-gen",
	Short: "Generate JSON for method call of TaskDetails (input-type: sonm.ID)",
	RunE:  grpccmd.TypeToJson("sonm.ID"),
}

var _Miner_TaskLogsCmd = &cobra.Command{
	Use:   "taskLogs",
	Short: "Make the TaskLogs method call, input-type: sonm.TaskLogsRequest output-type: sonm.TaskLogsChunk",
	RunE: grpccmd.RunE(
		"TaskLogs",
		"sonm.TaskLogsRequest",
		func(c io.Closer) interface{} {
			cc := c.(*grpc.ClientConn)
			return NewMinerClient(cc)
		},
	),
}

var _Miner_TaskLogsCmd_gen = &cobra.Command{
	Use:   "taskLogs-gen",
	Short: "Generate JSON for method call of TaskLogs (input-type: sonm.TaskLogsRequest)",
	RunE:  grpccmd.TypeToJson("sonm.TaskLogsRequest"),
}

var _Miner_DiscoverHubCmd = &cobra.Command{
	Use:   "discoverHub",
	Short: "Make the DiscoverHub method call, input-type: sonm.DiscoverHubRequest output-type: sonm.Empty",
	RunE: grpccmd.RunE(
		"DiscoverHub",
		"sonm.DiscoverHubRequest",
		func(c io.Closer) interface{} {
			cc := c.(*grpc.ClientConn)
			return NewMinerClient(cc)
		},
	),
}

var _Miner_DiscoverHubCmd_gen = &cobra.Command{
	Use:   "discoverHub-gen",
	Short: "Generate JSON for method call of DiscoverHub (input-type: sonm.DiscoverHubRequest)",
	RunE:  grpccmd.TypeToJson("sonm.DiscoverHubRequest"),
}

// Register commands with the root command and service command
func init() {
	grpccmd.RegisterServiceCmd(_MinerCmd)
	_MinerCmd.AddCommand(
		_Miner_PingCmd,
		_Miner_PingCmd_gen,
		_Miner_InfoCmd,
		_Miner_InfoCmd_gen,
		_Miner_HandshakeCmd,
		_Miner_HandshakeCmd_gen,
		_Miner_SaveCmd,
		_Miner_SaveCmd_gen,
		_Miner_LoadCmd,
		_Miner_LoadCmd_gen,
		_Miner_StartCmd,
		_Miner_StartCmd_gen,
		_Miner_StopCmd,
		_Miner_StopCmd_gen,
		_Miner_TasksStatusCmd,
		_Miner_TasksStatusCmd_gen,
		_Miner_TaskDetailsCmd,
		_Miner_TaskDetailsCmd_gen,
		_Miner_TaskLogsCmd,
		_Miner_TaskLogsCmd_gen,
		_Miner_DiscoverHubCmd,
		_Miner_DiscoverHubCmd_gen,
	)
}

// End grpccmd

func init() { proto.RegisterFile("miner.proto", fileDescriptor9) }

var fileDescriptor9 = []byte{
	// 743 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x54, 0xed, 0x4e, 0xe3, 0x46,
	0x14, 0x8d, 0x83, 0xd3, 0x90, 0x6b, 0x20, 0x70, 0x01, 0xe1, 0xba, 0xa8, 0x8d, 0xac, 0xb6, 0xa4,
	0xad, 0x1a, 0xd1, 0xb4, 0x42, 0x2d, 0xe5, 0x0f, 0x25, 0x54, 0xa0, 0x42, 0x37, 0x72, 0x78, 0x81,
	0x49, 0x3c, 0x1b, 0xbc, 0x49, 0x3c, 0xde, 0x99, 0x09, 0x52, 0xde, 0x61, 0x9f, 0x68, 0xff, 0xec,
	0x1b, 0xed, 0x33, 0xac, 0xe6, 0xc3, 0xb1, 0x93, 0xcd, 0xfe, 0x9b, 0xb9, 0xf7, 0xdc, 0x3b, 0xe7,
	0x9c, 0x3b, 0x33, 0xe0, 0xcd, 0x92, 0x94, 0xf2, 0x4e, 0xc6, 0x99, 0x64, 0xe8, 0x0a, 0x96, 0xce,
	0x02, 0x1c, 0x91, 0x8c, 0x0c, 0x93, 0x69, 0x22, 0x13, 0x2a, 0x4c, 0x26, 0x68, 0x8e, 0x58, 0x2a,
	0x49, 0x01, 0x0d, 0x9a, 0x49, 0xaa, 0xc0, 0x69, 0x42, 0x6c, 0xa0, 0x91, 0x12, 0x69, 0x96, 0xe1,
	0x2b, 0x38, 0x7e, 0x54, 0xd0, 0x3b, 0x92, 0xc6, 0xe2, 0x99, 0x4c, 0x68, 0x44, 0xdf, 0xce, 0xa9,
	0x90, 0xb8, 0x0f, 0x5b, 0xcf, 0xf3, 0xa1, 0xef, 0xb4, 0x9c, 0x76, 0x23, 0x52, 0x4b, 0xfc, 0x1e,
	0x6a, 0x92, 0x88, 0x89, 0xf0, 0xab, 0xad, 0xad, 0xb6, 0xd7, 0xdd, 0xeb, 0xa8, 0xa6, 0x9d, 0x27,
	0x22, 0x26, 0xf7, 0xe9, 0x6b, 0x16, 0x99, 0x64, 0xf8, 0xce, 0x81, 0xc3, 0xf5, 0x8e, 0xd9, 0x74,
	0x81, 0x47, 0x50, 0xd3, 0xf4, 0x6d, 0x47, 0xb3, 0xc1, 0x0b, 0xd8, 0x29, 0x2b, 0xf0, 0xab, 0x2d,
	0xa7, 0xed, 0x75, 0xd1, 0xb4, 0xbe, 0x29, 0x65, 0xa2, 0x15, 0x1c, 0x9e, 0x41, 0x3d, 0x25, 0xf2,
	0x69, 0x91, 0x51, 0x7f, 0xab, 0xe5, 0xb4, 0xf7, 0xba, 0xbb, 0xa6, 0xe4, 0xff, 0xeb, 0x27, 0x15,
	0x8c, 0xf2, 0x6c, 0xf8, 0xd1, 0x81, 0x03, 0x4d, 0x67, 0x20, 0x09, 0x97, 0xb9, 0xb8, 0x3d, 0xa8,
	0x26, 0xb1, 0x65, 0x52, 0x4d, 0x62, 0xfc, 0x15, 0x1a, 0x4b, 0xd3, 0x2c, 0x87, 0xa6, 0xe5, 0x90,
	0x87, 0xa3, 0x02, 0x81, 0xff, 0xc0, 0x2e, 0xa7, 0x42, 0x35, 0xec, 0xb3, 0x69, 0x32, 0x5a, 0x68,
	0x0e, 0x5e, 0xf7, 0x74, 0xbd, 0xa4, 0x8c, 0x89, 0x56, 0x4b, 0xf0, 0x0a, 0x1a, 0x9c, 0x0a, 0x36,
	0xe7, 0x23, 0x2a, 0x7c, 0x57, 0xd7, 0x7f, 0x5b, 0x38, 0x1a, 0xd9, 0x94, 0x22, 0x9c, 0x70, 0x3a,
	0xa3, 0xa9, 0x14, 0x51, 0x51, 0x80, 0x3e, 0xd4, 0x19, 0x8f, 0x29, 0xbf, 0x8f, 0xfd, 0x9a, 0x56,
	0x91, 0x6f, 0xc3, 0x3f, 0x00, 0x06, 0x6c, 0x34, 0xa1, 0xf2, 0x3a, 0x8e, 0x39, 0x22, 0xb8, 0x24,
	0x8e, 0x73, 0xd3, 0xf5, 0x5a, 0xc5, 0x32, 0xc6, 0xa5, 0xd6, 0xb9, 0x1b, 0xe9, 0x75, 0xf8, 0xde,
	0x81, 0x66, 0xd9, 0x26, 0x35, 0xb1, 0xd3, 0xb2, 0x29, 0xa6, 0x41, 0xc9, 0x83, 0x2b, 0xa8, 0xab,
	0xca, 0x47, 0x92, 0xd9, 0xfb, 0x10, 0x1a, 0xf6, 0x6b, 0x5d, 0x3a, 0x7d, 0x03, 0xba, 0x4d, 0x25,
	0x5f, 0x44, 0x79, 0x49, 0xf0, 0x1f, 0xec, 0x94, 0x13, 0xea, 0xb6, 0x4d, 0xe8, 0x22, 0xbf, 0x6d,
	0x13, 0xba, 0xc0, 0x1f, 0xa0, 0xf6, 0x42, 0xa6, 0x73, 0xba, 0x3a, 0x8e, 0xdb, 0x34, 0xce, 0x58,
	0xa2, 0xcc, 0x30, 0xd9, 0xcb, 0xea, 0x9f, 0x4e, 0xf8, 0x06, 0xb6, 0xf3, 0x5b, 0x88, 0xbf, 0x41,
	0x9d, 0x9b, 0x21, 0xeb, 0x66, 0x5e, 0xf7, 0xe4, 0x73, 0x5a, 0x3a, 0x1d, 0xe5, 0x38, 0xfc, 0x05,
	0x6a, 0x5c, 0x51, 0xb5, 0x27, 0x1d, 0x6f, 0xd4, 0x11, 0x19, 0x4c, 0xf8, 0x37, 0x34, 0x96, 0x1c,
	0xb0, 0x03, 0x0d, 0x9a, 0x6f, 0x7c, 0x47, 0xbb, 0xb0, 0x6f, 0xaa, 0x8b, 0x11, 0x44, 0x05, 0x24,
	0x3c, 0xb1, 0x8f, 0x6d, 0x20, 0x89, 0x9c, 0x8b, 0x47, 0x92, 0x59, 0x2e, 0xe1, 0x19, 0x78, 0x03,
	0xf2, 0xb2, 0x7c, 0x7b, 0x3e, 0xd4, 0x93, 0x19, 0x19, 0xd3, 0xfb, 0x9e, 0x75, 0x24, 0xdf, 0x76,
	0x3f, 0xb8, 0x50, 0xd3, 0x2d, 0xf0, 0x47, 0x70, 0xfb, 0x49, 0x3a, 0x46, 0xcf, 0x1a, 0x33, 0xcb,
	0xe4, 0x22, 0xb0, 0x2e, 0xa9, 0x84, 0x66, 0x1d, 0x56, 0x14, 0x4e, 0x1b, 0xb3, 0x09, 0xa7, 0xdf,
	0xad, 0xc5, 0xdd, 0x42, 0x63, 0xf9, 0x62, 0xf1, 0x9b, 0x92, 0x07, 0xeb, 0x3f, 0x43, 0xf0, 0xf5,
	0xe6, 0xa4, 0x69, 0xf3, 0x33, 0xb8, 0x4a, 0x09, 0x1e, 0x58, 0x1f, 0x0a, 0x55, 0x81, 0x65, 0x70,
	0xf3, 0x3c, 0x4f, 0x27, 0x61, 0xe5, 0xdc, 0xc1, 0x9f, 0xc0, 0x7d, 0x60, 0x24, 0xc6, 0x72, 0x22,
	0xb0, 0xdf, 0x4a, 0x9f, 0xb3, 0x31, 0xa7, 0x42, 0x84, 0x95, 0xb6, 0x73, 0xee, 0xe0, 0x5f, 0x50,
	0xd3, 0xb3, 0xc0, 0x2f, 0x8d, 0x33, 0xd8, 0x3c, 0xb6, 0xb0, 0x82, 0xdf, 0x81, 0x3b, 0x90, 0x2c,
	0xc3, 0x6d, 0xab, 0xb9, 0x17, 0x94, 0xad, 0x08, 0x2b, 0xf8, 0x2f, 0x78, 0xea, 0xfa, 0x08, 0x33,
	0x95, 0x15, 0xed, 0xeb, 0x83, 0x0a, 0x8e, 0xac, 0xac, 0x22, 0xae, 0x0f, 0xd1, 0x1c, 0xcf, 0x4d,
	0x9f, 0x1e, 0x95, 0x24, 0x99, 0x8a, 0xd2, 0x79, 0xc7, 0xc5, 0xbb, 0x36, 0x85, 0x39, 0xb5, 0x4b,
	0x73, 0x71, 0x1f, 0xd8, 0x58, 0x60, 0x09, 0xa4, 0xf6, 0xf9, 0x81, 0x87, 0xab, 0xe1, 0xc2, 0xbc,
	0x0b, 0xf0, 0x7a, 0x89, 0x18, 0xb1, 0x17, 0xca, 0xef, 0xe6, 0x43, 0xf4, 0x0d, 0xae, 0x14, 0x5a,
	0xb3, 0xdd, 0xaa, 0x1d, 0x7e, 0xa5, 0xff, 0xfd, 0xdf, 0x3f, 0x05, 0x00, 0x00, 0xff, 0xff, 0x99,
	0xe7, 0x01, 0x54, 0x4d, 0x06, 0x00, 0x00,
}
