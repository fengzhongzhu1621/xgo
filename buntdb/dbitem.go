package buntdb

import (
	"strconv"
	"time"

	"github.com/fengzhongzhu1621/xgo/datetime"
)

// dbItemOpts holds various meta information about an item.
type dbItemOpts struct {
	ex   bool      // does this item expire?
	exat time.Time // when does this item expire?
}
type dbItem struct {
	key, val string      // the binary key and value
	opts     *dbItemOpts // optional meta information
	keyless  bool        // keyless item for scanning 是否存在key，默认存在
}

// estAOFSetSize returns an estimated number of bytes that this item will use
// when stored in the aof file.
// 计算指令协议的字节数
func (dbi *dbItem) estAOFSetSize() int {
	var n int
	if dbi.opts != nil && dbi.opts.ex {
		n += estArraySize(5)
		n += estBulkStringSize("set")
		n += estBulkStringSize(dbi.key)
		n += estBulkStringSize(dbi.val)
		n += estBulkStringSize("ex")
		n += estBulkStringSize("99") // estimate two byte bulk string 延迟不能超过99秒
	} else {
		n += estArraySize(3)
		n += estBulkStringSize("set")
		n += estBulkStringSize(dbi.key)
		n += estBulkStringSize(dbi.val)
	}
	return n
}

// writeSetTo writes an item as a single SET record to the a bufio Writer.
func (dbi *dbItem) writeSetTo(buf []byte, now time.Time) []byte {
	if dbi.opts != nil && dbi.opts.ex {
		ex := dbi.opts.exat.Sub(now) / time.Second
		buf = appendArray(buf, 5)
		buf = appendBulkString(buf, "set")
		buf = appendBulkString(buf, dbi.key)
		buf = appendBulkString(buf, dbi.val)
		buf = appendBulkString(buf, "ex")
		buf = appendBulkString(buf, strconv.FormatUint(uint64(ex), 10))
	} else {
		buf = appendArray(buf, 3)
		buf = appendBulkString(buf, "set")
		buf = appendBulkString(buf, dbi.key)
		buf = appendBulkString(buf, dbi.val)
	}
	return buf
}

// writeSetTo writes an item as a single DEL record to the a bufio Writer.
func (dbi *dbItem) writeDeleteTo(buf []byte) []byte {
	buf = appendArray(buf, 2)
	buf = appendBulkString(buf, "del")
	buf = appendBulkString(buf, dbi.key)
	return buf
}

// expired evaluates id the item has expired. This will always return false when
// the item does not have `opts.ex` set to true.
func (dbi *dbItem) expired() bool {
	return dbi.opts != nil && dbi.opts.ex && time.Now().After(dbi.opts.exat)
}

// expiresAt will return the time when the item will expire. When an item does
// not expire `maxTime` is used.
func (dbi *dbItem) expiresAt() time.Time {
	if dbi.opts == nil || !dbi.opts.ex {
		return datetime.MaxTime
	}
	return dbi.opts.exat
}

// Less determines if a b-tree item is less than another. This is required
// for ordering, inserting, and deleting items from a b-tree. It's important
// to note that the ctx parameter is used to help with determine which
// formula to use on an item. Each b-tree should use a different ctx when
// sharing the same item.
// btree中元素的比较规则
func (dbi *dbItem) Less(dbi2 *dbItem, ctx interface{}) bool {
	switch ctx := ctx.(type) {
	case *exctx:
		// The expires b-tree formula
		if dbi2.expiresAt().After(dbi.expiresAt()) {
			return true
		}
		if dbi.expiresAt().After(dbi2.expiresAt()) {
			return false
		}
		// ttl相同，则继续比较key
	case *index:
		if ctx.less != nil {
			// Using an index
			if ctx.less(dbi.val, dbi2.val) {
				return true
			}
			if ctx.less(dbi2.val, dbi.val) {
				return false
			}
		}
	}
	// Always fall back to the key comparison. This creates absolute uniqueness.
	if dbi.keyless {
		return false
	} else if dbi2.keyless {
		return true
	}
	return dbi.key < dbi2.key
}

// Rect converts a string to a rectangle.
// An invalid rectangle will cause a panic.
// 获得索引的矩形表示（MBR）
func (dbi *dbItem) Rect(ctx interface{}) (min, max []float64) {
	switch ctx := ctx.(type) {
	case *index:
		return ctx.rect(dbi.val)
	}
	return nil, nil
}

func lessCtx(ctx interface{}) func(a, b interface{}) bool {
	return func(a, b interface{}) bool {
		return a.(*dbItem).Less(b.(*dbItem), ctx)
	}
}
