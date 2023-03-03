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

func Search(dir string, ignoreFiles []string, ignoreDirs []string, commentTypes []string) ([]Comment, error) {
	// Define regular expression to match the specified comment types
	commentRegex := regexp.MustCompile(fmt.Sprintf(`(?i)//\s*(%s)(?:\(([\w.-]+)\))?:\s*(.*)`, strings.Join(commentTypes, "|")))

	// Create a slice to hold the comments
	comments := []Comment{}

	// Walk the directory tree and search for comments in each file
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		// Ignore files with ignored file extensions
		for _, ext := range ignoreFiles {
			if strings.HasSuffix(info.Name(), ext) {
				return nil
			}
		}

		// Ignore directories with ignored names
		if info.IsDir() {
			for _, dir := range ignoreDirs {
				if info.Name() == dir {
					return filepath.SkipDir
				}
			}
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
