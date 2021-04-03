package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	JWT_secret_key        = []byte("my_secret_key")
	Migration_enabled     = false
	Authorization_enabled = false
	DB_init               = false
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}
	SetJWTSecretKey()
	SetMigrationEnabled()
	SetAuthorizationEnabled()
	SetDBInitEnabled()
}

func SetJWTSecretKey() {
	jwt_secret, is_setted := os.LookupEnv("JWT_SECRET")
	if is_setted {
		JWT_secret_key = []byte(jwt_secret)
	}
}

func SetMigrationEnabled() {
	var err error
	migration_enabled, is_setted := os.LookupEnv("MIGRATION_ENABLED")
	if is_setted {
		Migration_enabled, err = strconv.ParseBool(migration_enabled)
		if err != nil {
			log.Fatal("MIGRATION_ENABLED must be bool")
		}
	}
}

func SetAuthorizationEnabled() {
	var err error
	authorization_enabled, is_setted := os.LookupEnv("AUTHORIZATION_ENABLED")
	if is_setted {
		Authorization_enabled, err = strconv.ParseBool(authorization_enabled)
		if err != nil {
			log.Fatal("AUTHORIZATION_ENABLED must be bool")
		}
	}
}

func SetDBInitEnabled() {
	var err error
	db_init_enabled, is_setted := os.LookupEnv("DB_INIT")
	if is_setted {
		DB_init, err = strconv.ParseBool(db_init_enabled)
		if err != nil {
			log.Fatal("DB_INIT must be bool")
		}
	}
}

func GetURI() string {
	var uri string

	host := "0.0.0.0"
	port := "3000"

	host_env, is_setted := os.LookupEnv("API_HOST")
	if is_setted {
		host = host_env
	}
	port_env, is_setted := os.LookupEnv("API_PORT")
	if is_setted {
		port = port_env
	}

	uri = fmt.Sprintf("%s:%s", host, port)
	return uri
}
