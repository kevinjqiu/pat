package main

import (
	"github.com/hashicorp/go-multierror"
	"log"
	"path/filepath"
	"os"
	"flag"
	pat "github.com/kevinjqiu/pat/pkg"
)

const EnvVarTestFilePathsB64 = "TEST_FILE_PATHS_B64"

func collectTestFiles(globPatterns []string) ([]string, error) {
	if len(globPatterns) == 0 {
		return filepath.Glob("test_*.yaml")
	}

	var (
		err error
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

	for _, testFile := range testFiles {
		prt, err := pat.NewPromRuleTestFromFile(testFile)
		if err != nil {
			log.Fatal(err)
		}
		prt.Run()
	}
}
