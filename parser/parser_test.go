package parser_test

import (
	"testing"

	"git.bytecode.nl/foss/import-boundry-checker/lexer"
	"git.bytecode.nl/foss/import-boundry-checker/parser"
	"github.com/stretchr/testify/assert"
)

func TestParser_Parse_Correct(t *testing.T) {
	tests := []struct {
		input              string
		expectedLang       parser.Language
		expectedImportBase string
		expectedRules      []parser.Rule
	}{
		{`LANG "Go";

IMPORTRULE "git.bytecode.nl/single-projects/youngpwr/platform-backend/typings/entities"
CANNOTIMPORT "git.bytecode.nl/single-projects/youngpwr/platform-backend";

IMPORTRULE
  	"git.bytecode.nl/single-projects/youngpwr/platform-backend/domain"
CANNOTIMPORT
	"git.bytecode.nl/single-projects/youngpwr/platform-backend/infrastructure"
    "git.bytecode.nl/single-projects/youngpwr/platform-backend/data";`,
			"Go",
			"",
			[]parser.Rule{
				{"git.bytecode.nl/single-projects/youngpwr/platform-backend/typings/entities",
					[]string{"git.bytecode.nl/single-projects/youngpwr/platform-backend"}},
				{"git.bytecode.nl/single-projects/youngpwr/platform-backend/domain",
					[]string{
						"git.bytecode.nl/single-projects/youngpwr/platform-backend/infrastructure",
						"git.bytecode.nl/single-projects/youngpwr/platform-backend/data"}},
			}},
		{`LANG "Go";
IMPORTBASE "git.bytecode.nl/single-projects/youngpwr/platform-backend";

IMPORTRULE "git.bytecode.nl/single-projects/youngpwr/platform-backend/typings/entities"
CANNOTIMPORT "git.bytecode.nl/single-projects/youngpwr/platform-backend";

IMPORTRULE
  	"git.bytecode.nl/single-projects/youngpwr/platform-backend/domain"
CANNOTIMPORT
	"git.bytecode.nl/single-projects/youngpwr/platform-backend/infrastructure"
    "git.bytecode.nl/single-projects/youngpwr/platform-backend/data";`,
			"Go",
			"git.bytecode.nl/single-projects/youngpwr/platform-backend",
			[]parser.Rule{
				{"git.bytecode.nl/single-projects/youngpwr/platform-backend/typings/entities",
					[]string{"git.bytecode.nl/single-projects/youngpwr/platform-backend"}},
				{"git.bytecode.nl/single-projects/youngpwr/platform-backend/domain",
					[]string{
						"git.bytecode.nl/single-projects/youngpwr/platform-backend/infrastructure",
						"git.bytecode.nl/single-projects/youngpwr/platform-backend/data"}},
			}},
		{`LANG "Go";
IMPORTBASE "git.bytecode.nl/single-projects/youngpwr/platform-backend";

IMPORTRULE "[IMPORTBASE]/typings/entities"
CANNOTIMPORT "[IMPORTBASE]";

IMPORTRULE "[IMPORTBASE]/domain"
CANNOTIMPORT
	"[IMPORTBASE]/infrastructure"
    "[IMPORTBASE]/data";`,
			"Go",
			"git.bytecode.nl/single-projects/youngpwr/platform-backend",
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
		assert.Equal(t, test.expectedImportBase, p.ImportBase)
	}
}

func TestParser_Parse_Incorrect(t *testing.T) {
	incorrectInputs := []string{
		// Invalid language
		`LANG "COBOL";
IMPORTRULE "some/module"
CANNOTIMPORT "some/other/module";`,

		// Multiple importrules
		`LANG "Go";
IMPORTRULE "some/module1" "some/module2"
CANNOTIMPORT "some/other/module";`,

		// Not finishing the importrule
		`LANG "Go";
IMPORTRULE "some/module1" "some/module2";`,

		// Not setting the importrule, only the cannotimports
		`LANG "Go";
IMPORTRULE
CANNOTIMPORT "some/module2";`,

		// Not setting the language
		`IMPORTRULE "some/module"
CANNOTIMPORT "some/module2";`,

		// Only setting the language, and no importrules
		`LANG "Go";`,
	}

	for _, input := range incorrectInputs {
		// Lexer
		l := lexer.New(input)
		l.Exec()
		res, errs := l.Result()
		assert.Empty(t, errs)

		// Parser
		p := parser.New(res)
		p.Parse()
		assert.NotEmpty(t, p.Errors)
	}
}
