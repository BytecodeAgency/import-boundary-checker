package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"git.bytecode.nl/foss/import-boundry-checker/logging"
	"git.bytecode.nl/foss/import-boundry-checker/runner"
)

// TODO: Create cleaner entrypoint

func main() {
	// CLI flags
	configPath := flag.String("config", ".importrules", "Configuration path to be used when building import rule set")
	printHelp := flag.Bool("help", false, "Print CLI usage information")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	if *printHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	logger := logging.New(*verbose)

	c := getConfigString(*configPath)
	r := runner.New(c, logger)
	failed := r.Run()

	if failed {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func getConfigString(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	fullPath := abs
	config, err := ioutil.ReadFile(fullPath)
	if err != nil {
		panic(fmt.Sprintf("Could not read config file contents (%s), err %s", config, err))
	}
	return string(config)
}
