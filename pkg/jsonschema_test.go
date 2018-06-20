package pkg

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

type testcase struct {
	rules      string
	fixtures   string
	assertions string
	shouldPass bool
}

func (tc testcase) generateYAML() string {
	rules := tc.rules
	if rules == "" {
		rules = `rules:
  fromFile: foobar.yaml`
	}

	fixtures := tc.fixtures
	if fixtures == "" {
		fixtures = `fixtures: []`
	}

	assertions := tc.assertions
	if assertions == "" {
		assertions = `assertions: []`
	}

	return fmt.Sprintf(`
name: foobar
%s
%s
%s
`, rules, fixtures, assertions)
}

func (tc testcase) assert(t *testing.T) {
	testYaml := tc.generateYAML()
	_, err := NewPromRuleTestFromString([]byte(testYaml))
	if tc.shouldPass {
		assert.Nil(t, err)
	} else {
		assert.NotNil(t, err)
	}
}

func TestJSONSchemaValidationOnRules(t *testing.T) {
	testcases := []testcase{
		{
			rules: `rules:
  fromFile: foobar.yaml
`,
			shouldPass:true,
		},
		{
			rules: `rules:
  fromLiteral: HELLO
`,
			shouldPass: true,
		},
		{
			rules: `rules:
  fromFile: foobar.yaml
  fromLiteral: HELLO
`,
			shouldPass: false,
		},
	}

	for _, tc := range testcases {
		tc.assert(t)
	}
}

func TestJSONSchemaValidationOnFixture(t *testing.T) {
	testcases := []testcase{
		{
			fixtures:   `fixtures: []`,
			shouldPass: true,
		},
		{
			fixtures: `fixtures:
  - duration: 5d
    metrics: []`,
			shouldPass: false,
		},
		{
			fixtures: `fixtures:
  - duration: 1m
    metrics:
      - http_requests{server="foo"} 100
  - duration: 5m
    metrics:
      - http_requests{server="foo"} 100 200 300 400`,
			shouldPass: true,
		},
	}

	for _, tc := range testcases {
		tc.assert(t)
	}
}

func TestJSONSchemaValidationOnAssertions(t *testing.T) {
	testcases := []testcase{
		{
			assertions: `assertions:
  - at: 5m
    expected:
      - alertname: FOO
        alertstate: firing`,
        	shouldPass: true,
		},
		{
			assertions: `assertions:
  - at: 5d
    expected:
      - alertname: FOO
        alertstate: firing`,
			shouldPass: false,
		},
	}

	for _, tc := range testcases {
		tc.assert(t)
	}
}
