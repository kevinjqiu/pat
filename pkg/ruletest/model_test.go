package ruletest

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewPromRuleTestFromString(t *testing.T) {
	promRuleTest, err := FromString([]byte(`
name: Test HTTP Requests too low alert
rules:
  fromFile: rules.yaml
fixtures:
  5m:
    - http_requests{job="app-server", instance="0", group="canary", severity="overwrite-me"}	75 85  95 105 105  95  85
    - http_requests{job="app-server", instance="1", group="canary", severity="overwrite-me"}	80 90 100 110 120 130 140
assertions:
  - at: 0m
    expected:
      - ALERTS{alertname="HTTPRequestRateLow",alertstate="pending",group="canary",instance="0",job="app-server",severity="critical"} 1
      - ALERTS{alertname="HTTPRequestRateLow",alertstate="pending",group="canary",instance="1",job="app-server",severity="critical"} 1
  - at: 5m
    expected:
      - ALERTS{alertname="HTTPRequestRateLow",alertstate="firing",group="canary",instance="0",job="app-server",severity="critical"} 1
      - ALERTS{alertname="HTTPRequestRateLow",alertstate="firing",group="canary",instance="1",job="app-server",severity="critical"} 1
`))
	assert.Nil(t, err)

	assert.Equal(t, "Test HTTP Requests too low alert", promRuleTest.Name)
	assert.Equal(t, "rules.yaml", promRuleTest.Rules.FromFile)
	assert.Equal(t, 2, len(promRuleTest.Fixtures["5m"]))
	assert.Equal(t, 2, len(promRuleTest.Assertions))
}