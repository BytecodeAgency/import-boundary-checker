package runner

import (
	"fmt"

	"git.bytecode.nl/foss/import-boundry-checker/lexer"
	"git.bytecode.nl/foss/import-boundry-checker/parser"
)

func Run(config string) { // TODO: Extract CLI part to separate package
	// Lex and parse to get import ruleset
	l := doLex(config)
	p, lang := doParse(l)

	// DEV
	log(string(lang))
	for _, pp := range p {
		log(fmt.Sprintf("  ~> %+v", pp)) // TODO: Create printer/logging package
	}

}

func doLex(input string) []lexer.Result {
	lex := lexer.New(input)
	lex.Exec()
	lexRes, lexErr := lex.Result()
	if len(lexErr) > 0 {
		errStr := prettyprintErrs(lexErr)
		log("Lexer returned errors:\n" + errStr)
		fail("Lexing was not succesful")
	}
	return lexRes
}

func doParse(input []lexer.Result) ([]parser.Rule, parser.Language) {
	p := parser.New(input)
	p.Parse()
	if len(p.Errors) > 0 {
		errStr := prettyprintErrs(p.Errors)
		log("Parser returned errors:\n" + errStr)
		fail("Parsing was not succesful")
	}
	return p.Rules, p.Lang
}
