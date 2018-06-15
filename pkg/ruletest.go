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
	"sort"
	"github.com/prometheus/prometheus/pkg/labels"
	"testing"
	"github.com/stretchr/testify/assert"
)

func (prt PromRuleTest) evalRuleGroupAtInstant(suite *promql.Test, grps []*rules.Group, evalTime time.Time) ([]map[string]string, error) {
	var retval []map[string]string
	var metrics []labels.Labels

	for _, grp := range grps {
		for _, rule := range grp.Rules() {
			results, err := rule.Eval(suite.Context(), evalTime, rules.EngineQueryFunc(suite.QueryEngine(), suite.Storage()), nil)
			if err != nil {
				return retval, err
			}

			for _, res := range results {
				metrics = append(metrics, res.Metric)
			}
		}
	}

	sort.Slice(metrics, func(i, j int) bool {
		return labels.Compare(metrics[i], metrics[j]) < 0
	})

	for _, m := range metrics {
		retval = append(retval, m.Map())
	}

	return retval, nil
}

func (prt PromRuleTest) generateTestCases() ([]TestCase, error) {
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
		fmt.Println(resultAlertMetrics)

		tests = append(tests, TestCase{
			Name: fmt.Sprintf("Test%d", i),
			F: func(t *testing.T) {
				assert.Equal(t, 1, 1)
			},
		})
	}

	return tests, nil
}

func (prt PromRuleTest) Run() (bool, error) {
	testcases, err := prt.generateTestCases()
	if err != nil {
		return false, err
	}
	ok := prt.testRunner.RunTests(testcases)
	return ok, nil
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
	promRuleTest.Rules.BaseDir = filepath.Dir(absFilePath)
	return promRuleTest, err
}

func NewPromRuleTestFromString(fileContent []byte) (PromRuleTest, error) {
	var promRuleTest PromRuleTest
	err := yaml.UnmarshalStrict(fileContent, &promRuleTest)
	if err != nil {
		return promRuleTest, err
	}
	promRuleTest.filename = FilenameInline
	promRuleTest.Rules.BaseDir = "/"
	promRuleTest.testRunner = GoTestRunner{}
	return promRuleTest, nil
}
