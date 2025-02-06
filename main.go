// http backend server for kagari the personal playground
package main

import (
	"context"
	"fmt"
	"kagari/dataaccessor"
	dataaccessorimpl "kagari/dataaccessor/impl"
	"kagari/handler"
	handlerimpl "kagari/handler/impl"
	"kagari/setting"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func init() {
	envFileName := ".env"
	switch os.Getenv("ENV") {
	case "production":
		envFileName = ".env.prod"
	}
	err := godotenv.Load(envFileName)
	if err != nil {
		panic(fmt.Sprintf("error while loading .env %s", err))
	}
}

func main() {
	ctx := context.Background()
	dataaccessor.WithNeo4jConnection(ctx, dataaccessor.ConnectionInfo{
		ConnectionString: setting.Neo4jConnectionString(),
		User:             setting.Neo4jUser(),
		Password:         setting.Neo4jPassword(),
	}, func(neo4jDriver neo4j.DriverWithContext) {
		r := gin.Default()
		articleHandler := handlerimpl.NewArticleHandler(ctx, (dataaccessorimpl.NewArticleAccessor(ctx, neo4jDriver)))
		handler.BuildRoute(r, handler.Handlers{ArticleHandler: *articleHandler})
		r.Run()
	})
}
