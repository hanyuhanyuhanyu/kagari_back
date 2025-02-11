package model

import (
	"fmt"
	"io"
	"regexp"
	"time"
)

type UploadingArticle struct {
	Body             io.Reader
	OriginalFileName string
	Title            string
}

func (ua *UploadingArticle) CreateSavePath(id string) string {
	return fmt.Sprintf("%s/%s_%s.md", time.Now().UTC().Format("200601"), regexp.MustCompile(`\.md$`).ReplaceAllString(ua.OriginalFileName, ""), id)
}

type ArticleSearchParameter struct {
	Pager
}
