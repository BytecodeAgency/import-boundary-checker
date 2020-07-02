package main

import (
	"flag"
	"fmt"

	"git.bytecode.nl/foss/import-boundry-checker/lexer"
	"git.bytecode.nl/foss/import-boundry-checker/parser"
)

var verbose = false

// TODO: Create cleaner entrypoint

func main() {
	// CLI flags
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output of the main function (does not enable debugging for sub packages)")
	configPath := flag.String("config", "TODO", "Configuration path to be used when building import rule set")
	flag.Parse()

	logIfVerbose(fmt.Sprintf("verbose mode set to %b and using configPath %s", verbose, *configPath))

	tester := `LANG "Go";

IMPORTRULE "git.bytecode.nl/single-projects/youngpwr/platform-backend/typings/entities"
CANNOTIMPORT "git.bytecode.nl/single-projects/youngpwr/platform-backend";

IMPORTRULE
  	"git.bytecode.nl/single-projects/youngpwr/platform-backend/domain"
CANNOTIMPORT
	"git.bytecode.nl/single-projects/youngpwr/platform-backend/infrastructure"
    "git.bytecode.nl/single-projects/youngpwr/platform-backend/data";`

	l := doLex(tester)
	p, lang := doParse(l)

	// DEV
	log(string(lang))
	for _, pp := range p {
		log(fmt.Sprintf("  > %+v", pp))
	}
}



func doLex(input string) []lexer.Result {
	logIfVerbose("Building lexer")
	lex := lexer.New(input)
	logIfVerbose("Executing lexer")
	lex.Exec()
	logIfVerbose("Fetching lexer results")
	lexRes, lexErr := lex.Result()
	logIfVerbose("Checking lexer errors")
	if len(lexErr) > 0 {
		errStr := prettyprintErrs(lexErr)
		log("Lexer returned errors:\n" + errStr)
		fail("Lexing was not succesful")
	}
	logIfVerbose("Returning lexer results")
	return lexRes
}

func doParse(input []lexer.Result) ([]parser.Rule, parser.Language) {
	logIfVerbose("Building parser")
	p := parser.New(input)
	logIfVerbose("Executing parser")
	p.Parse()
	logIfVerbose("Checking parser errors")
	if len(p.Errors) > 0 {
		errStr := prettyprintErrs(p.Errors)
		log("Parser returned errors:\n" + errStr)
		fail("Parsing was not succesful")
	}
	logIfVerbose(fmt.Sprintf("Detected language %s", p.Lang))
	return p.Rules, p.Lang
}
