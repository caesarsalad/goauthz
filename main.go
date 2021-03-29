package main

import (
	"errors"
	"log"
	"strings"

	"github.com/caesarsalad/goauthz/config"
	"github.com/caesarsalad/goauthz/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func Login(c *fiber.Ctx) error {
	req_user := new(database.User)
	if err := c.BodyParser(req_user); err != nil {
		return c.SendStatus(400)
	}
	var user database.User
	err := database.DB.Where("username = ?", req_user.Username).Or("email = ?", req_user.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.Status(401)
		return c.JSON(map[string]string{"message": "username or password wrong"})
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req_user.Password))
	if err != nil {
		c.Status(401)
		return c.JSON(map[string]string{"message": "username or password wrong"})
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
	new_user := new(database.User)
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
		c.Status(400)
		return c.JSON(map[string]string{"message": "username or email already exist"})
	}
	bytes, _ := bcrypt.GenerateFromPassword([]byte(new_user.Password), 10)
	hashed_password := string(bytes)
	new_user.Password = hashed_password
	database.DB.Create(&new_user)
	return c.SendStatus(200)
}

func main() {
	app := fiber.New()

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
		return c.SendStatus(200)
	})

	app.Post("/login", Login)
	app.Post("/register", Register)

	database.ConnectDB()
	database.DB.AutoMigrate(&database.User{})

	uri := config.GetURI()
	app.Listen(uri)
}
