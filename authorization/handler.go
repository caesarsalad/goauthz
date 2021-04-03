package authorization

import (
	"bufio"
	"log"
	"regexp"

	"github.com/caesarsalad/goauthz/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func AddNewRule(c *fiber.Ctx) error {
	req_rule := new(database.Rule)
	if err := c.BodyParser(req_rule); err != nil {
		return c.SendStatus(400)
	}
	var r *regexp.Regexp
	if req_rule.MetaLocationID == MetaLocationUrl {
		var err error
		r, err = regexp.Compile(req_rule.Path)
		if err != nil {
			log.Println("error while compile regex ", req_rule.Path, err)
			return c.Status(400).
				JSON(map[string]string{"error": "error while compile regex. Check your regex path"})
		}
	}
	rule := database.Rule{Path: req_rule.Path, MetaKey: req_rule.MetaKey,
		MetaLocationID: req_rule.MetaLocationID, HTTPMethodID: req_rule.HTTPMethodID}
	err := database.DB.Create(&rule).Error
	if err != nil {
		log.Print("Error while create new rule ", err)
		return c.SendStatus(500)
	}
	if r != nil {
		Rule_regex_compiled[rule.ID] = r
	}

	c.Status(201)
	return c.JSON(rule)
}

func ListAllRules(c *fiber.Ctx) error {
	var rules []database.Rule
	database.DB.Find(&rules)

	return c.JSON(rules)
}

func AssignNewRule(c *fiber.Ctx) error {
	req := new(database.AssignedRules)
	if err := c.BodyParser(req); err != nil {
		return c.SendStatus(400)
	}
	tx := database.AssignedRules{RuleID: req.RuleID, UserID: req.UserID,
		MetaValue: req.MetaValue}
	err := database.DB.Create(&tx).Error
	if err != nil {
		log.Print("Error while assign rule ", err)
		return c.SendStatus(500)
	}
	delete(Cached_user_rules.Cache, tx.UserID)
	return c.SendStatus(201)
}

func ListAssignedRules(c *fiber.Ctx) error {
	var assigned_rules []database.AssignedRules
	database.DB.Preload(clause.Associations).Find(&assigned_rules)

	return c.JSON(assigned_rules)
}

func RuleSetFile(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["file"]
	for _, file := range files {
		log.Println("file info ", file.Filename, file.Size, file.Header["Content-Type"][0])
		file_raw, _ := file.Open()
		r := bufio.NewReader(file_raw)
		file_byte_array := make([]byte, file.Size)
		r.Read(file_byte_array)
		err := RuleSetYamlParser(file_byte_array)
		if err != nil {
			return c.Status(400).JSON(map[string]string{"error": err.Error()})
		}
	}
	return c.SendStatus(201)
}
