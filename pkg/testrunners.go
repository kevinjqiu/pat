package pkg

import "testing"

type TestCase struct {
	Name string
	F    func(*testing.T)
}

type TestRunner interface {
	RunTests([]TestCase)
}

type GoTestRunner struct{}

func (gtr GoTestRunner) RunTests(tc []TestCase) {
	// convert TestCases to testing.InternalTest
	testcases := []testing.InternalTest{}
	for _, test := range tc {
		testcases = append(testcases, testing.InternalTest{
			Name: test.Name,
			F:    test.F,
		})
	}

	testing.Main(func(pat, str string) (bool, error) { return true, nil }, testcases, nil, nil)
}
