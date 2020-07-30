package govalid2

import (
	"fmt"

	"github.com/BytecodeAgency/import-boundary-checker/examples/go-invalid-3/domain/user"
)

func main() {
	example := user.PrintUsername("exampleusername")
	fmt.Println(example)
}
