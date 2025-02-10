package impl

import (
	"context"
	"kagari/service"
	"kagari/service/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	service *service.ArticleService
}

func NewArticleHandler(ctx context.Context, acc service.ArticleAccessor) *ArticleHandler {
	return &ArticleHandler{service.NewArticleService(acc)}
}

func (ah *ArticleHandler) Get(c *gin.Context, id string) {
	article, err := ah.service.GetArticle(c.Request.Context(), id)
	switch {
	case err != nil:
		{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	case article == nil:
		{
			c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		}
	default:
		c.JSON(http.StatusOK, article)
	}
}
func (ah *ArticleHandler) Search(c *gin.Context, param model.ArticleSearchParameter) {
	article, err := ah.service.Search(c.Request.Context(), param)
	switch {
	case err != nil:
		{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	default:
		c.JSON(http.StatusOK, article)
	}
}
