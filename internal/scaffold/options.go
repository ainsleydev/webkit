package scaffold

// Option is a function that configures write options
type Option func(*writeOptions)

// writeOptions holds configuration for write operations
type writeOptions struct {
	mode WriteMode
}

// WithScaffoldMode sets the write mode for the operation
func WithScaffoldMode() Option {
	return func(opts *writeOptions) {
		opts.mode = ModeScaffold
	}
}

// defaultOptions returns the default write options
func defaultOptions() *writeOptions {
	return &writeOptions{
		mode: ModeGenerate,
	}
}

// applyOptions applies the given options to writeOptions
func applyOptions(opts ...Option) *writeOptions {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return o
}
