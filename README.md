# Import Boundry Checker DSL

_By Bytecode Digital Agency B.V._

## CLI options

_All CLI options are optional_

### Running the application

* `-config [filepath]`: set the configuration path (default: TODO)
* `-verbose`: set verbose output of main (does not enable debugging mode)

### Debugging

* `-debug_lexer`: enables debug output for the lexer
* `-debug_parser`: enables debug output for the parser

## Execution path

* `main` reads the config file and parses CLI options
* With the correct options `runner` is called, which executes the following steps
* The config file is passed through `lexer`, then through `parser` (with importrules as output)
* `filefinder` finds all files related to the language detected by `parser`
* If language is Go, `langs/golistimports` takes the file list and outputs the a map of imports per module
* TODO: `rulechecker` takes the imports map and ruleset from `parser` and checks whether all imports are valid
* TODO: `printer` takes the program results and prints them
