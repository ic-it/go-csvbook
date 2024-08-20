package csvbook

import "errors"

var (
	ErrInvalidLineNum  = errors.New("invalid line number")
	ErrorInvalidColNum = errors.New("invalid cell number")
	ErrInvalidMetadata = errors.New("invalid metadata")
)
