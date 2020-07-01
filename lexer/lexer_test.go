package lexer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"git.bytecode.nl/foss/import-boundry-checker/lexer"
	"git.bytecode.nl/foss/import-boundry-checker/token"
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
		{`"test1","test2";`, []lexer.Result{
			{token.STRING, "test1"},
			{token.COMMA, ""},
			{token.STRING, "test2"},
			{token.SEMICOLON, ""}}},
		{`LANG "Typescript";`, []lexer.Result{
			{token.KEYWORD_LANG, ""},
			{token.STRING, "Typescript"},
			{token.SEMICOLON, ""}}},
		{`LANG "Typescript"; IMPORTRULE "some/module/path" CANNOTIMPORT "other/module/path/1", "other/module/path/2";`, []lexer.Result{
			{token.KEYWORD_LANG, ""},
			{token.STRING, "Typescript"},
			{token.SEMICOLON, ""},
			{token.KEYWORD_IMPORTRULE, ""},
			{token.STRING, "some/module/path"},
			{token.KEYWORD_CANNOTIMPORT, ""},
			{token.STRING, "other/module/path/1"},
			{token.COMMA, ""},
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
		{`LANG "Typescript"; IMPORTRULE "some/module/path" CANNOTIMPORT "other/module/path/1", "other/module/path/2";`, false},
		{`LANG "Typescript"; IMPORTRULE "some/module/path" INVALIDKEYWORD "other/module/path/1", "other/module/path/2";`, true},
		{`LANG "Typescript"; IMPORTRULE "some/module/path" CannotImport "other/module/path/1", "other/module/path/2";`, true},
		{`LANG "Typescript"; importrule "some/module/path" CANNOTIMPORT "other/module/path/1", "other/module/path/2";`, true},
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
