package service

import (
	"Blog/models"
	"Blog/pkg/repository"
	"mime/multipart"
)

type Authorization interface {
	SignIn(username, password string) (int, error)
	SignUp(username, password string) error
	GetUser(id int) (models.User, error)
}

type Article interface {
	ConvertMarkdownToHTML(markdown string) (string, error)
	GetArticleByID(id int) (models.Article, error)
	CreateArticle(title, body string, authorId int, file *multipart.FileHeader) (int, error)
	DeleteArticle()
	UpdateArticle()
	GetLastArticles(count int) ([]models.Article, error)
}

type Service struct {
	Authorization
	Article
}

func NewService(repos *repository.Repository) *Service {
	return &Service{Authorization: NewAuthService(repos), Article: NewArticleService(repos)}
}
