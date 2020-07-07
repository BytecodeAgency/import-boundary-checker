package runner

import (
	"fmt"

	"git.bytecode.nl/foss/import-boundry-checker/filefinder"
	"git.bytecode.nl/foss/import-boundry-checker/lexer"
	"git.bytecode.nl/foss/import-boundry-checker/parser"
)

type Runner struct {
	configFile string
}

func New(configFile string) Runner {
	return Runner{
		configFile: configFile,
	}
}

func (r Runner) Run() { // TODO: Extract CLI part to separate package
	// Lex and parse to get import ruleset
	l := doLex(r.configFile)
	p := doParse(l)
	f := doFilewalker(p.Lang)

	// DEV
	log("LANG: " + string(p.Lang))
	log("IMPORTBASE: " + p.ImportBase)
	log("\nIMPORTRULES:")
	for _, pp := range p.Rules { // TODO: Move this print function to parser method
		log(fmt.Sprintf("  ~> %s CANNOTIMPORT", pp.RuleFor)) // TODO: Create printer/logging package
		for _, noimport := range pp.CannotImport {
			log(fmt.Sprintf("     - %+v", noimport)) // TODO: Create printer/logging package
		}
	}
	log("\nFILES:")
	for _, ff := range f {
		log(fmt.Sprintf("  ~> %+v", ff)) // TODO: Create printer/logging package
	}

}

func doLex(input string) []lexer.Result {
	lex := lexer.New(input)
	lex.Exec()
	lexRes, lexErr := lex.Result()
	if len(lexErr) > 0 {
		errStr := prettyprintErrs(lexErr)
		log("Lexer returned errors:\n" + errStr)
		fail("Lexing was not succesful")
	}
	return lexRes
}

func doParse(input []lexer.Result) parser.Parser {
	p := parser.New(input)
	p.Parse()
	if len(p.Errors) > 0 {
		errStr := prettyprintErrs(p.Errors)
		log("Parser returned errors:\n" + errStr)
		fail("Parsing was not succesful")
	}
	return p
}

func doFilewalker(lang parser.Language) []string {
	if lang == parser.LangGo {
		files, err := filefinder.FindFilesWithExtInDir(".", []string{"go"}) // TODO: Make the directory editable via config
		if err != nil {
			fail(err.Error())
		}
		return files
	}
	fail(fmt.Sprintf("language %s is not supported", lang))
	return []string{} // Won't ever reach this code due to panic call one line up
}
