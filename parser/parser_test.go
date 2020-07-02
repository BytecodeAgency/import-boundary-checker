package parser_test

import (
	"testing"

	"git.bytecode.nl/foss/import-boundry-checker/lexer"
	"git.bytecode.nl/foss/import-boundry-checker/parser"
	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		input         string
		expectedLang  parser.Language
		expectedRules []parser.Rule
	}{
		{`LANG "Go";

IMPORTRULE "git.bytecode.nl/single-projects/youngpwr/platform-backend/typings/entities"
CANNOTIMPORT "git.bytecode.nl/single-projects/youngpwr/platform-backend";

IMPORTRULE
  	"git.bytecode.nl/single-projects/youngpwr/platform-backend/domain"
CANNOTIMPORT
	"git.bytecode.nl/single-projects/youngpwr/platform-backend/infrastructure"
    "git.bytecode.nl/single-projects/youngpwr/platform-backend/data";`, "Go",
			[]parser.Rule{
				{"git.bytecode.nl/single-projects/youngpwr/platform-backend/typings/entities",
					[]string{"git.bytecode.nl/single-projects/youngpwr/platform-backend"}},
				{"git.bytecode.nl/single-projects/youngpwr/platform-backend/domain",
					[]string{
						"git.bytecode.nl/single-projects/youngpwr/platform-backend/infrastructure",
						"git.bytecode.nl/single-projects/youngpwr/platform-backend/data"}},
			}},
	}

	for _, test := range tests {
		// Lexer
		l := lexer.New(test.input)
		l.Exec()
		res, errs := l.Result()
		assert.Empty(t, errs)

		// Parser
		p := parser.New(res)
		p.Parse()
		assert.Empty(t, p.Errors)
		assert.Equal(t, test.expectedLang, p.Lang)
		assert.Equal(t, test.expectedRules, p.Rules)
	}
}
