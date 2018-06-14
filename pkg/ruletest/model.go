package ruletest
import (
	"gopkg.in/yaml.v2"
	"os"
	"io/ioutil"
)

type RuleLoader struct {
	FromFile    string `yaml:"fromFile"`
	FromLiteral string `yaml:"fromLiteral"`
}

type Instant string

type Duration string

type Metric string

type MetricFixtures map[Duration][]Metric

type Assertion struct {
	At       Instant  `yaml:"at"`
	Expected []Metric `yaml:"expected"`
	Comment  string   `yaml:"comment,omitempty"`
}

type PromRuleTest struct {
	Name       string         `yaml:"name"`
	Rules      RuleLoader     `yaml:"rules"`
	Fixtures   MetricFixtures `yaml:"fixtures"`
	Assertions []Assertion    `yaml:"assertions"`
}

func FromFile(filepath string) (PromRuleTest, error) {
	_, err := os.Stat(filepath)
	if err != nil {
		return PromRuleTest{}, err
	}

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return PromRuleTest{}, err
	}

	return FromString(content)
}

func FromString(fileContent []byte) (PromRuleTest, error) {
	var promRuleTest PromRuleTest
	err := yaml.UnmarshalStrict(fileContent, &promRuleTest)
	if err != nil {
		return promRuleTest, err
	}
	return promRuleTest, nil
}
