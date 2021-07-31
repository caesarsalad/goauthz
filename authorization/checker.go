package authorization

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type userRules struct {
	ID             uint
	Path           string
	MetaKey        sql.NullString
	MetaValue      sql.NullString
	MetaLocationID sql.NullInt32
	HTTPMethodID   sql.NullInt32
	PathPrefix     bool
}

func regexValidation(request_path string, rule userRules) bool {
	var re *regexp.Regexp
	var ok bool
	if re, ok = Rule_regex_compiled[rule.ID]; !ok {
		re, _ = regexp.Compile(rule.Path)
		Rule_regex_compiled[rule.ID] = re
	}
	results := re.FindAllStringSubmatch(request_path, -1)
	if len(results) == 0 {
		return false
	}
	// We need handle with more generic way.
	if results[0][1] == rule.MetaValue.String {
		return true
	}
	return false
}

func metaValidation(c *fiber.Ctx, rule userRules) bool {
	request_body := c.Body()
	if rule.MetaKey.Valid {
		switch uint(rule.MetaLocationID.Int32) {
		case MetaLocationQuery:
			if c.Query(rule.MetaKey.String) == rule.MetaValue.String {
				return true
			}
			return false
		case MetaLocationBody:
			request_body_json := make(map[string]interface{})
			err := json.Unmarshal(request_body, &request_body_json)
			if err != nil {
				return false
			}
			request_body_value := request_body_json[rule.MetaKey.String]
			request_body_value_string := fmt.Sprintf("%v", request_body_value)
			if request_body_value_string == rule.MetaValue.String {
				return true
			}
			return false
		}
	}
	return true
}

func CheckRules(c *fiber.Ctx, user_id uint) bool {
	user_rules := getUserRules(user_id)
	matched := false
	request_path := c.Path()
	for _, rule := range user_rules {
		if rule.HTTPMethodID.Valid && rule.HTTPMethodID.Int32 != 0 {
			if HttpMethodIDMap[c.Method()] != uint(rule.HTTPMethodID.Int32) {
				continue
			}
		}
		if rule.MetaLocationID.Valid && uint(rule.MetaLocationID.Int32) == MetaLocationUrl {
			if regexValidation(request_path, rule) {
				matched = true
				break
			}
		} else if strings.HasPrefix(request_path, rule.Path) {
			if rule.PathPrefix {
				if metaValidation(c, rule) {
					matched = true
					break
				}
			} else {
				if request_path == rule.Path {
					if metaValidation(c, rule) {
						matched = true
						break
					}
				}
			}
		}
	}
	return matched
}
