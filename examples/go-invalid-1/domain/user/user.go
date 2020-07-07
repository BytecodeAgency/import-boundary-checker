package user

import "git.bytecode.nl/foss/import-boundry-checker/examples/go-invalid-1/data/database"

func GetTheUser() string {
	return database.GetUser()
}
