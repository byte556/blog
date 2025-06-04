package article

import (
	"net/http"
	"strconv"

	"Blog/pkg/service"

	"github.com/gin-gonic/gin"
)

// Handler provides HTTP handlers for article microservice.
type Handler struct {
	service service.Article
}

func NewHandler(s service.Article) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()
	r.GET("/health", h.health)
	r.GET("/articles", h.getLastArticles)
	r.GET("/articles/:id", h.getArticleByID)
	return r
}

func (h *Handler) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) getLastArticles(c *gin.Context) {
	countStr := c.DefaultQuery("count", "10")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid count"})
		return
	}
	articles, err := h.service.GetLastArticles(count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, articles)
}

func (h *Handler) getArticleByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	article, err := h.service.GetArticleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, article)
}
