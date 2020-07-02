package main

import "fmt"

// TODO: Create standardized logging system
// TODO: Add support for logging from Lexer and Parser
var stacktrace []string

// TODO: Add better (colorized) logger
func log(line string) {
	stacktrace = append(stacktrace, "INFO: " + line)
	fmt.Println(line)
}

func logIfVerbose(line string) {
	stacktrace = append(stacktrace, "VERBOSE: " + line)
	if verbose {
		log(line)
	}
}

func prettyprintErrs(errs []error) string {
	var errStr string
	for _, err := range errs {
		errStr += fmt.Sprintf(" -> %s\n", err)
	}
	return errStr
}

func fail(message string) {
	log("Exiting process, stacktrace:")
	for _, line := range stacktrace {
		log(line)
	}
	panic(message)
}
