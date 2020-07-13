package logging

import (
	"fmt"

	"github.com/BytecodeAgency/import-boundary-checker/parser"
	"github.com/BytecodeAgency/import-boundary-checker/rulechecker"
)

type Entry struct {
	Level    LogLevel
	Contents string
}

type Logger struct {
	Entries          []Entry
	Rules            []parser.Rule
	ImportChart      map[string][]string
	ImportViolations []rulechecker.Violation
	Verbose          bool
}

func New(verbose bool) *Logger {
	logger := Logger{
		Verbose: verbose,
	}
	return &logger
}

/*
 * Adders and settings for information
 */

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
		logInfo("Rules is nil")
		return
	}
	if len(l.Rules) == 0 {
		logInfo("Rules length is 0")
		return
	}
	logInfo("Printing parser rules:")
	for _, rule := range l.Rules {
		logCont(fmt.Sprintf("%s cannot import:", rule.RuleFor))
		for _, cannotImport := range rule.CannotImport {
			logCont(fmt.Sprintf("  - %s", cannotImport))
		}
	}
}

func (l *Logger) printImportChart() {
	if l.ImportChart == nil {
		logInfo("Import chart is nil")
		return
	}
	logInfo("Printing import chart:")
	for file, imports := range l.ImportChart {
		logCont(fmt.Sprintf("%s imports:", file))
		for _, imp := range imports {
			logCont(fmt.Sprintf("  - %s", imp))
		}
	}
}

func (l *Logger) printImportViolations() {
	if l.ImportViolations == nil || len(l.ImportViolations) == 0 {
		logInfo("No import violations found")
		return
	}
	for _, v := range l.ImportViolations {
		logError(fmt.Sprintf("File '%s' violated import rule, cannot import", v.Filename))
		logCont(fmt.Sprintf("'%s' but imported", v.CannotImport))
		logCont(fmt.Sprintf("'%s'", v.ImportLine))
	}
}

/*
 * PRINTING ENTRIES
 */

func (l *Logger) printAlwaysOutput() {
	for _, entry := range l.Entries {
		if entry.Level == LogLevelError {
			logError(entry.Contents)
		}
	}
}

func (l *Logger) printVerboseOutput() {
	for _, entry := range l.Entries {
		switch entry.Level {
		case LogLevelError:
			logError(entry.Contents)
		case LogLevelWarn:
			logWarn(entry.Contents)
		case LogLevelDebug:
			logInfo(entry.Contents)
		case LogLevelTrace:
			logTrace(entry.Contents)
		default:
			logTrace(entry.Contents)
		}
	}
	l.printRules()
	l.printImportChart()
}

func (l *Logger) Print() {
	if l.Verbose {
		l.printVerboseOutput()
	} else {
		l.printAlwaysOutput()
	}
	l.printImportViolations() // TODO: Only print if arrived at the checking stage, f.e. not if there was an error while parsing
}

/*
 * EXITING METHODS - DO NOT CALL os.Exit HERE
 */

func (l *Logger) FailWithMessage(message string) {
	l.Print()
	logError("Error occurred: " + message)
}

func (l *Logger) FailWithError(message string, err error) {
	l.Print()
	logError(fmt.Sprintf("Error occurred: %s (%s)", message, err.Error()))
}

func (l *Logger) FailWithErrors(message string, errs []error) {
	l.Print()
	logError("Error occurred: " + message)
	for _, err := range errs {
		logError(err.Error())
	}
}

func (l *Logger) Success() {
	l.Print()
	logInfo("No errors seem to have occurred!")
	logInfo("Exiting gracefully")
}
