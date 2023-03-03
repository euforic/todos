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
	FilePath   string `json:"filePath"`
	LineNumber int    `json:"lineNumber"`
	Type       string `json:"type"`
	Text       string `json:"text"`
	Username   string `json:"username,omitempty"`
}

func Search(dir string, commentTypes []string, ignores []string) ([]Comment, error) {
	// Define regular expression to match the specified comment types
	commentRegex := regexp.MustCompile(fmt.Sprintf(`(?i)\s*(%s)(?:\(([\w.-]+)\))?:\s*(.*)`, strings.Join(commentTypes, "|")))

	// Create a slice to hold the comments
	comments := []Comment{}
	ignorePatterns := []string{}

	// Read the contents of the .gitignore file
	ignoreFilePath := filepath.Join(dir, ".gitignore")
	ignoreContent, err := os.ReadFile(ignoreFilePath)
	if err == nil {
		// Parse the .gitignore file to obtain a list of ignored files and directories
		ignoreScanner := bufio.NewScanner(strings.NewReader(string(ignoreContent)))
		for ignoreScanner.Scan() {
			pattern := ignoreScanner.Text()
			if !strings.HasPrefix(pattern, "#") && len(pattern) > 0 {
				ignorePatterns = append(ignorePatterns, pattern)
			}
		}
	}

	// Add the additional ignores to the ignore patterns
	ignorePatterns = append(ignorePatterns, ignores...)

	// Walk the directory tree and search for comments in each file
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		// Check if file or directory should be ignored based on patterns from .gitignore and additional ignores
		for _, pattern := range ignorePatterns {
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

		if info.IsDir() {
			return nil
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
				username := matches[2]
				commentText := strings.TrimSpace(matches[3])
				comment := Comment{
					FilePath:   path,
					LineNumber: i,
					Type:       commentType,
					Text:       commentText,
					Username:   username,
				}
				comments = append(comments, comment)
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return comments, nil
}
