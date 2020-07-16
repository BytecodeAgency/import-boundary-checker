package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/BytecodeAgency/import-boundary-checker/logging"
	"github.com/BytecodeAgency/import-boundary-checker/runner"
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

	// Parse config file
	c, err := getConfigString(*configPath)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	// Create logger and run app
	logger := logging.New(*verbose)
	r := runner.New(c, logger)
	failed := r.Run()

	// Get and print the logs
	fmt.Print(logger.Logs.String())

	if failed {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func getConfigString(path string) (string, error) {
	fullPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("parse absolute config path: %v", err)
	}
	config, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("read config file (%s): %v", path, err)
	}
	return string(config), nil
}
