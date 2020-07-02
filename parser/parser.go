package parser

import (
	"fmt"

	"git.bytecode.nl/foss/import-boundry-checker/keyword"
	"git.bytecode.nl/foss/import-boundry-checker/lexer"
	"git.bytecode.nl/foss/import-boundry-checker/token"
)

type Language string

const (
	LangUnset      = Language("")
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
	currentKeyword keyword.Keyword
	currentRule    Rule
}

func New(input []lexer.Result) Parser {
	return Parser{
		input: input,
		Lang:  LangUnset,
	}
}

func (p *Parser) logError(details string) {
	err := fmt.Errorf("%s with data currentKeyword %s, currentRule %s",
		details, p.currentKeyword, p.currentRule)
	p.Errors = append(p.Errors, err)
}

func (p *Parser) Parse() {
	for _, lexRes := range p.input {
		switch lexRes.Token {

		// Set currentKeyword
		case token.KEYWORD_LANG:
			p.currentKeyword = keyword.Lang
		case token.KEYWORD_IMPORTRULE:
			p.currentKeyword = keyword.ImportRule
		case token.KEYWORD_CANNOTIMPORT:
			p.currentKeyword = keyword.CannotImport

		// Handle expression end
		case token.SEMICOLON:
			p.endExpression()

		// Handle string input, save to location based on currentKeyword
		case token.STRING:
			p.saveStringToParserData(lexRes.Contents)
		}
	}
	if p.Lang == LangUnset {
		p.logError("language has not been set")
	}
	if len(p.Rules) == 0 {
		p.logError("no rules have been given")
	}
}

func (p *Parser) endExpression() {
	// Setting the language should not save the rule, so return
	if p.currentKeyword == keyword.Lang {
		return
	}

	// Validate that data for the rule is set
	if p.currentRule.RuleFor == "" || len(p.currentRule.CannotImport) == 0 {
		p.logError("not all required rule data is set")
	}
	p.Rules = append(p.Rules, p.currentRule)
	p.currentRule = Rule{} // Reset to default values
}

func (p *Parser) saveStringToParserData(ruleContents string) {
	switch p.currentKeyword {
	case keyword.ImportRule:
		if p.currentRule.RuleFor != "" {
			p.logError(fmt.Sprintf("RuleFor has already been set to %s", p.currentRule.RuleFor))
		}
		p.currentRule.RuleFor = ruleContents
	case keyword.CannotImport:
		p.currentRule.CannotImport = append(p.currentRule.CannotImport, ruleContents)

	// Set the language of the Parser instance
	case keyword.Lang:
		// Validate that the given language is correct
		switch Language(ruleContents) {
		case LangGo:
			p.Lang = LangGo
		case LangTypescript:
			p.Lang = LangTypescript
		default:
			p.logError(fmt.Sprintf("language '%s' is not supported", ruleContents))
		}
	}
}
