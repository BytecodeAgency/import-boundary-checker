package database

import "github.com/BytecodeAgency/import-boundary-checker/examples/go-invalid-5/data/interactions"

func ValidateUser() error {
	return interactions.Validate()
}
func GetUser() string {
	return "admin"
}
