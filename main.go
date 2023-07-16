package main

import (
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	r := gin.Default()
	store, err := redis.NewStore(10, "tcp", redisHost+":"+redisPort, "", []byte("secret"))
	if err != nil {
		panic("failed to create redis store")
	}

	store.Options(sessions.Options{Path: "/", Domain: "localhost", MaxAge: 3600, Secure: false, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	r.Use(sessions.Sessions("mysession", store))

	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	authorized.GET("/login", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		session := sessions.Default(c)
		session.Set("user", user)
		session.Save()
	})

	r.GET("/admin/secrets", func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.Redirect(http.StatusSeeOther, "/admin/login")
		}

		userString, ok := user.(string)
		if !ok {
			c.Redirect(http.StatusSeeOther, "/admin/login")
		}
		if secret, ok := secrets[userString]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user.(string), "secret": secret})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"user": "", "secret": ""})
		}
	})

	r.Run(":8080")
}
