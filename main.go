package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/euforic/todos/todos"
)

func main() {
	// Define command line flags
	dir := flag.String("dir", ".", "Directory to search for comments")
	ignores := flag.String("ignore", "", "Comma-separated list of files and directories to ignore")
	sortBy := flag.String("sortby", "", "Sort results by field")
	commentTypesStr := flag.String("types", "TODO,FIXME", "Comma-separated list of comment types to search for")
	searchHidden := flag.Bool("hidden", false, "Search hidden files and directories")
	validateMax := flag.Int("validate-max", 0, "Validate that the number of comments is less than or equal to the max")
	outputStyle := flag.String("output", "file", "Output style (table, file, json)")
	flag.Parse()

	// Parse the ignore flag into a slice of strings
	var ignoreList []string
	if *ignores != "" {
		ignoreList = strings.Split(*ignores, ",")
	}

	// Add the hidden files and directories ignore
	if !*searchHidden {
		ignoreList = append(ignoreList, ".*")
	}

	ignorePatterns, err := todos.ParseGitignore(*dir)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	ignoreList = append(ignoreList, ignorePatterns...)

	// Parse the comment types into a slice of strings
	commentTypes := strings.Split(*commentTypesStr, ",")

	// Search for comments
	comments, err := todos.Search(*dir, commentTypes, ignoreList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	// Validate the max number of comments
	if *validateMax > 0 && len(comments) > *validateMax {
		fmt.Fprintf(os.Stderr, "Error: %d comments found, max is %d", len(comments), *validateMax)
		os.Exit(1)
	}

	// Sort the comments
	if *sortBy == "" {
		sort.Slice(comments, func(i, j int) bool {
			switch *sortBy {
			case "username":
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
	}

	// Output the comments
	switch *outputStyle {
	case "table":
		outputTable(comments)
	case "file":
		outputFile(comments)
	case "json":
		outputJSON(comments)
	default:
		outputTable(comments)
	}
}

// outputFile outputs the comments in a file format
func outputFile(comments []todos.Comment) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	output := make(map[string][]todos.Comment)
	for _, comment := range comments {
		output[comment.FilePath] = append(output[comment.FilePath], comment)
	}

	for file, comments := range output {
		fmt.Fprintf(w, "%s [%d Comments]:\n", file, len(comments))
		for i, comment := range comments {
			author := ""
			if comment.Author != "" {
				author = "(" + comment.Author + ")"
			}

			fmt.Fprintf(w, "%d\t|\t%s%s:\t%s\t\n", comment.Line, comment.Type, author, comment.Text)
			if i == len(comments)-1 {
				fmt.Fprintf(w, "\n")
			}
		}
	}
	w.Flush()
}

// outputJSON outputs the comments in JSON format
func outputJSON(comments []todos.Comment) {
	commentsJSON, err := json.MarshalIndent(comments, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(string(commentsJSON))
}

// outputTable outputs the comments in a tabular format
func outputTable(comments []todos.Comment) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	for _, comment := range comments {
		username := "unknown"
		if comment.Author != "" {
			username = comment.Author
		}
		commentString := fmt.Sprintf("%s\t%s\t%s:%d\t%s", username, comment.Type, comment.FilePath, comment.Line, comment.Text)
		fmt.Fprintln(w, commentString)
	}
	w.Flush()
}
