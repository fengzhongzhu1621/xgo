package buntdb

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fengzhongzhu1621/xgo"
	"github.com/tidwall/btree"
)

// exctx is a simple b-tree context for ordering by expiration.
type exctx struct {
	db *DB
}

// DB represents a collection of key-value pairs that persist on disk.
// Transactions are used for all forms of data access to the DB.
type DB struct {
	mu   sync.RWMutex // the gatekeeper for all fields
	file *os.File     // the underlying file，数据库文件
	buf  []byte       // a buffer to write to
	// 下面3个属性在执行flushdb时需要重置
	keys      *btree.BTree      // a tree of all item ordered by key 存放指令
	exps      *btree.BTree      // a tree of items ordered by expiration 存放过期的指令
	idxs      map[string]*index // the index trees.
	insIdxs   []*index          // a reuse buffer for gathering indexes
	flushes   int               // a count of the number of disk flushes
	closed    bool              // set when the database has been closed
	config    Config            // the database configuration
	persist   bool              // do we write to disk
	shrinking bool              // when an aof shrink is in-process.
	lastaofsz int               // the size of the last shrink aof size
}

// ReadConfig returns the database configuration.
func (db *DB) ReadConfig(config *Config) error {
	// 加读锁
	db.mu.RLock()
	defer db.mu.RUnlock()
	// 数据库已关闭不能修改配置
	if db.closed {
		return xgo.ErrDatabaseClosed
	}
	// 读取结果
	*config = db.config
	return nil
}

// SetConfig updates the database configuration.
func (db *DB) SetConfig(config Config) error {
	// 加写锁
	db.mu.Lock()
	defer db.mu.Unlock()

	// 数据库已关闭不能修改配置
	if db.closed {
		return xgo.ErrDatabaseClosed
	}
	// 判断同步策略，未知策略返回错误
	switch config.SyncPolicy {
	default:
		return xgo.ErrInvalidSyncPolicy
	case Never, EverySecond, Always:
	}

	db.config = config
	return nil
}

// Close releases all database resources.
// All transactions must be closed before closing the database.
func (db *DB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// 重复关闭报错
	if db.closed {
		return xgo.ErrDatabaseClosed
	}
	// 标记关闭
	db.closed = true

	// 持久化
	if db.persist {
		db.file.Sync() // do a sync but ignore the error
		if err := db.file.Close(); err != nil {
			return err
		}
	}
	// Let's release all references to nil. This will help both with debugging
	// late usage panics and it provides a hint to the garbage collector
	db.keys, db.exps, db.idxs, db.file = nil, nil, nil, nil
	return nil
}

// Load loads commands from reader. This operation blocks all reads and writes.
// Note that this can only work for fully in-memory databases opened with
// Open(":memory:").
func (db *DB) Load(rd io.Reader) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.persist {
		// cannot load into databases that persist to disk
		return xgo.ErrPersistenceActive
	}
	_, err := db.readLoad(rd, time.Now())
	return err
}

// load reads entries from the append only database file and fills the database.
// The file format uses the Redis append only file format, which is and a series
// of RESP commands. For more information on RESP please read
// http://redis.io/topics/protocol. The only supported RESP commands are DEL and
// SET.
func (db *DB) load() error {
	// 获得文件的修改时间
	fi, err := db.file.Stat()
	if err != nil {
		return err
	}
	// 读取文件，返回读取的字节数
	n, err := db.readLoad(db.file, fi.ModTime())
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			// The db file has ended mid-command, which is allowed but the
			// data file should be truncated to the end of the last valid
			// command
			//
			// Truncate改变文件的大小，对文件进行截断，它不会改变I/O的当前位置。
			// 如果截断文件，多出的部分就会被丢弃。如果出错，错误底层类型是*PathError。
			//
			// 如果数据文件中间遇到结束符，后续的无效内容需要删除
			if err := db.file.Truncate(n); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	// 指针指向文件结尾
	if _, err := db.file.Seek(n, 0); err != nil {
		return err
	}

	// 计算指令转为为协议表示需要的字节数
	var estaofsz int
	db.keys.Walk(func(items []interface{}) {
		for _, v := range items {
			estaofsz += v.(*dbItem).estAOFSetSize()
		}
	})
	db.lastaofsz += estaofsz
	return nil
}

// readLoad reads from the reader and loads commands into the database.
// modTime is the modified time of the reader, should be no greater than
// the current time.Now().
// Returns the number of bytes of the last command read and the error if any.
func (db *DB) readLoad(rd io.Reader, modTime time.Time) (n int64, err error) {
	defer func() {
		// 修改错误
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	// 指令 set a b ex 10的协议示例如下
	// *5\r\n
	// $3\r\nset\r\n
	// $1\r\na\r\n
	// $1\r\nb\r\n
	// $2\r\nex\r\n
	// $2\r\n10\r\n

	totalSize := int64(0)
	data := make([]byte, 4096)
	parts := make([]string, 0, 8)

	// 读文件
	r := bufio.NewReader(rd)
	for {
		// peek at the first byte. If it's a 'nul' control character then
		// ignore it and move to the next byte.
		// 读取并返回一个字节，如果没有字节可读，则返回错误信息
		c, err := r.ReadByte()
		if err != nil {
			// 返回读取的字节数
			if err == io.EOF {
				err = nil
			}
			return totalSize, err
		}
		// 忽略控制字符
		if c == 0 {
			// ignore nul control characters
			n += 1
			continue
		}
		// 对读取字节的操作进行一个撤销操作
		// 取消已读取的最后一个字节（即把字节重新放回读取缓冲区的前部）。只有最近一次读取的单个字节才能取消读取。
		// UnreadByte方法的意思是：将上一次ReadByte的字节还原，使得再次调用ReadByte返回的结果和上一次调用相同，
		// 也就是说，UnreadByte是重置上一次的ReadByte。
		// 注意：UnreadByte调用之前必须调用了ReadByte,且不能连续调用UnreadByte
		if err := r.UnreadByte(); err != nil {
			return totalSize, err
		}

		// read a single command. 获得parts的数量
		// first we should read the number of parts that the of the command
		cmdByteSize := int64(0)
		// 读的时候从缓冲器里找第一个出现的分隔符（delim），找到后，把从缓冲器头部开始到分隔符之间的所有byte进行返回
		// 它和 ReadSlice 具有相同的签名，但是 ReadSlice 是一个低级别的函数，ReadBytes 的实现使用了 ReadSlice。 那么两者之间有什么不同呢?
		// 在分隔符找不到的情况下，ReadBytes 可以多次调用 ReadSlice，而且可以累积返回的数据。 这意味着 ReadBytes 将不再受到 缓存大小的限制
		// ReadBytes 会将分隔符一起返回，所以需要额外的一些工作来重新处理数据(除非返回分隔符是有用的)。
		line, err := r.ReadBytes('\n')
		if err != nil {
			return totalSize, err
		}
		// 命令必须以*开头
		if line[0] != '*' {
			return totalSize, xgo.ErrInvalid
		}
		cmdByteSize += int64(len(line))

		// convert the string number to and int
		var n int
		if len(line) == 4 && line[len(line)-2] == '\r' {
			// eg: *1\r\n => 1
			if line[1] < '0' || line[1] > '9' {
				return totalSize, xgo.ErrInvalid
			}
			n = int(line[1] - '0')
		} else {
			if len(line) < 5 || line[len(line)-2] != '\r' {
				return totalSize, xgo.ErrInvalid
			}
			// eg: *123\r\n => 123
			for i := 1; i < len(line)-2; i++ {
				if line[i] < '0' || line[i] > '9' {
					return totalSize, xgo.ErrInvalid
				}
				n = n*10 + int(line[i]-'0')
			}
		}

		// 清空切片parts
		parts = parts[:0]

		// read each part of the command.
		// 读取n个parts
		for i := 0; i < n; i++ {
			// read the number of bytes of the part.
			// 必须是$开头
			line, err := r.ReadBytes('\n')
			if err != nil {
				return totalSize, err
			}
			if line[0] != '$' {
				return totalSize, xgo.ErrInvalid
			}
			cmdByteSize += int64(len(line))

			// convert the string number to and int
			var n int
			if len(line) == 4 && line[len(line)-2] == '\r' {
				// eg: $1\r\n => 1
				if line[1] < '0' || line[1] > '9' {
					return totalSize, xgo.ErrInvalid
				}
				n = int(line[1] - '0')
			} else {
				// eg: $123\r\n => 123
				if len(line) < 5 || line[len(line)-2] != '\r' {
					return totalSize, xgo.ErrInvalid
				}
				for i := 1; i < len(line)-2; i++ {
					if line[i] < '0' || line[i] > '9' {
						return totalSize, xgo.ErrInvalid
					}
					n = n*10 + int(line[i]-'0')
				}
			}
			// resize the read buffer
			// 扩容缓冲区
			if len(data) < n+2 {
				dataln := len(data)
				for dataln < n+2 {
					dataln *= 2
				}
				data = make([]byte, dataln)
			}

			// 读取指定数量的字节到数组中
			if _, err = io.ReadFull(r, data[:n+2]); err != nil {
				return totalSize, err
			}
			if data[n] != '\r' || data[n+1] != '\n' {
				return totalSize, xgo.ErrInvalid
			}

			// copy string
			parts = append(parts, string(data[:n]))
			cmdByteSize += int64(n + 2)
		}

		// finished reading the command

		if len(parts) == 0 {
			continue
		}
		if (parts[0][0] == 's' || parts[0][0] == 'S') &&
			(parts[0][1] == 'e' || parts[0][1] == 'E') &&
			(parts[0][2] == 't' || parts[0][2] == 'T') {
			// SET
			if len(parts) < 3 || len(parts) == 4 || len(parts) > 5 {
				return totalSize, xgo.ErrInvalid
			}
			if len(parts) == 5 {
				// set key value ex timeout
				if strings.ToLower(parts[3]) != "ex" {
					return totalSize, xgo.ErrInvalid
				}
				ex, err := strconv.ParseUint(parts[4], 10, 64)
				if err != nil {
					return totalSize, err
				}
				now := time.Now()
				dur := (time.Duration(ex) * time.Second) - now.Sub(modTime)
				if dur > 0 {
					// 未过期
					db.insertIntoDatabase(&dbItem{
						key: parts[1],
						val: parts[2],
						opts: &dbItemOpts{
							ex:   true,
							exat: now.Add(dur), // 过期时间
						},
					})
				}
			} else {
				// set key value
				db.insertIntoDatabase(&dbItem{key: parts[1], val: parts[2]})
			}
		} else if (parts[0][0] == 'd' || parts[0][0] == 'D') &&
			(parts[0][1] == 'e' || parts[0][1] == 'E') &&
			(parts[0][2] == 'l' || parts[0][2] == 'L') {
			// DEL key
			if len(parts) != 2 {
				return totalSize, xgo.ErrInvalid
			}
			db.deleteFromDatabase(&dbItem{key: parts[1]})
		} else if (parts[0][0] == 'f' || parts[0][0] == 'F') &&
			strings.ToLower(parts[0]) == "flushdb" {
			// flushdb
			db.keys = btreeNew(lessCtx(nil))
			db.exps = btreeNew(lessCtx(&exctx{db}))
			db.idxs = make(map[string]*index)
		} else {
			// 无效指令
			return totalSize, xgo.ErrInvalid
		}
		totalSize += cmdByteSize
	}
}

// insertIntoDatabase performs inserts an item in to the database and updates
// all indexes. If a previous item with the same key already exists, that item
// will be replaced with the new one, and return the previous item.
func (db *DB) insertIntoDatabase(item *dbItem) *dbItem {
	var pdbi *dbItem
	// Generate a list of indexes that this item will be inserted in to.
	// 获得key命中的索引集
	idxs := db.insIdxs
	for _, idx := range db.idxs {
		if idx.match(item.key) {
			idxs = append(idxs, idx)
		}
	}

	// 将指令对象放到btree中，返回旧的值
	prev := db.keys.Set(item)
	if prev != nil {
		// A previous item was removed from the keys tree. Let's
		// fully delete this item from all indexes.
		// 从过期btree中删除旧的指令
		pdbi = prev.(*dbItem)
		if pdbi.opts != nil && pdbi.opts.ex {
			// Remove it from the expires tree.
			db.exps.Delete(pdbi)
		}

		// 从索引中删除旧指令
		for _, idx := range idxs {
			if idx.btr != nil {
				// Remove it from the btree index.
				idx.btr.Delete(pdbi)
			}
			if idx.rtr != nil {
				// Remove it from the rtree index.
				idx.rtr.Remove(pdbi)
			}
		}
	}

	// 添加新的指令到带有过期时间的btree中
	if item.opts != nil && item.opts.ex {
		// The new item has eviction options. Add it to the
		// expires tree
		db.exps.Set(item)
	}

	// 加入索引
	for i, idx := range idxs {
		if idx.btr != nil {
			// Add new item to btree index.
			idx.btr.Set(item)
		}
		if idx.rtr != nil {
			// Add new item to rtree index.
			idx.rtr.Insert(item)
		}
		// clear the index
		idxs[i] = nil
	}

	// reuse the index list slice
	db.insIdxs = idxs[:0]

	// we must return the previous item to the caller.
	return pdbi
}

// deleteFromDatabase removes and item from the database and indexes. The input
// item must only have the key field specified thus "&dbItem{key: key}" is all
// that is needed to fully remove the item with the matching key. If an item
// with the matching key was found in the database, it will be removed and
// returned to the caller. A nil return value means that the item was not
// found in the database
func (db *DB) deleteFromDatabase(item *dbItem) *dbItem {
	var pdbi *dbItem
	// 删除指令
	prev := db.keys.Delete(item)
	if prev != nil {
		pdbi = prev.(*dbItem)
		if pdbi.opts != nil && pdbi.opts.ex {
			// Remove it from the exipres tree.
			db.exps.Delete(pdbi)
		}
		for _, idx := range db.idxs {
			if !idx.match(pdbi.key) {
				continue
			}
			if idx.btr != nil {
				// Remove it from the btree index.
				idx.btr.Delete(pdbi)
			}
			if idx.rtr != nil {
				// Remove it from the rtree index.
				idx.rtr.Remove(pdbi)
			}
		}
	}
	return pdbi
}

// Save writes a snapshot of the database to a writer. This operation blocks all
// writes, but not reads. This can be used for snapshots and backups for pure
// in-memory databases using the ":memory:". Database that persist to disk
// can be snapshotted by simply copying the database file.
func (db *DB) Save(wr io.Writer) error {
	var err error
	db.mu.RLock()
	defer db.mu.RUnlock()

	// use a buffered writer and flush every 4MB
	var buf []byte
	now := time.Now()

	// 遍历所有指令
	// iterated through every item in the database and write to the buffer
	btreeAscend(db.keys, func(item interface{}) bool {
		// 将指令转换为协议放到buf中
		dbi := item.(*dbItem)
		buf = dbi.writeSetTo(buf, now)
		if len(buf) > 1024*1024*4 {
			// flush when buffer is over 4MB
			_, err = wr.Write(buf)
			if err != nil {
				// 写失败停止遍历
				return false
			}
			buf = buf[:0]
		}
		return true
	})
	if err != nil {
		return err
	}

	// one final flush
	// 保存不足4M的数据
	if len(buf) > 0 {
		_, err = wr.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

// get return an item or nil if not found.
func (db *DB) get(key string) *dbItem {
	item := db.keys.Get(&dbItem{key: key})
	if item != nil {
		return item.(*dbItem)
	}
	return nil
}

// managed calls a block of code that is fully contained in a transaction.
// This method is intended to be wrapped by Update and View

// View executes a function within a managed read-only transaction.
// When a non-nil error is returned from the function that error will be return
// to the caller of View().
//
// Executing a manual commit or rollback from inside the function will result
// in a panic.
func (db *DB) View(fn func(tx *Tx) error) error {
	return db.managed(false, fn)
}

func (db *DB) managed(writable bool, fn func(tx *Tx) error) (err error) {
	var tx *Tx
	// 创建一个事务对象
	tx, err = db.Begin(writable)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			// The caller returned an error. We must rollback.
			// 指令写文件失败，则对数据进行回滚
			// 如果回滚失败则可能会导致数据不一致，所以 rollbackInner 操作尽可能不出错，这里忽略硬件故障和第三方包出错的场景
			_ = tx.Rollback()
			return
		}
		if writable {
			// Everything went well. Lets Commit()
			// 将指令写入到文件中
			err = tx.Commit()
		} else {
			// read-only transaction can only roll back.
			// 读操作的回滚仅仅是做了解锁操作
			err = tx.Rollback()
		}
	}()

	tx.funcd = true
	defer func() {
		// 事务处理函数执行完毕后设置为 false
		// 用于保证在执行的过程中不能调用 commit / rollback
		tx.funcd = false
	}()

	// 执行事务处理函数
	err = fn(tx)
	return
}

// Begin opens a new transaction.
// Multiple read-only transactions can be opened at the same time but there can
// only be one read/write transaction at a time. Attempting to open a read/write
// transactions while another one is in progress will result in blocking until
// the current read/write transaction is completed.
//
// All transactions must be closed by calling Commit() or Rollback() when done.
// 创建一个事务对象
func (db *DB) Begin(writable bool) (*Tx, error) {
	tx := &Tx{
		db:       db,
		writable: writable,
	}
	tx.lock()
	if db.closed {
		tx.unlock()
		return nil, xgo.ErrDatabaseClosed
	}
	if writable {
		// writable transactions have a writeContext object that
		// contains information about changes to the database.
		tx.wc = &txWriteContext{}
		tx.wc.rollbackItems = make(map[string]*dbItem)
		tx.wc.rollbackIndexes = make(map[string]*index)
		if db.persist {
			// 需要持久化的指令集
			tx.wc.commitItems = make(map[string]*dbItem)
		}
	}
	return tx, nil
}

// Update executes a function within a managed read/write transaction.
// The transaction has been committed when no error is returned.
// In the event that an error is returned, the transaction will be rolled back.
// When a non-nil error is returned from the function, the transaction will be
// rolled back and the that error will be return to the caller of Update().
//
// Executing a manual commit or rollback from inside the function will result
// in a panic.
func (db *DB) Update(fn func(tx *Tx) error) error {
	return db.managed(true, fn)
}

// DropIndex removes an index.
func (db *DB) DropIndex(name string) error {
	return db.Update(func(tx *Tx) error {
		return tx.DropIndex(name)
	})
}

// Indexes returns a list of index names.
func (db *DB) Indexes() ([]string, error) {
	var names []string
	var err = db.View(func(tx *Tx) error {
		var err error
		names, err = tx.Indexes()
		return err
	})
	return names, err
}

// CreateIndex builds a new index and populates it with items.
// The items are ordered in an b-tree and can be retrieved using the
// Ascend* and Descend* methods.
// An error will occur if an index with the same name already exists.
//
// When a pattern is provided, the index will be populated with
// keys that match the specified pattern. This is a very simple pattern
// match where '*' matches on any number characters and '?' matches on
// any one character.
// The less function compares if string 'a' is less than string 'b'.
// It allows for indexes to create custom ordering. It's possible
// that the strings may be textual or binary. It's up to the provided
// less function to handle the content format and comparison.
// There are some default less function that can be used such as
// IndexString, IndexBinary, etc.
func (db *DB) CreateIndex(name, pattern string,
	less ...func(a, b string) bool) error {
	return db.Update(func(tx *Tx) error {
		return tx.CreateIndex(name, pattern, less...)
	})
}

// ReplaceIndex builds a new index and populates it with items.
// The items are ordered in an b-tree and can be retrieved using the
// Ascend* and Descend* methods.
// If a previous index with the same name exists, that index will be deleted.
func (db *DB) ReplaceIndex(name, pattern string,
	less ...func(a, b string) bool) error {
	return db.Update(func(tx *Tx) error {
		err := tx.CreateIndex(name, pattern, less...)
		if err != nil {
			if err == xgo.ErrIndexExists {
				err := tx.DropIndex(name)
				if err != nil {
					return err
				}
				return tx.CreateIndex(name, pattern, less...)
			}
			return err
		}
		return nil
	})
}

// CreateSpatialIndex builds a new index and populates it with items.
// The items are organized in an r-tree and can be retrieved using the
// Intersects method.
// An error will occur if an index with the same name already exists.
//
// The rect function converts a string to a rectangle. The rectangle is
// represented by two arrays, min and max. Both arrays may have a length
// between 1 and 20, and both arrays must match in length. A length of 1 is a
// one dimensional rectangle, and a length of 4 is a four dimension rectangle.
// There is support for up to 20 dimensions.
// The values of min must be less than the values of max at the same dimension.
// Thus min[0] must be less-than-or-equal-to max[0].
// The IndexRect is a default function that can be used for the rect
// parameter.
func (db *DB) CreateSpatialIndex(name, pattern string,
	rect func(item string) (min, max []float64)) error {
	return db.Update(func(tx *Tx) error {
		return tx.CreateSpatialIndex(name, pattern, rect)
	})
}

// ReplaceSpatialIndex builds a new index and populates it with items.
// The items are organized in an r-tree and can be retrieved using the
// Intersects method.
// If a previous index with the same name exists, that index will be deleted.
func (db *DB) ReplaceSpatialIndex(name, pattern string,
	rect func(item string) (min, max []float64)) error {
	return db.Update(func(tx *Tx) error {
		err := tx.CreateSpatialIndex(name, pattern, rect)
		if err != nil {
			if err == xgo.ErrIndexExists {
				err := tx.DropIndex(name)
				if err != nil {
					return err
				}
				return tx.CreateSpatialIndex(name, pattern, rect)
			}
			return err
		}
		return nil
	})
}
