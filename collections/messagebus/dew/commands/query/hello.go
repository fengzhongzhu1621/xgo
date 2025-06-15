package query

// HelloQuery is a simple query that returns a greeting message.
type HelloQuery struct {
	// Name is the name of the user.
	Name string

	// Result is the output of the query.
	// You can define any struct as the result.
	Result string
}
