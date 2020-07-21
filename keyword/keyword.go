package keyword

type Keyword string

const (
	Lang         = Keyword("LANG")
	ImportRule   = Keyword("IMPORTRULE")
	CannotImport = Keyword("CANNOTIMPORT")
	Allow        = Keyword("ALLOW")
	ImportBase   = Keyword("IMPORTBASE")
)
