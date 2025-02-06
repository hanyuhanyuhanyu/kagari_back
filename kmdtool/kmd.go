package main

import (
	"context"
	"fmt"
	"kagari/dataaccessor"
	dataaccessorimpl "kagari/dataaccessor/impl"
	"kagari/kmdtool/handler"
	"kagari/setting"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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

func buildHandlers(ctx context.Context, neo4jDriver neo4j.DriverWithContext) (ss *handlers, err error) {
	ss = &handlers{}
	articleAccessor, err := dataaccessorimpl.NewArticleAccessor(ctx, neo4jDriver)
	if err != nil {
		return nil, err
	}
	ss.article = handler.NewArticleHandler(ctx, articleAccessor)
	return
}
func main() {
	ctx := context.Background()
	err := dataaccessor.WithNeo4jConnection(ctx, dataaccessor.ConnectionInfo{
		ConnectionString: setting.Neo4jConnectionString(),
		User:             setting.Neo4jUser(),
		Password:         setting.Neo4jPassword(),
	}, func(neo4jDriver neo4j.DriverWithContext) {
		handlers, err := buildHandlers(ctx, neo4jDriver)
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
