package utility

import "errors"

var (
	ErrorBedinTime error = errors.New("TimeError: you don't sleep in p.m 0 ~ p.m 9")
)
