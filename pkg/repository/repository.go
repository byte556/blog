package repository

import (
	"Blog/models"
)

type Authorization interface {
	GetUserById(id int) (models.User, error)
	GetUserByName(username string) (models.User, error)
	AddUser(user *models.User) error
}
type Article interface {
	CreateArticle(title, body string, authorId int, name string) (int, error)
	GetArticleById(id int) (models.Article, error)
	GetAllArticles() ([]models.Article, error)
	UpdateArticle(article models.Article) error
	DeleteArticle(id int) error
	GetLastArticles(count int) ([]models.Article, error)
}

type Repository struct {
	Authorization
	Article
}

func NewRepository(db *SQLite) *Repository {
	return &Repository{db, db}
}
