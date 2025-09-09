package codec

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"

	errs "github.com/fengzhongzhu1621/xgo/xerror"
)

var _ IMsg = (*msg)(nil)
var msgPool = sync.Pool{
	New: func() interface{} {
		return &msg{}
	},
}

// msg is the context of rpc.
type msg struct {
	context             context.Context
	frameHead           interface{}
	requestTimeout      time.Duration
	serializationType   int
	compressType        int
	streamID            uint32
	dyeing              bool
	dyeingKey           string
	serverRPCName       string
	clientRPCName       string
	serverMetaData      MetaData
	clientMetaData      MetaData
	callerServiceName   string
	calleeServiceName   string
	calleeContainerName string
	serverRspErr        error
	clientRspErr        error
	serverReqHead       interface{}
	serverRspHead       interface{}
	clientReqHead       interface{}
	clientRspHead       interface{}
	localAddr           net.Addr
	remoteAddr          net.Addr
	logger              interface{}
	callerApp           string
	callerServer        string
	callerService       string
	callerMethod        string
	calleeApp           string
	calleeServer        string
	calleeService       string
	calleeMethod        string
	namespace           string
	setName             string
	envName             string
	envTransfer         string
	requestID           uint32
	calleeSetName       string
	streamFrame         interface{}
	commonMeta          CommonMeta
	callType            RequestType
}

// resetDefault reset all fields of msg to default value.
func (m *msg) resetDefault() {
	m.context = nil
	m.frameHead = nil
	m.requestTimeout = 0
	m.serializationType = 0
	m.compressType = 0
	m.dyeing = false
	m.dyeingKey = ""
	m.serverRPCName = ""
	m.clientRPCName = ""
	m.serverMetaData = nil
	m.clientMetaData = nil
	m.callerServiceName = ""
	m.calleeServiceName = ""
	m.calleeContainerName = ""
	m.serverRspErr = nil
	m.clientRspErr = nil
	m.serverReqHead = nil
	m.serverRspHead = nil
	m.clientReqHead = nil
	m.clientRspHead = nil
	m.localAddr = nil
	m.remoteAddr = nil
	m.logger = nil
	m.callerApp = ""
	m.callerServer = ""
	m.callerService = ""
	m.callerMethod = ""
	m.calleeApp = ""
	m.calleeServer = ""
	m.calleeService = ""
	m.calleeMethod = ""
	m.namespace = ""
	m.setName = ""
	m.envName = ""
	m.envTransfer = ""
	m.requestID = 0
	m.streamFrame = nil
	m.streamID = 0
	m.calleeSetName = ""
	m.commonMeta = nil
	m.callType = 0
}

// Context restores old context when create new msg.
func (m *msg) Context() context.Context {
	return m.context
}

// WithNamespace set server's namespace.
func (m *msg) WithNamespace(namespace string) {
	m.namespace = namespace
}

// Namespace returns namespace.
func (m *msg) Namespace() string {
	return m.namespace
}

// WithEnvName sets environment.
func (m *msg) WithEnvName(envName string) {
	m.envName = envName
}

// WithSetName sets set name.
func (m *msg) WithSetName(setName string) {
	m.setName = setName
}

// SetName returns set name.
func (m *msg) SetName() string {
	return m.setName
}

// WithCalleeSetName sets the callee set name.
func (m *msg) WithCalleeSetName(s string) {
	m.calleeSetName = s
}

// CalleeSetName returns the callee set name.
func (m *msg) CalleeSetName() string {
	return m.calleeSetName
}

// EnvName returns environment.
func (m *msg) EnvName() string {
	return m.envName
}

// WithEnvTransfer sets environment transfer value.
func (m *msg) WithEnvTransfer(envTransfer string) {
	m.envTransfer = envTransfer
}

// EnvTransfer returns environment transfer value.
func (m *msg) EnvTransfer() string {
	return m.envTransfer
}

// WithRemoteAddr sets remote address.
func (m *msg) WithRemoteAddr(addr net.Addr) {
	m.remoteAddr = addr
}

// WithLocalAddr set local address.
func (m *msg) WithLocalAddr(addr net.Addr) {
	m.localAddr = addr
}

// RemoteAddr returns remote address.
func (m *msg) RemoteAddr() net.Addr {
	return m.remoteAddr
}

// LocalAddr returns local address.
func (m *msg) LocalAddr() net.Addr {
	return m.localAddr
}

// RequestTimeout returns request timeout set by
// upstream business protocol.
func (m *msg) RequestTimeout() time.Duration {
	return m.requestTimeout
}

// WithRequestTimeout sets request timeout.
func (m *msg) WithRequestTimeout(t time.Duration) {
	m.requestTimeout = t
}

// FrameHead returns frame head.
func (m *msg) FrameHead() interface{} {
	return m.frameHead
}

// WithFrameHead sets frame head.
func (m *msg) WithFrameHead(f interface{}) {
	m.frameHead = f
}

// SerializationType returns the value of body serialization, which is
// defined in serialization.go.
func (m *msg) SerializationType() int {
	return m.serializationType
}

// WithSerializationType sets body serialization type of body.
func (m *msg) WithSerializationType(t int) {
	m.serializationType = t
}

// CompressType returns compress type value, which is defined in compress.go.
func (m *msg) CompressType() int {
	return m.compressType
}

// WithCompressType sets compress type.
func (m *msg) WithCompressType(t int) {
	m.compressType = t
}

// ServerRPCName returns server rpc name.
func (m *msg) ServerRPCName() string {
	return m.serverRPCName
}

// ServerMetaData returns server meta data, which is passed to server.
func (m *msg) ServerMetaData() MetaData {
	return m.serverMetaData
}

// WithServerMetaData sets server meta data.
func (m *msg) WithServerMetaData(d MetaData) {
	if d == nil {
		d = MetaData{}
	}
	m.serverMetaData = d
}

// ClientMetaData returns client meta data, which will pass to downstream.
func (m *msg) ClientMetaData() MetaData {
	return m.clientMetaData
}

// WithClientMetaData set client meta data.
func (m *msg) WithClientMetaData(d MetaData) {
	if d == nil {
		d = MetaData{}
	}
	m.clientMetaData = d
}

// CalleeServiceName returns callee service name.
func (m *msg) CalleeServiceName() string {
	return m.calleeServiceName
}

// CalleeContainerName returns callee container name.
func (m *msg) CalleeContainerName() string {
	return m.calleeContainerName
}

// WithCalleeContainerName sets callee container name.
func (m *msg) WithCalleeContainerName(s string) {
	m.calleeContainerName = s
}

// WithStreamFrame sets stream frame.
func (m *msg) WithStreamFrame(i interface{}) {
	m.streamFrame = i
}

// StreamFrame returns stream frame.
func (m *msg) StreamFrame() interface{} {
	return m.streamFrame
}

// CallerServiceName returns caller service name.
func (m *msg) CallerServiceName() string {
	return m.callerServiceName
}

// WithServerRspErr sets server response error.
func (m *msg) WithServerRspErr(e error) {
	m.serverRspErr = e
}

// WithStreamID sets stream id.
func (m *msg) WithStreamID(streamID uint32) {
	m.streamID = streamID
}

// StreamID returns stream id.
func (m *msg) StreamID() uint32 {
	return m.streamID
}

// ClientRspErr returns client response error, which created when client call downstream.
func (m *msg) ClientRspErr() error {
	return m.clientRspErr
}

// WithClientRspErr sets client response err, this method will called
// when client parse response package.
func (m *msg) WithClientRspErr(e error) {
	m.clientRspErr = e
}

// ServerReqHead returns the package head of request
func (m *msg) ServerReqHead() interface{} {
	return m.serverReqHead
}

// WithServerReqHead sets the package head of request
func (m *msg) WithServerReqHead(h interface{}) {
	m.serverReqHead = h
}

// ServerRspHead returns the package head of response
func (m *msg) ServerRspHead() interface{} {
	return m.serverRspHead
}

// WithServerRspHead sets the package head returns to upstream
func (m *msg) WithServerRspHead(h interface{}) {
	m.serverRspHead = h
}

// ClientReqHead returns the request package head of client,
// this is set only when cross protocol call.
func (m *msg) ClientReqHead() interface{} {
	return m.clientReqHead
}

// WithClientReqHead sets the request package head of client.
func (m *msg) WithClientReqHead(h interface{}) {
	m.clientReqHead = h
}

// ClientRspHead returns the request package head of client.
func (m *msg) ClientRspHead() interface{} {
	return m.clientRspHead
}

// WithClientRspHead sets the response package head of client.
func (m *msg) WithClientRspHead(h interface{}) {
	m.clientRspHead = h
}

// Dyeing return the dyeing mark.
func (m *msg) Dyeing() bool {
	return m.dyeing
}

// WithDyeing sets the dyeing mark.
func (m *msg) WithDyeing(dyeing bool) {
	m.dyeing = dyeing
}

// DyeingKey returns the dyeing key.
func (m *msg) DyeingKey() string {
	return m.dyeingKey
}

// WithDyeingKey sets the dyeing key.
func (m *msg) WithDyeingKey(key string) {
	m.dyeingKey = key
}

// CallerApp returns caller app.
func (m *msg) CallerApp() string {
	return m.callerApp
}

// WithCallerApp sets caller app.
func (m *msg) WithCallerApp(app string) {
	m.callerApp = app
}

// CallerServer returns caller server.
func (m *msg) CallerServer() string {
	return m.callerServer
}

// WithCallerServer sets caller server.
func (m *msg) WithCallerServer(s string) {
	m.callerServer = s
}

// CallerService returns caller service.
func (m *msg) CallerService() string {
	return m.callerService
}

// WithCallerService sets caller service.
func (m *msg) WithCallerService(s string) {
	m.callerService = s
}

// WithCallerMethod sets caller method.
func (m *msg) WithCallerMethod(s string) {
	m.callerMethod = s
}

// CallerMethod returns caller method.
func (m *msg) CallerMethod() string {
	return m.callerMethod
}

// CalleeApp returns caller app.
func (m *msg) CalleeApp() string {
	return m.calleeApp
}

// WithCalleeApp sets callee app.
func (m *msg) WithCalleeApp(app string) {
	m.calleeApp = app
}

// CalleeServer returns callee server.
func (m *msg) CalleeServer() string {
	return m.calleeServer
}

// WithCalleeServer sets callee server.
func (m *msg) WithCalleeServer(s string) {
	m.calleeServer = s
}

// CalleeService returns callee service.
func (m *msg) CalleeService() string {
	return m.calleeService
}

// WithCalleeService sets callee service.
func (m *msg) WithCalleeService(s string) {
	m.calleeService = s
}

// WithCalleeMethod sets callee method.
func (m *msg) WithCalleeMethod(s string) {
	m.calleeMethod = s
}

// CalleeMethod returns callee method.
func (m *msg) CalleeMethod() string {
	return m.calleeMethod
}

// WithLogger sets logger into context message. Generally, the logger is
// created from WithFields() method.
func (m *msg) WithLogger(l interface{}) {
	m.logger = l
}

// Logger returns logger from context message.
func (m *msg) Logger() interface{} {
	return m.logger
}

// WithRequestID sets request id.
func (m *msg) WithRequestID(id uint32) {
	m.requestID = id
}

// RequestID returns request id.
func (m *msg) RequestID() uint32 {
	return m.requestID
}

// WithCommonMeta sets common meta data.
func (m *msg) WithCommonMeta(c CommonMeta) {
	m.commonMeta = c
}

// CommonMeta returns common meta data.
func (m *msg) CommonMeta() CommonMeta {
	return m.commonMeta
}

// WithCallType sets type of call.
func (m *msg) WithCallType(t RequestType) {
	m.callType = t
}

// CallType returns type of call.
func (m *msg) CallType() RequestType {
	return m.callType
}

// ClientRPCName returns client rpc name.
func (m *msg) ClientRPCName() string {
	return m.clientRPCName
}

// WithServerRPCName sets server rpc name.
func (m *msg) WithServerRPCName(s string) {
	if m.serverRPCName == s {
		return
	}
	m.serverRPCName = s
	m.updateMethodNameUsingRPCName(s)
}

// WithCallerServiceName sets caller service name.
func (m *msg) WithCallerServiceName(s string) {
	if m.callerServiceName == s {
		return
	}
	m.callerServiceName = s
	if s == "*" {
		return
	}
	app, server, service := getAppServerService(s)
	m.WithCallerApp(app)
	m.WithCallerServer(server)
	m.WithCallerService(service)
}

// WithCalleeServiceName sets callee service name.
func (m *msg) WithCalleeServiceName(s string) {
	if m.calleeServiceName == s {
		return
	}
	m.calleeServiceName = s
	if s == "*" {
		return
	}
	app, server, service := getAppServerService(s)
	m.WithCalleeApp(app)
	m.WithCalleeServer(server)
	m.WithCalleeService(service)
}

// ServerRspErr returns server response error, which is created by handler.
func (m *msg) ServerRspErr() *errs.Error {
	if m.serverRspErr == nil {
		return nil
	}
	e, ok := m.serverRspErr.(*errs.Error)
	if !ok {
		return &errs.Error{
			Type: errs.ErrorTypeBusiness,
			Code: errs.RetUnknown,
			Msg:  m.serverRspErr.Error(),
		}
	}

	return e
}

// WithClientRPCName sets client rpc name, which will be called
// by client stub.
func (m *msg) WithClientRPCName(s string) {
	if m.clientRPCName == s {
		return
	}
	m.clientRPCName = s
	m.updateMethodNameUsingRPCName(s)
}

func (m *msg) updateMethodNameUsingRPCName(s string) {
	// 检查给定的字符串是否符合 trpc 格式规范
	if rpcNameIsTRPCForm(s) {
		m.WithCalleeMethod(methodFromRPCName(s))
		return
	}
	if m.CalleeMethod() == "" {
		m.WithCalleeMethod(s)
	}
}

// methodFromRPCName returns the method parsed from rpc string.
func methodFromRPCName(s string) string {
	return s[strings.LastIndex(s, "/")+1:]
}

// getAppServerService 根据服务名称字符串解析出应用名、服务名和方法名
// 服务名称示例: trpc.app.server.service
func getAppServerService(s string) (app, server, service string) {
	// 检查字符串中点的数量是否足够（至少2个点，即3个部分）
	if strings.Count(s, ".") >= ServiceSectionLength-1 {
		// 找到第一个点的位置并加1，得到应用名的起始位置
		i := strings.Index(s, ".") + 1
		// 从应用名后开始找第二个点的位置，并计算服务名的起始位置
		j := strings.Index(s[i:], ".") + i + 1
		// 从服务名后开始找第三个点的位置，并计算方法名的起始位置
		k := strings.Index(s[j:], ".") + j + 1
		// 提取应用名（第一个点和第二个点之间的内容）
		app = s[i : j-1]
		// 提取服务名（第二个点和第三个点之间的内容）
		server = s[j : k-1]
		// 提取方法名（第三个点之后的所有内容）
		service = s[k:]
		return
	}
	// 处理点数量不足的情况：先尝试解析应用名

	// app
	i := strings.Index(s, ".")
	if i == -1 {
		app = s
		return
	}
	app = s[:i]

	// server
	i++
	j := strings.Index(s[i:], ".")
	if j == -1 {
		server = s[i:]
		return
	}
	j += i + 1
	server = s[i : j-1]

	// service
	service = s[j:]

	return
}

// rpcNameIsTRPCForm 检查给定的字符串是否符合 trpc 格式规范。
// 其功能等价于以下正则表达式（但出于性能考虑，使用字符串操作实现）：
//
//	var r = regexp.MustCompile(`^/[^/.]+\.[^/]+/[^/.]+$`)
//
//	func rpcNameIsTRPCForm(s string) bool {
//		return r.MatchString(s)
//	}
//
// 经测试，正则表达式版本比当前实现慢很多。
// 参考 message_bench_test.go 中的 BenchmarkRPCNameIsTRPCForm 基准测试结果。
func rpcNameIsTRPCForm(s string) bool {
	// 检查字符串是否为空
	if len(s) == 0 {
		return false // 空字符串不符合格式
	}
	// 检查字符串是否以斜杠 '/' 开头 (对应正则中的 ^/)
	if s[0] != '/' {
		return false // 首字符不是斜杠，不符合格式
	}
	// 定义起始位置常量，跳过开头的 '/'
	const start = 1
	// 在起始位置之后的子串中查找第一个点号 '.' 的位置
	firstDot := strings.Index(s[start:], ".")
	// 检查是否找到点号，且点号不能是第一个字符 (对应正则中的 [^.]+\.)
	// firstDot == -1 表示没找到点号，firstDot == 0 表示点号紧接在斜杠后
	if firstDot == -1 || firstDot == 0 {
		return false
	}
	// 检查从起始位置到第一个点号之间的子串（应用名部分）是否包含斜杠 '/'
	// 此部分应不包含斜杠 (对应正则中的 [^/]+\.)
	if strings.Contains(s[start:start+firstDot], "/") {
		return false
	}
	// 在第一个点号之后的子串中查找第二个斜杠 '/' 的位置
	secondSlash := strings.Index(s[start+firstDot:], "/")
	// 检查是否找到第二个斜杠，且第二个斜杠不能紧跟在点号之后 (对应正则中的 [^/]+/)
	// secondSlash == -1 表示没找到，secondSlash == 1 表示斜杠紧接在点号后
	if secondSlash == -1 || secondSlash == 1 {
		return false
	}
	// 计算第二个斜杠在原始字符串中的绝对位置
	// 检查该位置是否是字符串的最后一个字符（即方法名部分不能为空）
	if start+firstDot+secondSlash == len(s)-1 {
		return false
	}
	// 定义偏移量常量，用于跳过第二个斜杠
	const offset = 1
	// 检查第二个斜杠之后的部分（方法名部分）是否包含任何斜杠 '/' 或点号 '.'
	// 此部分应不包含斜杠或点号 (对应正则中的 [^/.]+$)
	if strings.ContainsAny(s[start+firstDot+secondSlash+offset:], "/.") {
		return false
	}
	// 所有检查均通过，字符串符合 trpc 格式
	return true
}
