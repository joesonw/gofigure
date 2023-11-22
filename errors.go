package gofigure

import (
	"errors"
	"fmt"
)

var (
	ErrPathNotFound     = errors.New("path not found")
	ErrConfigParseError = errors.New("config parse error")
	ErrInvalidPath      = errors.New("invalid path")
)

type nodeError struct {
	node *Node
	err  error
}

func (e *nodeError) Error() string {
	return fmt.Sprintf("%s.yaml@%d:%d %s", e.node.Filepath(), e.node.Line(), e.node.Column(), e.err.Error())
}

func (e *nodeError) Unwrap() error {
	return e.err
}

func NewNodeError(node *Node, err error) error {
	return &nodeError{
		node: node,
		err:  err,
	}
}

func IsNodeError(err error) bool {
	_, ok := err.(*nodeError)
	return ok
}
