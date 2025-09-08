package codec

import (
	"context"
	"net"
	"time"

	"trpc.group/trpc-go/trpc-go/errs"
)

// RequestType is the type of client request, such as SendAndRecv，SendOnly.
type RequestType int

// MetaData is request penetrate message.
type MetaData map[string][]byte

// CommonMeta is common meta message.
type CommonMeta map[interface{}]interface{}

// Msg defines core message data for multi protocol, business protocol
// should set this message when packing and unpacking data.
type IMsg interface {
	// Context returns rpc context
	Context() context.Context

	// WithRemoteAddr sets upstream address for server,
	// or downstream address for client.
	WithRemoteAddr(addr net.Addr)

	// WithLocalAddr sets server local address.
	WithLocalAddr(addr net.Addr)

	// RemoteAddr returns upstream address for server,
	// or downstream address for client.
	RemoteAddr() net.Addr

	// LocalAddr returns server local address.
	LocalAddr() net.Addr

	// WithNamespace sets server namespace.
	WithNamespace(string)

	// Namespace returns server namespace.
	Namespace() string

	// WithEnvName sets server environment.
	WithEnvName(string)

	// EnvName returns server environment.
	EnvName() string

	// WithSetName sets server set name.
	WithSetName(string)

	// SetName returns server set name.
	SetName() string

	// WithEnvTransfer sets environment message for transfer.
	WithEnvTransfer(string)

	// EnvTransfer returns environment message for transfer.
	EnvTransfer() string

	// WithRequestTimeout sets the upstream timeout for server,
	// or downstream timeout for client.
	WithRequestTimeout(time.Duration)

	// RequestTimeout returns the upstream timeout for server,
	// or downstream timeout for client.
	RequestTimeout() time.Duration

	// WithSerializationType sets serialization type.
	WithSerializationType(int)

	// SerializationType returns serialization type.
	SerializationType() int

	// WithCompressType sets compress type.
	WithCompressType(int)

	// CompressType returns compress type.
	CompressType() int

	// WithServerRPCName sets server handler method name.
	WithServerRPCName(string)

	// WithClientRPCName sets client rpc name for downstream.
	WithClientRPCName(string)

	// ServerRPCName returns method name of current server handler name,
	// such as /trpc.app.server.service/method.
	ServerRPCName() string

	// ClientRPCName returns method name of downstream interface.
	ClientRPCName() string

	// WithCallerServiceName sets caller service name.
	WithCallerServiceName(string)

	// WithCalleeServiceName sets callee service name.
	WithCalleeServiceName(string)

	// WithCallerApp sets caller app. For server this app is upstream app,
	// but for client, is its own app.
	WithCallerApp(string)

	// WithCallerServer sets caller server. For server this server is upstream server,
	// but for client, is its own server.
	WithCallerServer(string)

	// WithCallerService sets caller service, For server this service is upstream service,
	// but for client, is its own service.
	WithCallerService(string)

	// WithCallerMethod sets caller method, For server this mothod is upstream mothod,
	// but for client, is its own method.
	WithCallerMethod(string)

	// WithCalleeApp sets callee app. For server, this app is its own app,
	// but for client, is downstream's app.
	WithCalleeApp(string)

	// WithCalleeServer sets callee server. For server, this server is its own server,
	// but for client, is downstream's server.
	WithCalleeServer(string)

	// WithCalleeService sets callee service. For server, this service is its own service,
	// but for client, is downstream's service.
	WithCalleeService(string)

	// WithCalleeMethod sets callee method. For server, this method is its own method,
	// but for client, is downstream's method.
	WithCalleeMethod(string)

	// CallerServiceName returns caller service name, such as trpc.app.server.service.
	// For server, this name is upstream's service name, but for client, is its own service name.
	CallerServiceName() string

	// CallerApp returns caller app. For server, this app is upstream's app,
	// but for client, is its own app.
	CallerApp() string

	// CallerServer returns caller server. For server, this is upstream's server,
	// but for client, is its own server.
	CallerServer() string

	// CallerService returns caller service. For server, this service is upstream's service,
	// but for client, is its own service.
	CallerService() string

	// CallerMethod returns caller method. For server, this method is upstream's method,
	// but for client, is its own method.
	CallerMethod() string

	// CalleeServiceName returns callee service name. For server, this name is its own service name,
	// but for client, is downstream's service name.
	CalleeServiceName() string

	// CalleeApp returns callee app. For server, this app is its own app,
	// but for client, is downstream's app.
	CalleeApp() string

	// CalleeServer returns callee server. For server, this server name is its own name,
	// but for client, is downstream's server name.
	CalleeServer() string

	// CalleeService returns callee service. For server, this service is its own service,
	// but for client, is downstream's service.
	CalleeService() string

	// CalleeMethod returns callee method. For server, this method is its own method,
	// but for client, is downstream's method.
	CalleeMethod() string

	// CalleeContainerName sets callee container name.
	CalleeContainerName() string

	// WithCalleeContainerName return callee container name.
	WithCalleeContainerName(string)

	// WithServerMetaData sets server meta data.
	WithServerMetaData(MetaData)

	// ServerMetaData returns server meta data.
	ServerMetaData() MetaData

	// WithFrameHead sets frame head.
	WithFrameHead(interface{})

	// FrameHead returns frame head.
	FrameHead() interface{}

	// WithServerReqHead sets server request head.
	WithServerReqHead(interface{})

	// ServerReqHead returns server request head.
	ServerReqHead() interface{}

	// WithServerRspHead sets server response head, this head will return to upstream.
	WithServerRspHead(interface{})

	// ServerRspHead returns server response head, this head will return to upstream.
	ServerRspHead() interface{}

	// WithDyeing sets dyeing mark.
	WithDyeing(bool)

	// Dyeing returns dyeing mark.
	Dyeing() bool

	// WithDyeingKey sets dyeing key.
	WithDyeingKey(string)

	// DyeingKey returns dyeing key.
	DyeingKey() string

	// WithServerRspErr sets response error for server.
	WithServerRspErr(error)

	// ServerRspErr returns response error for server.
	ServerRspErr() *errs.Error

	// WithClientMetaData sets client meta data.
	WithClientMetaData(MetaData)

	// ClientMetaData returns client meta data.
	ClientMetaData() MetaData

	// WithClientReqHead sets client request head.
	WithClientReqHead(interface{})

	// ClientReqHead returns client request head.
	ClientReqHead() interface{}

	// WithClientRspErr sets response error for client.
	WithClientRspErr(error)

	// ClientRspErr returns response error for client.
	ClientRspErr() error

	// WithClientRspHead sets response head for client.
	WithClientRspHead(interface{})

	// ClientRspHead returns response head for client.
	ClientRspHead() interface{}

	// WithLogger sets logger into context.
	WithLogger(interface{})

	// Logger returns logger from context.
	Logger() interface{}

	// WithRequestID sets request id.
	WithRequestID(uint32)

	// RequestID returns request id.
	RequestID() uint32

	// WithStreamID sets stream id.
	WithStreamID(uint32)

	// StreamID return stream id.
	StreamID() uint32

	// StreamFrame sets stream frame.
	StreamFrame() interface{}

	// WithStreamFrame returns stream frame.
	WithStreamFrame(interface{})

	// WithCalleeSetName sets callee set name.
	WithCalleeSetName(string)

	// CalleeSetName returns callee set name.
	CalleeSetName() string

	// WithCommonMeta sets common meta data.
	WithCommonMeta(CommonMeta)

	// CommonMeta returns common meta data.
	CommonMeta() CommonMeta

	// WithCallType sets call type.
	WithCallType(RequestType)

	// CallType returns call type.
	CallType() RequestType
}
