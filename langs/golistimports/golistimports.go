package golistimports

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func ExtractForFileList(relativeFilePaths []string) (map[string][]string, error) {
	imports := make(map[string][]string, len(relativeFilePaths))
	for _, path := range relativeFilePaths {
		abs, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		source, err := ioutil.ReadFile(abs)
		if err != nil {
			return nil, err
		}
		importsForFile, err := ExtractForSourceFile(string(source))
		if err != nil {
			return nil, err
		}
		imports[path] = importsForFile
	}
	return imports, nil
}

func ExtractForSourceFile(source string) ([]string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "sourcefile.go", source, parser.ParseComments)
	if err != nil {
		return []string{}, err
	}

	if len(file.Imports) == 0 {
		return []string{}, err
	}

	var imports []string
	for _, imp := range file.Imports {
		if imp.Path != nil {
			p := *imp.Path
			// Output of imports is `\"module/path\"`, change to `module/path`
			importLine := strings.ReplaceAll(p.Value, "\"", "")
			imports = append(imports, importLine)
		}
	}
	return imports, nil
}
