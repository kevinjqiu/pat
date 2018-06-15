package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	pat "github.com/kevinjqiu/pat/pkg"
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
	flag.Parse()
	testFiles, err := collectTestFiles(flag.Args())

	if err != nil {
		log.Fatal(err)
	}

	allTestCases := []pat.TestCase{}
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

	testRunner := pat.GoTestRunner{}
	testRunner.RunTests(allTestCases)
}
