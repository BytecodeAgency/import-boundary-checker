package lexer

import (
	"git.bytecode.nl/foss/import-boundry-checker/token"
)

type Result struct {
	Token    token.Token
	Contents string
}

type Lexer struct {
	// Input and result
	input  []rune
	result []Result

	// Intermediate values
	buffer          []rune
	bufferTokenType token.Token
	position        int
	currentChar     rune
	nextChar        rune
}

func New(input string) Lexer {
	return Lexer{
		input:           []rune(input),
		result:          []Result{},
		buffer:          []rune{},
		bufferTokenType: token.UNKNOWN,
		position:        0,
		currentChar:     rune(input[0]),
		nextChar:        rune(input[1]),
	}
}

func (l Lexer) Result() []Result {
	return l.result
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
	l.buffer = []rune{}
	l.bufferTokenType = token.UNKNOWN
}

// Recursive loop that keeps running until we have reached the end of the input
func (l *Lexer) Exec() { // TODO: Return error
	l.execStep()
	done := l.next()
	if done {
		return
	}
	l.Exec()
}

// Executes a single step, called by the Exec loop, but does not loop itself
func (l *Lexer) execStep() {
	// If buffer tokenType is unknown, do type recognition
	if l.bufferTokenType == token.UNKNOWN {
		switch l.currentChar {
		case ' ': // TODO: Better handle whitespace
			return
		case '"': // If the current char is ", we have encountered a string
			l.bufferTokenType = token.STRING
			return
		default:
			panic("Should not reach this code (1)") // TODO: Return error message somehow
		}
	}

	// Buffer is filled, so continue until we encounter the end of the current Token content
	switch l.bufferTokenType {
	case token.STRING:
		// End of string
		if l.currentChar == '"' {
			l.finishBuffer()
		} else {
			l.currentCharToBuffer()
		}
		return
	}
	panic("Should not reach this code (2)") // TODO: Return error message somehow
}
