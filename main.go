package main

import (
	"Blog/pkg/handler"
	"Blog/pkg/repository"
	"Blog/pkg/service"
	_ "modernc.org/sqlite"
)

func main() {
	db := repository.NewSQLite("sqlite", "./database.db")
	rep := repository.NewRepository(db)
	serv := service.NewService(rep)
	handl := handler.NewHandler(serv)
	engine := handl.InitRoutes()
	err := engine.Run(":8080")
	if err != nil {
		return
	}
}
