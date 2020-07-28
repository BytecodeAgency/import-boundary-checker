package database

import "github.com/BytecodeAgency/import-boundary-checker/examples/go-valid-3/data/interactions"

func ValidateUser() error {
	return interactions.Validate()
}
func GetUser() string {
	return "admin"
}
