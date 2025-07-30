package assign

type Option interface {
	apply(*options)
}

// DryRun specifies that a dry-run should be executed and no individuals should
// be assigned.
func DryRun(ok bool) Option {
	return option(func(o *options) {
		o.dryRun = ok
	})
}

type options struct {
	dryRun bool
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt.apply(o)
	}
}

type option func(*options)

func (o option) apply(opts *options) {
	o(opts)
}
