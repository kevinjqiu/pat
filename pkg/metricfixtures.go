package pkg

import (
	"fmt"
	"strings"
	"github.com/prometheus/prometheus/promql"
)

func (mf MetricFixtures) generatePromQLTestInstructions() (string, error) {
	instructions := []string{}
	instructions = append(instructions,"clear")

	for _, fixtures := range mf {
		instructions = append(instructions, fmt.Sprintf("load %s", fixtures.Duration))
		for _, metric := range fixtures.Metrics {
			instructions = append(instructions, fmt.Sprintf("    %s", metric))
		}
	}

	return strings.Join(instructions, "\n"), nil
}

func (mf MetricFixtures) Load() (*promql.Test, error) {
	instructions, err := mf.generatePromQLTestInstructions()
	if err != nil {
		return nil, err
	}

	suite, err := promql.NewTest(StubTestCase{}, instructions)
	if err != nil {
		return nil, err
	}
	return suite, nil
}
