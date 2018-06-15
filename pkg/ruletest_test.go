package pkg

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
)

func TestGetTestCaseName(t *testing.T) {
	prt := PromRuleTest{Name: "Test HTTP Requests too low alert"}
	assert.Equal(t, "Test_HTTP_Requests_too_low_alert_1", prt.getTestCaseName(1))
}

func TestNewPromRuleTestFromFile(t *testing.T) {
	promRuleTest, err := NewPromRuleTestFromFile("testdata/test.yaml")
	assert.Nil(t, err)
	assert.Equal(t, "Test HTTP Requests too low alert", promRuleTest.Name)
	assert.Equal(t, "rules.yaml", promRuleTest.Rules.FromFile)
	assert.Equal(t, 2, len(promRuleTest.Fixtures[0].Metrics))
	assert.Equal(t, 2, len(promRuleTest.Assertions))
	assert.Equal(t, "testdata/test.yaml", promRuleTest.filename)
}

func TestNewPromRuleTestFromString(t *testing.T) {
	fileContent, err := ioutil.ReadFile("testdata/test.yaml")
	assert.Nil(t, err)

	promRuleTest, err := NewPromRuleTestFromString(fileContent)
	assert.Nil(t, err)

	assert.Equal(t, "Test HTTP Requests too low alert", promRuleTest.Name)
	assert.Equal(t, "rules.yaml", promRuleTest.Rules.FromFile)
	assert.Equal(t, 2, len(promRuleTest.Fixtures[0].Metrics))
	assert.Equal(t, 2, len(promRuleTest.Assertions))
	assert.Equal(t, FilenameInline, promRuleTest.filename)
}

