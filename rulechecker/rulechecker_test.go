package rulechecker_test

import (
	"testing"

	"github.com/BytecodeAgency/import-boundry-checker/parser"
	"github.com/BytecodeAgency/import-boundry-checker/rulechecker"
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
			[]parser.Rule{{"package/domain", []string{"package/data"}}},
			map[string][]string{"package/domain/file.go": {"fmt", "package/typings/entities"}},
		},
		{
			false,
			[]parser.Rule{{"package/domain", []string{"package/data", "ioutil"}}},
			map[string][]string{"package/domain/file.go": {"fmt", "ioutil"}},
		},
		{
			false,
			[]parser.Rule{{"package/domain", []string{"package/data"}}},
			map[string][]string{"package/domain/file.go": {"fmt", "package/data"}},
		},
	}
	for _, test := range tests {
		rc := rulechecker.New(test.rules, test.importChart)
		res := rc.Check()
		assert.Equal(t, res, test.expectedResult)
	}
}
