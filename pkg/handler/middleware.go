package handler

import (
	"Blog/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (h *Handler) UserMiddleware(context *gin.Context) {

	user, _ := h.getUser(context) // Получаем пользователя

	context.Set("user", user)

	context.Next()

}

func (h *Handler) getUser(context *gin.Context) (*models.User, error) {
	session := sessions.Default(context)
	idValue := session.Get("id") // Получаем ID из сессии

	if idValue == nil {
		return nil, nil // Нет сессии, но это не ошибка
	}

	// Приводим `idValue` к int (так как оно может быть float64)
	if idValue == "" {
		return nil, nil // Нет сессии, но это не ошибка
	}
	id, ok := idValue.(int)
	if !ok {
		return nil, nil
	}

	// Проверяем, что ID > 0 (валидный пользователь)
	if id <= 0 {
		return nil, nil
	}

	// Получаем пользователя из базы
	user, err := h.service.GetUser(id)
	if err != nil {
		return nil, nil
	}

	return &user, nil
}
