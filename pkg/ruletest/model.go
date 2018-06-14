package ruletest
import (
	"gopkg.in/yaml.v2"
	"os"
	"io/ioutil"
	"github.com/prometheus/prometheus/rules"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"github.com/prometheus/prometheus/promql"
	"time"
	"github.com/prometheus/prometheus/pkg/labels"
	"fmt"
	"strings"
	"log"
)

type StubTestCase struct{}

func (stc StubTestCase) Fatal(args ...interface{}) {
	log.Fatal(args)
}
func (stc StubTestCase) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args)
}

type RuleLoader struct {
	FromFile    string `yaml:"fromFile,omitempty"`
	FromLiteral string `yaml:"fromLiteral,omitempty"`
}

func (ruleLoader RuleLoader) Load() ([]*rules.Group, error) {
	var (
		ruleGroups *rulefmt.RuleGroups
		errs []error
		filename string
	)

	switch {
	case ruleLoader.FromFile != "":
		filename = ruleLoader.FromFile
		ruleGroups, errs = rulefmt.ParseFile(ruleLoader.FromFile)
		// TODO: add validation?
		if len(errs) != 0 {
			return nil, errs[0] // TODO: multi-error
		}
	case ruleLoader.FromLiteral != "":
		filename = "__inline__"
		ruleGroups, errs = rulefmt.Parse([]byte(ruleLoader.FromLiteral))
		if len(errs) != 0 {
			return nil, errs[0]
		}
	}

	if ruleGroups == nil {
		return nil, fmt.Errorf("Must specify a rule loading strategy")
	}

	retval := []*rules.Group{}
	for _, rg := range ruleGroups.Groups {
		rls := make([]rules.Rule, 0, len(rg.Rules))
		for _, r := range rg.Rules {
			expr, err := promql.ParseExpr(r.Expr)
			if err != nil {
				return nil, err
			}

			if r.Alert != "" {
				rule := rules.NewAlertingRule(
					r.Alert,
					expr,
					time.Duration(r.For),
					labels.FromMap(r.Labels),
					labels.FromMap(r.Annotations),
					nil,    // TODO: stub logger
				)
				rls = append(rls, rule)
			}
		}
		group := rules.NewGroup(rg.Name, filename, time.Duration(rg.Interval), rls, &rules.ManagerOptions{})
		retval = append(retval, group)
	}
	return retval, nil
}

type Instant string

type Duration string

type Metric string

type DurationMetrics struct {
	Duration Duration `yaml:"duration"`
	Metrics  []Metric `yaml:"metrics"`
}

type MetricFixtures []DurationMetrics

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

func (prt PromRuleTest) Run() error {
	suite, err := prt.Fixtures.Load()
	if err != nil {
		return err
	}
	fmt.Println(suite)
	return nil
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
