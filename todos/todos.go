package todos

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/euforic/todos/pkg/gitignore"
)

// Comment represents a comment
type Comment struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Type   string `json:"type"`
	Text   string `json:"text"`
	Author string `json:"author"`
}

// Search searches a directory for comments
func Search(dir string, commentTypes []string, ignores []string) ([]Comment, error) {
	searchHidden, ignores := removeHiddenIgnore(ignores)

	comments := []Comment{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if info.IsDir() {
			if !searchHidden && strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}

		if shouldIgnoreFile(info, ignores, path, searchHidden) {
			return nil
		}

		// Open the file and search for comments
		file, openErr := os.Open(path)
		if openErr != nil {
			// Ignore directories that can't be opened
			if os.IsPermission(openErr) || os.IsNotExist(openErr) {
				return nil
			}
			return openErr
		}
		defer file.Close()

		fileComments, parseErr := Parse(file, path, commentTypes)
		if parseErr != nil {
			return parseErr
		}

		comments = append(comments, fileComments...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return comments, nil
}

// ParseFile parses a file for comments
func removeHiddenIgnore(ignores []string) (bool, []string) {
	searchHidden := true
	for i, ignore := range ignores {
		if ignore == ".*" {
			searchHidden = false
			ignores = append(ignores[:i], ignores[i+1:]...)
		}
	}
	return searchHidden, ignores
}

// shouldIgnoreFile returns true if the file should be ignored.
func shouldIgnoreFile(info os.FileInfo, ignores []string, path string, searchHidden bool) bool {
	if !searchHidden && strings.HasPrefix(info.Name(), ".") {
		return true
	}

	for _, pattern := range ignores {
		matched, err := gitignore.Match(pattern, path)
		if err != nil {
			return true
		}

		matchedPath, err := filepath.Match(pattern, info.Name())
		if err != nil {
			return true
		}

		if matched || matchedPath {
			return !info.IsDir()
		}
	}
	return false
}

// Parse parses the specified file and returns a slice of comments.
func Parse(r io.Reader, path string, commentTypes []string) ([]Comment, error) {
	// Define regular expression to match the specified comment types
	commentRegex := regexp.MustCompile(fmt.Sprintf(`(?i)\s*(%s)(?:\(([\w.-]+)\))?:\s*(.*)`, strings.Join(commentTypes, "|")))

	// Create a slice to hold the comments
	var comments []Comment

	scanner := bufio.NewScanner(r)
	for i := 1; scanner.Scan(); i++ {
		line := scanner.Text()
		if matches := commentRegex.FindStringSubmatch(line); matches != nil {
			commentType := strings.ToUpper(matches[1])
			author := matches[2]
			commentText := strings.TrimSpace(matches[3])
			comment := Comment{
				File:   path,
				Line:   i,
				Type:   commentType,
				Text:   commentText,
				Author: author,
			}
			comments = append(comments, comment)
		}
	}
	if err := scanner.Err(); err != nil {
		// Ignore lines that are too long
		if scanner.Err() != bufio.ErrTooLong {
			return nil, err
		}
	}
	return comments, nil
}

// ParseGitignore parses the .gitignore file in the specified directory and returns a slice of
// patterns to ignore.
func ParseGitignore(dir string) ([]string, error) {
	// Open the .gitignore file
	file, err := os.Open(filepath.Join(dir, ".gitignore"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a slice to hold the patterns
	patterns := []string{}

	// Read the patterns from the .gitignore file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Ignore comments
		if strings.HasPrefix(line, "#") {
			continue
		}
		// Ignore blank lines
		if line == "" {
			continue
		}
		patterns = append(patterns, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return patterns, nil
}
