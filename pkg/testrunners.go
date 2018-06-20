package pkg

import (
	"testing"
	"io"
	"errors"
)

type TestCase struct {
	Name string
	F    func(*testing.T)
}

type TestRunner interface {
	RunTests([]TestCase) int
}

type GoTestRunner struct{}

var errMain = errors.New("testing: unexpected use of func Main")
type matcher func(pat, str string) (bool, error)
func (f matcher) MatchString(pat, str string) (bool, error)   { return f(pat, str) }
func (f matcher) StartCPUProfile(w io.Writer) error           { return errMain }
func (f matcher) StopCPUProfile()                             {}
func (f matcher) WriteHeapProfile(w io.Writer) error          { return errMain }
func (f matcher) WriteProfileTo(string, io.Writer, int) error { return errMain }
func (f matcher) ImportPath() string                          { return "" }
func (f matcher) StartTestLog(io.Writer)                      {}
func (f matcher) StopTestLog() error                          { return errMain }

func (gtr GoTestRunner) RunTests(tc []TestCase) int {
	// convert TestCases to testing.InternalTest
	var testcases []testing.InternalTest
	for _, test := range tc {
		testcases = append(testcases, testing.InternalTest{
			Name: test.Name,
			F:    test.F,
		})
	}

	matchAllFunc := matcher(func(pat, str string) (bool, error) { return true, nil })
	return testing.MainStart(matchAllFunc, testcases, nil, nil).Run()
}
