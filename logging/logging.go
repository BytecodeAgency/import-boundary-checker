package logging

import (
	"fmt"
	"strings"

	"github.com/BytecodeAgency/import-boundary-checker/parser"
	"github.com/BytecodeAgency/import-boundary-checker/rulechecker"
)

type Entry struct {
	Level    LogLevel
	Contents string
}

type Logger struct {
	Logs             strings.Builder
	Entries          []Entry
	Rules            []parser.Rule
	ImportChart      map[string][]string
	ImportViolations []rulechecker.Violation
	Verbose          bool
	Completed        bool
}

func New(verbose bool) *Logger {
	logger := Logger{
		Verbose:   verbose,
		Logs:      strings.Builder{},
		Completed: false,
	}
	return &logger
}

/*
 * Adders and setters for information
 */

func (l *Logger) log(contents string) {
	_, err := l.Logs.Write([]byte(contents))
	if err != nil {
		panic(err) // Could not write to logs, TODO: Clean up
	}
}

func (l *Logger) AddAlways(contents string) {
	l.Entries = append(l.Entries, Entry{
		Level:    LogLevelError,
		Contents: contents,
	})
}

func (l *Logger) AddDebug(contents string) {
	l.Entries = append(l.Entries, Entry{
		Level:    LogLevelDebug,
		Contents: contents,
	})
}

func (l *Logger) SetRules(rules []parser.Rule) {
	l.Rules = rules
}

func (l *Logger) SetImportChart(imports map[string][]string) {
	l.ImportChart = imports
}

func (l *Logger) SetImportViolations(violations []rulechecker.Violation) {
	l.ImportViolations = violations
}

/*
 * PRINTING SPECIFIC LOGGER FIELDS
 */

func (l *Logger) printRules() {
	if l.Rules == nil {
		l.log(logInfo("Rules is nil"))
		return
	}
	if len(l.Rules) == 0 {
		l.log(logInfo("Rules length is 0"))
		return
	}
	l.log(logInfo("Printing parser rules:"))
	for _, rule := range l.Rules {
		l.log(logCont(fmt.Sprintf("%s cannot import:", rule.RuleFor)))
		for _, cannotImport := range rule.CannotImport {
			l.log(logCont(fmt.Sprintf("  - %s", cannotImport)))
		}
	}
}

func (l *Logger) printImportChart() {
	if l.ImportChart == nil {
		l.log(logInfo("Import chart is nil"))
		return
	}
	l.log(logInfo("Printing import chart:"))
	for file, imports := range l.ImportChart {
		logCont(fmt.Sprintf("%s imports:", file))
		for _, imp := range imports {
			logCont(fmt.Sprintf("  - %s", imp))
		}
	}
}

func (l *Logger) printImportViolations() {
	if l.ImportViolations == nil || len(l.ImportViolations) == 0 {
		l.log(logInfo("No import violations found"))
		return
	}
	for _, v := range l.ImportViolations {
		l.log(logError(fmt.Sprintf("File '%s' violated import rule, cannot import", v.Filename)))
		l.log(logCont(fmt.Sprintf("'%s' but imported", v.CannotImport)))
		l.log(logCont(fmt.Sprintf("'%s'", v.ImportLine)))
	}
}

/*
 * PRINTING ENTRIES
 */

func (l *Logger) printAlwaysOutput() {
	for _, entry := range l.Entries {
		if entry.Level == LogLevelError {
			l.log(logError(entry.Contents))
		}
	}
}

func (l *Logger) printVerboseOutput() {
	for _, entry := range l.Entries {
		switch entry.Level {
		case LogLevelError:
			l.log(logError(entry.Contents))
		case LogLevelWarn:
			l.log(logWarn(entry.Contents))
		case LogLevelDebug:
			l.log(logInfo(entry.Contents))
		case LogLevelTrace:
			l.log(logTrace(entry.Contents))
		default:
			l.log(logTrace(entry.Contents))
		}
	}
	l.printRules()
	l.printImportChart()
}

func (l *Logger) Print(completed bool) {
	// Print output from logging entries (parser, lexer)
	if l.Verbose {
		l.printVerboseOutput()
	} else {
		l.printAlwaysOutput()
	}

	// Print import violations
	if completed {
		l.printImportViolations()
	} else {
		l.log(logWarn("Import violations were not checked due to error (run with `-verbose` for debug info)"))
	}
}

/*
 * EXITING METHODS - DO NOT CALL os.Exit HERE
 */

func (l *Logger) FailWithMessage(message string) {
	l.Print(false)
	l.log(logError("Error occurred: " + message))
}

func (l *Logger) FailWithMessageCompleted(message string) {
	l.Print(true)
	l.log(logError("Error occurred: " + message))
}

func (l *Logger) FailWithError(message string, err error) {
	l.Print(false)
	l.log(logError(fmt.Sprintf("Error occurred: %s (%s)", message, err.Error())))
}

func (l *Logger) FailWithErrors(message string, errs []error) {
	l.Print(false)
	l.log(logError("Error occurred: " + message))
	for _, err := range errs {
		l.log(logError(err.Error()))
	}
}

func (l *Logger) Success() {
	l.Print(true)
	l.log(logInfo("No errors seem to have occurred!"))
	l.log(logInfo("Exiting gracefully"))
}
