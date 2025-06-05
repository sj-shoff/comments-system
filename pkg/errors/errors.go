package errors

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrParentNotFound = errors.New("parent comment not found")
)
