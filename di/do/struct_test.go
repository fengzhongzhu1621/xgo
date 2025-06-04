package do

import (
	"fmt"
	"testing"

	"github.com/samber/do"
)

// ////////////////////////////////////////////////////////////////
type Engine struct {
	Started bool // 表示引擎是否已启动
}

// 构造函数：创建 Engine 实例
func NewEngine(i *do.Injector) (*Engine, error) {
	return &Engine{Started: false}, nil
}

// [Optional] Implements do.Healthcheckable.
func (c *Engine) HealthCheck() error {
	return fmt.Errorf("engine broken")
}

// [Optional] Implements do.Shutdownable.
func (c *Engine) Shutdown() error {
	println("engine stopped")
	return nil
}

// ////////////////////////////////////////////////////////////////
type Car struct {
	Engine *Engine // Car 依赖一个 Engine 实例
}

// 构造函数：创建 Car 实例时自动注入 Engine
func NewCar(i *do.Injector) (*Car, error) {
	// 从容器中获取 Engine 实例（依赖注入的核心操作）
	// MustInvoke 会阻塞直到 Engine 可用，如果失败则 panic（生产环境建议用 Invoke + 错误处理）。
	engine := do.MustInvoke[*Engine](i)
	return &Car{Engine: engine}, nil
}

func (c *Car) Start() {
	c.Engine.Started = true   // 启动引擎
	println("car is running") // 打印运行状态
}

// [Optional] Implements do.Healthcheckable.
func (c *Car) HealthCheck() error {
	return fmt.Errorf("car broken")
}

// [Optional] Implements do.Shutdownable.
func (c *Car) Shutdown() error {
	println("car stopped")
	return nil
}

// ////////////////////////////////////////////////////////////////
func TestStruct(t *testing.T) {
	// 初始化一个空的 DI 容器，创建依赖注入容器
	injector := do.New()

	// 注册服务
	do.Provide(injector, NewEngine) // 注册 Engine 服务（懒加载单例）
	do.Provide(injector, NewCar)    // 注册 Car 服务，依赖 Engine

	// 获取 Car 实例（自动创建 Car 及其依赖 Engine）
	car := do.MustInvoke[*Car](injector)

	// 使用服务
	car.Start() // 输出: "Car is running"

	fmt.Println(do.HealthCheck[*Car](injector)) // 输出: "car broken"

	// shutdown all services in reverse order
	// 关闭容器（触发 Shutdown 钩子，如果实现）
	// injector.ShutdownOnSIGTERM()    // will block until receiving sigterm signal
	injector.Shutdown()
	// 输出:
	// car stopped
	// engine stopped
}
