package config

import (
	"fmt"
	"os"
)

var JWT_secret_key = []byte("my_secret_key")

func SetJWTSecretKey() {
	jwt_secret, is_setted := os.LookupEnv("JWT_SECRET")
	if is_setted {
		JWT_secret_key = []byte(jwt_secret)
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
