package repository

import (
	"Blog/models"
	"time"
)

/*
CREATE TABLE IF NOT EXISTS comments (

	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL UNIQUE,
	article_id INTEGER NOT NULL,
	text TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE

);`
*/
func (a *SQLite) CreateComment(userId, articleId int, authorName, text string) (int, error) {
	query := "INSERT INTO comments (user_id, article_id, author_name, text, created_at) VALUES (?, ?, ?, ?, ?)"
	e, err := a.db.Exec(query, userId, articleId, authorName, text, time.Now().UTC())
	if err != nil {
		return 0, err
	}
	id, err := e.LastInsertId()
	return int(id), err
}
func (a *SQLite) GetAllCommentByArticleId(articleId int) ([]models.Comment, error) {
	query := "SELECT id, user_id, article_id, author_name, text, created_at FROM comments WHERE article_id = ?"
	rows, err := a.db.Query(query, articleId)
	if err != nil {
		return []models.Comment{}, err
	}
	defer rows.Close() // Добавим закрытие rows для избежания утечек

	ret := []models.Comment{}
	for rows.Next() {
		var comment models.Comment
		err = rows.Scan(&comment.Id, &comment.AuthorId, &comment.PostId, &comment.AuthorName, &comment.Text, &comment.CreatedAt)
		if err != nil {
			return []models.Comment{}, err
		}
		ret = append(ret, comment)
	}

	// Проверяем ошибки после итерации
	if err = rows.Err(); err != nil {
		return []models.Comment{}, err
	}

	return ret, nil
}
