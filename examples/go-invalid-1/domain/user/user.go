package user

import "github.com/BytecodeAgency/import-boundary-checker/examples/go-invalid-1/data/database"

func GetTheUser() string {
	return database.GetUser()
}
