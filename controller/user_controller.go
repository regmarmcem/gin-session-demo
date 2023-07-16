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
	}
	session := sessions.Default(c)
	session.Set("user", user)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/home")
}
