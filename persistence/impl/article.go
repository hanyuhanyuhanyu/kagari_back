package impl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"kagari/entity"
	"kagari/persistence"
	"kagari/service"
	"kagari/service/model"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
)

func NewArticleAccessor(ctx context.Context, conn *pgx.Conn) (service.ArticleAccessor, error) {
	s3cli, err := persistence.InitS3Client(ctx)
	if err != nil {
		return nil, err
	}
	if s3cli == nil {
		return nil, errors.New("init s3 client but it somehow became null")
	}
	return &ArticleAccessor{conn, s3cli}, nil
}

type ArticleAccessor struct {
	conn  *pgx.Conn
	s3cli *s3.Client
}

func (aa *ArticleAccessor) Search(ctx context.Context, param model.ArticleSearchParameter) ([]entity.Article, error) {
	rows, err := aa.conn.Query(ctx, "SELECT id, title, url, created_at FROM articles LIMIT $1 OFFSET $2", param.GetLimit(), param.GetOffset())
	if err != nil {
		return nil, err
	}
	articles := make([]entity.Article, 0)
	for rows.Next() {
		var id, title, url string
		var createdAt time.Time
		if err := rows.Scan(&id, &title, &url, &createdAt); err != nil {
			return nil, err
		}
		articles = append(articles, entity.Article{
			ID:        id,
			Title:     title,
			URL:       url,
			CreatedAt: createdAt,
		})
	}
	return articles, nil
}
func (aa *ArticleAccessor) GetOne(ctx context.Context, id string) (*entity.Article, error) {
	rows, err := aa.conn.Query(ctx, "SELECT id, title, url, created_at FROM articles WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	var _id, title, url string
	var createdAt time.Time
	if err := rows.Scan(&id, &title, &url, &createdAt); err != nil {
		return nil, err
	}
	return &entity.Article{
		ID:        _id,
		Title:     title,
		URL:       url,
		CreatedAt: createdAt,
	}, nil
}
func createSavePath(originalFileName string) (id, keyPath string) {
	id = uuid.NewString()
	ext := filepath.Ext(originalFileName)
	keyPath = fmt.Sprintf("%s/%s_%s%s", time.Now().UTC().Format("200601"), originalFileName[:len(originalFileName)-len(ext)], id, ext)
	return

}
func (aa *ArticleAccessor) uploadImages(ctx context.Context, article *model.UploadingArticle) error {
	for key, val := range article.ImageSources {
		_, savePath := createSavePath(filepath.Base(key))
		bucket := persistence.KagariMarkdownBucket
		contentType := fmt.Sprintf("image/%s", filepath.Ext(key))
		_, err := aa.s3cli.PutObject(ctx, &s3.PutObjectInput{Bucket: &bucket, Key: &savePath, Body: val, ContentType: &contentType})
		if err != nil {
			return err
		}
		savedPath := filepath.Join(fmt.Sprintf("https://%s.s3.ap-northeast-1.amazonaws.com/", bucket), savePath)
		article.ReplaceImageSource(key, savedPath)
	}
	return nil
}
func (aa *ArticleAccessor) Upload(ctx context.Context, article *model.UploadingArticle) (string, error) {
	id, key := createSavePath(article.OriginalFileName)
	bucket := persistence.KagariMarkdownBucket
	contentType := "text/markdown; charset=UTF-8"
	tx, err := aa.conn.Begin(ctx)
	if err != nil {
		return "", err
	}
	_, err = tx.Exec(ctx, "INSERT INTO articles (id, title, url, created_at) VALUES ($1, $2, $3, $4)", id, article.Title, key, time.Now().UTC())
	if err != nil {
		tx.Rollback(ctx)
		return "", err
	}
	err = aa.uploadImages(ctx, article)
	if err != nil {
		tx.Rollback(ctx)
		return "", err
	}
	_, err = aa.s3cli.PutObject(ctx, &s3.PutObjectInput{Bucket: &bucket, Key: &key, Body: bytes.NewReader(article.Body), ContentType: &contentType})
	tx.Commit(ctx)
	return id, err
}
