package pkg

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
)

func TestNewPromRuleTestFromFile(t *testing.T) {
	promRuleTest, err := NewPromRuleTestFromFile("testdata/test.yaml")
	assert.Nil(t, err)
	assert.Equal(t, "Test HTTP Requests too low alert", promRuleTest.Name)
	assert.Equal(t, "rules.yaml", promRuleTest.Rules.FromFile)
	assert.Equal(t, 2, len(promRuleTest.Fixtures[0].Metrics))
	assert.Equal(t, 2, len(promRuleTest.Assertions))
	assert.Equal(t, "testdata/test.yaml", promRuleTest.Filename)
}

func TestNewPromRuleTestFromString(t *testing.T) {
	fileContent, err := ioutil.ReadFile("testdata/test.yaml")
	assert.Nil(t, err)

	promRuleTest, err := NewPromRuleTestFromString(fileContent)
	assert.Nil(t, err)

	assert.Equal(t, "Test HTTP Requests too low alert", promRuleTest.Name)
	assert.Equal(t, "rules.yaml", promRuleTest.Rules.FromFile)
	assert.Equal(t, 2, len(promRuleTest.Fixtures[0].Metrics))
	assert.Equal(t, 2, len(promRuleTest.Assertions))
	assert.Equal(t, FilenameInline, promRuleTest.Filename)
}

func TestPromRuleTestRun(t *testing.T) {
	promRuleTest, err := NewPromRuleTestFromFile("testdata/test.yaml")
	assert.Nil(t, err)
	err = promRuleTest.Run()
	assert.Nil(t, err)
	//fmt.Println(err)
	//err = prt.Run()
	//fmt.Print(err)
}
