// http backend server for kagari the personal playground
package main

import (
	"context"
	"kagari/dataaccessor"
	dataaccessorimpl "kagari/dataaccessor/impl"
	"kagari/handler"
	handlerimpl "kagari/handler/impl"
	"kagari/setting"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("error while loading .env")
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
