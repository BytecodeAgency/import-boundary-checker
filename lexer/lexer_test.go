package lexer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"git.bytecode.nl/foss/import-boundry-checker/lexer"
	"git.bytecode.nl/foss/import-boundry-checker/token"
)

func TestLexer_SingleLine(t *testing.T) {
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
	}
	for _, test := range tests {
		l := lexer.New(test.input)
		l.Exec()
		res, errs := l.Result()
		assert.Empty(t, errs)
		assert.Equal(t, test.expected, res)
	}
}
