package internal

import (
	"dario.cat/mergo"
)

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

func ApplyDefaults[T any](defaults *T, supplied *T) *T {
	if supplied == nil {
		return defaults
	}

	res := new(T)
	_ = mergo.Merge(res, supplied)
	_ = mergo.Merge(res, defaults)
	return res
}
