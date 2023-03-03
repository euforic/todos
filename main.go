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
	ignoreFiles := flag.String("ignore-files", "", "Comma-separated list of file extensions to ignore")
	ignoreDirs := flag.String("ignore-dirs", "", "Comma-separated list of directory names to ignore")
	outputJSON := flag.Bool("json", false, "Output results in JSON format")
	sortBy := flag.String("sortby", "", "Sort results by field")
	commentTypesStr := flag.String("comment-types", "TODO,FIXME", "Comma-separated list of comment types to search for")
	flag.Parse()

	// Parse the ignore flags into slices of strings
	var ignoredFiles []string
	if *ignoreFiles != "" {
		ignoredFiles = strings.Split(*ignoreFiles, ",")
	}

	var ignoredDirs []string
	if *ignoreDirs != "" {
		ignoredDirs = strings.Split(*ignoreDirs, ",")
	}

	// Parse the comment types into a slice of strings
	commentTypes := strings.Split(*commentTypesStr, ",")

	// Search for comments
	comments, err := todos.Search(*dir, ignoredFiles, ignoredDirs, commentTypes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	sort.Slice(comments, func(i, j int) bool {
		switch *sortBy {
		case "username":
			return comments[i].Username < comments[j].Username
		case "file":
			return comments[i].FilePath < comments[j].FilePath
		case "line":
			return comments[i].LineNumber < comments[j].LineNumber
		case "type":
			return comments[i].Type < comments[j].Type
		case "text":
			return comments[i].Text < comments[j].Text
		default:
			return comments[i].FilePath < comments[j].FilePath
		}
	})

	// Output the comments
	if *outputJSON {
		commentsJSON, err := json.MarshalIndent(comments, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Println(string(commentsJSON))
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "USERNAME\tTYPE\tFILE\tTEXT")

	for _, comment := range comments {
		username := "unknown"
		if comment.Username != "" {
			username = comment.Username
		}
		commentString := fmt.Sprintf("%s\t%s\t%s:%d\t%s", username, comment.Type, comment.FilePath, comment.LineNumber, comment.Text)
		fmt.Fprintln(w, commentString)
	}
	w.Flush()

}
