package gofigure

import (
	"context"
)

type Feature interface {
	Name() string
	Resolve(ctx context.Context, loader *Loader, node *Node) (*Node, error)
}
