package p

import "golang.org/x/exp/constraints"

func NotNil(v interface{}) bool {
	return v != nil
}

func Int[T constraints.Integer](v T) *T {
	return &v
}

func String(v string) *string {
	return &v
}

func Bool(v bool) *bool {
	return &v
}
