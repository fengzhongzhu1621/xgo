package xorm

import (
	"bytes"

	"xorm.io/xorm"
)

// ColumnsSearch ... 用于多个字段的模糊查询匹配
func ColumnsSearch(session *xorm.Session, keyword string, args ...string) *xorm.Session {
	if keyword == "" {
		return session
	}

	var (
		buffer bytes.Buffer
		keys   []interface{}
	)

	like_keyword := "%" + keyword + "%"

	for i, column := range args {
		if i != 0 {
			buffer.WriteString(" OR ")
		}
		buffer.WriteString(column + " LIKE ?")
		keys = append(keys, like_keyword)
	}

	return session.And(buffer.String(), keys...)
}

// ColumnsEnumsIntSearch  用于多个INT字段的精确匹配
func ColumnsEnumsIntSearch(session *xorm.Session, array []int, column string) *xorm.Session {
	var (
		buffer bytes.Buffer
		keys   []interface{}
	)

	for i, v := range array {
		if i != 0 {
			buffer.WriteString(" OR ")
		}
		buffer.WriteString(column + " = ? ")
		keys = append(keys, v)
	}

	return session.And(buffer.String(), keys...)
}
