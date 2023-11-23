package gofigure

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Loader struct {
	features []Feature

	root        *Node
	loadedFiles map[string]bool
}

func New() *Loader {
	return &Loader{
		loadedFiles: map[string]bool{},
	}
}

func (l *Loader) WithFeatures(features ...Feature) *Loader {
	l.features = append(l.features, features...)
	return l
}

func (l *Loader) Load(name string, contents []byte) error {
	name = filepath.Clean(name)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	var yamlNode yaml.Node
	if err := yaml.Unmarshal(contents, &yamlNode); err != nil {
		return fmt.Errorf("unable to unmarshal file %q: %w", name, errors.Join(err, ErrConfigParseError))
	}

	fileNode := NewNode(yamlNode.Content[0], NodeFilepath(name))
	name = strings.TrimSuffix(name, filepath.Ext(name))
	names := strings.Split(name, string(filepath.Separator))
	fileNode = PackNodeInNestedKeys(fileNode, names...)

	if l.root == nil {
		l.root = fileNode
	} else {
		newNode, err := mergeToNode(l.root, fileNode)
		if err != nil {
			return fmt.Errorf("unable to merge file %q: %w", name, errors.Join(err, ErrConfigParseError))
		}
		l.root = newNode
	}

	l.loadedFiles[name] = true
	return nil
}

func (l *Loader) Get(ctx context.Context, path string, target any) error {
	node, err := l.GetNode(ctx, path)
	if err != nil {
		return err
	}
	if node == nil {
		return nil
	}
	return node.ToYAMLNode().Decode(target)
}

func (l *Loader) GetNode(ctx context.Context, path string) (*Node, error) {
	checkFile := strings.ReplaceAll(path, ".", string(filepath.Separator))
	loaded := false
	for file := range l.loadedFiles {
		if strings.HasPrefix(checkFile, file) {
			loaded = true
			break
		}
	}

	if !loaded {
		return nil, fmt.Errorf("no file contains path %q: %w", path, ErrPathNotFound)
	}

	current := l.root
	if len(path) > 0 {
		paths, err := parseDotPath(path)
		if err != nil {
			return nil, fmt.Errorf("unable to parse path %q: %w", path, err)
		}

		for _, p := range paths {
			if p.key != "" { // map
				current, err = current.GetMappingChild(p.key)
				if err != nil {
					return nil, err
				}
			} else { // slice
				current, err = current.GetSequenceChild(p.index)
				if err != nil {
					return nil, err
				}
			}

			if current == nil {
				break
			}

			if !current.resolved {
				current, err = l.resolve(ctx, current)
				if err != nil {
					return nil, err
				}
			}
			if current == nil {
				break
			}
		}
	}

	if current == nil {
		return nil, nil
	}

	return l.resolve(ctx, current)
}

func (l *Loader) resolve(ctx context.Context, node *Node) (result *Node, reterr error) {
	if node.resolved {
		if node.resolvedNode != nil {
			return node.resolvedNode, nil
		}
		return node, nil
	}
	node.resolved = true
	defer func() {
		if reterr != nil {
			node.resolved = false
		}
	}()

	if node.kind == yaml.MappingNode {
		for key := range node.mappingNodes {
			childNode, err := l.resolve(ctx, node.mappingNodes[key])
			if err != nil {
				return nil, err
			}
			node.mappingNodes[key] = childNode
		}
	}

	if node.kind == yaml.SequenceNode {
		for i := range node.sequenceNodes {
			childNode, err := l.resolve(ctx, node.sequenceNodes[i])
			if err != nil {
				return nil, err
			}
			node.sequenceNodes[i] = childNode
		}
	}

	if node.style&yaml.TaggedStyle == 0 {
		return node, nil
	}

	for _, feature := range l.features {
		if feature.Name() == node.tag {
			result, err := feature.Resolve(ctx, l, node)
			if err != nil {
				node.resolved = false
				return nil, err
			}
			node.resolvedNode = result
			return result, nil
		}
	}

	return node, nil
}
