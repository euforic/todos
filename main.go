package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/euforic/todos/todos"
)

func main() {
	// Define command line flags
	ignores := flag.String("ignore", "", "Comma-separated list of files and directories to ignore")
	sortBy := flag.String("sortby", "", "Sort results by field (author, file, line, type, text) to sort descending, postfix with ':desc' (e.g. author:desc)")
	commentTypesStr := flag.String("types", "TODO,FIXME", "Comma-separated list of comment types to search for")
	searchHidden := flag.Bool("hidden", false, "Search hidden files and directories")
	validateMax := flag.Int("validate-max", 0, "Validate that the number of comments is less than or equal to the max")
	outputStyle := flag.String("output", "group", "Output style (table, group, json, md)")
	format := flag.String("format", "", "Go template string to use for output style (-output will be ignored if format is set)")
	flag.Parse()

	dir := flag.Arg(0)
	if dir == "" {
		dir = "."
	}

	// Parse the ignore flag into a slice of strings
	var ignoreList []string
	if *ignores != "" {
		ignoreList = strings.Split(*ignores, ",")
	}

	// Add the hidden files and directories ignore
	if !*searchHidden {
		ignoreList = append(ignoreList, ".*")
	}

	ignorePatterns, err := todos.ParseGitignore(dir)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	ignoreList = append(ignoreList, ignorePatterns...)

	// Parse the comment types into a slice of strings
	commentTypes := strings.Split(*commentTypesStr, ",")

	// Search for comments
	comments, err := todos.Search(dir, commentTypes, ignoreList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	// Validate the max number of comments
	if *validateMax > 0 && len(comments) > *validateMax {
		fmt.Fprintf(os.Stderr, "Error: %d comments found, max is %d", len(comments), *validateMax)
		os.Exit(1)
	}

	var sortDesc bool

	sortbyParts := strings.Split(*sortBy, ":")
	*sortBy = sortbyParts[0]

	if len(sortbyParts) > 1 && sortbyParts[1] == "desc" {
		sortDesc = true
	}

	if *format != "" {
		*outputStyle = "format"
	}

	// Output the comments
	var outputErr error
	switch *outputStyle {
	case "table":
		outputErr = todos.WriteTable(os.Stdout, comments, *sortBy, sortDesc)
	case "group":
		outputErr = todos.WriteFileGroup(os.Stdout, comments, *sortBy, sortDesc)
	case "json":
		outputErr = todos.WriteJSON(os.Stdout, comments, *sortBy, sortDesc)
	case "md":
		outputErr = todos.WriteMarkdown(os.Stdout, comments, *sortBy, sortDesc)
	case "format":
		outputErr = todos.WriteTemplate(os.Stdout, comments, *sortBy, sortDesc, *format)
	default:
		outputErr = todos.WriteTable(os.Stdout, comments, *sortBy, sortDesc)
	}

	if outputErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", outputErr.Error())
		os.Exit(1)
	}
}
