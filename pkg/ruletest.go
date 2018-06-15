package pkg

import (
	"time"
	"fmt"
	"github.com/prometheus/prometheus/rules"
	"os"
	"io/ioutil"
	"path/filepath"
	"gopkg.in/yaml.v2"
)

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
		for _, grp := range grps {
			for _, rule := range grp.Rules() {
				res, err := rule.Eval(suite.Context(), evalTime, rules.EngineQueryFunc(suite.QueryEngine(), suite.Storage()), nil)
				if err != nil {
					return err
				}

				if len(res) != len(assertion.Expected) {
					return fmt.Errorf("")
				}
			}
		}
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
