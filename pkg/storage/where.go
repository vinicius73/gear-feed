package storage

type WhereOptions struct {
	Is          *Status
	Not         *Status
	AllowMissed *bool
}

type WhereOption func(opts *WhereOptions)

func Where(opts ...WhereOption) WhereOptions {
	var options WhereOptions

	for _, opt := range opts {
		opt(&options)
	}

	return options
}

func WhereIs(status Status) WhereOption {
	return func(opts *WhereOptions) {
		opts.Is = &status
	}
}

func WhereNot(status Status) WhereOption {
	return func(opts *WhereOptions) {
		opts.Not = &status
	}
}

func WhereAllowMissed(allow bool) WhereOption {
	return func(opts *WhereOptions) {
		opts.AllowMissed = &allow
	}
}
