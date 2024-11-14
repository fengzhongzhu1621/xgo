package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

const (
	minPageSize = 1
	maxPageSize = 50
)

// GetPage 从请求参数获得页码
func GetPage(c *gin.Context) int {
	page := cast.ToInt(c.Query("page"))

	return max(1, page)
}

// GetLimit 从请求参数获得每页的大小
func GetPageSize(c *gin.Context) int {
	pageSize := cast.ToInt(c.Query("limit"))
	pageSize = min(maxPageSize, pageSize)
	pageSize = max(minPageSize, pageSize)

	return pageSize
}

// GetOffset 获得数据的偏移量
func GetOffset(c *gin.Context) int {
	page, limit := GetPage(c), GetPageSize(c)

	return (page - 1) * limit
}
