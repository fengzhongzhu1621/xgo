package do

import (
	"fmt"
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
)

// ////////////////////////////////////////////////////////////////
// 用于注入接口
type IEngineService interface{}

type engineServiceImplem struct{}

// [Optional] Implements do.Healthcheckable.
func (c *engineServiceImplem) HealthCheck() error {
	return fmt.Errorf("engine broken")
}

func NewEngineService(i *do.Injector) (IEngineService, error) {
	return &engineServiceImplem{}, nil
}

// ////////////////////////////////////////////////////////////////
// 用于注入结构体
type CarService struct {
	Engine IEngineService
}

func (c *CarService) Start() {
	println("car starting")
}

// [Optional] Implements do.Shutdownable.
func (c *CarService) Shutdown() error {
	println("car stopped")
	return nil
}

func NewCarService(i *do.Injector) (*CarService, error) {
	engine := do.MustInvoke[IEngineService](i)
	car := CarService{Engine: engine}
	return &car, nil
}

// 注入接口示例
func TestInterface(t *testing.T) {
	injector := do.New()

	// provides CarService
	do.Provide(injector, NewCarService)
	// provides EngineService
	do.Provide(injector, NewEngineService)

	car := do.MustInvoke[*CarService](injector)
	car.Start()
	// prints "car starting"

	car2 := do.MustInvoke[*CarService](injector)
	car2.Start()
	// prints "car starting"

	// 测试多次从缓存中获取的对象是否是同一个
	assert.Equal(t, car, car2)

	fmt.Println(do.HealthCheck[IEngineService](injector))
	// returns "engine broken"

	statuses := injector.HealthCheck()
	fmt.Println(statuses)
	// map[*do.CarService:<nil> *do.IEngineService:engine broken]

	// injector.ShutdownOnSIGTERM()    // will block until receiving sigterm signal
	injector.Shutdown()
	// prints "car stopped"
}
