// Package stack provides a non-thread-safe stack.
// 非并发安全的栈实现
package linkstack

// Stack is a non-thread-safe stack.
// 基于链表实现的非线程安全栈
type Stack[T any] struct {
	top  *node[T] // 栈顶节点指针
	size int      // 栈中元素数量
}

type node[T any] struct {
	value T        // 节点存储的值
	prev  *node[T] // 指向前一个节点的指针
}

// New creates a stack.
// 创建新的空栈实例
func New[T any]() *Stack[T] {
	return &Stack[T]{} // 返回初始化的空栈
}

// Size returns the stack size.
// 返回栈中元素的数量
func (st *Stack[T]) Size() int {
	return st.size // 直接返回size字段
}

// Reset resets the stack.
// 重置栈，清空所有元素
func (st *Stack[T]) Reset() {
	st.top = nil // 将栈顶指针设为nil
	st.size = 0  // 将元素数量重置为0
}

// Push pushes an element onto the stack.
// 将元素压入栈顶
func (st *Stack[T]) Push(value T) {
	newNode := &node[T]{
		value: value,  // 设置新节点的值
		prev:  st.top, // 新节点的prev指向当前栈顶
	}
	st.top = newNode // 更新栈顶指针为新节点
	st.size++        // 增加栈大小计数
}

// Pop pops an element from the stack.
// 弹出栈顶元素并返回
func (st *Stack[T]) Pop() (T, bool) {
	if st.size == 0 {
		var zero T
		return zero, false // 栈为空时返回零值和false
	}
	topNode := st.top          // 保存当前栈顶节点
	st.top = topNode.prev      // 将栈顶指针指向下一个节点
	topNode.prev = nil         // 断开原栈顶节点的prev指针
	st.size--                  // 减少栈大小计数
	return topNode.value, true // 返回栈顶元素的值和true
}

// Peek looks at the top element of the stack.
// 查看栈顶元素但不弹出
func (st *Stack[T]) Peek() (T, bool) {
	if st.size == 0 {
		var zero T
		return zero, false // 栈为空时返回零值和false
	}
	return st.top.value, true // 返回栈顶元素的值和true
}
