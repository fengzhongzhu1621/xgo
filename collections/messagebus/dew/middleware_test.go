package dew

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/fengzhongzhu1621/xgo/collections/messagebus/dew/commands/action"
	"github.com/fengzhongzhu1621/xgo/collections/messagebus/dew/commands/query"
	"github.com/fengzhongzhu1621/xgo/collections/messagebus/dew/handlers"
	"github.com/fengzhongzhu1621/xgo/collections/messagebus/dew/middlewares"
	"github.com/go-dew/dew"
)

func TestMiddleware(t *testing.T) {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	// 初始化消息总线 bus
	bus := initializeBus()

	fmt.Println("--- Authorization Example ---")

	// 模拟普通成员（Member）的操作
	if err := runMemberScenario(bus); err != nil {
		return fmt.Errorf("member scenario failed: %w", err)
	}

	// 模拟管理员（Admin）的操作
	if err := runAdminScenario(bus); err != nil {
		return fmt.Errorf("admin scenario failed: %w", err)
	}

	fmt.Println("\n--- Authorization Example finished ---")
	return nil
}

func initializeBus() dew.Bus {
	bus := dew.New()
	bus.Group(func(bus dew.Bus) {
		// 对所有 ACTION 类型的消息（即命令）应用 AdminOnly 中间件。这个中间件的作用是只允许管理员执行命令。
		// 查询无此验证中间件
		bus.Use(dew.ACTION, middlewares.AdminOnly)
		// 对所有类型的消息（包括 ACTION 和 QUERY）应用 LogCommand 中间件。这个中间件的作用是记录每条消息的执行日志。
		bus.Use(dew.ALL, middlewares.LogCommand)
		// 注册一个处理组织相关操作的消息处理器（Handler），比如查询组织详情或更新组织信息。
		bus.Register(handlers.NewOrgHandler())
	})
	return bus
}

// 模拟了一个普通成员的操作
func runMemberScenario(bus dew.Bus) error {
	// 创建一个带有成员身份的上下文 memberContext（通过 middlewares.AuthContext 设置当前用户为 MemberID）
	busContext := dew.NewContext(context.Background(), bus)
	memberContext := middlewares.AuthContext(
		busContext,
		&middlewares.CurrentUser{ID: middlewares.MemberID},
	)

	// 查询组织详情
	// 使用 dew.Query 发送一个查询请求 GetOrgDetailsQuery。
	// 预期这个查询应该成功，因为查询通常不受权限限制（除非特别指定）。
	fmt.Println("\n1. Execute a query to get the organization profile (should succeed for member).")
	orgProfile, err := dew.Query(memberContext, &query.GetOrgDetailsQuery{})
	if err != nil {
		return fmt.Errorf("unexpected error in GetOrgDetailsQuery: %w", err)
	}
	fmt.Printf("Organization Profile: %s\n", orgProfile.Result)

	// 尝试更新组织信息
	// 使用 dew.Dispatch 发送一个更新操作 UpdateOrgAction。
	fmt.Println(
		"\n2. Dispatch an action to update the organization profile (should fail for member).",
	)
	_, err = dew.Dispatch(memberContext, &action.UpdateOrgAction{Name: "Foo"})
	if err == nil {
		return fmt.Errorf("expected unauthorized error, got nil")
	}
	if err != middlewares.ErrUnauthorized {
		return fmt.Errorf("expected unauthorized error, got: %w", err)
	}
	fmt.Printf("Expected unauthorized error: %v\n", err)

	return nil
}

func runAdminScenario(bus dew.Bus) error {
	busContext := dew.NewContext(context.Background(), bus)
	adminContext := middlewares.AuthContext(
		busContext,
		&middlewares.CurrentUser{ID: middlewares.AdminID},
	)

	fmt.Println(
		"\n3. Dispatch an action to update the organization profile (should succeed for admin).",
	)
	err := dew.DispatchMulti(adminContext, dew.NewAction(&action.UpdateOrgAction{Name: "Foo"}))
	if err != nil {
		return fmt.Errorf("unexpected error in UpdateOrgAction: %w", err)
	}
	fmt.Println("\nOrganization profile updated successfully.")

	return nil
}
