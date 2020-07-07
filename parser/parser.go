package parser

import (
	"flag"
	"fmt"

	"git.bytecode.nl/foss/import-boundry-checker/keyword"
	"git.bytecode.nl/foss/import-boundry-checker/lexer"
	"git.bytecode.nl/foss/import-boundry-checker/token"
)

type Language string

var DEBUG = false

func init() {
	flag.BoolVar(&DEBUG, "debug_parser", false, "Enable debugging for parser")
}

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
	input      []lexer.Result
	Lang       Language
	ImportBase string
	Rules      []Rule
	Errors     []error

	// Intermediate
	currentKeyword keyword.Keyword
	currentRule    Rule
}

func New(input []lexer.Result) Parser {
	return Parser{
		input:      input,
		Lang:       LangUnset,
		ImportBase: "",
	}
}

func (p *Parser) logError(details string) {
	err := fmt.Errorf("%s with data currentKeyword %s, currentRule %s",
		details, p.currentKeyword, p.currentRule)
	p.Errors = append(p.Errors, err)
}

func (p *Parser) Parse() {
	if DEBUG {
		fmt.Printf("starting the parsing of %+v", p.input)
	}

	for _, lexRes := range p.input {
		switch lexRes.Token {

		// Set currentKeyword
		case token.KEYWORD_LANG:
			p.currentKeyword = keyword.Lang
		case token.KEYWORD_IMPORTRULE:
			p.currentKeyword = keyword.ImportRule
		case token.KEYWORD_CANNOTIMPORT:
			p.currentKeyword = keyword.CannotImport
		case token.KEYWORD_IMPORTBASE:
			p.currentKeyword = keyword.ImportBase

		// Handle expression end
		case token.SEMICOLON:
			p.endExpression()

		// Handle string input, save to location based on currentKeyword
		case token.STRING:
			p.saveStringToParserData(lexRes.Contents)

		// Should not reach this code
		default:
			p.logError(fmt.Sprintf("received unsupported token %s with contents %s", lexRes.Token, lexRes.Contents))
		}

		if DEBUG {
			fmt.Printf("after parsing %s with contents %s, got currentKeyword %s, currentRule %s with lang %s and rules %+v and errors %s",
				lexRes.Token, lexRes.Contents, p.currentKeyword, p.currentRule, p.Lang, p.Rules, p.Errors)
		}
	}
	if p.Lang == LangUnset {
		p.logError("language has not been set")
	}
	if len(p.Rules) == 0 {
		p.logError("no rules have been given")
	}

	if DEBUG {
		fmt.Printf("after parsing, got currentKeyword %s, currentRule %s with lang %s and rules %+v and errors %s",
			p.currentKeyword, p.currentRule, p.Lang, p.Rules, p.Errors)
	}
}

func (p *Parser) endExpression() {
	// Setting the language or importbase should not save the rule, so return
	if p.currentKeyword == keyword.Lang || p.currentKeyword == keyword.ImportBase {
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
	case keyword.ImportBase:
		p.ImportBase = ruleContents

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
