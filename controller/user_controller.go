package controller

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/regmarcmem/gin-session-demo/service"
)

type UserController struct {
	s *service.UserService
}

func NewUserController(s *service.UserService) *UserController {
	return &UserController{s: s}
}

func (ctr *UserController) GetSignup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}

func (ctr *UserController) PostSignup(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := ctr.s.Signup(email, password)
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusSeeOther, "/signup")
		return
	}
	session := sessions.Default(c)
	session.Set("user", user.Email)
	err = session.Save()
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusSeeOther, "/signup")
		return
	}
	c.Redirect(http.StatusSeeOther, "/home")
}

func (ctr *UserController) GetSignin(c *gin.Context) {
	c.HTML(http.StatusOK, "signin.html", nil)
}

func (ctr *UserController) PostSignin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := ctr.s.Signin(email, password)
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusSeeOther, "/signin")
		return
	}
	session := sessions.Default(c)
	session.Set("user", user.Email)
	err = session.Save()
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusSeeOther, "/signin")
		return
	}

	c.Redirect(http.StatusSeeOther, "/home")
}

func (ctr *UserController) GetSignout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user")
	err := session.Save()
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusSeeOther, "/home")
		return
	}
	c.Redirect(http.StatusSeeOther, "/signin")
}
