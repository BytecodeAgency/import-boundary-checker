package user

import "github.com/BytecodeAgency/import-boundary-checker/examples/go-invalid-3/data/database"

func PrintUsername(username string) string {
	db := database.New(username)
	return db.Username()
}
