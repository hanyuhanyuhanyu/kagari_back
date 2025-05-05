package main

import (
	"context"
	"fmt"
	"kagari/kmdtool/handler"
	"kagari/persistence"
	persistenceimpl "kagari/persistence/impl"
	"log"
	"os"

	pgx "github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func init() {
	envFileName := ".env"
	switch os.Getenv("ENV") {
	case "production":
		envFileName = ".env.prod"
	}
	_, err := os.Stat(envFileName)
	if os.IsNotExist(err) {
		return
	}
	err = godotenv.Load(envFileName)
	if err != nil {
		panic(fmt.Sprintf("error while loading .env %s", err))
	}
}

type handlers struct {
	article *handler.ArticleHandler
}

func buildHandlers(ctx context.Context, conn *pgx.Conn) (ss *handlers, err error) {
	ss = &handlers{}
	articleAccessor, err := persistenceimpl.NewArticleAccessor(ctx, conn)
	if err != nil {
		return nil, err
	}
	ss.article = handler.NewArticleHandler(ctx, articleAccessor)
	return
}
func main() {
	ctx := context.Background()
	err := persistence.WithPsqlConnection(ctx, func(conn *pgx.Conn) {
		handlers, err := buildHandlers(ctx, conn)
		if err != nil {
			log.Fatal(err)
		}
		commands := &cli.Command{
			Commands: []*cli.Command{
				handlers.article.Handlers(),
			},
		}
		err = commands.Run(context.Background(), os.Args)
		if err != nil {
			log.Fatal(err)
		}
	})
	if err != nil {
		log.Fatal(err)
		return
	}
}
