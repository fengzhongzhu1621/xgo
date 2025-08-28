package action

import (
	"context"
	"fmt"
)

// HelloAction 表示问候某人的动作。
type HelloAction struct {
	Name string
}

// Validate 检查 HelloAction 是否有效。
func (c HelloAction) Validate(_ context.Context) error {
	if c.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}
