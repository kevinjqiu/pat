package pkg

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMetricFixturesGeneratePromQLTestInstructions(t *testing.T) {
	mf := MetricFixtures{
		DurationMetricsFixture{
			Duration: "1m",
			Metrics: []Metric{
				`http_requests{job="app-server", instance="0", group="canary", severity="overwrite-me"}	75`,
				`http_requests{job="app-server", instance="1", group="canary", severity="overwrite-me"}	75`,
			},
		},
		DurationMetricsFixture{
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

