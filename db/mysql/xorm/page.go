package xorm

import (
	"github.com/fengzhongzhu1621/xgo/db/mysql"
)

const (
	DefaultPage     = 1    // 当前页数
	DefaultPageSize = 20   // 每页多少条数据
	MaxPageSize     = 1000 // 每次最多取多少条
)

// BaseModel 分页model
type Page struct {
	// 当前页码， 1-based
	Page int `xorm:"-"`
	// 每页显示的记录数
	PageSize int `xorm:"-"`
}

// parsePageAndPageSize 从查询条件中获取当前页码和每页的数量
func (model *Page) ParsePageAndPageSize(params mysql.CommonQueryConditionMap) {
	var (
		tested bool
		p1, p2 int
		p3, p4 int64
	)

	// 从查询条件中获取当前页码
	page, ok := params["page"]
	if ok {
		if p1, tested = page.(int); !tested {
			if p3, tested = page.(int64); !tested {
				model.Page = DefaultPage
			} else {
				model.Page = int(p3)
			}
		} else {
			model.Page = p1
		}

	}

	// 从查询条件中获取每页显示的记录数
	pageSize, ok := params["pagesize"]
	if ok {
		if p2, tested = pageSize.(int); !tested {
			if p4, tested = pageSize.(int64); !tested {
				model.PageSize = DefaultPageSize
			} else {
				model.PageSize = int(p4)
			}
		} else {
			model.PageSize = int(p2)
		}
	}
	if model.Page <= 0 {
		model.Page = DefaultPage
	}
	if model.PageSize <= 0 {
		model.PageSize = DefaultPageSize
	}
}

func (model *Page) PageLimitOffset() int {
	return (model.Page - 1) * model.PageSize
}
