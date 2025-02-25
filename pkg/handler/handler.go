package handler

import (
	"Blog/pkg/service"
	"Blog/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"html/template"
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
	router.Use(sessions.Sessions("session", store))

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

		post.GET("/add", h.articleAddPage)
		post.GET("/:id", h.articlePage)
	}

	return router
}
