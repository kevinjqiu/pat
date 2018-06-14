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

func TestRuleLoaderFromFile(t *testing.T) {
	ruleLoader := RuleLoader{FromFile: "testdata/rules.yaml"}

	ruleGroups, err := ruleLoader.Load()
	assert.Nil(t, err)

	assert.Equal(t, 1, len(ruleGroups))
	assert.Equal(t, "prometheus.rules", ruleGroups[0].Name())
	assert.Equal(t, "testdata/rules.yaml", ruleGroups[0].File())
	assert.Equal(t, 1, len(ruleGroups[0].Rules()))
	assert.Equal(t, "HTTPRequestRateLow", ruleGroups[0].Rules()[0].Name())
}

func TestRuleLoaderFromLiteral(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/rules.yaml")
	assert.Nil(t, err)
	ruleLoader := RuleLoader{FromLiteral: string(content)}

	ruleGroups, err := ruleLoader.Load()
	assert.Nil(t, err)

	assert.Equal(t, 1, len(ruleGroups))
	assert.Equal(t, "prometheus.rules", ruleGroups[0].Name())
	assert.Equal(t, "__inline__", ruleGroups[0].File())
	assert.Equal(t, 1, len(ruleGroups[0].Rules()))
	assert.Equal(t, "HTTPRequestRateLow", ruleGroups[0].Rules()[0].Name())
}

func TestRuleLoaderNoLoadingStrategySpecified(t *testing.T) {
	ruleLoader := RuleLoader{}
	_, err := ruleLoader.Load()
	assert.NotNil(t, err)
}