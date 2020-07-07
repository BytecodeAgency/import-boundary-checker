package runner

import (
	"flag"
	"fmt"
	"os"

	"git.bytecode.nl/foss/import-boundry-checker/lexer"
	"git.bytecode.nl/foss/import-boundry-checker/parser"
)

var verbose = false

func Run(config string) { // TODO: Extract CLI part to separate package
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output of the main function (does not enable debugging for sub packages)")

	printHelp := flag.Bool("help", false, "Print CLI usage information")

	if *printHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	logIfVerbose(fmt.Sprintf("verbose mode set to %t", verbose))

	// Lex and parse to get import ruleset
	l := doLex(config)
	p, lang := doParse(l)

	// DEV
	log(string(lang))
	for _, pp := range p {
		log(fmt.Sprintf("  > %+v", pp)) // TODO: Create printer/logging package
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
