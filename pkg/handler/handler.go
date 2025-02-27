package handler

import (
	"Blog/pkg/service"
	"Blog/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/utrack/gin-csrf"
	"html/template"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	// Настройка cookie-сессий
	store := cookie.NewStore([]byte("top-secret"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   12 * 60 * 60,
		Secure:   false,
		HttpOnly: true,
	})

	// Добавляем middleware для сессий ПЕРЕД CSRF
	router.Use(sessions.Sessions("session", store))

	// Настройка CSRF
	router.Use(csrf.Middleware(csrf.Options{
		Secret: "your-csrf-secret",
		ErrorFunc: func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{"error": "CSRF token mismatch"})
			c.Abort()
			return
		},
	}))

	// Подключаем наши функции для шаблонов
	router.SetFuncMap(template.FuncMap{
		"stripHTML": utils.StripHTMLTags,
		"shorten":   utils.Shorten,
	})

	// Грузим шаблоны
	router.LoadHTMLGlob("static/templates/*")
	// Статичные файлы
	router.Static("/public", "./public")

	// Ваш middleware
	router.Use(h.UserMiddleware)

	// Роуты
	router.GET("/", h.homePage)
	router.GET("/sign-in", h.signInPage)
	router.POST("/sign-in", h.signIn)
	router.GET("/sign-up", h.signUpPage)
	router.POST("/sign-up", h.signUp)
	router.POST("/logout", h.logout)

	post := router.Group("/article")
	{
		post.POST("/preview", h.articlePreview)
		post.POST("/create", h.articleCreate)
		post.POST("/:id/create", h.commentCreate) // Исправлено!
		post.GET("/add", h.AuthRequired, h.articleAddPage)
		post.GET("/:id", h.articlePage)
	}
	return router
}
