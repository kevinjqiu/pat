package pkg

import (
	"testing"
	"path/filepath"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
)

func TestRuleLoaderFromFile(t *testing.T) {
	currentDir, err := filepath.Abs(".")
	assert.Nil(t, err)

	ruleLoader := RuleLoader{FromFile: "testdata/rules.yaml", baseDir: currentDir}

	ruleGroups, err := ruleLoader.Load()
	assert.Nil(t, err)

	assert.Equal(t, 1, len(ruleGroups))
	assert.Equal(t, "prometheus.rules", ruleGroups[0].Name())
	assert.Equal(t, "testdata/rules.yaml", ruleGroups[0].File())
	assert.Equal(t, 2, len(ruleGroups[0].Rules()))
	assert.Equal(t, "HTTPRequestRateLow", ruleGroups[0].Rules()[0].Name())
	assert.Equal(t, "httprequestratelow_record", ruleGroups[0].Rules()[1].Name())
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
	assert.Equal(t, 2, len(ruleGroups[0].Rules()))
	assert.Equal(t, "HTTPRequestRateLow", ruleGroups[0].Rules()[0].Name())
	assert.Equal(t, "httprequestratelow_record", ruleGroups[0].Rules()[1].Name())
}

func TestRuleLoaderNoLoadingStrategySpecified(t *testing.T) {
	ruleLoader := RuleLoader{}
	_, err := ruleLoader.Load()
	assert.NotNil(t, err)
}

