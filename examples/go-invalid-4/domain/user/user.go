package user

import (
	"fmt"

	"github.com/BytecodeAgency/import-boundary-checker/examples/go-invalid-4/data/interactors"
)

func PrintUsername(db interactors.DatabaseInteractor) string {
	return fmt.Sprintf("The username is %s", db.Username())
}
