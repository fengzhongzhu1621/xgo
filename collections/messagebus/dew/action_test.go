package dew

import (
	"context"
	"fmt"
	"testing"

	"github.com/fengzhongzhu1621/xgo/collections/messagebus/dew/commands/action"
	"github.com/fengzhongzhu1621/xgo/collections/messagebus/dew/handlers"
	"github.com/go-dew/dew"
)

func TestRegister(t *testing.T) {
	// 初始化命令总线。
	bus := dew.New()

	// 注册 HelloHandler。
	bus.Register(&handlers.HelloHandler{})

	// 创建带有总线的上下文。
	ctx := dew.NewContext(context.Background(), bus)

	// 从命令行参数获取名称或使用默认值。
	name := "Dew"

	// 创建并调度 HelloAction。
	if _, err := dew.Dispatch(ctx, &action.HelloAction{Name: name}); err != nil {
		fmt.Println("failed to dispatch HelloAction: %w", err)
	}
}
