package lexer

import (
	"fmt"

	"git.bytecode.nl/foss/import-boundry-checker/keyword"
	"git.bytecode.nl/foss/import-boundry-checker/token"
)

const DEBUG = false // TODO: Support flag `-debug`

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
		return true
	}

	// Increment position and set characters
	l.position++
	l.currentChar = l.input[l.position]
	if l.position < len(l.input)-1 {
		l.nextChar = l.input[l.position+1]
	} else {
		l.nextChar = 0
	}
	return false
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
	err := fmt.Errorf("could not parse %q (next char %q) at position %d, with buffer %s and tokentype '%s' (error location: %s)", l.currentChar, l.nextChar, l.position, string(l.buffer), l.bufferTokenType, errorLocation)
	l.errors = append(l.errors, err)
}

func (l *Lexer) Format() string {
	return fmt.Sprintf("position: %d, currentChar: %q, nextChar: %q, buffer: %s, bufferType: %s, errors: %s, input: %s, results: %s",
		l.position, l.currentChar, l.nextChar, string(l.buffer), l.bufferTokenType, l.errors, string(l.input), l.result)
}

// Recursive loop that keeps running until we have reached the end of the input
func (l *Lexer) Exec() {
	if DEBUG {
		fmt.Printf("Start Exec with %s\n", l.Format()) // TODO: Add trace functionality (log in struct?)
	}
	l.execStep()
	done := l.next()
	if done {
		if DEBUG {
			fmt.Printf("Done Exec with %s\n", l.Format())
		}
		return
	}
	l.Exec()
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
	case keyword.Lang:
		l.bufferTokenType = token.KEYWORD_LANG
	case keyword.ImportRule:
		l.bufferTokenType = token.KEYWORD_IMPORTRULE
	case keyword.CannotImport:
		l.bufferTokenType = token.KEYWORD_CANNOTIMPORT
	default:
		l.logErrorAtPosition("recognizeKeywordFromBuffer")
	}
	l.buffer = []byte{}
}
