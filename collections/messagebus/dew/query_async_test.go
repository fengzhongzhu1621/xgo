package dew

import (
	"context"
	"fmt"
	"testing"

	"github.com/fengzhongzhu1621/xgo/collections/messagebus/dew/commands/query"
	"github.com/fengzhongzhu1621/xgo/collections/messagebus/dew/handlers"
	"github.com/go-dew/dew"
)

// Use QueryAsync for handling multiple queries concurrently:
func TestQueryAsync(t *testing.T) {
	// Initialize the Command Bus.
	bus := dew.New()

	// Register the handler for the HelloAction.
	bus.Register(new(handlers.HelloHandler))

	// Create a context with the bus.
	ctx := dew.NewContext(context.Background(), bus)

	helloQuery := &query.HelloQuery{Name: "Dew"}
	orgQuery := &query.HelloQuery{Name: "Dew"}

	err := dew.QueryAsync(ctx, dew.NewQuery(helloQuery), dew.NewQuery(orgQuery))
	if err != nil {
		fmt.Println("Error executing queries:", err)
	} else {
		fmt.Printf("Result: %+v\n", helloQuery.Result)
	}
}
