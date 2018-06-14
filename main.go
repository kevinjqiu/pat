package main

import (
	"github.com/hashicorp/go-multierror"
	"log"
	"path/filepath"
	"os"
	"flag"
	"fmt"
	"os/exec"
	"encoding/json"
	"encoding/base64"
)

const EnvVarTestFilePathsB64 = "TEST_FILE_PATHS_B64"

func collectTestFiles(globPatterns []string) ([]string, error) {
	if len(globPatterns) == 0 {
		return filepath.Glob("*.test")
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

func encodeTestFilePaths(testFiles []string) (string, error) {
	jsonBody, err := json.Marshal(testFiles)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(jsonBody), err
}

func decodeTestFilePaths(encodedTestFilePaths string) ([]string, error) {
	jsonBody, err := base64.StdEncoding.DecodeString(encodedTestFilePaths)
	if err != nil {
		return []string{}, err
	}

	var filePaths []string
	err = json.Unmarshal(jsonBody, &filePaths)

	if err != nil {
		return []string{}, err
	}
	return filePaths, nil
}

func main() {
	flag.Parse()
	testFiles, err := collectTestFiles(flag.Args())

	encodedTestFilePaths, err := encodeTestFilePaths(testFiles)

	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("go", "test")

	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", EnvVarTestFilePathsB64, encodedTestFilePaths))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
