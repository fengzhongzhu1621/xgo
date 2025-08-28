package handlers

import (
	"context"
	"fmt"

	"github.com/fengzhongzhu1621/xgo/collections/messagebus/dew/commands/action"
	"github.com/fengzhongzhu1621/xgo/collections/messagebus/dew/commands/query"
)

// HelloHandler 处理 HelloAction。
type HelloHandler struct{}

// HandleHello 是 HelloAction 的处理函数。
func (h *HelloHandler) HandleHelloAction(ctx context.Context, cmd *action.HelloAction) error {
	fmt.Printf("Hello, %s!\n", cmd.Name)
	return nil
}

func (h *HelloHandler) HandleHelloQuery(ctx context.Context, cmd *query.HelloQuery) error {
	cmd.Result = fmt.Sprintf("Hello, %s!", cmd.Name)
	return nil
}
