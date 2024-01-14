package buntdb

import (
	"strconv"

	"github.com/tidwall/btree"
	"github.com/tidwall/match"
)

// estIntSize returns the string representions size.
// Has the same result as len(strconv.Itoa(x)).
// 计算数字字符串的长度
func estIntSize(x int) int {
	n := 1
	if x < 0 {
		n++
		x *= -1
	}
	for x >= 10 {
		n++
		x /= 10
	}
	return n
}

// estArraySize 计算指令 *5\r\n 的长度
func estArraySize(count int) int {
	return 1 + estIntSize(count) + 2
}

// estBulkStringSize 计算指令$3\r\nset\r\n的长度
func estBulkStringSize(s string) int {
	return 1 + estIntSize(len(s)) + 2 + len(s) + 2
}

// appendArray 构造一个指令的组成部分的数量协议
func appendArray(buf []byte, count int) []byte {
	buf = append(buf, '*')
	buf = strconv.AppendInt(buf, int64(count), 10)
	buf = append(buf, '\r', '\n')
	return buf
}

// appendBulkString 构造指令每个部分的协议
func appendBulkString(buf []byte, s string) []byte {
	buf = append(buf, '$')
	buf = strconv.AppendInt(buf, int64(len(s)), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, s...)
	buf = append(buf, '\r', '\n')
	return buf
}

//// Wrappers around btree Ascend/Descend

func bLT(tr *btree.BTree, a, b interface{}) bool { return tr.Less(a, b) }
func bGT(tr *btree.BTree, a, b interface{}) bool { return tr.Less(b, a) }

// func bLTE(tr *btree.BTree, a, b interface{}) bool { return !tr.Less(b, a) }
// func bGTE(tr *btree.BTree, a, b interface{}) bool { return !tr.Less(a, b) }

// Ascend

// btreeAscend 升序遍历 [8, 8]
func btreeAscend(tr *btree.BTree, iter func(item interface{}) bool) {
	tr.Ascend(nil, iter)
}

// btreeAscendLessThan 生序遍历，直到有个元素大于pivot [8, pivot)
func btreeAscendLessThan(tr *btree.BTree, pivot interface{},
	iter func(item interface{}) bool,
) {
	tr.Ascend(nil, func(item interface{}) bool {
		return bLT(tr, item, pivot) && iter(item)
	})
}

// btreeAscendGreaterOrEqual 从pivot开始生序遍历 [pivot, 8)
func btreeAscendGreaterOrEqual(tr *btree.BTree, pivot interface{},
	iter func(item interface{}) bool,
) {
	tr.Ascend(pivot, iter)
}

// btreeAscendRange 升序遍历一个范围 [greaterOrEqual, lessThan)
func btreeAscendRange(tr *btree.BTree, greaterOrEqual, lessThan interface{},
	iter func(item interface{}) bool,
) {
	tr.Ascend(greaterOrEqual, func(item interface{}) bool {
		return bLT(tr, item, lessThan) && iter(item)
	})
}

// Descend

// btreeDescend 逆序遍历 [8, 8]
func btreeDescend(tr *btree.BTree, iter func(item interface{}) bool) {
	tr.Descend(nil, iter)
}

// btreeDescendGreaterThan [pivot, 8]
func btreeDescendGreaterThan(tr *btree.BTree, pivot interface{},
	iter func(item interface{}) bool,
) {
	tr.Descend(nil, func(item interface{}) bool {
		return bGT(tr, item, pivot) && iter(item)
	})
}

// btreeDescendRange [greaterThan, lessOrEqual]
func btreeDescendRange(tr *btree.BTree, lessOrEqual, greaterThan interface{},
	iter func(item interface{}) bool,
) {
	tr.Descend(lessOrEqual, func(item interface{}) bool {
		return bGT(tr, item, greaterThan) && iter(item)
	})
}

// btreeDescendLessOrEqual [8, pivot]
func btreeDescendLessOrEqual(tr *btree.BTree, pivot interface{},
	iter func(item interface{}) bool,
) {
	tr.Descend(pivot, iter)
}

// 创建一个btree（不支持并发写）
func btreeNew(less func(a, b interface{}) bool) *btree.BTree {
	// Using NewNonConcurrent because we're managing our own locks.
	return btree.NewNonConcurrent(less)
}

// Match returns true if the specified key matches the pattern. This is a very
// simple pattern matcher where '*' matches on any number characters and '?'
// matches on any one character.
func Match(key, pattern string) bool {
	return match.Match(key, pattern)
}
