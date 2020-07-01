package token

type Token string

const (
	UNSET                = Token("unset")   // TODO: Make sure this does not make it out of the lexer
	UNKNOWN              = Token("unknown") // TODO: Make sure this does not make it out of the lexer
	STRING               = Token("string")
	SEMICOLON            = Token("semicolon")
	COMMA                = Token("comma")
	KEYWORD_LANG         = Token("kw_lang")
	KEYWORD_IMPORTRULE   = Token("kw_importrule")
	KEYWORD_CANNOTIMPORT = Token("kw_cannotimport")
)
