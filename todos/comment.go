package todos

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"
)

// WriteFileGroup writes the comments to the io.Writer as a file list
func WriteFileGroup(w io.Writer, comments []Comment, sortby string, desc bool) error {
	if w == nil {
		return fmt.Errorf("output is nil")
	}

	tabW := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fileGroups := make(map[string][]Comment)
	for _, comment := range comments {
		fileGroups[comment.FilePath] = append(fileGroups[comment.FilePath], comment)
	}
	keys := make([]string, 0, len(fileGroups))
	sort.Strings(keys)

	for file, comments := range fileGroups {
		fmt.Fprintf(tabW, "%s [%d Comments]:\n", file, len(comments))
		if sortby != "" {
			sort.Slice(comments, func(i, j int) bool {
				switch sortby {
				case "author":
					return comments[i].Author < comments[j].Author
				case "file":
					return comments[i].FilePath < comments[j].FilePath
				case "line":
					return comments[i].Line < comments[j].Line
				case "type":
					return comments[i].Type < comments[j].Type
				case "text":
					return comments[i].Text < comments[j].Text
				default:
					return comments[i].FilePath < comments[j].FilePath
				}
			})

			if desc {
				for i, j := 0, len(comments)-1; i < j; i, j = i+1, j-1 {
					comments[i], comments[j] = comments[j], comments[i]
				}
			}
		}

		for i, comment := range comments {
			author := ""
			if comment.Author != "" {
				author = "(" + comment.Author + ")"
			}

			fmt.Fprintf(tabW, "%d\t|\t%s%s:\t%s\t\n", comment.Line, comment.Type, author, comment.Text)
			if i == len(comments)-1 {
				fmt.Fprintf(tabW, "\n")
			}
		}
	}

	return nil
}

// WriteJSON writes the comments to the io.Writer as JSON
func WriteJSON(w io.Writer, comments []Comment, sortby string, desc bool) error {
	if len(comments) == 0 {
		return nil
	}

	sortComments(comments, sortby, desc)

	commentsJSON, err := json.MarshalIndent(comments, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		return err
	}

	if _, err := w.Write(commentsJSON); err != nil {
		return err
	}

	return nil
}

// WriteTable writes the comments to the io.Writer as a table
func WriteTable(w io.Writer, comments []Comment, sortby string, desc bool) error {
	if len(comments) == 0 {
		return nil
	}

	sortComments(comments, sortby, desc)

	tabW := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	for _, comment := range comments {
		commentString := fmt.Sprintf("%s\t%s\t%s:%d\t%s", comment.Author, comment.Type, comment.FilePath, comment.Line, comment.Text)
		fmt.Fprintln(tabW, commentString)
	}

	return nil
}

// WriteMarkdown writes the comments to the io.Writer as a markdown table
func WriteMarkdown(w io.Writer, comments []Comment, sortby string, desc bool) error {
	if len(comments) == 0 {
		return nil
	}

	sortComments(comments, sortby, desc)

	fmt.Println("## TODOs")
	fmt.Println("")
	fmt.Println("| Author | Type | File:Line | Text |")
	fmt.Println("| --- | --- | --- | --- |")

	for _, comment := range comments {
		fmt.Printf("| %s | %s | %s:%d | %s |\n", comment.Author, comment.Type, comment.FilePath, comment.Line, comment.Text)
	}

	return nil
}

// sortComments sorts the comments by the given sortby string and
// reverses the order if the sortby string ends with ":desc"
func sortComments(comments []Comment, sortby string, desc bool) {
	sort.Slice(comments, func(i, j int) bool {
		switch sortby {
		case "author":
			return comments[i].Author < comments[j].Author
		case "file":
			return comments[i].FilePath < comments[j].FilePath
		case "line":
			return comments[i].Line < comments[j].Line
		case "type":
			return comments[i].Type < comments[j].Type
		case "text":
			return comments[i].Text < comments[j].Text
		default:
			return comments[i].FilePath < comments[j].FilePath
		}
	})

	if desc {
		for i, j := 0, len(comments)-1; i < j; i, j = i+1, j-1 {
			comments[i], comments[j] = comments[j], comments[i]
		}
	}
}
