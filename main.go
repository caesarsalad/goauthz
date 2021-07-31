package main

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/caesarsalad/goauthz/authorization"
	"github.com/caesarsalad/goauthz/config"
	"github.com/caesarsalad/goauthz/cutils"
	"github.com/caesarsalad/goauthz/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" validate:"required"`
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func Login(c *fiber.Ctx) error {
	req_user := new(User)
	if err := c.BodyParser(req_user); err != nil {
		return c.SendStatus(400)
	}
	validation_err := cutils.ValidateStruct(*req_user)
	if validation_err != nil {
		return c.Status(400).JSON(validation_err)
	}
	var user database.User
	err := database.DB.Where("username = ?", req_user.Username).Or("email = ?", req_user.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(401).JSON(map[string]string{"message": "username or password wrong"})
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req_user.Password))
	if err != nil {
		return c.Status(401).JSON(map[string]string{"message": "username or password wrong"})
	}

	claims := &Claims{
		Username: user.Username,
		Email:    user.Email,
		UserID:   user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(config.JWT_secret_key)
	if err != nil {
		return c.SendStatus(500)
	}
	return c.JSON(map[string]string{"token": tokenString})
}

func Register(c *fiber.Ctx) error {
	new_user := new(User)
	if err := c.BodyParser(new_user); err != nil {
		return c.SendStatus(400)
	}
	var user database.User
	err := database.DB.Where("username = ?", new_user.Username).Or("email = ?", new_user.Email).First(&user).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}
	if user.ID > 0 {
		return c.Status(400).JSON(map[string]string{"message": "username or email already exist"})
	}
	bytes, _ := bcrypt.GenerateFromPassword([]byte(new_user.Password), 10)
	hashed_password := string(bytes)
	new_user.Password = hashed_password
	database.DB.Create(&new_user)
	return c.SendStatus(200)
}

func main() {
	app := fiber.New()

	app.Use(requestid.New(requestid.Config{
		Generator: utils.UUIDv4,
	},
	))
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${time} ${locals:requestid} ${ip} ${status} - ${method} ${path}​\n​",
	}))

	app.Use("/", func(c *fiber.Ctx) error {
		switch c.Path() {
		case "/login":
			{
				return c.Next()
			}
		case "/register":
			{
				return c.Next()
			}
		}
		token_header := string(c.Request().Header.Peek("Authorization"))
		token := strings.Split(token_header, "Bearer")
		if len(token) != 2 {
			return c.SendStatus(401)
		}
		tknStr := strings.TrimSpace(token[1])
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JWT_secret_key, nil
		})
		if err != nil {
			return c.SendStatus(401)
		}
		if !tkn.Valid {
			return c.SendStatus(401)
		}
		if config.Authorization_enabled {
			if strings.HasPrefix(c.Path(), "/api/goauthz/internal/") {
				return c.Next()
			}
			if !authorization.CheckRules(c, claims.UserID) {
				return c.SendStatus(403)
			}
		}
		return c.SendStatus(200)
	})

	app.Post("/login", Login)
	app.Post("/register", Register)

	internal_api := app.Group("/api/goauthz/internal/")
	internal_api.Get("/rule", authorization.ListAllRules)
	internal_api.Post("/rule", authorization.AddNewRule)
	internal_api.Get("/assigned_rules", authorization.ListAssignedRules)
	internal_api.Post("/assign_rule", authorization.AssignNewRule)
	internal_api.Post("/rule_set_file", authorization.RuleSetFile)
	internal_api.Get("/rule_set_file", authorization.RuleSetDump)

	database.ConnectDB()
	if config.Migration_enabled {
		log.Println("Auto Migration Enabled")
		database.DB.AutoMigrate(&database.User{},
			&database.Rule{},
			&database.AssignedRules{})
		if config.DB_init {
			authorization.InitStaticTypesDB()
		}
	}

	go func() {
		for {
			authorization.CompileAllRegexRules()
			authorization.ReCacheUserRules()
			time.Sleep(30 * time.Second)
		}
	}()

	uri := config.GetURI()
	app.Listen(uri)
}
