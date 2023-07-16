package main

import (
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/joho/godotenv"
	"github.com/regmarcmem/gin-session-demo/api"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	store, err := redis.NewStore(10, "tcp", redisHost+":"+redisPort, "", []byte("secret"))
	if err != nil {
		panic("failed to create redis store")
	}
	store.Options(sessions.Options{Path: "/", Domain: "localhost", MaxAge: 3600, Secure: false, HttpOnly: true, SameSite: http.SameSiteLaxMode})

	r := api.NewRouter()
	r.Use(sessions.Sessions("mysession", store))

	r.Run(":8080")
}
