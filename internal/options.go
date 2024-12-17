package internal

func ResolveOptions[T any](opts []func(*T)) *T {
	opt := new(T)
	for _, o := range opts {
		o(opt)
	}
	return opt
}
