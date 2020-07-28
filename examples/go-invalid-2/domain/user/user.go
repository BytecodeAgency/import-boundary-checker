package user

import (
	"github.com/BytecodeAgency/import-boundary-checker/examples/go-invalid-2/data/database"
	"github.com/BytecodeAgency/import-boundary-checker/examples/go-invalid-2/data/interactions"
)

func Validate() error {
	return interactions.Validate()
}

func GetTheUser() string {
	return database.GetUser()
}
