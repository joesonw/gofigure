package gofigure

import (
	"fmt"
)

func newNodeError(node *Node, err error) error {
	return fmt.Errorf("%d:%d@%s: %w", node.line, node.column, node.Filepath(), err)
}
