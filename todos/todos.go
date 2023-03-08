package todos

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Comment struct {
	FilePath string `json:"filePath"`
	Line     int    `json:"line"`
	Type     string `json:"type"`
	Text     string `json:"text"`
	Author   string `json:"author"`
}

func Search(dir string, commentTypes []string, ignores []string) ([]Comment, error) {
	// FIXME: For some reason the hidden files and folder ignore pattern is not working
	searchHidden := true
	for i, ignore := range ignores {
		// remove hidden files and directories regex
		if ignore == ".*" {
			searchHidden = false
			// remove the hidden files and directories regex from the ignores
			ignores = append(ignores[:i], ignores[i+1:]...)
		}
	}

	// Define regular expression to match the specified comment types
	commentRegex := regexp.MustCompile(fmt.Sprintf(`(?i)\s*(%s)(?:\(([\w.-]+)\))?:\s*(.*)`, strings.Join(commentTypes, "|")))

	// Create a slice to hold the comments
	comments := []Comment{}

	// Walk the directory tree and search for comments in each file
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		// FIXME: For some reason the hidden files and folder ignore pattern is not working
		if info.IsDir() {
			if !searchHidden && strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
				return filepath.SkipDir
			}

			return nil
		}
		// Ignore hidden files
		if !searchHidden && strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// Check if file or directory should be ignored based on patterns from .gitignore and additional ignores
		for _, pattern := range ignores {
			matched, err := filepath.Match(pattern, info.Name())
			if err != nil {
				return err
			}
			if matched {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Open the file and search for comments
		file, outErr := os.Open(path)
		if outErr != nil {
			// Ignore directories that can't be opened
			if os.IsPermission(err) || os.IsNotExist(err) {
				return nil
			}
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for i := 1; scanner.Scan(); i++ {
			line := scanner.Text()
			if matches := commentRegex.FindStringSubmatch(line); matches != nil {
				commentType := strings.ToUpper(matches[1])
				author := matches[2]
				commentText := strings.TrimSpace(matches[3])
				comment := Comment{
					FilePath: path,
					Line:     i,
					Type:     commentType,
					Text:     commentText,
					Author:   author,
				}
				comments = append(comments, comment)
			}
		}
		if err := scanner.Err(); err != nil {
			if scanner.Err() == bufio.ErrTooLong {
				return nil
			}
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return comments, nil
}

// Parses the .gitignore file in the specified directory and returns a slice of
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
