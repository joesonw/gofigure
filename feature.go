package gofigure

import (
	"context"
)

type Feature interface {
	Name() string
	Resolve(ctx context.Context, loader *Loader, node *Node) (*Node, error)
}

type ResolveFunc func(ctx context.Context, loader *Loader, node *Node) (*Node, error)

type featureFunc struct {
	name    string
	resolve ResolveFunc
}

func (f *featureFunc) Name() string {
	return f.name
}

func (f *featureFunc) Resolve(ctx context.Context, loader *Loader, node *Node) (*Node, error) {
	return f.resolve(ctx, loader, node)
}

func FeatureFunc(name string, resolve ResolveFunc) Feature {
	return &featureFunc{
		name:    name,
		resolve: resolve,
	}
}
