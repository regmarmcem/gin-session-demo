package api

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/regmarcmem/gin-session-demo/controller"
	"github.com/regmarcmem/gin-session-demo/service"
	"gorm.io/gorm"
)

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

	signinCheckGroup := r.Group("/", checkSignin())
	{
		signinCheckGroup.GET("/home", func(c *gin.Context) {
			session := sessions.Default(c)
			user := session.Get("user")
			c.HTML(http.StatusOK, "home.html", gin.H{"user": user})
		})
		signinCheckGroup.GET("/signout", c.GetSignout)
	}
	signoutCheckGroup := r.Group("/", checkSignout())
	{
		signoutCheckGroup.GET("/signup", c.GetSignup)
		signoutCheckGroup.POST("/signup", c.PostSignup)
		signoutCheckGroup.GET("/signin", c.GetSignin)
		signoutCheckGroup.POST("/signin", c.PostSignin)
	}

	return r
}

func checkSignin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		_, ok := user.(string)
		if !ok {
			c.Redirect(http.StatusFound, "/signin")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func checkSignout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		_, ok := user.(string)
		if ok {
			c.Redirect(http.StatusFound, "/home")
			c.Abort()
		} else {
			c.Next()
		}
	}
}
