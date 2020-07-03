package golistimports

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

func ExtractForSourceFile(source string) []string {
	// TODO: Parse source file, find line beginning with `import (` or `import "`
	// TODO: Implement
	//return []string{source}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "foo.go", source, parser.ParseComments)
	if err != nil {
		panic(err) // TODO: Error handling
	}

	if len(file.Imports) == 0 {
		return []string{}
	}

	var imports []string
	for _, imp := range file.Imports {
		fmt.Printf("%+v", imp.Path)
		if imp.Path != nil {
			p := *imp.Path
			importLine := strings.Replace(p.Value, "\"", "", -1)
			imports = append(imports, importLine)
		}
	}
	return imports
	//var buf bytes.Buffer
	//printer.Fprint(&buf, fset, file)
	//fmt.Println(buf.String())
}
