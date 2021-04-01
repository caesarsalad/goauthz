package authorization

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/caesarsalad/goauthz/database"
	"github.com/gofiber/fiber/v2"
)

type userRules struct {
	Path           string
	MetaKey        sql.NullString
	MetaValue      sql.NullString
	MetaLocationID sql.NullInt32
}

func metaValidation(c *fiber.Ctx, rule userRules) bool {
	request_body := c.Body()
	if rule.MetaKey.Valid {
		switch int(rule.MetaLocationID.Int32) {
		case 1:
			if c.Query(rule.MetaKey.String) == rule.MetaValue.String {
				return true
			}
			return false
		case 2:
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
	var user_rules []userRules
	database.DB.Table("assigned_rules").
		Select("rules.path", "assigned_rules.meta_value", "rules.meta_key", "rules.meta_location_id").
		Joins("INNER JOIN rules ON rules.id = assigned_rules.rule_id").
		Where("assigned_rules.user_id = ?", user_id).Scan(&user_rules)
	matched := false
	request_path := c.Path()
	for _, rule := range user_rules {
		if strings.HasPrefix(request_path, rule.Path) {
			if metaValidation(c, rule) {
				matched = true
				break
			}
		}
	}
	return matched
}
