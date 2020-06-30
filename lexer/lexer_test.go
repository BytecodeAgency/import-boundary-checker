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
		//{`LANG`, []lexer.Result{{token.KEYWORD_LANG, ""}}},
		//{`LANG "Typescript"`, []lexer.Result{{token.KEYWORD_LANG, "Typescript"}}},
	}
	for _, test := range tests {
		l := lexer.New(test.input)
		l.Exec()
		res := l.Result()
		assert.Equal(t, test.expected, res)
	}
}
