package paginate

import (
	"math"
)

type Paging struct {
	Page      int // 当前页码
	PageSize  int // 每页条数
	Total     int // 总条数
	PageCount int // 总页数
	Offset    int // 偏移量
}

// CreatePaging 创建分页对象
func CreatePaging(Page, pageSize, total int) *Paging {
	// 默认 1-based
	if Page < 1 {
		Page = 1
	}
	// 每页最小 10 条记录
	if pageSize < 1 {
		pageSize = 10
	}

	// 计算总页数
	pageCount := math.Ceil(float64(total) / float64(pageSize))
	// 计算偏移量
	offset := pageSize * (Page - 1)

	paging := new(Paging)
	paging.Page = Page
	paging.PageSize = pageSize
	paging.Total = total
	paging.PageCount = int(pageCount)
	paging.Offset = offset

	return paging
}
