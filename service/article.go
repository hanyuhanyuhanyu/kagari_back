package service

import (
	"context"
	"kagari/entity"
)

type ArticleAccessor interface {
	GetOne(ctx context.Context, id string) (*entity.Article, error)
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
