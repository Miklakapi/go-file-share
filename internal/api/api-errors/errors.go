package apierrors

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrInvalidFile    = errors.New("invalid file")
)
