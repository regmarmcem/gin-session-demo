package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/joho/godotenv"
	"github.com/regmarcmem/gin-session-demo/api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	store, err := redis.NewStore(10, "tcp", os.Getenv("REDIS_HOST")+":"+os.Getenv("REDIS_PORT"), "", []byte("secret"))
	if err != nil {
		panic("failed to create redis store")
	}
	store.Options(sessions.Options{Path: "/", Domain: "localhost", MaxAge: 3600, Secure: false, HttpOnly: true, SameSite: http.SameSiteLaxMode})

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Tokyo", os.Getenv("DB_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	r := api.NewRouter(db)
	r.Use(sessions.Sessions("mysession", store))

	r.Run(":8080")
}
