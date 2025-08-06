package mock

// Option set options.
type Option func(*options)

// WithMock set mock data.
func WithMock(mock *Item) Option {
	return func(opts *options) {
		opts.mocks = append(opts.mocks, mock)
	}
}
