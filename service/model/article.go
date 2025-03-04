package model

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
)

type UploadingArticle struct {
	Body             []byte
	ImageSources     map[string]io.Reader
	OriginalFileName string
	Title            string
}

var imgFileRegexpTemplate = `!\[[^\]]+\]\(%s(\s+[^\)]+)?\)`

func getImageFileRegexp(imgPath *string) *regexp.Regexp {
	if imgPath == nil {
		return regexp.MustCompile(fmt.Sprintf(imgFileRegexpTemplate, `([^\s)]+)`))
	}
	return regexp.MustCompile(fmt.Sprintf(imgFileRegexpTemplate, *imgPath))
}

func (ua *UploadingArticle) ReplaceImageSource(from string, to string) {
	ua.Body = getImageFileRegexp(&from).ReplaceAllFunc(ua.Body, func(found []byte) []byte {
		return bytes.ReplaceAll(found, []byte(from), []byte(to))
	})
	ua.ImageSources[to] = ua.ImageSources[from]
	delete(ua.ImageSources, from)
}

func (ua *UploadingArticle) FindImagePaths() (found [][]byte) {
	imgs := getImageFileRegexp(nil).FindAllSubmatch(ua.Body, -1)

	for _, matches := range imgs {
		if len(matches) < 2 {
			continue
		}
		found = append(found, matches[1])
	}
	return
}

type ArticleSearchParameter struct {
	Pager
}
