package handler

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"Blog/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

// ArticlePreviewRequest – структура для валидации запроса превью статьи.
type ArticlePreviewRequest struct {
	Markdown string `json:"markdown" binding:"required"`
}

// ArticleForm – структура для валидации формы создания статьи.
type ArticleForm struct {
	Title            string `form:"title" binding:"required"`
	Body             string `form:"body" binding:"required"`
	IdempotencyToken string `form:"idempotencyToken" binding:"required"`
}

// articlePreview обрабатывает запрос на предварительный просмотр Markdown.
func (h *Handler) articlePreview(c *gin.Context) {
	var req ArticlePreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неверный JSON",
			"details": err.Error(),
		})
		return
	}

	htmlContent, err := h.service.ConvertMarkdownToHTML(req.Markdown)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Не удалось преобразовать Markdown",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"html": htmlContent})
}

// articleAddPage рендерит страницу создания статьи и генерирует idempotency-токен.
func (h *Handler) articleAddPage(c *gin.Context) {
	user, err := h.getUser(c)
	if err != nil {
		c.HTML(http.StatusOK, "error.html", gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	// Генерируем токен и сохраняем его в сессии.
	idempotencyToken := utils.GenerateToken()
	session := sessions.Default(c)
	session.Set("idempotencyToken", idempotencyToken)
	if err := session.Save(); err != nil {
		log.Printf("Ошибка сохранения сессии: %v", err)
	}

	c.HTML(http.StatusOK, "add_post.html", gin.H{
		"user":             user,
		"csrf_token":       csrf.GetToken(c),
		"idempotencyToken": idempotencyToken,
	})
}

// articleCreate создаёт статью с проверкой idempotency-токена и валидацией входных данных.
func (h *Handler) articleCreate(c *gin.Context) {
	// Валидация формы создания статьи.
	var form ArticleForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неверные входные данные",
			"details": err.Error(),
		})
		return
	}

	user, err := h.getUser(c)
	if err != nil {
		c.HTML(http.StatusOK, "error.html", gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.HTML(http.StatusOK, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}

	// Проверка idempotency-токена.
	session := sessions.Default(c)
	sessionToken := session.Get("idempotencyToken")
	if !utils.ValidateIdempotencyToken(sessionToken, form.IdempotencyToken) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка токена"})
		return
	}
	// Удаляем токен для предотвращения повторного использования.
	session.Delete("idempotencyToken")
	session.Save()

	// Получаем файл обложки.
	coverFile, err := c.FormFile("cover")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Вы должны передать обложку",
			"details": err.Error(),
		})
		return
	}

	// Передаём валидированные данные в сервис для создания статьи.
	articleID, err := h.service.CreateArticle(form.Title, form.Body, user.Id, coverFile)
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

// articlePage отображает страницу отдельной статьи.
func (h *Handler) articlePage(c *gin.Context) {
	user, err := h.getUser(c)
	if err != nil {
		c.HTML(http.StatusOK, "error.html", gin.H{"error": err.Error()})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Неверный ID"})
		return
	}

	article, err := h.service.GetArticleByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Статья не найдена"})
		return
	}

	// Преобразуем Markdown статьи в HTML.
	article.Body, err = h.service.ConvertMarkdownToHTML(article.Body)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	bodyHTML := template.HTML(article.Body)

	comments, err := h.service.GetAllCommentsFromArticle(article.Id)
	token := csrf.GetToken(c)
	c.HTML(http.StatusOK, "view_post.html", gin.H{
		"Article":    article,
		"csrf_token": token,
		"BodyHTML":   bodyHTML,
		"user":       user,
		"comments":   comments,
	})
}
