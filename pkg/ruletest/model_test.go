package ruletest

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"time"
)

func TestNewPromRuleTestFromFile(t *testing.T) {
	promRuleTest, err := FromFile("testdata/test.yaml")
	assert.Nil(t, err)
	assert.Equal(t, "Test HTTP Requests too low alert", promRuleTest.Name)
	assert.Equal(t, "rules.yaml", promRuleTest.Rules.FromFile)
	assert.Equal(t, 2, len(promRuleTest.Fixtures[0].Metrics))
	assert.Equal(t, 2, len(promRuleTest.Assertions))
}

func TestNewPromRuleTestFromString(t *testing.T) {
	fileContent, err := ioutil.ReadFile("testdata/test.yaml")
	assert.Nil(t, err)

	promRuleTest, err := FromString(fileContent)
	assert.Nil(t, err)

	assert.Equal(t, "Test HTTP Requests too low alert", promRuleTest.Name)
	assert.Equal(t, "rules.yaml", promRuleTest.Rules.FromFile)
	assert.Equal(t, 2, len(promRuleTest.Fixtures[0].Metrics))
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

func TestMetricFixturesGeneratePromQLTestInstructions(t *testing.T) {
	mf := MetricFixtures{
		DurationMetrics{
			Duration: "1m",
			Metrics: []Metric{
				`http_requests{job="app-server", instance="0", group="canary", severity="overwrite-me"}	75`,
				`http_requests{job="app-server", instance="1", group="canary", severity="overwrite-me"}	75`,
			},
		},
		DurationMetrics{
			Duration: "5m",
			Metrics: []Metric{
				`http_errors{job="app-server", instance="0", group="canary", severity="overwrite-me"}	75`,
				`http_errors{job="app-server", instance="1", group="canary", severity="overwrite-me"}	75`,
			},
		},
	}

	instructions, err := mf.generatePromQLTestInstructions()
	assert.Nil(t, err)
	assert.Equal(t, `clear
load 1m
    http_requests{job="app-server", instance="0", group="canary", severity="overwrite-me"}	75
    http_requests{job="app-server", instance="1", group="canary", severity="overwrite-me"}	75
load 5m
    http_errors{job="app-server", instance="0", group="canary", severity="overwrite-me"}	75
    http_errors{job="app-server", instance="1", group="canary", severity="overwrite-me"}	75`, instructions)
}

func TestInstantToDuration(t *testing.T) {
	instant := Instant("5m")
	d, err := instant.ToDuration()
	assert.Nil(t, err)
	assert.Equal(t, 5 * time.Minute, d)

	instant = Instant("2h")
	d, err = instant.ToDuration()
	assert.Nil(t, err)
	assert.Equal(t, 2 * time.Hour, d)

	instant = Instant("ab")
	_, err = instant.ToDuration()
	assert.NotNil(t, err)
}
