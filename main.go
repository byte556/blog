package main

import (
	"Blog/pkg/handler"
	"Blog/pkg/repository"
	"Blog/pkg/service"
	_ "modernc.org/sqlite"
	"os"
)

func main() {
	port := os.Getenv("PORT") // Получаем порт из окружения
	if port == "" {
		port = "8080" // Значение по умолчанию
	}
	db := repository.NewSQLite("sqlite", "./database.db")
	rep := repository.NewRepository(db)
	serv := service.NewService(rep)
	handl := handler.NewHandler(serv)
	engine := handl.InitRoutes()
	err := engine.Run(":" + port)
	if err != nil {
		return
	}
}
