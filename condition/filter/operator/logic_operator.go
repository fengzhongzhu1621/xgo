package operator

import "fmt"

const (
	// And logic operator
	And LogicOperator = "AND"
	// Or logic operator
	Or LogicOperator = "OR"
)

// LogicOperator defines the logic operator
type LogicOperator string

// Validate the logic operator is valid or not.
func (lo LogicOperator) Validate() error {
	switch lo {
	case And:
	case Or:
	default:
		return fmt.Errorf("unsupported expression's logic operator: %s", lo)
	}

	return nil
}
