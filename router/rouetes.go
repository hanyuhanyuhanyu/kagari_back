package router

import (
	"fmt"
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
	r.GET("/article/:id", func(c *gin.Context) {
		fmt.Println(c.Request.Host, c.Request.RequestURI, c.Request.TLS, c.Request.URL, c.Request.PostForm, c.Request.ProtoMinor)
		res, _ := handlers.ArticleHandler.Get(c, c.Param("id"))
		response(c, res)
	})
	r.GET("/article/search", func(c *gin.Context) {
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
