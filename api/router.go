package api

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("static/*.html")
	r.Static("/assets", "./static/assets")
	r.Static("/dist", "./static/dist")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/signin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signin.html", nil)
	})

	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	authorized.GET("/signin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		session := sessions.Default(c)
		session.Set("user", user)
		session.Save()
	})

	r.GET("/admin/secrets", func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.Redirect(http.StatusSeeOther, "/admin/signin")
		}

		userString, ok := user.(string)
		if !ok {
			c.Redirect(http.StatusSeeOther, "/admin/signin")
		}
		if secret, ok := secrets[userString]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user.(string), "secret": secret})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"user": "", "secret": ""})
		}
	})

	return r
}