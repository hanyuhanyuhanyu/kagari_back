package handler

import (
	"kagari/handler/handlerimpl"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	ArticleHandler handlerimpl.ArticleHandler
}

func BuildRoute(r *gin.Engine, handlers Handlers) {
	r.GET("/article/:id", handlers.ArticleHandler.Get)
}
