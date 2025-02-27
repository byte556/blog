package handler

import (
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
)

func (h *Handler) homePage(context *gin.Context) {
	user, err := h.getUser(context)
	if err != nil {
		context.HTML(http.StatusOK, "error.html", gin.H{"error": err.Error()})
		return
	}

	posts, err := h.service.GetLastArticles(10)

	context.HTML(http.StatusOK, "index.html", gin.H{"user": user, "posts": posts, "csrf_token": csrf.GetToken(context)})
}
