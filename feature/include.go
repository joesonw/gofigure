package feature

import (
	"context"
	"fmt"
	iofs "io/fs"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/joesonw/gofigure"
)

type includeFeature struct {
	fs          iofs.FS
	loadedNodes map[string]bool
}

func Include(fs iofs.FS) gofigure.Feature {
	return &includeFeature{
		fs:          fs,
		loadedNodes: map[string]bool{},
	}
}

func (*includeFeature) Name() string {
	return "!include"
}

func (f *includeFeature) Resolve(ctx context.Context, loader *gofigure.Loader, node *gofigure.Node) (*gofigure.Node, error) {
	if node.Kind() != yaml.MappingNode {
		return nil, fmt.Errorf("!include requires a mapping node")
	}

	fileNode, err := node.GetMappingChild("file")
	if err != nil {
		return nil, gofigure.NewNodeError(node, err)
	}

	if fileNode != nil {
		pathNode, err := fileNode.GetMappingChild("path")
		if err != nil {
			return nil, gofigure.NewNodeError(fileNode, fmt.Errorf("unable to get path: %w", err))
		}
		if pathNode == nil {
			return nil, gofigure.NewNodeError(fileNode, fmt.Errorf("key \"path\" is missing for file"))
		}

		parseNode, err := fileNode.GetMappingChild("parse")
		if err != nil {
			return nil, gofigure.NewNodeError(fileNode, fmt.Errorf("unable to get parse: %w", err))
		}

		keyNode, err := fileNode.GetMappingChild("key")
		if err != nil {
			return nil, gofigure.NewNodeError(fileNode, fmt.Errorf("unable to get key: %w", err))
		}

		parse := false
		if parseNode != nil {
			parse, err = parseNode.BoolValue()
			if err != nil {
				return nil, gofigure.NewNodeError(parseNode, err)
			}
		}

		path := strings.TrimSpace(pathNode.Value())
		contents, err := iofs.ReadFile(f.fs, path)
		if err != nil {
			return nil, gofigure.NewNodeError(pathNode, fmt.Errorf("unable to read file %q: %w", path, err))
		}

		if !parse {
			return gofigure.NewScalarNode(string(contents)), nil
		}

		if !f.loadedNodes[path] {
			if err := loader.Load(path, contents); err != nil {
				return nil, gofigure.NewNodeError(pathNode, fmt.Errorf("unable to load file %q: %w", path, err))
			}
			f.loadedNodes[path] = true
		}

		dotPath := filepath.Clean(path)
		dotPath = strings.TrimSuffix(dotPath, filepath.Ext(dotPath))
		dotPath = strings.ReplaceAll(dotPath, string(filepath.Separator), ".")

		if keyNode != nil {
			dotPath += "." + strings.TrimSpace(keyNode.Value())
		}

		return loader.GetNode(ctx, dotPath)
	}

	return gofigure.NewScalarNode(""), nil
}
