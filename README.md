# demog

 A command-line utility that retrieves demographic data for a specified set of U.S. states from a public API and outputs that data in the requested format.

 If the specified output-format parameter is "CSV," the data is output in CSV format, sorted alphabetically by state.

 If the specified output-parameter format is "averages," a single integer is output representing the weighted average of "income below poverty" across all the specified input states.

## Assumptions

1. States will not be input by their abbreviations
2. Whitespace between commas counts as part of the field

## Building

1. Make sure you have [dep](https://github.com/golang/dep) installed.
2. Ruj `dep ensure` to get the requirements
3. Build the application via `go build`
4. Run the application `./demog`

## Usage

```
NAME:
   demog - A CLI for retrieving demographic data for sets of US states

USAGE:
   demog [global options] command [command options] [arguments...]

VERSION:
   0.0.0

AUTHOR:
   Philip Gardner

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --format value, -f value  output format for demographic data [csv,averages]
   --help, -h                show help
   --version, -v             print the version
   ```