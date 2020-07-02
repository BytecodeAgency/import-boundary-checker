package parser

import "git.bytecode.nl/foss/import-boundry-checker/lexer"

type Language string

const (
	LangGo         = Language("Go")
	LangTypescript = Language("Typescript")
)

type Rule struct {
	RuleFor      string
	CannotImport []string
}

type Parser struct {
	// Input and output
	input  []lexer.Result
	Lang   Language
	Rules  []Rule
	Errors []error

	// Intermediate
	// TODO: Implement
}

func New(input []lexer.Result) Parser {
	return Parser{
		input: input,
	}
}

func (p *Parser) Parse() {
	// TODO: Implement
}
