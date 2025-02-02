package entity

import (
	"github.com/gin-gonic/gin"
)

type Article struct {
	ID  string
	URL string
}

func (a *Article) FromMap(m map[string]any) *Article {
	a.ID = m["id"].(string)
	a.URL = m["url"].(string)
	return a
}
func (a *Article) Json() gin.H {
	return gin.H{
		"id":  a.ID,
		"url": a.URL,
	}
}
