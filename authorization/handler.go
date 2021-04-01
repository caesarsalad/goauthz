package authorization

import (
	"log"

	"github.com/caesarsalad/goauthz/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func AddNewRule(c *fiber.Ctx) error {
	req_rule := new(database.Rule)
	if err := c.BodyParser(req_rule); err != nil {
		return c.SendStatus(400)
	}
	rule := database.Rule{Path: req_rule.Path, MetaKey: req_rule.MetaKey,
		MetaLocationID: req_rule.MetaLocationID}
	err := database.DB.Create(&rule).Error
	if err != nil {
		log.Print("Error while create new rule ", err)
		return c.SendStatus(500)
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
	return c.SendStatus(201)
}

func ListAssignedRules(c *fiber.Ctx) error {
	var assigned_rules []database.AssignedRules
	database.DB.Preload(clause.Associations).Find(&assigned_rules)

	return c.JSON(assigned_rules)
}
