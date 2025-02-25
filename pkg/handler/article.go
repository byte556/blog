package handler

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) articlePreview(c *gin.Context) {
	// 1) Получаем JSON {"markdown": "..."}
	var req struct {
		Markdown string `json:"markdown"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// 2) Конвертируем Markdown -> HTML через сервис
	html, err := h.service.ConvertMarkdownToHTML(req.Markdown)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert"})
		return
	}

	// 3) Возвращаем JSON: { "html": "..."}
	c.JSON(http.StatusOK, gin.H{"html": html})
}

func (h *Handler) articleAddPage(c *gin.Context) {
	user, err := h.getUser(c)
	if err != nil {
		c.HTML(http.StatusOK, "error.html", gin.H{"error": err.Error()})
		return
	}
	// Рендерим шаблон "add_post.html", передавая user
	c.HTML(http.StatusOK, "add_post.html", gin.H{"user": user})
}
func (h *Handler) articleCreate(c *gin.Context) {
	user, err := h.getUser(c)
	if err != nil {
		c.HTML(http.StatusOK, "error.html", gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.HTML(http.StatusOK, ".html", gin.H{"user": nil}) //you need login
	}
	title := c.PostForm("title")
	body := c.PostForm("body")

	if title == "" || body == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Заголовок и текст статьи не должны быть пустыми",
		})
		return
	}

	coverFile, err := c.FormFile("cover")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Вы должны передать обложку",
		})
		return
	}

	articleID, err := h.service.CreateArticle(title, body, user.Id, coverFile)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Не удалось сохранить файл: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":     "/article/" + strconv.Itoa(articleID),
		"message": "Статья успешно создана",
	})
}

func (h *Handler) articlePage(context *gin.Context) {
	user, err := h.getUser(context)
	if err != nil {
		context.HTML(http.StatusOK, "error.html", gin.H{"error": err.Error()})
	}
	// id из URL: /article/:id
	idStr := context.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		context.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Неверный ID"})
		return
	}

	// Получаем статью
	article, err := h.service.GetArticleByID(id)
	if err != nil {
		context.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Статья не найдена"})
		return
	}

	// Если article.Body уже содержит HTML (после конвертации из Markdown + санитизация),
	// делаем template.HTML, чтобы Go-шаблоны не экранировали теги.
	bodyHTML := template.HTML(article.Body)
	// Рендерим шаблон
	context.HTML(http.StatusOK, "view_post.html", gin.H{
		"Article":  article,
		"BodyHTML": bodyHTML,
		"user":     user,
	})

}
