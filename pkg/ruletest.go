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
	"github.com/stretchr/testify/assert"
)

func (prt PromRuleTest) evalRuleGroupAtInstant(suite *promql.Test, grps []*rules.Group, evalTime time.Time) ([]map[string]string, error) {
	var retval []map[string]string
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

func (prt PromRuleTest) Run() error {
	suite, err := prt.Fixtures.Load()
	if err != nil {
		return err
	}

	grps, err := prt.Rules.Load()
	if err != nil {
		return err
	}

	baseTime := time.Unix(0, 0)
	for _, assertion := range prt.Assertions {
		duration, err := assertion.At.ToDuration()
		if err != nil {
			return err
		}
		evalTime := baseTime.Add(duration)
		resultAlertMetrics, err := prt.evalRuleGroupAtInstant(suite, grps, evalTime)
		if err != nil {
			return err
		}
		fmt.Println(resultAlertMetrics)
		assert.Contains()
	}
	return nil
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
	promRuleTest.Filename = filename
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
	promRuleTest.Filename = FilenameInline
	promRuleTest.Rules.BaseDir = "/"
	return promRuleTest, nil
}
