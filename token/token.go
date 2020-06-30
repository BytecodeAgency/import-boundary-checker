package token

type Token string

const (
	UNKNOWN      = Token("unknown") // TODO: Make sure this does not make it out of the lexer
	KEYWORD_LANG = Token("kw_lang")
	IDENTIFIER   = Token("ident")
	STRING       = Token("string")
)
