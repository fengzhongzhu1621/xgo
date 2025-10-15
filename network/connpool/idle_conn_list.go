package connpool

// connList 维护空闲连接，使用栈结构来管理连接
//
// 栈方法相比队列方法有一个优势：当请求量相对较小但请求分布仍然相对均匀时，
// 队列方法会导致被占用的连接延迟释放，而栈方法可以更快地重用最近使用的连接
type connList struct {
	count      int       // 当前列表中的连接数量
	head, tail *PoolConn // 列表的头节点和尾节点指针
}

// pushHead 将连接添加到列表头部（栈顶）
// 参数:
//
//	pc: 要添加的池连接
func (l *connList) pushHead(pc *PoolConn) {
	// 标记连接在池中
	pc.inPool = true

	// 添加到列表头部（栈顶）
	pc.next = l.head // 新连接的next指向当前头节点
	pc.prev = nil    // 新连接是头节点，prev为nil

	if l.count == 0 {
		l.tail = pc // 如果列表为空，尾节点也是新连接
	} else {
		// 添加双向链表支持
		l.head.prev = pc // 当前头节点的prev指向新连接
	}

	l.count++   // 连接计数增加
	l.head = pc // 更新头节点为新连接
}

// popHead 从列表头部（栈顶）移除并返回连接
func (l *connList) popHead() {
	pc := l.head // 获取当前头节点
	l.count--    // 连接计数减少
	if l.count == 0 {
		l.head, l.tail = nil, nil // 如果列表为空，重置头尾节点
	} else {
		pc.next.prev = nil // 新头节点的prev设为nil
		l.head = pc.next   // 更新头节点为下一个节点
	}
	pc.next, pc.prev = nil, nil // 断开移除连接的指针
	pc.inPool = false           // 标记连接不在池中
}

// pushTail 将连接添加到列表尾部
// 参数:
//
//	pc: 要添加的池连接
func (l *connList) pushTail(pc *PoolConn) {
	pc.inPool = true // 标记连接在池中
	pc.next = nil    // 新连接是尾节点，next为nil
	pc.prev = l.tail // 新连接的prev指向当前尾节点
	if l.count == 0 {
		l.head = pc // 如果列表为空，头节点也是新连接
	} else {
		l.tail.next = pc // 当前尾节点的next指向新连接
	}
	l.count++   // 连接计数增加
	l.tail = pc // 更新尾节点为新连接
}

// popTail 从列表尾部移除并返回连接
func (l *connList) popTail() {
	pc := l.tail // 获取当前尾节点
	l.count--    // 连接计数减少
	if l.count == 0 {
		l.head, l.tail = nil, nil // 如果列表为空，重置头尾节点
	} else {
		pc.prev.next = nil // 新尾节点的next设为nil
		l.tail = pc.prev   // 更新尾节点为前一个节点
	}
	pc.next, pc.prev = nil, nil // 断开移除连接的指针
	pc.inPool = false           // 标记连接不在池中
}
