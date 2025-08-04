package do

import (
	"fmt"
	"testing"

	"github.com/samber/do"
)

/**
 * Wheel
 */
type Wheel struct{}

/**
 * Engine
 */
type Engine2 struct{}

func (e *Engine2) HealthCheck() error {
	return fmt.Errorf("engine broken")
}

/**
 * Car
 */
type Car2 struct {
	Engine *Engine2
	Wheels []*Wheel
}

func (c *Car2) Start() {
	println("vroooom")
}

func TestMultipleObject(t *testing.T) {
	injector := do.New()

	// provide wheels
	do.ProvideNamedValue(injector, "wheel-1", &Wheel{})
	do.ProvideNamedValue(injector, "wheel-2", &Wheel{})
	do.ProvideNamedValue(injector, "wheel-3", &Wheel{})
	do.ProvideNamedValue(injector, "wheel-4", &Wheel{})

	// provide car
	do.Provide(injector, func(i *do.Injector) (*Car2, error) {
		car := Car2{
			Engine: do.MustInvoke[*Engine2](i),
			Wheels: []*Wheel{
				do.MustInvokeNamed[*Wheel](i, "wheel-1"),
				do.MustInvokeNamed[*Wheel](i, "wheel-2"),
				do.MustInvokeNamed[*Wheel](i, "wheel-3"),
				do.MustInvokeNamed[*Wheel](i, "wheel-4"),
			},
		}

		return &car, nil
	})

	// provide engine
	do.Provide(injector, func(i *do.Injector) (*Engine2, error) {
		return &Engine2{}, nil
	})

	// start car
	car := do.MustInvoke[*Car2](injector)
	car.Start() // vroooom

	// check single service
	fmt.Println(do.HealthCheck[*Engine](injector))
	// check all services
	fmt.Println(injector.HealthCheck())
}
