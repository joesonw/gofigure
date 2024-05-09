package gofigure

import (
	"fmt"
	"strconv"
)

type DotPath struct {
	Key   string
	Index int
}

//nolint:gocyclo
func ParseDotPath(path string) ([]*DotPath, error) {
	if path == "" {
		return nil, nil
	}
	var paths []*DotPath
	curr := ""
	isInBracket := false
	if path[0] == '.' || path[len(path)-1] == '.' { // cannot start with dot
		return nil, fmt.Errorf("%s: %w", path, ErrInvalidPath)
	}
	for i := 0; i < len(path); i++ {
		if path[i] == '.' {
			if isInBracket { // unterminated bracket
				return nil, fmt.Errorf("%s: %w", path, ErrInvalidPath)
			}

			if curr == "" { // no key specified
				if len(paths) > 0 && paths[len(paths)-1].Key == "" {
					continue
				}
				return nil, fmt.Errorf("%s: %w", path, ErrInvalidPath)
			}
			paths = append(paths, &DotPath{
				Key: curr,
			})
			curr = ""
			continue
		}

		if path[i] == '[' {
			if isInBracket { // unterminated bracket
				return nil, fmt.Errorf("%s: %w", path, ErrInvalidPath)
			}
			if curr != "" {
				paths = append(paths, &DotPath{
					Key: curr,
				})
				curr = ""
			}
			isInBracket = true
			continue
		}

		if path[i] == ']' {
			if !isInBracket { // not in bracket
				return nil, fmt.Errorf("%s: %w", path, ErrInvalidPath)
			}

			if curr == "" { // no index specified
				return nil, fmt.Errorf("%s: %w", path, ErrInvalidPath)
			}

			index, err := strconv.Atoi(curr)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", path, ErrInvalidPath)
			}

			isInBracket = false
			paths = append(paths, &DotPath{
				Index: index,
			})
			curr = ""
			continue
		}

		curr += string(path[i])
	}

	if isInBracket { // unterminated bracket
		return nil, fmt.Errorf("%s: %w", path, ErrInvalidPath)
	}

	if curr != "" {
		paths = append(paths, &DotPath{
			Key: curr,
		})
	}

	return paths, nil
}
