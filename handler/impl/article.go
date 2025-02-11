package impl

import (
	"context"
	"kagari/handler"
	"kagari/service"
	"kagari/service/model"
	"net/http"
)

type ArticleHandler struct {
	service *service.ArticleService
}

func NewArticleHandler(ctx context.Context, acc service.ArticleAccessor) *ArticleHandler {
	return &ArticleHandler{service.NewArticleService(acc)}
}

func (ah *ArticleHandler) Get(c context.Context, id string) (handler.Response, error) {
	article, err := ah.service.GetArticle(c, id)
	switch {
	case err != nil:
		{
			return handler.Response{
				Status: http.StatusInternalServerError,
				Body: map[string]any{
					"error": err.Error(),
				},
			}, err
		}
	case article == nil:
		{
			return handler.Response{
				Status: http.StatusNotFound,
				Body: map[string]any{
					"error": "article not found",
				},
			}, nil

		}
	default:
		return handler.Response{
			Status: http.StatusOK,
			Body:   article,
		}, nil
	}
}
func (ah *ArticleHandler) Search(c context.Context, param model.ArticleSearchParameter) (handler.Response, error) {
	articles, err := ah.service.Search(c, param)
	switch {
	case err != nil:
		{
			return handler.Response{
				Status: http.StatusInternalServerError,
				Body: map[string]any{
					"error": err.Error(),
				},
			}, err
		}
	default:
		return handler.Response{
			Status: http.StatusOK,
			Body:   articles,
		}, nil
	}
}
