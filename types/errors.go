package types

import "errors"

var (
	ErrInvalidFormat = errors.New("invalid-format")
	ErrInvalidType   = errors.New("invalid-type")
)
