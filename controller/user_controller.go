package controller

import (
	"net/http"

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
	email := c.PostForm("inputEmail")
	password := c.PostForm("inputPassword")

	user, err := ctr.s.Signup(email, password)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/signup")
		return
	}
	c.HTML(http.StatusOK, "home.html", gin.H{"user": user})
}
