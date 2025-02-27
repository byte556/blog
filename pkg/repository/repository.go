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
type Comment interface {
	CreateComment(userId, articleId int, authorName, text string) (int, error)
	GetAllCommentByArticleId(id int) ([]models.Comment, error)
	//UpdateArticle(article models.Comment) error
	//DeleteArticle(id int) error
}

type Repository struct {
	Authorization
	Article
	Comment
}

func NewRepository(db *SQLite) *Repository {
	return &Repository{db, db, db}
}
