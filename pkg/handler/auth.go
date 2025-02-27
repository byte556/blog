package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

// SignInForm – структура для валидации данных формы входа.
type SignInForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// SignUpForm – структура для валидации данных формы регистрации.
type SignUpForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// signInPage рендерит страницу входа с CSRF-токеном.
func (h *Handler) signInPage(c *gin.Context) {
	user, err := h.getUser(c)
	if err != nil {
		c.HTML(http.StatusOK, "error.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "sign-in.html", gin.H{
		"user":       user,
		"csrf_token": csrf.GetToken(c),
	})
}

// signIn обрабатывает вход пользователя с валидацией данных.
func (h *Handler) signIn(c *gin.Context) {
	var form SignInForm
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "sign-in.html", gin.H{
			"error":   "Неверные входные данные",
			"details": err.Error(),
		})
		return
	}

	id, err := h.service.SignIn(form.Username, form.Password)
	if err != nil {
		c.HTML(http.StatusOK, "sign-in.html", gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)
	session.Set("id", id)
	session.Save()
	c.Redirect(http.StatusFound, "/")
}

// signUpPage рендерит страницу регистрации с CSRF-токеном.
func (h *Handler) signUpPage(c *gin.Context) {
	user, err := h.getUser(c)
	if err != nil {
		c.HTML(http.StatusOK, "error.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "sign-up.html", gin.H{
		"user":       user,
		"csrf_token": csrf.GetToken(c),
	})
}

// signUp обрабатывает регистрацию с валидацией входных данных.
func (h *Handler) signUp(c *gin.Context) {
	var form SignUpForm
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "sign-up.html", gin.H{
			"error":   "Неверные входные данные",
			"details": err.Error(),
		})
		return
	}

	err := h.service.SignUp(form.Username, form.Password)
	if err != nil {
		c.HTML(http.StatusOK, "sign-up.html", gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/sign-in")
}

// logout завершает сессию и перенаправляет пользователя на главную страницу.
func (h *Handler) logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/")
}
