package govalid2_test

import (
	"testing"

	"github.com/BytecodeAgency/import-boundary-checker/examples/go-valid-2/data/database"
	"github.com/BytecodeAgency/import-boundary-checker/examples/go-valid-2/domain/user"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	testData := "this is some username"
	db := database.New(testData)
	example := user.PrintUsername(db)

	assert.Equal(t, testData, example)
}
