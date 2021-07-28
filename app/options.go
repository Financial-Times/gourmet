package app

type appOptions struct {
	terminator Terminator
}

func newDefaultOptions() *appOptions {
	return &appOptions{
		terminator: OSSignalTerminator,
	}
}

// Option - app configuration
type Option func(*appOptions)

// WithTerminator - define custom app terminator
func WithTerminator(t Terminator) Option {
	return func(o *appOptions) {
		o.terminator = t
	}
}
