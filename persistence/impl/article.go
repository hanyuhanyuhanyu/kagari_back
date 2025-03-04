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
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/samber/lo"
)

func NewArticleAccessor(ctx context.Context, driver neo4j.DriverWithContext) (service.ArticleAccessor, error) {
	s3cli, err := persistence.InitS3Client(ctx)
	if err != nil {
		return nil, err
	}
	if s3cli == nil {
		return nil, errors.New("init s3 client but it somehow became null")
	}
	return &ArticleAccessor{driver, s3cli}, nil
}

type ArticleAccessor struct {
	driver neo4j.DriverWithContext
	s3cli  *s3.Client
}

func (aa *ArticleAccessor) Search(ctx context.Context, param model.ArticleSearchParameter) ([]entity.Article, error) {
	result, err := neo4j.ExecuteQuery(ctx, aa.driver, "MATCH (a:Article) ORDER BY a.created_at DESC LIMIT $limit SKIP $offset RETURN a", map[string]any{
		"limit":  param.GetLimit(),
		"offset": param.GetOffset(),
	}, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	first := result.Records[0]
	if first == nil {
		return nil, nil
	}
	return lo.Map(result.Records, func(r *neo4j.Record, _ int) entity.Article {
		return *(&entity.Article{}).FromMap(r.AsMap()["a"].(dbtype.Node).GetProperties())
	}), nil
}
func (aa *ArticleAccessor) GetOne(ctx context.Context, id string) (*entity.Article, error) {
	result, err := neo4j.ExecuteQuery(ctx, aa.driver, "MATCH (a:Article {id: $id}) RETURN a", map[string]any{
		"id": id,
	}, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	if len(result.Records) == 0 {
		return nil, nil
	}
	first := result.Records[0]
	return (&entity.Article{}).FromMap(first.AsMap()["a"].(dbtype.Node).GetProperties()), nil
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
	sess := aa.driver.NewSession(ctx, neo4j.SessionConfig{})
	_, err := sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, "MERGE (a:Article {id: $id, url: $url, title: $title, created_at: $created_at})", map[string]any{
			"id":         id,
			"title":      article.Title,
			"url":        key,
			"created_at": time.Now().UTC(),
		})
		if err != nil {
			return nil, err
		}
		err = aa.uploadImages(ctx, article)
		if err != nil {
			return nil, err
		}
		_, err = aa.s3cli.PutObject(ctx, &s3.PutObjectInput{Bucket: &bucket, Key: &key, Body: bytes.NewReader(article.Body), ContentType: &contentType})
		return nil, err
	})
	return id, err
}
