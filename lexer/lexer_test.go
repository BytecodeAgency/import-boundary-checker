package lexer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/BytecodeAgency/import-boundry-checker/lexer"
	"github.com/BytecodeAgency/import-boundry-checker/token"
)

func TestLexer_SingleLine_Correct(t *testing.T) {
	tests := []struct {
		input    string
		expected []lexer.Result
	}{
		{`"test"`, []lexer.Result{{token.STRING, "test"}}},
		{`"test";`, []lexer.Result{
			{token.STRING, "test"},
			{token.SEMICOLON, ""}}},
		{`"test"         ;`, []lexer.Result{
			{token.STRING, "test"},
			{token.SEMICOLON, ""}}},
		{`"test1""test2";`, []lexer.Result{
			{token.STRING, "test1"},
			{token.STRING, "test2"},
			{token.SEMICOLON, ""}}},
		{`LANG "Typescript";`, []lexer.Result{
			{token.KEYWORD_LANG, ""},
			{token.STRING, "Typescript"},
			{token.SEMICOLON, ""}}},
		{`LANG "Typescript"; IMPORTRULE "some/module/path" CANNOTIMPORT "other/module/path/1" "other/module/path/2";`, []lexer.Result{
			{token.KEYWORD_LANG, ""},
			{token.STRING, "Typescript"},
			{token.SEMICOLON, ""},
			{token.KEYWORD_IMPORTRULE, ""},
			{token.STRING, "some/module/path"},
			{token.KEYWORD_CANNOTIMPORT, ""},
			{token.STRING, "other/module/path/1"},
			{token.STRING, "other/module/path/2"},
			{token.SEMICOLON, ""}}},
	}
	for _, test := range tests {
		l := lexer.New(test.input)
		l.Exec()
		res, errs := l.Result()
		assert.Empty(t, errs)
		assert.Equal(t, test.expected, res)
	}
}

func TestLexer_SingleLine_Failure(t *testing.T) {
	tests := []struct {
		input     string
		shouldErr bool
	}{
		{`LANG "Typescript";`, false},
		{`LANG "Typescript"; IMPORTRULE "some/module/path" CANNOTIMPORT "other/module/path/1" "other/module/path/2";`, false},
		{`LANG "Typescript"; IMPORTRULE "some/module/path" INVALIDKEYWORD "other/module/path/1" "other/module/path/2";`, true},
		{`LANG "Typescript"; IMPORTRULE "some/module/path" CannotImport "other/module/path/1" "other/module/path/2";`, true},
		{`LANG "Typescript"; importrule "some/module/path" CANNOTIMPORT "other/module/path/1" "other/module/path/2";`, true},
		{`"test1", "test2";`, true},
		{`"test1" "test2";`, false},
		{`LANG "Typescript"; IMPORTRULE "some/module/path" CANNOTIMPORT "other/module/path/1", "other/module/path/2";`, true},
		{`LANG "Typescript"; IMPORTRULE "some/module/path" CANNOTIMPORT "other/module/path/1" "other/module/path/2";`, false},
	}
	for _, test := range tests {
		l := lexer.New(test.input)
		l.Exec()
		_, errs := l.Result()
		if test.shouldErr {
			assert.NotEmpty(t, errs)
		} else {
			assert.Empty(t, errs)
		}
	}
}

func TestLexer_MultiLine_Correct(t *testing.T) {
	tests := []struct {
		input    string
		expected []lexer.Result
	}{
		{`"test1"
"test2";`, []lexer.Result{
			{token.STRING, "test1"},
			{token.STRING, "test2"},
			{token.SEMICOLON, ""}}},
		{`
LANG
"Typescript"
;`, []lexer.Result{
			{token.KEYWORD_LANG, ""},
			{token.STRING, "Typescript"},
			{token.SEMICOLON, ""}}},
		{`LANG "Typescript";
IMPORTRULE "some/module/path"
CANNOTIMPORT "other/module/path/1" "other/module/path/2";`, []lexer.Result{
			{token.KEYWORD_LANG, ""},
			{token.STRING, "Typescript"},
			{token.SEMICOLON, ""},
			{token.KEYWORD_IMPORTRULE, ""},
			{token.STRING, "some/module/path"},
			{token.KEYWORD_CANNOTIMPORT, ""},
			{token.STRING, "other/module/path/1"},
			{token.STRING, "other/module/path/2"},
			{token.SEMICOLON, ""}}},
		{`LANG "Go";
IMPORTBASE "git.bytecode.nl/single-projects/youngpwr/platform-backend";

IMPORTRULE "git.bytecode.nl/single-projects/youngpwr/platform-backend/typings/entities"
CANNOTIMPORT "git.bytecode.nl/single-projects/youngpwr/platform-backend";

IMPORTRULE
  	"git.bytecode.nl/single-projects/youngpwr/platform-backend/domain"
CANNOTIMPORT
	"git.bytecode.nl/single-projects/youngpwr/platform-backend/infrastructure"
    "git.bytecode.nl/single-projects/youngpwr/platform-backend/data";`, []lexer.Result{
			{token.KEYWORD_LANG, ""},
			{token.STRING, "Go"},
			{token.SEMICOLON, ""},
			{token.KEYWORD_IMPORTBASE, ""},
			{token.STRING, "git.bytecode.nl/single-projects/youngpwr/platform-backend"},
			{token.SEMICOLON, ""},
			{token.KEYWORD_IMPORTRULE, ""},
			{token.STRING, "git.bytecode.nl/single-projects/youngpwr/platform-backend/typings/entities"},
			{token.KEYWORD_CANNOTIMPORT, ""},
			{token.STRING, "git.bytecode.nl/single-projects/youngpwr/platform-backend"},
			{token.SEMICOLON, ""},
			{token.KEYWORD_IMPORTRULE, ""},
			{token.STRING, "git.bytecode.nl/single-projects/youngpwr/platform-backend/domain"},
			{token.KEYWORD_CANNOTIMPORT, ""},
			{token.STRING, "git.bytecode.nl/single-projects/youngpwr/platform-backend/infrastructure"},
			{token.STRING, "git.bytecode.nl/single-projects/youngpwr/platform-backend/data"},
			{token.SEMICOLON, ""}}},
	}
	for _, test := range tests {
		l := lexer.New(test.input)
		l.Exec()
		res, errs := l.Result()
		assert.Empty(t, errs)
		assert.Equal(t, test.expected, res)
	}
}
