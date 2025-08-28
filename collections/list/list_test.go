package list

import (
	"container/list"
	"fmt"
	"testing"
)

func TestList(t *testing.T) {
	// 创建一个新的双向链表
	l := list.New()

	// 在链表尾部插入元素
	l.PushBack("A")
	l.PushBack("B")
	l.PushBack("C")

	// 在链表头部插入元素
	l.PushFront("D")

	// 遍历链表并打印元素
	// D
	// A
	// B
	// C
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

	// 删除链表中的元素
	l.Remove(l.Front())

	// 再次遍历链表并打印元素
	// 	A
	// B
	// C
	fmt.Println("After removing the first element:")
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}
