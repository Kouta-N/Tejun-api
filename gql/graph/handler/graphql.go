package graph

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"main.go/gql/graph"
)

const defaultPort = "8084"

func main() {
      port := os.Getenv("PORT")
      if port == "" {
          port = defaultPort
      }

      r := gin.Default()
      r.Use(cors.New(cors.Config{ // Use() を通してアタッチされたミドルウェアは、 すべてのリクエストのハンドラチェーンに含まれる
            AllowOrigins:     []string{"http://localhost:3000"},
            AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
            AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
            AllowCredentials: true,
      }))

      srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

      // playground.Handlerは、ブラウザ上でGraphQLクエリを試すことができる"GraphQL playground"というツールを提供します。"/query"はGraphQLリクエストを送信するエンドポイントを示します。
      r.GET("/", gin.WrapF(playground.Handler("GraphQL playground", "/query")))
      r.POST("/query", gin.WrapH(srv))

						// WrapF is a helper function for wrapping http.HandlerFunc and returns a Gin middleware.
      // func WrapF(f http.HandlerFunc) HandlerFunc {
      // 	return func(c *Context) {
      // 		f(c.Writer, c.Request)
      // 	}
      // }

      // WrapH is a helper function for wrapping http.Handler and returns a Gin middleware.
      // func WrapH(h http.Handler) HandlerFunc {
      // 	return func(c *Context) {
      // 		h.ServeHTTP(c.Writer, c.Request)
      // 	}
      // }

      log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
      r.Run(":" + port) // デフォルトでPanicが発生するので、log.Fatalは不要
}