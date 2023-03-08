package gitignore

// Package gitignore implements pattern matching for .gitignore files.

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

const (
	dblAsterisks = "**"
	comment      = "#"
	negate       = "!"
)

var (
	errGlob = errors.New("unable to glob pattern")
)

// Match matches patterns in the same manner that gitignore does.
// Reference https://git-scm.com/docs/gitignore.
func Match(pattern, value string) (bool, error) {
	if pattern == "" {
		return false, nil
	}

	if strings.HasPrefix(pattern, comment) {
		return false, nil
	}

	pattern = strings.TrimSuffix(pattern, " ")

	negated := strings.HasPrefix(pattern, negate)
	if negated {
		pattern = strings.TrimPrefix(pattern, negate)
	}

	pattern = strings.TrimSuffix(pattern, string(filepath.Separator))

	if strings.Contains(pattern, dblAsterisks) {
		return matchDblAsterisk(pattern, value, negated)
	}

	if !strings.Contains(pattern, string(filepath.Separator)) {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return false, fmt.Errorf("%w: %s", errGlob, pattern)
		}

		for _, match := range matches {
			if match == value {
				return !negated, nil
			}
		}

		return negated, nil
	}

	pattern = filepath.ToSlash(pattern)
	value = filepath.ToSlash(value)

	matched := false
	err := filepath.Walk(filepath.Dir(value), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		if matched {
			return filepath.SkipDir
		}

		match, err := filepath.Match(pattern, path)
		if err != nil {
			return fmt.Errorf("%w: %s", errGlob, pattern)
		}

		if match {
			matched = true
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return matched != negated, nil
}

func matchDblAsterisk(pattern, value string, negated bool) (bool, error) {
	pattern = filepath.ToSlash(pattern)
	value = filepath.ToSlash(value)

	if strings.HasPrefix(pattern, dblAsterisks) {
		pattern = strings.TrimPrefix(pattern, dblAsterisks)
		return strings.HasSuffix(value, pattern) != negated, nil
	}

	if strings.HasSuffix(pattern, dblAsterisks) {
		pattern = strings.TrimSuffix(pattern, dblAsterisks)
		return strings.HasPrefix(value, pattern) != negated, nil
	}

	parts := strings.Split(pattern, dblAsterisks)
	for i, part := range parts {
		part = filepath.ToSlash(part)

		switch i {
		case 0:
			if !strings.HasPrefix(value, part) {
				return false, nil
			}
		case len(parts) - 1:
			if !strings.HasSuffix(value, part) {
				return false, nil
			}
		default:
			index := strings.Index(value, part)
			if index == -1 {
				return false, nil
			}

			value = value[index+len(part):]
		}
	}

	return !negated, nil
}
