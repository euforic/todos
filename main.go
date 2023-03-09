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
	outputStyle := flag.String("output", "table", "Output style (table, group, json, md)")
	format := flag.String("format", "", "Go template string to use for output style (-output will be ignored if format is set)")
	flag.Parse()

	dir := flag.Arg(0)
	if dir == "" {
		dir = "."
	}

	ignoreList := parseIgnoreList(ignores, searchHidden, dir)

	commentTypes := strings.Split(*commentTypesStr, ",")

	comments, err := todos.Search(dir, commentTypes, ignoreList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	validateCommentsCount(*validateMax, comments)

	sortField, sortDesc := parseSortBy(*sortBy)

	formatStr := ""
	if *format != "" {
		*outputStyle = "format"
		formatStr = *format
	}

	outputComments(*outputStyle, comments, sortField, sortDesc, formatStr)
}

// parseIgnoreList parses the ignore list from the command line
func parseIgnoreList(ignores *string, searchHidden *bool, dir string) []string {
	ignoreList := []string{}

	if *ignores != "" {
		ignoreList = strings.Split(*ignores, ",")
	}

	if !*searchHidden {
		ignoreList = append(ignoreList, ".*")
	}

	ignorePatterns, err := todos.ParseGitignore(dir)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	ignoreList = append(ignoreList, ignorePatterns...)
	return ignoreList
}

// validateCommentsCount validates that the number of comments is less than or equal to the max
func validateCommentsCount(validateMax int, comments []todos.Comment) {
	if validateMax > 0 && len(comments) > validateMax {
		fmt.Fprintf(os.Stderr, "Error: %d comments found, max is %d", len(comments), validateMax)
		os.Exit(1)
	}
}

// parseSortBy parses the sortby flag from the command line
func parseSortBy(sortBy string) (string, bool) {
	sortField := ""
	sortDesc := false

	sortbyParts := strings.Split(sortBy, ":")
	if len(sortbyParts) > 0 {
		sortField = sortbyParts[0]
	}

	if len(sortbyParts) > 1 && sortbyParts[1] == "desc" {
		sortDesc = true
	}

	return sortField, sortDesc
}

// outputComments outputs the comments in the specified format
func outputComments(outputStyle string, comments []todos.Comment, sortField string, sortDesc bool, formatStr string) {
	var outputErr error
	switch outputStyle {
	case "table":
		outputErr = todos.WriteTable(os.Stdout, comments, sortField, sortDesc)
	case "group":
		outputErr = todos.WriteFileGroup(os.Stdout, comments, sortField, sortDesc)
	case "json":
		outputErr = todos.WriteJSON(os.Stdout, comments, sortField, sortDesc)
	case "md":
		outputErr = todos.WriteMarkdown(os.Stdout, comments, sortField, sortDesc)
	case "format":
		outputErr = todos.WriteTemplate(os.Stdout, comments, sortField, sortDesc, formatStr)
	default:
		outputErr = todos.WriteTable(os.Stdout, comments, sortField, sortDesc)
	}

	if outputErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", outputErr.Error())
		os.Exit(1)
	}
}
