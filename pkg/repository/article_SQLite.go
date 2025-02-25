package repository

import (
	"Blog/models"
	"time"
)

func (a *SQLite) CreateArticle(title, body string, authorId int, url string) (int, error) {
	query := "INSERT INTO articles (author_id, title, body, cover_url, created_at,  is_posted) VALUES (?, ?, ?, ?, ?,?)"
	e, err := a.db.Exec(query, authorId, title, body, url, time.Now().UTC(), true)
	if err != nil {
		return 0, err
	}
	id, err := e.LastInsertId()
	return int(id), err
}

func (a *SQLite) GetArticleById(id int) (models.Article, error) {
	var article models.Article
	query := "SELECT id, author_id, title, body, cover_url, is_posted, created_at FROM articles WHERE id = ?"
	err := a.db.QueryRow(query, id).Scan(&article.Id, &article.AuthorId, &article.Title, &article.Body, &article.CoverURL, &article.IsPosted, &article.CreatedAt)
	return article, err
}

func (a *SQLite) GetAllArticles() ([]models.Article, error) {
	query := "SELECT id, author_id, title, body, cover_url, is_posted, created_at FROM articles"
	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.Id, &article.AuthorId, &article.Title, &article.Body, &article.CoverURL, &article.IsPosted, &article.CreatedAt); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (a *SQLite) UpdateArticle(article models.Article) error {
	query := "UPDATE articles SET title = ?, body = ?, is_posted = ?, cover_url =? WHERE id = ?"
	_, err := a.db.Exec(query, article.Title, article.Body, article.IsPosted, article.CoverURL, article.Id)
	return err
}

func (a *SQLite) DeleteArticle(id int) error {
	query := "DELETE FROM articles WHERE id = ?"
	_, err := a.db.Exec(query, id)
	return err
}

func (a *SQLite) GetLastArticles(count int) ([]models.Article, error) {
	query := "SELECT id, author_id, title, body, cover_url, is_posted, created_at FROM articles ORDER BY created_at DESC LIMIT ?"
	rows, err := a.db.Query(query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.Id, &article.AuthorId, &article.Title, &article.Body, &article.CoverURL, &article.IsPosted, &article.CreatedAt); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}
