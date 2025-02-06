package entity

import (
	"fmt"
	"io"
	"regexp"
	"time"

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

type UploadingArticle struct {
	Body             io.Reader
	OriginalFileName string
	Title            string
}

func (ua *UploadingArticle) CreateSavePath(id string) string {
	return fmt.Sprintf("%s/%s_%s.md", time.Now().UTC().Format("200601"), regexp.MustCompile(`\.md$`).ReplaceAllString(ua.OriginalFileName, ""), id)
}
