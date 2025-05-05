// http backend server for kagari the personal playground
package main

import (
	"context"
	"fmt"
	handlerimpl "kagari/handler/impl"
	"kagari/persistence"
	persistenceimpl "kagari/persistence/impl"
	"kagari/router"
	"log"
	"os"
	"time"

	pgx "github.com/jackc/pgx/v5"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	persistence.WithPsqlConnection(ctx, func(conn *pgx.Conn) {
		r := gin.Default()
		acc, err := persistenceimpl.NewArticleAccessor(ctx, conn)
		if err != nil {
			log.Fatalf("create accessor fail %v", err)
		}
		articleHandler := handlerimpl.NewArticleHandler(ctx, acc)
		r.Use(cors.New(cors.Config{
			AllowOrigins: []string{
				"http://localhost:3001",
				"http://kagari-frontend-static.s3-website-ap-northeast-1.amazonaws.com",
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
