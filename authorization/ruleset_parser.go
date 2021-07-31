package authorization

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/caesarsalad/goauthz/database"
	"gopkg.in/yaml.v2"
)

type RuleSets struct {
	RuleSet []database.Rule `yaml:"RuleSets,flow"`
}

func ruleSetValidation(rule database.Rule) error {
	var count int64
	database.DB.Model(&database.Rule{}).Where("path=? AND http_method_id=?", rule.Path, rule.HTTPMethodID).Count(&count)
	if count > 0 {
		return fmt.Errorf("Path and method must be unique together. Path: %s already exist in DB", rule.Path)
	}
	if rule.MetaLocationID == MetaLocationUrl {
		_, err := regexp.Compile(rule.Path)
		if err != nil {
			return fmt.Errorf("error while compile path regex. %s", rule.Path)
		}
	}
	return nil
}

func RuleSetYamlParser(file []byte) error {
	var rule_sets RuleSets
	err := yaml.Unmarshal([]byte(file), &rule_sets)
	if err != nil {
		return errors.New("error while ruleset yaml parsing")
	}
	for _, rule := range rule_sets.RuleSet {
		err = ruleSetValidation(rule)
		if err != nil {
			return err
		}
	}
	database.DB.Create(rule_sets.RuleSet)
	return nil
}

func RulseSetYamlDump() []byte {
	rule_sets := RuleSets{}
	database.DB.Find(&rule_sets.RuleSet)
	file, err := yaml.Marshal(rule_sets)
	if err != nil {
		log.Println("error while dump rulesets ", err)
	}
	return file
}
