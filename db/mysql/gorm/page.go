package gorm_2

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

type Pagination struct {
	page     int
	pageSize int
}

func (p *Pagination) GetPage() int     { return p.page }
func (p *Pagination) GetPageSize() int { return p.pageSize }
func (p *Pagination) ModifyStatement(stm *gorm.Statement) {
	// 修改语句以添加分页
	db := stm.DB
	db.Limit(p.pageSize).Offset((p.page - 1) * p.pageSize)
}

func (p *Pagination) Build(_ clause.Builder) {
	// Build方法留空，因为分页不需要额外的SQL子句
}
