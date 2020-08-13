package parser

import (
	"flag"
	"fmt"
	"strings"

	"github.com/BytecodeAgency/import-boundary-checker/keyword"
	"github.com/BytecodeAgency/import-boundary-checker/lexer"
	"github.com/BytecodeAgency/import-boundary-checker/token"
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
	RuleFor               string
	CannotImport          []string
	AllowImportExceptions []string
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
	err := fmt.Errorf("%s", details)
	errDebug := fmt.Errorf("Debugdata: currentKeyword %s, currentRule %s", p.currentKeyword, p.currentRule)
	p.Errors = append(p.Errors, errDebug, err)
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
		case token.KEYWORD_ALLOW:
			p.currentKeyword = keyword.Allow
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

	// Check for language and if there are any rules
	if p.Lang == LangUnset {
		p.logError("language has not been set")
	}
	if len(p.Rules) == 0 {
		p.logError("No rules have been found while parsing (maybe you forgot to use a semicolon?)")
	}

	// Replace the `[IMPORTBASE]` variable in the strings
	importbaseVarString := fmt.Sprintf("[%s]", keyword.ImportBase) // Make it variable so we don't hardcode, all code depends on the keyword package
	for i, rule := range p.Rules {
		rule.RuleFor = strings.ReplaceAll(rule.RuleFor, importbaseVarString, p.ImportBase)
		for j, ci := range rule.CannotImport {
			rule.CannotImport[j] = strings.ReplaceAll(ci, importbaseVarString, p.ImportBase)
		}
		for j, ci := range rule.AllowImportExceptions {
			rule.AllowImportExceptions[j] = strings.ReplaceAll(ci, importbaseVarString, p.ImportBase)
		}
		p.Rules[i] = rule
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
	if p.currentRule.RuleFor == "" || (len(p.currentRule.CannotImport) == 0 && len(p.currentRule.AllowImportExceptions) == 0) {
		p.logError("not all required rule data is set")
	}
	if p.currentRule.AllowImportExceptions == nil {
		p.currentRule.AllowImportExceptions = []string{}
	}
	p.Rules = append(p.Rules, p.currentRule)
	p.currentRule = Rule{} // Reset to default values
}

func (p *Parser) saveStringToParserData(ruleContents string) {
	switch p.currentKeyword {
	case keyword.ImportRule:
		if p.currentRule.RuleFor != "" {
			p.logError(fmt.Sprintf("RuleFor has already been set to %s\n  -> Did you forget to add a semicolon at the end of your previous rule?", p.currentRule.RuleFor)) // TODO: Log this cleaner
		}
		p.currentRule.RuleFor = ruleContents
	case keyword.CannotImport:
		p.currentRule.CannotImport = append(p.currentRule.CannotImport, ruleContents)
	case keyword.Allow:
		p.currentRule.AllowImportExceptions = append(p.currentRule.AllowImportExceptions, ruleContents)
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
