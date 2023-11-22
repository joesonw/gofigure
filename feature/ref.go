package feature

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/joesonw/gofigure"
)

type referenceFeature struct {
}

func Reference() gofigure.Feature {
	return &referenceFeature{}
}

func (referenceFeature) Name() string {
	return "!ref"
}

func (t *referenceFeature) Resolve(ctx context.Context, loader *gofigure.Loader, node *gofigure.Node) (*gofigure.Node, error) {
	if node.Kind() != yaml.ScalarNode {
		return nil, fmt.Errorf("!ref only supports scalar node")
	}

	path := node.Value()
	result, err := loader.GetNode(ctx, path)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("!ref %s not found", node.Value())
	}
	return result, nil
}
