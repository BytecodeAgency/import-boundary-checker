package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"git.bytecode.nl/foss/import-boundry-checker/lexer"
	"git.bytecode.nl/foss/import-boundry-checker/parser"
)

var verbose = false

// TODO: Create cleaner entrypoint

func main() {
	// CLI flags
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output of the main function (does not enable debugging for sub packages)")
	configPath := flag.String("config", ".importrules", "Configuration path to be used when building import rule set")
	printHelp := flag.Bool("help", false, "Print CLI usage information")
	flag.Parse()

	if *printHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	logIfVerbose(fmt.Sprintf("verbose mode set to %b and using configPath %s", verbose, *configPath))

	c := getConfigString(*configPath)
	l := doLex(c)
	p, lang := doParse(l)

	// DEV
	log(string(lang))
	for _, pp := range p {
		log(fmt.Sprintf("  > %+v", pp))
	}
}

func getConfigString(path string) string {
	logIfVerbose("Creating fullPath")
	abs, err := filepath.Abs(path)
	if err != nil {
		fail("Could not create absolute config file path " + err.Error())
	}
	fullPath := abs
	logIfVerbose("fullPath = " + fullPath)
	config, err := ioutil.ReadFile(fullPath)
	logIfVerbose("config = " + string(config))
	if err != nil {
		fail(fmt.Sprintf("Could not read config file contents (%s), err %s", config, err))
	}
	return string(config)
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
