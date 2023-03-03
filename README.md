# Todos

This is a command-line program written in Go that searches for `// TODO: xxxx` style comments in a directory and outputs them to the console or in JSON format.

## Usage

The program accepts the following command-line arguments:

- `-dir`: Specifies the directory to search for comments. Defaults to the current directory.
- `-ignore-files`: A comma-separated list of file extensions to ignore.
- `-ignore-dirs`: A comma-separated list of directory names to ignore.
- `-json`: Outputs the results in JSON format.
- `-comment-types`: A comma-separated list of comment types to search for. The default is "TODO,FIXME".

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

To search a specific directory, use the -dir flag followed by the directory path. For example, to search the directory ~/projects/myproject, run the following command:

```bash
todos -dir=~/projects/myproject
```

### Ignore Files

To ignore files with certain file extensions, use the -ignore-files flag followed by a comma-separated list of file extensions. For example, to ignore files with the extensions .txt and .log, run the following command:

```bash
todos -ignore-files=txt,log
```

### Ignore Directories

To ignore directories with certain names, use the -ignore-dirs flag followed by a comma-separated list of directory names. For example, to ignore directories named node_modules and vendor, run the following command:

```bash
todos -ignore-dirs=node_modules,vendor
```

### Output in JSON Format

To output the results in JSON format, use the -json flag. For example, to output the results in JSON format, run the following command:

```bash
todos -json
```

### Search for Different Comment Types

To search for different types of comments, use the -comment-types flag followed by a comma-separated list of comment types. For example, to search for comments with the types TODO, FIXME, and NOTE, run the following command:

```bash
todos -comment-types=TODO,FIXME,NOTE
```
