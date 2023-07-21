package api

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/regmarcmem/gin-session-demo/controller"
	"github.com/regmarcmem/gin-session-demo/service"
	"gorm.io/gorm"
)

var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func NewRouter(db *gorm.DB, store sessions.Store) *gin.Engine {
	s := service.NewUserService(db)
	c := controller.NewUserController(s)

	r := gin.Default()
	r.Use(sessions.Sessions("mysession", store))

	r.LoadHTMLGlob("static/*.html")
	r.Static("/assets", "./static/assets")
	r.Static("/dist", "./static/dist")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	loginCheckGroup := r.Group("/", checkLogin())
	{
		loginCheckGroup.GET("/home", func(c *gin.Context) {
			session := sessions.Default(c)
			user := session.Get("user")
			c.HTML(http.StatusOK, "home.html", gin.H{"user": user})
		})
	}
	logoutCheckGroup := r.Group("/", checkLogout())
	{
		logoutCheckGroup.GET("/signup", c.GetSignup)
		logoutCheckGroup.POST("/signup", c.PostSignup)
		logoutCheckGroup.GET("/signin", c.GetSignin)
		logoutCheckGroup.POST("/signin", c.PostSignin)
	}

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

func checkLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.Redirect(http.StatusFound, "/signin")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func checkLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user != nil {
			c.Redirect(http.StatusFound, "/home")
			c.Abort()
		} else {
			c.Next()
		}
	}
}
