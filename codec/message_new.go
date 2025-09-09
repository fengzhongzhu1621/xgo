package codec

import "context"

// Message returns the message of context.
func Message(ctx context.Context) IMsg {
	val := ctx.Value(ContextKeyMessage)
	m, ok := val.(*msg)
	if !ok {
		return &msg{context: ctx}
	}
	return m
}

// WithNewMessage create a new empty message, and put it into ctx,
func WithNewMessage(ctx context.Context) (context.Context, IMsg) {
	m := msgPool.Get().(*msg)
	// 在上下文传递 msg
	ctx = context.WithValue(ctx, ContextKeyMessage, m)
	m.context = ctx
	return ctx, m
}

// PutBackMessage return struct Message to sync pool,
// and reset all the members of Message to default
func PutBackMessage(sourceMsg IMsg) {
	m, ok := sourceMsg.(*msg)
	if !ok {
		return
	}
	m.resetDefault()
	msgPool.Put(m)
}

// EnsureMessage returns context and message, if there is a message in context,
// returns the original one, if not, returns a new one.
func EnsureMessage(ctx context.Context) (context.Context, IMsg) {
	// 先从上下文获取消息
	val := ctx.Value(ContextKeyMessage)
	if m, ok := val.(*msg); ok {
		return ctx, m
	}

	// 如果上下文中没有消息，创建一个新消息（默认值）
	return WithNewMessage(ctx)
}

// WithCloneContextAndMessage creates a new context, then copy the message of current context
// into new context, this method will return the new context and message for stream mod.
func WithCloneContextAndMessage(ctx context.Context) (context.Context, IMsg) {
	newMsg := msgPool.Get().(*msg)

	// 从当前上下文中获取消息对象
	val := ctx.Value(ContextKeyMessage)

	// 创建新的上下文
	newCtx := context.Background()
	newCtx = context.WithValue(newCtx, ContextKeyMessage, newMsg)
	newMsg.context = newCtx

	m, ok := val.(*msg)
	if !ok {
		return newCtx, newMsg
	}

	copyCommonMessage(m, newMsg)
	copyServerToServerMessage(m, newMsg)

	return newCtx, newMsg
}

// WithCloneMessage copy a new message and put into context, each rpc call should
// create a new message, this method will be called by client stub.
func WithCloneMessage(ctx context.Context) (context.Context, IMsg) {
	newMsg := msgPool.Get().(*msg)

	// 从当前上下文中获取消息对象
	val := ctx.Value(ContextKeyMessage)

	// 更新上下文中的消息对象
	ctx = context.WithValue(ctx, ContextKeyMessage, newMsg)
	newMsg.context = ctx

	m, ok := val.(*msg)
	if !ok {
		return ctx, newMsg
	}

	copyCommonMessage(m, newMsg)
	copyServerToClientMessage(m, newMsg)

	return ctx, newMsg
}

// copyCommonMessage copy common data of message.
func copyCommonMessage(m *msg, newMsg *msg) {
	// Do not copy compress type here, as it will cause subsequence RPC calls to inherit the upstream
	// compress type which is not the expected behavior. Compress type should not be propagated along
	// the entire RPC invocation chain.
	newMsg.frameHead = m.frameHead
	newMsg.requestTimeout = m.requestTimeout
	newMsg.serializationType = m.serializationType
	newMsg.serverRPCName = m.serverRPCName
	newMsg.clientRPCName = m.clientRPCName
	newMsg.serverReqHead = m.serverReqHead
	newMsg.serverRspHead = m.serverRspHead
	newMsg.dyeing = m.dyeing
	newMsg.dyeingKey = m.dyeingKey
	newMsg.serverMetaData = m.serverMetaData.Clone()
	newMsg.logger = m.logger
	newMsg.namespace = m.namespace
	newMsg.envName = m.envName
	newMsg.setName = m.setName
	newMsg.envTransfer = m.envTransfer
	newMsg.commonMeta = m.commonMeta.Clone()
}

// copyClientMessage copy the message transferred from server to client.
func copyServerToClientMessage(m *msg, newMsg *msg) {
	newMsg.clientMetaData = m.serverMetaData.Clone()
	// clone this message for downstream client, so caller is equal to callee.
	newMsg.callerServiceName = m.calleeServiceName
	newMsg.callerApp = m.calleeApp
	newMsg.callerServer = m.calleeServer
	newMsg.callerService = m.calleeService
	newMsg.callerMethod = m.calleeMethod
}

func copyServerToServerMessage(m *msg, newMsg *msg) {
	newMsg.callerServiceName = m.callerServiceName
	newMsg.callerApp = m.callerApp
	newMsg.callerServer = m.callerServer
	newMsg.callerService = m.callerService
	newMsg.callerMethod = m.callerMethod

	newMsg.calleeServiceName = m.calleeServiceName
	newMsg.calleeService = m.calleeService
	newMsg.calleeApp = m.calleeApp
	newMsg.calleeServer = m.calleeServer
	newMsg.calleeMethod = m.calleeMethod
}
