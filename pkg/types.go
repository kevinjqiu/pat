package pkg

import (
	"log"
	"time"
)

const FilenameInline = "__inline__"

type StubTestCase struct{}

func (stc StubTestCase) Fatal(args ...interface{}) {
	log.Fatal(args)
}

func (stc StubTestCase) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args)
}

type RuleLoader struct {
	FromFile    string `yaml:"fromFile,omitempty" json:"fromFile,omitempty"`
	FromLiteral string `yaml:"fromLiteral,omitempty" json:"fromLiteral,omitempty"`
	baseDir     string
}

type Assertion struct {
	At       Instant `yaml:"at" json:"at"`
	Expected []Alert `yaml:"expected" json:"expected"`
	Comment  string  `yaml:"comment,omitempty" json:"comment,omitempty"`
}

type Alert = map[string]string

type PromRuleTest struct {
	Name       string         `yaml:"name" json:"name"`
	Rules      RuleLoader     `yaml:"rules" json:"rules"`
	Fixtures   MetricFixtures `yaml:"fixtures" json:"fixtures"`
	Assertions []Assertion    `yaml:"assertions" json:"assertions"`
	filename   string
	testRunner TestRunner
}

type Instant string

func (i Instant) ToDuration() (time.Duration, error) {
	duration, err := time.ParseDuration(string(i))
	if err != nil {
		return duration, err
	}
	return duration, nil
}

type Duration = string

type Metric = string

type DurationMetricsFixture struct {
	Duration Duration `yaml:"duration" json:"duration"`
	Metrics  []Metric `yaml:"metrics" json:"metrics"`
}

type MetricFixtures []DurationMetricsFixture
