package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type CommentForm struct {
	Text string `form:"text" binding:"required"`
}

func (h *Handler) commentCreate(c *gin.Context) {
	articleId := c.Param("id")

	form := new(CommentForm)
	if err := c.ShouldBind(&form); err != nil {
		log.Printf("Ошибка привязки формы: %v", err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Ошибка в форме комментария",
		})
		return
	}

	user, err := h.getUser(c)
	if err != nil {
		log.Printf("Ошибка получения пользователя: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка сервера"})
		return
	}
	if user == nil {
		log.Printf("Пользователь не авторизован")
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}

	id, err := strconv.Atoi(articleId)
	if err != nil {
		log.Printf("Неверный ID статьи: %v", err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Неверный ID статьи"})
		return
	}

	_, err = h.service.CreateComment(form.Text, user.Name, id, user.Id)
	if err != nil {
		log.Printf("Ошибка создания комментария: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Не удалось создать комментарий"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/article/"+articleId)
}
