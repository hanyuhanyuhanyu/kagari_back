package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"kagari/service"
	"kagari/service/model"
	"kagari/util/jsonutil"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

type ArticleHandler struct {
	service *service.ArticleService
}

func NewArticleHandler(ctx context.Context, acc service.ArticleAccessor) *ArticleHandler {
	return &ArticleHandler{service.NewArticleService(acc)}
}

func (ah *ArticleHandler) Handlers() *cli.Command {
	return &cli.Command{
		Name:    "article",
		Usage:   "manage articles",
		Aliases: []string{"a"},
		Commands: []*cli.Command{
			ah.upload(),
			ah.search(),
		},
	}
}

func (ah *ArticleHandler) upload() *cli.Command {
	var title string

	return &cli.Command{
		Name:    "up",
		Usage:   "upload markdown file",
		Aliases: []string{"u"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "title",
				Aliases:     []string{"t"},
				Usage:       "uploading file title",
				Destination: &title,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			pwd, err := os.Getwd()
			if err != nil {
				return err
			}
			body := cmd.Args().First()
			if body == "" {
				return errors.New("file path must be given")
			}
			filePath := body
			if filepath.IsLocal(body) {
				filePath = filepath.Join(pwd, body)
			}
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			text, err := io.ReadAll(file)
			if err != nil {
				return err
			}
			article := &model.UploadingArticle{
				Body:             text,
				Title:            title,
				OriginalFileName: filepath.Base(filePath),
			}
			images := article.FindImagePaths()
			sources, err := getImageSources(filepath.Dir(filePath), images)
			if err != nil {
				return err
			}
			article.ImageSources = sources
			id, err := ah.service.Post(ctx, article)
			if err != nil {
				return err
			}
			fmt.Printf("file uploaded successfully with id %s\n", id)
			return nil
		},
	}
}

func getImageSources(parent string, imgPaths [][]byte) (maps map[string]io.Reader, err error) {
	maps = make(map[string]io.Reader)
	for _, i := range imgPaths {
		str := string(i)
		// / or .で始まっていなければローカルを指していないとみなす
		if !bytes.HasPrefix(i, []byte("/")) && !bytes.HasPrefix(i, []byte(".")) {
			continue
		}
		original := filepath.Join(parent, str)
		reader, err := os.Open(original)
		if err != nil {
			return maps, err
		}
		maps[str] = reader
	}
	return
}

func (ah *ArticleHandler) search() *cli.Command {
	var limit uint64
	var offset uint64

	return &cli.Command{
		Name:    "search",
		Usage:   "search for files",
		Aliases: []string{"s"},
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:        "limit",
				Aliases:     []string{"l"},
				Usage:       "fetching amout",
				Destination: &limit,
				Value:       0,
			},
			&cli.UintFlag{
				Name:        "offset",
				Aliases:     []string{"o"},
				Usage:       "fetching offset",
				Destination: &offset,
				Value:       0,
			},
		},
		Action: func(ctx context.Context, _ *cli.Command) error {
			articles, err := ah.service.Search(ctx, model.ArticleSearchParameter{
				Pager: model.Pager{
					Limit:  limit,
					Offset: offset,
				},
			})
			if err != nil {
				return err
			}
			marshaled, err := jsonutil.Json(articles, jsonutil.Option{Indent: 2})
			if err != nil {
				return err
			}
			fmt.Println(string(marshaled))
			return nil
		},
	}
}
