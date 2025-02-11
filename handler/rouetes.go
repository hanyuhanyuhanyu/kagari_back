package handler

import (
	"kagari/handler/impl"
	"kagari/service/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	ArticleHandler impl.ArticleHandler
}

func getPagerFromQuery(c *gin.Context) (*model.Pager, error) {
	limit := c.DefaultQuery("limit", "0")
	offset := c.DefaultQuery("offset", "0")
	return model.PagerFromString(limit, offset)
}
func BuildRoute(r *gin.Engine, handlers Handlers) {
	r.GET("search", func(c *gin.Context) {
		handlers.ArticleHandler.Get(c, c.Param("id"))
	})
	r.GET("/article/:id", func(c *gin.Context) {
		pager, err := getPagerFromQuery(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		param := model.ArticleSearchParameter{
			Pager: *pager,
		}
		handlers.ArticleHandler.Search(c, param)
	})
}
