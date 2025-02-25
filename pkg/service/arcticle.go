package service

import (
	"Blog/models"
	"Blog/pkg/repository"
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday" // для санитизации
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func (serv *ArticleService) GetArticleByID(id int) (models.Article, error) {
	return serv.repos.GetArticleById(id)
}
func (serv *ArticleService) CreateArticle(title, body string, authorId int, file *multipart.FileHeader) (int, error) {
	fileName, err := serv.saveCoverFile(file)
	if err != nil {
		return 0, err
	}
	markdownBody, err := serv.ConvertMarkdownToHTML(body)
	if err != nil {
		return 0, err
	}
	return serv.repos.CreateArticle(title, markdownBody, authorId, fileName)
}
func (serv *ArticleService) saveCoverFile(file *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	if !allowedExts[ext] {
		return "", fmt.Errorf("файлы типа %s не разрешены", ext)
	}

	uniqueName := uuid.New().String() + ext
	coverPath := "./public/" + uniqueName

	// Открываем исходный файл
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(coverPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Копируем содержимое
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return coverPath, nil
}
func (serv *ArticleService) DeleteArticle() {
	return
}
func (serv *ArticleService) UpdateArticle() {
	return
}
func (serv *ArticleService) ConvertMarkdownToHTML(markdown string) (string, error) {
	// Создаем новый Goldmark, указывая GFM-расширения
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,     // Добавляет поддержку таблиц, strikethrough, autolink и т.д.
			extension.Linkify, // Преобразует ссылки

		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(), // Переносы строк
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return "", err // или models.ErrConvert
	}

	rawHTML := buf.String()

	// Bluemonday санитизирует потенциально опасный HTML
	policy := bluemonday.UGCPolicy()
	safeHTML := policy.Sanitize(rawHTML)

	return safeHTML, nil
}
func (serv *ArticleService) GetLastArticles(count int) ([]models.Article, error) {
	return serv.repos.GetLastArticles(count)
}

type ArticleService struct {
	repos repository.Article
}

func NewArticleService(article repository.Article) *ArticleService {
	return &ArticleService{article}
}
