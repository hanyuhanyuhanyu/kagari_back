package service

import (
	"context"
	"kagari/entity"
)

type ArticleAccessor interface {
	Search(ctx context.Context, condition entity.ArticleSearchParameter) ([]entity.Article, error)
	GetOne(ctx context.Context, id string) (*entity.Article, error)
	Upload(ctx context.Context, article *entity.UploadingArticle) (id string, err error)
}
type ArticleService struct {
	accessor ArticleAccessor
}

func NewArticleService(accessor ArticleAccessor) *ArticleService {
	return &ArticleService{accessor: accessor}
}

func (as *ArticleService) GetArticle(ctx context.Context, id string) (*entity.Article, error) {
	return as.accessor.GetOne(ctx, id)
}
func (as *ArticleService) Post(ctx context.Context, article *entity.UploadingArticle) (id string, err error) {
	return as.accessor.Upload(ctx, article)
}
func (as *ArticleService) Search(ctx context.Context, param entity.ArticleSearchParameter) ([]entity.Article, error) {
	return as.accessor.Search(ctx, param)
}
