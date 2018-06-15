package pkg

import (
	"path"
	"github.com/prometheus/prometheus/rules"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"fmt"
	"github.com/prometheus/prometheus/promql"
	"time"
	"github.com/prometheus/prometheus/pkg/labels"
)

func (ruleLoader RuleLoader) getRuleFilePath() string {
	if path.IsAbs(ruleLoader.FromFile) {
		return ruleLoader.FromFile
	}
	return path.Join(ruleLoader.BaseDir, ruleLoader.FromFile)
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
		ruleGroups, errs = rulefmt.ParseFile(ruleLoader.getRuleFilePath())
		if len(errs) != 0 {
			return nil, errs[0] // TODO: multi-error
		}
		// TODO: add validation?
	case ruleLoader.FromLiteral != "":
		filename = FilenameInline
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
