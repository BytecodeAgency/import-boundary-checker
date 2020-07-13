package runner

import (
	"fmt"

	"github.com/BytecodeAgency/import-boundry-checker/filefinder"
	"github.com/BytecodeAgency/import-boundry-checker/langs/golistimports"
	"github.com/BytecodeAgency/import-boundry-checker/lexer"
	"github.com/BytecodeAgency/import-boundry-checker/logging"
	"github.com/BytecodeAgency/import-boundry-checker/parser"
	"github.com/BytecodeAgency/import-boundry-checker/rulechecker"
)

type Runner struct {
	configFile string
	logger     *logging.Logger
	failed     bool
}

func New(configFile string, logger *logging.Logger) Runner {
	return Runner{
		configFile: configFile,
		logger:     logger,
		failed:     false,
	}
}

func (r *Runner) Run() (failureOccurred bool) {
	// Lex and parse to get import ruleset
	l := r.doLex(r.configFile)
	if r.failed {
		return r.failed
	}
	p := r.doParse(l)
	r.logger.AddDebug("Detected language: " + string(p.Lang))
	r.logger.AddDebug("Detected importbase: " + p.ImportBase)
	if r.failed {
		return r.failed
	}
	imps := r.doGetImports(p.Lang, p.ImportBase)

	// Do rule check
	r.doRuleCheck(p.Rules, imps)
	if r.failed {
		return r.failed
	}

	// Return if there are failures
	return r.failed
}

func (r *Runner) setFailed() {
	r.failed = true
}

func (r *Runner) doLex(input string) []lexer.Result {
	lex := lexer.New(input)
	lex.Exec()
	lexRes, lexErr := lex.Result()
	if len(lexErr) > 0 {
		r.setFailed()
		r.logger.FailWithErrors("Lexing was not successful", lexErr)
	}
	return lexRes
}

func (r *Runner) doParse(input []lexer.Result) parser.Parser {
	p := parser.New(input)
	p.Parse()
	if len(p.Errors) > 0 {
		r.setFailed()
		r.logger.FailWithErrors("Lexing was not successful", p.Errors)
	}
	r.logger.SetRules(p.Rules)
	return p
}

func (r *Runner) doGetImports(lang parser.Language, importbase string) map[string][]string {
	if lang == parser.LangGo {
		files, err := filefinder.FindFilesWithExtInDir(".", []string{"go"}, []string{"_test.go"}) // TODO: Make the directory and extensions editable via config
		if err != nil {
			r.setFailed()
			r.logger.FailWithError("Could not run file finder", err)
		}
		importmap, err := golistimports.ExtractForFileList(files, importbase)
		r.logger.SetImportChart(importmap)
		if err != nil {
			r.setFailed()
			r.logger.FailWithError("Could not run import map extractor", err)
		}
		return importmap
	}
	r.logger.FailWithMessage(fmt.Sprintf("language %s is not supported", lang))
	return nil // Won't ever reach this code due to panic call one line up
}

func (r *Runner) doRuleCheck(rules []parser.Rule, imports map[string][]string) {
	rc := rulechecker.New(rules, imports)
	if success := rc.Check(); !success { // Means there are violations
		r.setFailed()
		r.logger.SetImportViolations(rc.Violations)
		r.logger.FailWithMessage("Import violations were found")
		return
	}
	r.logger.Success()
}
