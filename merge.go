package gofigure

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func PackNodeInNestedKeys(node *Node, keys ...string) *Node {
	result := node
	for i := len(keys) - 1; i >= 0; i-- {
		key := keys[i]
		childNode := NewMappingNode(map[string]*Node{
			key: result,
		})

		result.parent = childNode
		result.mappingKey = key
		result.hasMappingKey = true
		result = childNode
	}
	return result
}

func MergeNodes(nodes ...*Node) (*Node, error) {
	var rootNode *Node
	for _, node := range nodes {
		if rootNode == nil {
			rootNode = node
		} else {
			result, err := mergeToNode(rootNode, node)
			if err != nil {
				return nil, err
			}
			rootNode = result
		}
	}

	return rootNode, nil
}

func mergeToNode(n, another *Node) (*Node, error) {
	var err error

	shouldAppend := false
	if another.style == yaml.TaggedStyle {
		if another.tag != "!append" {
			return another, nil
		}
		another.style = n.style
		another.tag = ""
		shouldAppend = true
	}

	if n.kind != another.kind {
		return nil, fmt.Errorf("cannot merge %d with %d", n.kind, another.kind)
	}

	switch n.kind {
	case yaml.AliasNode:
		return nil, fmt.Errorf("alias node cannot be merged")
	case yaml.DocumentNode:
		return nil, fmt.Errorf("document node cannot be merged")
	case yaml.MappingNode:
		for key, value := range another.mappingNodes {
			if destNode, ok := n.mappingNodes[key]; ok {
				n.mappingNodes[key], err = mergeToNode(destNode, value)
				if err != nil {
					return nil, err
				}
			} else {
				n.mappingNodes[key] = value
			}
		}
	case yaml.SequenceNode:
		if shouldAppend {
			n.sequenceNodes = append(n.sequenceNodes, another.sequenceNodes...)
			return n, nil
		}
		return another, nil
	case yaml.ScalarNode:
		return another, nil
	}
	return n, nil
}
