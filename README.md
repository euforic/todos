# Todos

This is a command-line program written in Go that searches for `// TODO: xxxx` style comments in a directory and outputs them to the console or in JSON format.

## Usage

The program accepts the following command-line arguments:

- `-dir`: Specifies the directory to search for comments. Defaults to the current directory.
- `-ignore`: A comma-separated list of files and directories to ignore, in gitignore format.
- `-sortby`: Sort results by field (`author`, `file`, `line`, `type`, or `text`)
- `-output`: Output style (table, file, json). Default: table
- `-types`: A comma-separated list of comment types to search for. The default is "TODO,FIXME".
- `-hidden`: Search hidden files and directories.
- `-validate-max`: Validate that the number of comments is less than or equal to the max.

## Install

To install the program, run the following command:

```bash
go install github.com/euforic/todos@latest
```

## Running the Program

To run the program for the current dir, execute the following command:

```bash
todos
```

You can also pass command-line arguments to the program. For example, to search for comments in the directory "myproject" and output the results in JSON format, run the following command:

```bash
todos -output json ./myproject
```

### Search a Specific Directory

To search a specific directory, use `todos [options] <dir>` the first arg is the directory path. For example, to search the directory `~/projects/myproject`, run the following command:

```bash
todos ~/projects/myproject
```

### Ignore Files and Directories

To ignore files and directories, use the `-ignore` flag followed by a comma-separated list of files and directories in gitignore format. For example, to ignore files with the extensions `.txt` and `.log`, and directories named `vendor` and `node_modules`, run the following command:

```bash
todos -ignore "*.txt,*.log,vendor/,node_modules/"
```

### Output in Format Style

To output the results in the chosen format (json, file, table), use the `-output` flag. For example, to output the results in json format, run the following command:

```bash
todos -output json
```

### Sort Results

To sort the results, use the `-sortby` flag followed by the field to sort by. The valid fields are `author`, `file`, `line`, `type`, and `text`. For example, to sort the results by comment type, run the following command:

```bash
todos -sortby type
```

### Search for Different Comment Types

To search for different types of comments, use the `-types` flag followed by a comma-separated list of comment types. For example, to search for comments with the types `TODO`, `FIXME`, and `NOTE`, run the following command:

```bash
todos -types=TODO,FIXME,NOTE
```

### Validate Max

To validate the maximum amount of comments does not exceed a certin value, use the `-validate-max` flag followed by the maximium number of comments. For example, to limit the maximum amount of comments to `20`, run the following command:

```bash
todos -validate-max 20
```
