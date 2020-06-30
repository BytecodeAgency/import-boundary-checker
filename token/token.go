package token

type Token string

const (
	UNKNOWN   = Token("unknown") // TODO: Make sure this does not make it out of the lexer
	STRING    = Token("string")
	SEMICOLON = Token("semicolon")
	// TODO: Add keyword support (lang, importrule, cannotimport)
)
