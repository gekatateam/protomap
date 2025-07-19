package protomap

import "errors"

var (
	ErrNoSuchFile    = errors.New("no such file")
	ErrNoMessages    = errors.New("descriptor does not contains messages")
	ErrNoSuchMessage = errors.New("no such message in descriptor")
)
