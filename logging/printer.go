package logging

import "fmt"

var (
	logError = Color("\033[1;31mERROR\033[0m")
	logWarn  = Color("\033[1;33mWARN\033[0m")
	logInfo  = Color("\033[1;36mINFO\033[0m")
	logTrace = Color("\033[1;32mTRACE\033[0m") // Change to 37 for white
	logCont  = Color("\033[1;32m  ->\033[0m")  // Change to 37 for white
)

func Color(colorString string) func(message string) {
	printer := func(message string) {
		fmt.Printf("%5s: %s\n", colorString, message)
	}
	return printer
}
