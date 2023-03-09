package todos

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"sort"
	"strings"
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
		fileGroups[comment.File] = append(fileGroups[comment.File], comment)
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
					return comments[i].File < comments[j].File
				case "line":
					return comments[i].Line < comments[j].Line
				case "type":
					return comments[i].Type < comments[j].Type
				case "text":
					return comments[i].Text < comments[j].Text
				default:
					return comments[i].File < comments[j].File
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
		commentString := fmt.Sprintf("%s\t%s\t%s:%d\t%s", comment.Author, comment.Type, comment.File, comment.Line, comment.Text)
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

	fmt.Println("| Type | Author | File:Line | Text |")
	fmt.Println("| --- | --- | --- | --- |")

	for _, comment := range comments {
		fmt.Printf("| %s | %s | %s | %s:%d |\n", comment.Type, comment.Author, comment.Text, comment.File, comment.Line)
	}

	return nil
}

// WriteTemplate writes the comments to the io.Writer using the given template
func WriteTemplate(w io.Writer, comments []Comment, sortby string, desc bool, sourceStr string) error {
	if len(comments) == 0 {
		return nil
	}

	sortComments(comments, sortby, desc)

	sourceStr = strings.Replace(sourceStr, `\n`, "\n", -1)
	sourceStr = strings.Replace(sourceStr, `\t`, "\t", -1)

	t, err := template.New("comments").Parse(sourceStr)
	if err != nil {
		return err
	}

	return t.Execute(w, comments)
}

// sortComments sorts the comments by the given sortby string and
// reverses the order if the sortby string ends with ":desc"
func sortComments(comments []Comment, sortby string, desc bool) {
	sort.Slice(comments, func(i, j int) bool {
		switch sortby {
		case "author":
			return comments[i].Author < comments[j].Author
		case "file":
			return comments[i].File < comments[j].File
		case "line":
			return comments[i].Line < comments[j].Line
		case "type":
			return comments[i].Type < comments[j].Type
		case "text":
			return comments[i].Text < comments[j].Text
		default:
			return comments[i].File < comments[j].File
		}
	})

	if desc {
		for i, j := 0, len(comments)-1; i < j; i, j = i+1, j-1 {
			comments[i], comments[j] = comments[j], comments[i]
		}
	}
}
