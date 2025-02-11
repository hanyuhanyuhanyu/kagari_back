package router

import (
	"kagari/handler"
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
func response(c *gin.Context, res handler.Response) {
	switch res.ContentType {
	default:
		c.JSON(res.Status, res.Body)
	}
}
func BuildRoute(r *gin.Engine, handlers Handlers) {
	r.GET("/article/search", func(c *gin.Context) {
		res, _ := handlers.ArticleHandler.Get(c, c.Param("id"))
		response(c, res)
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
		res, _ := handlers.ArticleHandler.Search(c, param)
		response(c, res)
	})
}
