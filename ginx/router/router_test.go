package router

import (
	"testing"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	t.Parallel()

	cfg := config.GetGlobalConfig()

	router := gin.Default()
	rootGroup := router.Group("/")
	RegisterRouter(cfg, router)
	RegisterRouterGroup(cfg, rootGroup)

	assert.NotNil(t, router)
	assert.NotNil(t, rootGroup)
}
