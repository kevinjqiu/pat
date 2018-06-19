package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	pat "github.com/kevinjqiu/pat/pkg"
	"fmt"
)

func collectTestFiles(globPatterns []string) ([]string, error) {
	if len(globPatterns) == 0 {
		return filepath.Glob("test_*.yaml")
	}

	var (
		err       error
		filePaths []string
	)

	for _, pattern := range globPatterns {
		files, newErr := filepath.Glob(pattern)
		if err != nil {
			err = multierror.Append(err, newErr)
		}

		for _, file := range files {
			if _, err := os.Stat(file); err == nil {
				filePaths = append(filePaths, file)
			}
		}
	}

	return filePaths, err
}

func main() {
	fs := flag.NewFlagSet("pat", flag.ContinueOnError)
	err := fs.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	testFiles, err := collectTestFiles(fs.Args()[1:])

	if err != nil {
		log.Fatal(err)
	}

	var allTestCases []pat.TestCase
	for _, testFile := range testFiles {
		prt, err := pat.NewPromRuleTestFromFile(testFile)
		if err != nil {
			log.Fatal(err)
		}

		testCasesForFile, err := prt.GenerateTestCases()
		for _, tc := range testCasesForFile {
			allTestCases = append(allTestCases, tc)
		}
	}

	if len(allTestCases) == 0 {
		fmt.Println("WARNING: No tests discovered. Exiting...")
		os.Exit(0)
	}

	testRunner := pat.GoTestRunner{}
	os.Exit(testRunner.RunTests(allTestCases))
}
