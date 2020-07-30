package govalid2

import (
	"fmt"

	"github.com/BytecodeAgency/import-boundary-checker/examples/go-invalid-4/data/database"
	"github.com/BytecodeAgency/import-boundary-checker/examples/go-invalid-4/domain/user"
)

func main() {
	db := database.New("username")
	example := user.PrintUsername(db)
	fmt.Println(example)
}
