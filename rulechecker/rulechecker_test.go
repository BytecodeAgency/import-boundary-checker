package rulechecker_test

import (
	"testing"

	"github.com/BytecodeAgency/import-boundary-checker/parser"
	"github.com/BytecodeAgency/import-boundary-checker/rulechecker"
	"github.com/stretchr/testify/assert"
)

func TestRuleChecker_Check(t *testing.T) {
	tests := []struct {
		expectedResult bool
		rules          []parser.Rule
		importChart    rulechecker.ImportChart
	}{
		{
			true,
			[]parser.Rule{{"package/domain", []string{"package/data"}, []string{}}},
			map[string][]string{"package/domain/file.go": {"fmt", "package/typings/entities"}},
		},
		{
			true,
			[]parser.Rule{{"package/domain", []string{"package/data"}, []string{"package/data/detail"}}},
			map[string][]string{"package/domain/file.go": {"fmt", "package/typings/entities", "package/data/detail"}},
		},
		{
			false,
			[]parser.Rule{{"package/domain", []string{"package/data", "ioutil"}, []string{}}},
			map[string][]string{"package/domain/file.go": {"fmt", "ioutil"}},
		},
		{
			false,
			[]parser.Rule{{"package/domain", []string{"package/data"}, []string{}}},
			map[string][]string{"package/domain/file.go": {"fmt", "package/data"}},
		},
		{
			false,
			[]parser.Rule{{"package/domain", []string{"package/data"}, []string{"package/data/detail"}}},
			map[string][]string{"package/domain/file.go": {"fmt", "package/typings/entities", "package/data/detail", "package/data"}},
		},
	}
	for _, test := range tests {
		rc := rulechecker.New(test.rules, test.importChart)
		res := rc.Check()
		assert.Equal(t, res, test.expectedResult)
	}
}
