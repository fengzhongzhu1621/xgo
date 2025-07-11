package validation

// Validator is the interface for automatic validation.
type Validator interface {
	Validate() error
}
