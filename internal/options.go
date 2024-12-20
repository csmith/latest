package internal

import (
	"dario.cat/mergo"
)

func ApplyDefaults[T any](defaults *T, supplied *T) *T {
	if supplied == nil {
		return defaults
	}

	res := new(T)
	_ = mergo.Merge(res, supplied)
	_ = mergo.Merge(res, defaults)
	return res
}
