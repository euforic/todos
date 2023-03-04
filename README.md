# Todos

This is a command-line program written in Go that searches for `// TODO: xxxx` style comments in a directory and outputs them to the console or in JSON format.

## Usage

The program accepts the following command-line arguments:

- `-dir`: Specifies the directory to search for comments. Defaults to the current directory.
- `-ignore`: A comma-separated list of files and directories to ignore, in gitignore format.
- `-json`: Outputs the results in JSON format.
- `-sortby`: Sort results by field (`username`, `file`, `line`, `type`, or `text`)
- `-comment-types`: A comma-separated list of comment types to search for. The default is "TODO,FIXME".
- `-hidden`: Search hidden files and directories.

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
todos -dir=myproject -json
```

### Search a Specific Directory

To search a specific directory, use the `-dir` flag followed by the directory path. For example, to search the directory `~/projects/myproject`, run the following command:

```bash
todos -dir=~/projects/myproject
```

### Ignore Files and Directories

To ignore files and directories, use the `-ignore` flag followed by a comma-separated list of files and directories in gitignore format. For example, to ignore files with the extensions `.txt` and `.log`, and directories named `vendor` and `node_modules`, run the following command:

```bash
todos -ignore="*.txt,*.log,vendor/,node_modules/"
```

### Output in JSON Format

To output the results in JSON format, use the `-json` flag. For example, to output the results in JSON format, run the following command:

```bash
todos -json
```

### Sort Results

To sort the results, use the `-sortby` flag followed by the field to sort by. The valid fields are `username`, `file`, `line`, `type`, and `text`. For example, to sort the results by comment type, run the following command:

```bash
todos -sortby=type
```

### Search for Different Comment Types

To search for different types of comments, use the `-comment-types` flag followed by a comma-separated list of comment types. For example, to search for comments with the types `TODO`, `FIXME`, and `NOTE`, run the following command:

```bash
todos -comment-types=TODO,FIXME,NOTE
```
