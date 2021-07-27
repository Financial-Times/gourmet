package structloader

const (
	fallbackLoaderTagName   = "data"
	fallbackDefaultTagName  = "default"
	fallbackRequiredTagName = "required"
)

type loaderOptions struct {
	loaderTagName   string
	defaultTagName  string
	requiredTagName string
}

func newLoaderOptions() *loaderOptions {
	return &loaderOptions{
		loaderTagName:   fallbackLoaderTagName,
		defaultTagName:  fallbackDefaultTagName,
		requiredTagName: fallbackRequiredTagName,
	}
}

func (o *loaderOptions) applyOptions(opts ...LoaderOption) {
	for _, opt := range opts {
		opt(o)
	}
}

type LoaderOption func(*loaderOptions)

func WithLoaderTagName(name string) LoaderOption {
	return func(o *loaderOptions) {
		o.loaderTagName = name
	}
}

func WithDefaultTagName(name string) LoaderOption {
	return func(o *loaderOptions) {
		o.defaultTagName = name
	}
}

func WithRequiredTagName(name string) LoaderOption {
	return func(o *loaderOptions) {
		o.requiredTagName = name
	}
}
