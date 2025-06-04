package main

import (
	"Blog/microservices/article"
	"Blog/pkg/repository"
	"Blog/pkg/service"
	_ "modernc.org/sqlite"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	db := repository.NewSQLite("sqlite", "./database.db")
	repos := repository.NewRepository(db)
	services := service.NewService(repos)

	handler := article.NewHandler(services.Article)
	router := handler.InitRoutes()

	router.Run(":" + port)
}
