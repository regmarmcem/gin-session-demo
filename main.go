package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func main() {

	r := gin.Default()

	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	authorized.GET("/login", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("user", user, 3600, "/", "localhost", false, true)
	})

	r.GET("/admin/secrets", func(c *gin.Context) {
		user, err := c.Cookie("user")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/admin/login")
		}

		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		}
	})

	r.Run("localhost:8080")
}
