package dew

import (
	"context"
	"fmt"
	"testing"

	"github.com/fengzhongzhu1621/xgo/collections/messagebus/dew/commands/query"
	"github.com/fengzhongzhu1621/xgo/collections/messagebus/dew/handlers"
	"github.com/go-dew/dew"
)

func TestQuery(t *testing.T) {
	// Initialize the Command Bus.
	bus := dew.New()

	// Register the handler for the HelloAction.
	bus.Register(new(handlers.HelloHandler))

	// Create a context with the bus.
	ctx := dew.NewContext(context.Background(), bus)

	// Execute the query.
	query, err := dew.Query(ctx, &query.HelloQuery{Name: "Dew"})
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Result: %+v\n", query.Result)
	}
}
