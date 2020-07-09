package runner

import (
	"fmt"

	"git.bytecode.nl/foss/import-boundry-checker/filefinder"
	"git.bytecode.nl/foss/import-boundry-checker/langs/golistimports"
	"git.bytecode.nl/foss/import-boundry-checker/lexer"
	"git.bytecode.nl/foss/import-boundry-checker/parser"
	"git.bytecode.nl/foss/import-boundry-checker/rulechecker"
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
	imps := doGetImports(p.Lang, p.ImportBase)

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
	log("\nIMPORTMAP:")
	for file, imports := range imps {
		log(fmt.Sprintf("  ~> %s", file)) // TODO: Create printer/logging package
		for _, imp := range imports {
			log(fmt.Sprintf("    -> %s", imp)) // TODO: Create printer/logging package
		}
	}

	doRuleCheck(p.Rules, imps)
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

func doGetImports(lang parser.Language, importbase string) map[string][]string {
	if lang == parser.LangGo {
		files, err := filefinder.FindFilesWithExtInDir(".", []string{"go"}, []string{"_test.go"}) // TODO: Make the directory and extensions editable via config
		if err != nil {
			fail(err.Error())
		}
		importmap, err := golistimports.ExtractForFileList(files, importbase)
		if err != nil {
			fail(err.Error())
		}
		return importmap
	}
	fail(fmt.Sprintf("language %s is not supported", lang))
	return nil // Won't ever reach this code due to panic call one line up
}

func doRuleCheck(rules []parser.Rule, imports map[string][]string) {
	rc := rulechecker.New(rules, imports)
	if success := rc.Check(); success {
		succeed()
	}
	// There are failures, pretty print them
	for _, v := range rc.Violations {
		log(fmt.Sprintf("ERR: %s cannot import %s, but imports %s", v.Filename, v.CannotImport, v.ImportLine))
	}
	fail("Import violations were found.")
}
