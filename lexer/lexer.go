package lexer

import (
	"flag"
	"fmt"

	"github.com/BytecodeAgency/import-boundary-checker/keyword"
	"github.com/BytecodeAgency/import-boundary-checker/token"
)

var DEBUG = false

func init() {
	flag.BoolVar(&DEBUG, "debug_lexer", false, "Enable debugging for lexer")
}

type Result struct {
	Token    token.Token
	Contents string
}

type Lexer struct {
	// Input and result
	input  []byte
	result []Result
	errors []error

	// Intermediate values
	buffer          []byte
	bufferTokenType token.Token
	position        int
	line            int
	currentChar     byte
	nextChar        byte
}

func New(input string) Lexer {
	return Lexer{
		input:           []byte(input),
		result:          []Result{},
		buffer:          []byte{},
		bufferTokenType: token.UNSET,
		position:        0,
		line:            1,
		currentChar:     input[0],
		nextChar:        input[1],
	}
}

func (l Lexer) Result() ([]Result, []error) {
	return l.result, l.errors
}

func (l *Lexer) next() (done bool) {
	// Return if we are at the end of the input
	if l.position+1 == len(l.input) {
		return false
	}

	// Increment line if we encounter a linebreak
	if l.currentChar == '\n' {
		l.line++
	}

	// Increment position and set characters
	l.position++
	l.currentChar = l.input[l.position]
	if l.position < len(l.input)-1 {
		l.nextChar = l.input[l.position+1]
	} else {
		l.nextChar = 0
	}
	return true
}

func (l *Lexer) currentCharToBuffer() {
	l.buffer = append(l.buffer, l.currentChar)
}

func (l *Lexer) finishBuffer() {
	l.result = append(l.result, Result{l.bufferTokenType, string(l.buffer)})
	l.buffer = []byte{}
	l.bufferTokenType = token.UNSET
}

func (l *Lexer) logErrorAtPosition(errorLocation string) {
	errDebug := fmt.Errorf("debug info: buffer %s and tokentype '%s', absolute location %d (error location: %s)", string(l.buffer), l.bufferTokenType, l.position, errorLocation)
	err := fmt.Errorf("could not parse %q (next char %q) on line %d", l.currentChar, l.nextChar, l.line)
	l.errors = append(l.errors, err, errDebug)
}

func (l *Lexer) Format() string {
	return fmt.Sprintf("position: %d, currentChar: %q, nextChar: %q, buffer: %s, bufferType: %s, errors: %s, input: %s, results: %s",
		l.position, l.currentChar, l.nextChar, string(l.buffer), l.bufferTokenType, l.errors, string(l.input), l.result)
}

// Recursive loop that keeps running until we have reached the end of the input
func (l *Lexer) Exec() {
	if DEBUG {
		fmt.Printf("Called Exec with %s\n", l.Format()) // TODO: Add trace functionality (log in struct?)
	}
	l.execStep()   // Call l.execStep once to handle the first character, then enter into loop to finish the tokenization
	for l.next() { // l.next() returns true when should continue, false when done
		if DEBUG {
			fmt.Printf("Doing exepStep with %s\n", l.Format())
		}
		l.execStep()
	}
	if DEBUG {
		fmt.Printf("Done Exec with %s\n", l.Format())
	}
}

// Executes a single step, called by the Exec loop, but does not loop itself
func (l *Lexer) execStep() {
	if l.bufferTokenType == token.UNSET {
		l.execStepUnknownTokenType()
	} else {
		l.execStepKnownTokenType()
	}
}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\n' || char == '\t' || char == '\r'
}

func (l *Lexer) execStepUnknownTokenType() {
	// Ignore all whitespace
	if isWhitespace(l.currentChar) {
		return // l.next() is called in l.Exec() after returning
	}

	// Recognize the type
	switch l.currentChar {
	case '"': // If the current char is ", we have encountered a string
		l.bufferTokenType = token.STRING
	case ';':
		l.bufferTokenType = token.SEMICOLON
		l.finishBuffer()
	default:
		l.bufferTokenType = token.UNKNOWN
		l.currentCharToBuffer()
	}
}

func (l *Lexer) execStepKnownTokenType() {
	switch l.bufferTokenType {
	case token.STRING:
		if l.currentChar == '"' { // End of string
			l.finishBuffer()
		} else {
			l.currentCharToBuffer()
		}
	case token.UNKNOWN: // In current version, can only be keyword
		if isWhitespace(l.currentChar) { // End of keyword
			l.recognizeKeywordFromBuffer()
			l.finishBuffer()
		} else {
			l.currentCharToBuffer()
		}
	default:
		l.logErrorAtPosition("execStepKnownTokenType")
	}
}

func (l *Lexer) recognizeKeywordFromBuffer() {
	switch string(l.buffer) {
	case string(keyword.Lang):
		l.bufferTokenType = token.KEYWORD_LANG
	case string(keyword.ImportRule):
		l.bufferTokenType = token.KEYWORD_IMPORTRULE
	case string(keyword.CannotImport):
		l.bufferTokenType = token.KEYWORD_CANNOTIMPORT
	case string(keyword.ImportBase):
		l.bufferTokenType = token.KEYWORD_IMPORTBASE
	default:
		l.logErrorAtPosition("recognizeKeywordFromBuffer")
	}
	l.buffer = []byte{}
}
