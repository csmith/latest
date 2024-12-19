package internal

func ResolveOptions[T any](opts []func(*T)) *T {
	return ResolveOptionsWithDefaults(opts, new(T))
}

func ResolveOptionsWithDefaults[T any](opts []func(*T), defaults *T) *T {
	opt := defaults
	for _, o := range opts {
		o(opt)
	}
	return opt
}
