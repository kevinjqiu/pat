package ruletest

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
)

func TestNewPromRuleTestFromFile(t *testing.T) {
	promRuleTest, err := FromFile("testdata/test.yaml")
	assert.Nil(t, err)
	assert.Equal(t, "Test HTTP Requests too low alert", promRuleTest.Name)
	assert.Equal(t, "rules.yaml", promRuleTest.Rules.FromFile)
	assert.Equal(t, 2, len(promRuleTest.Fixtures["5m"]))
	assert.Equal(t, 2, len(promRuleTest.Assertions))
}

func TestNewPromRuleTestFromString(t *testing.T) {
	fileContent, err := ioutil.ReadFile("testdata/test.yaml")
	assert.Nil(t, err)

	promRuleTest, err := FromString(fileContent)
	assert.Nil(t, err)

	assert.Equal(t, "Test HTTP Requests too low alert", promRuleTest.Name)
	assert.Equal(t, "rules.yaml", promRuleTest.Rules.FromFile)
	assert.Equal(t, 2, len(promRuleTest.Fixtures["5m"]))
	assert.Equal(t, 2, len(promRuleTest.Assertions))
}