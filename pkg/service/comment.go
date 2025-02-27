package service

import (
	"Blog/models"
	"Blog/pkg/repository"
)

type CommentService struct {
	repos repository.Comment
}

func (serv *CommentService) CreateComment(title, authorName string, articleId, authorId int) (int, error) {
	return serv.repos.CreateComment(authorId, articleId, authorName, title)
}

func (serv *CommentService) GetAllCommentsFromArticle(articleId int) ([]models.Comment, error) {
	return serv.repos.GetAllCommentByArticleId(articleId)
}
func NewCommentService(authorization repository.Comment) Comment {
	return &CommentService{authorization}
}
