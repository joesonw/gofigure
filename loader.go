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

	root *Node
}

func New() *Loader {
	return &Loader{}
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

	// root node is a document node, and the first child is a map holds all the values
	fileNode := NewNode(yamlNode.Content[0], NodeFilepath(name))
	name = strings.TrimSuffix(name, filepath.Ext(name))
	names := strings.Split(name, string(filepath.Separator))
	// nest the file with its path, e.g, config/app.yaml -> config.app
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

			// always try to resolve the node, so if it has resolvedNode, it will be used instead
			current, err = l.resolve(ctx, current)
			if err != nil {
				return nil, err
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

func (l *Loader) resolve(ctx context.Context, node *Node) (resultNode *Node, reterr error) {
	if node.resolved {
		// if the node is resolved by tagged resolver, the result is stored in resolvedNode (so the original value can be preserved)
		if node.resolvedNode != nil {
			return node.resolvedNode, nil
		}
		return node, nil
	}

	// set it to true first to avoid infinite loop
	node.resolved = true
	defer func() {
		if reterr != nil {
			node.resolved = false
		}
	}()

	// resolve children first for mapping and sequence nodes

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

	// if not tagged, no further action is needed
	if node.style&yaml.TaggedStyle == 0 {
		return node, nil
	}

	// resolve the node with the feature if matched
	for _, feature := range l.features {
		if feature.Name() == node.tag {
			result, err := feature.Resolve(ctx, l, node)
			if err != nil {
				node.resolved = false
				return nil, err
			}
			node.resolvedNode = result
			return node, nil
		}
	}

	return node, nil
}
