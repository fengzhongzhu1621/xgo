package pool

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Student struct {
	Name   string
	Age    int32
	Remark [1024]byte
}

var (
	buf, _      = json.Marshal(Student{Name: "bob", Age: 25})
	studentPool = sync.Pool{
		New: func() interface{} {
			return new(Student)
		},
	}
)

type Person struct {
	Name string
	Age  int
}

type personPool struct {
	pool sync.Pool
}

// Get Get 之后需要进行断言和一些初始化操作，Get 之后的数据状态是不确定的。
func (pp *personPool) Get(name string, age int) (p *Person, err error) {
	// 从池中获取一个对象
	p, ok := pp.pool.Get().(*Person)
	if !ok {
		return nil, err
	}

	// 初始化对象的值
	p.Name = name
	p.Age = age

	return p, nil
}

// Put 在 Put 前要对临时对象做一些清理工作，以免影响下一次复用。
func (pp *personPool) Put(p *Person) {
	// 归还前需要清理状态
	p.Name = ""
	p.Age = 0

	// 归还到缓存池
	pp.pool.Put(p)
}

var PersonPool = &personPool{
	pool: sync.Pool{
		New: func() interface{} {
			return new(Person)
		},
	},
}

func TestPersonPool(t *testing.T) {
	// 从池中获取一个对象
	p1, err := PersonPool.Get("mark", 18)
	if err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, "mark", p1.Name)

	// 使用完毕后放回池中
	PersonPool.Put(p1)
}

func BenchmarkUnmarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		stu := &Student{}
		_ = json.Unmarshal(buf, stu)
	}
}

func BenchmarkUnmarshalWithPool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		stu := studentPool.Get().(*Student)
		_ = json.Unmarshal(buf, stu)
		studentPool.Put(stu)
	}
}
