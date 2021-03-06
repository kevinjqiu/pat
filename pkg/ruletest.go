package pkg

import (
	"time"
	"github.com/prometheus/prometheus/rules"
	"os"
	"io/ioutil"
	"path/filepath"
	"gopkg.in/yaml.v2"
	"github.com/prometheus/prometheus/promql"
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"regexp"
	"log"
	"sort"
	"github.com/xeipuuv/gojsonschema"
	"github.com/hashicorp/go-multierror"
)

var jsonSchemaLoader gojsonschema.JSONLoader

func init() {
	asset, err := Asset("schema/schema.json")
	if err != nil {
		panic(err)
	}
	jsonSchemaLoader = gojsonschema.NewBytesLoader(asset)
}

func (prt PromRuleTest) evalRuleGroupAtInstant(suite *promql.Test, grps []*rules.Group, evalTime time.Time) ([]map[string]string, error) {
	retval := make([]map[string]string, 0)

	for _, grp := range grps {
		for _, rule := range grp.Rules() {
			results, err := rule.Eval(suite.Context(), evalTime, rules.EngineQueryFunc(suite.QueryEngine(), suite.Storage()), nil)
			if err != nil {
				return retval, err
			}

			for _, res := range results {
				retval = append(retval, res.Metric.Map())
			}
		}
	}
	return retval, nil
}

func (prt PromRuleTest) GenerateTestCases() ([]TestCase, error) {
	var tests []TestCase
	suite, err := prt.Fixtures.Load()
	if err != nil {
		return tests, err
	}

	grps, err := prt.Rules.Load()
	if err != nil {
		return tests, err
	}

	baseTime := time.Unix(0, 0)
	for i, assertion := range prt.Assertions {
		duration, err := assertion.At.ToDuration()
		if err != nil {
			return tests, err
		}
		evalTime := baseTime.Add(duration)
		resultAlertMetrics, err := prt.evalRuleGroupAtInstant(suite, grps, evalTime)
		if err != nil {
			return tests, err
		}
		tests = append(tests, prt.newTestCase(i, assertion, resultAlertMetrics))
	}

	return tests, nil
}

func (prt PromRuleTest) getTestCaseName(assertionIdx int) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	assertion := prt.Assertions[assertionIdx]
	return fmt.Sprintf("%s_at_%s", reg.ReplaceAllString(prt.Name, "_"), assertion.At)
}

func (prt PromRuleTest) newTestCase(idx int, assertion Assertion, resultAlerts []map[string]string) TestCase {
	return TestCase{
		Name: prt.getTestCaseName(idx),
		F: func(t *testing.T) {
			assert.Equal(t, len(assertion.Expected), len(resultAlerts), "Alert count does not match expected")
			assertAlertsEqual(t, assertion.Expected, resultAlerts)
		},
	}
}

func (prt PromRuleTest) validateJSONSchema() error {
	documentLoader := gojsonschema.NewGoLoader(prt)

	res, err := gojsonschema.Validate(jsonSchemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if res.Valid() {
		return nil
	}

	var multiErr error
	for _, e := range res.Errors() {
		multiErr = multierror.Append(multiErr, fmt.Errorf("validation error: %s, field=%s, value=%s", e.Description(), e.Field(), e.Value()))
	}
	return multiErr
}

func NewPromRuleTestFromFile(filename string) (PromRuleTest, error) {
	_, err := os.Stat(filename)
	if err != nil {
		return PromRuleTest{}, err
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return PromRuleTest{}, err
	}

	promRuleTest, err := NewPromRuleTestFromString(content)
	if err != nil {
		return promRuleTest, err
	}
	promRuleTest.filename = filename
	absFilePath, err := filepath.Abs(filename)
	if err != nil {
		return PromRuleTest{}, err
	}
	promRuleTest.Rules.baseDir = filepath.Dir(absFilePath)
	return promRuleTest, err
}

func NewPromRuleTestFromString(fileContent []byte) (PromRuleTest, error) {
	var promRuleTest PromRuleTest
	err := yaml.UnmarshalStrict(fileContent, &promRuleTest)
	if err != nil {
		return promRuleTest, err
	}

	err = promRuleTest.validateJSONSchema()
	if err != nil {
		return promRuleTest, err
	}

	promRuleTest.filename = FilenameInline
	promRuleTest.Rules.baseDir = "/"
	promRuleTest.testRunner = GoTestRunner{}
	return promRuleTest, nil
}

func sortedKeys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func alertLessThan(a, b map[string]string) bool {
	aKeys := sortedKeys(a)
	bKeys := sortedKeys(b)

	if len(aKeys) < len(bKeys) {
		return true
	}
	if len(aKeys) > len(bKeys) {
		return false
	}

	for i, _ := range aKeys {
		aKey := aKeys[i]
		bKey := bKeys[i]

		if aKey < bKey {
			return true
		}
		if aKey > bKey {
			return false
		}
		aVal := a[aKey]
		bVal := b[bKey]
		if aVal < bVal {
			return true
		}
	}

	return false
}

func assertAlertsEqual(t *testing.T, expected, actual []map[string]string) bool {
	// Add __name__ attribute to the expected map
	for _, e := range expected {
		if _, ok := e["__name__"]; !ok {
			e["__name__"] = "ALERTS"
		}
	}

	sort.SliceStable(expected, func(i, j int) bool {
		return alertLessThan(expected[i], expected[j])
	});

	sort.SliceStable(actual, func(i, j int) bool {
		return alertLessThan(actual[i], actual[j])
	});

	return assert.EqualValues(t, expected, actual)
}
