package user

import "github.com/BytecodeAgency/import-boundary-checker/examples/go-valid-2/data/interactors"

func PrintUsername(db interactors.DatabaseInteractor) string {
	return db.Username()
}
