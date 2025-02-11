package impl

import (
	"context"
	"errors"
	"kagari/dataaccessor"
	"kagari/entity"
	"kagari/service"
	"kagari/service/model"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/samber/lo"
)

func NewArticleAccessor(ctx context.Context, driver neo4j.DriverWithContext) (service.ArticleAccessor, error) {
	s3cli, err := dataaccessor.InitS3Client(ctx)
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
	result, err := neo4j.ExecuteQuery(ctx, aa.driver, "MATCH (a:Article) ORDER BY a.created_at LIMIT $limit SKIP $offset RETURN a", map[string]any{
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
	first := result.Records[0]
	if first == nil {
		return nil, nil
	}
	return (&entity.Article{}).FromMap(first.AsMap()["a"].(dbtype.Node).GetProperties()), nil
}
func (aa *ArticleAccessor) Upload(ctx context.Context, article *model.UploadingArticle) (string, error) {
	id := uuid.NewString()
	bucket := dataaccessor.KagariMarkdownBucket
	key := article.CreateSavePath(id)
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
		_, err = aa.s3cli.PutObject(ctx, &s3.PutObjectInput{Bucket: &bucket, Key: &key, Body: article.Body, ContentType: &contentType})
		return nil, err
	})
	return id, err
}
