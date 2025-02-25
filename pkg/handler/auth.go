package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) logout(context *gin.Context) {
	session := sessions.Default(context)
	session.Clear()
	session.Save()
	context.Redirect(http.StatusFound, "/")

}
func (h *Handler) signInPage(context *gin.Context) {
	user, err := h.getUser(context)
	if err != nil {
		context.HTML(http.StatusOK, "error.html", gin.H{"error": err.Error()})
		return
	}
	context.HTML(http.StatusOK, "sign-in.html", gin.H{"user": user})
}

func (h *Handler) signIn(context *gin.Context) {

	username := context.PostForm("username")
	password := context.PostForm("password")

	id, err := h.service.SignIn(username, password)
	if err != nil {
		context.HTML(http.StatusOK, "sign-in.html", gin.H{"error": err.Error()})
		return
	}
	session := sessions.Default(context)
	session.Set("id", id)

	session.Save()
	context.Redirect(http.StatusFound, "/")
}

func (h *Handler) signUpPage(context *gin.Context) {
	user, err := h.getUser(context)
	if err != nil {
		context.HTML(http.StatusOK, "error.html", gin.H{"error": err.Error()})
		return
	}
	context.HTML(http.StatusOK, "sign-up.html", gin.H{"user": user})

}

func (h *Handler) signUp(context *gin.Context) {
	username := context.PostForm("username")
	password := context.PostForm("password")

	err := h.service.SignUp(username, password)
	if err != nil {
		context.HTML(http.StatusOK, "sign-in.html", gin.H{"error": err.Error()})
		return
	}
	context.Redirect(http.StatusFound, "/sign-in")
}
