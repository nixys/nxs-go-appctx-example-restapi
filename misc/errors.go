package misc

import "errors"

var (
	ErrArgSuccessExit = errors.New("arg success exit")
	ErrNotFound       = errors.New("entity not found")
)
