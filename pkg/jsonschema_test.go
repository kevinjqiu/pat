package pkg

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestJSONSchemaValidationOnRules(t *testing.T) {
	type testcase struct {
		rules string
		shouldPass bool
	}

	testcases := []testcase{
		{
			`rules:
  fromFile: foobar.yaml
`, true,
		},
		{
			`rules:
  fromLiteral: HELLO
`, true,
		},
		{
			`rules:
  fromFile: foobar.yaml
  fromLiteral: HELLO
`, false,
		},
	}

	for _, tc := range testcases {
		testYaml := fmt.Sprintf(`
name: foobar
%s
fixtures: []
assertions: []
`, tc.rules)
		_, err := NewPromRuleTestFromString([]byte(testYaml))
		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}

