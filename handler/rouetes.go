package handler

import (
	"kagari/handler/impl"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	ArticleHandler impl.ArticleHandler
}

func BuildRoute(r *gin.Engine, handlers Handlers) {
	r.GET("/article/:id", handlers.ArticleHandler.Get)
}
