// http backend server for kagari the personal playground
package main

import (
	"context"
	"fmt"
	handlerimpl "kagari/handler/impl"
	"kagari/persistence"
	persistenceimpl "kagari/persistence/impl"
	"kagari/router"
	"kagari/setting"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
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
	_, err := os.Stat(envFileName)
	if os.IsNotExist(err) {
		return
	}
	err = godotenv.Load(envFileName)
	if err != nil {
		panic(fmt.Sprintf("error while loading .env %s", err))
	}
}

func main() {
	ctx := context.Background()
	persistence.WithNeo4jConnection(ctx, persistence.ConnectionInfo{
		ConnectionString: setting.Neo4jConnectionString(),
		User:             setting.Neo4jUser(),
		Password:         setting.Neo4jPassword(),
	}, func(neo4jDriver neo4j.DriverWithContext) {
		r := gin.Default()
		acc, err := persistenceimpl.NewArticleAccessor(ctx, neo4jDriver)
		if err != nil {
			log.Fatalf("create accessor fail %v", err)
		}
		articleHandler := handlerimpl.NewArticleHandler(ctx, acc)
		r.Use(cors.New(cors.Config{
			AllowOrigins: []string{
				"http://localhost:3001",
				"https://kagari-frontend-static.s3.ap-northeast-1.amazonaws.com",
			},
			AllowMethods: []string{
				"GET",
				"POST",
				"PATCH",
				"PUT",
				"DELETE",
				"OPTION",
			},
			AllowHeaders: []string{
				"Access-Control-Allow-Credentials",
				"Access-Control-Allow-Headers",
				"Content-Type",
				"Content-Length",
				"Accept-Encoding",
				"Authorization"},
			AllowCredentials: true,
			MaxAge:           8 * time.Hour,
		}))
		router.BuildRoute(r, router.Handlers{ArticleHandler: *articleHandler})
		r.Run()
	})
}
