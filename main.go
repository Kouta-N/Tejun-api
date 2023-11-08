// 参考 https://zenn.dev/shimpo/articles/setup-go-mysql-with-docker-compose
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"main.go/gql/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID            int64          `db:"id" json:"id"`
	Email         string         `db:"email" json:"email"`
	Password      string         `db:"password" json:"password"`
	Name          string         `db:"name" json:"name"`
	ProfileImage  sql.NullString `db:"profile_image" json:"profileImage"`
	EmailVerified bool           `db:"email_verified" json:"emailVerified"`
	CreatedAt     *time.Time     `db:"created_at" json:"createdAt"`
	UpdatedAt     *time.Time     `db:"updated_at" json:"updatedAt"`
	DeletedAt     *time.Time     `db:"deleted_at" json:"deletedAt"`
}

const defaultPort = "8082"

func getUsers() []*User {
	db, err := sql.Open("mysql", "tester:password@tcp(db:3306)/test?charset=utf8&parseTime=true") //parseTime=trueは、DATEおよびDATETIME値の出力タイプを[]byte/stringの代わりにtime.Timeに変更
	if err != nil {
		panic(err)
	}
	defer db.Close()

	results, err := db.Query("SELECT * FROM users")
	fmt.Println("⭐️", results)
	if err != nil {
		panic(err)
	}

	var users []*User
	for results.Next() {
		var u User
		err := results.Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.ProfileImage, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
		if err != nil {
			panic(err)
		}
		users = append(users, &u)
	}
	return users
}

func usersPage(w http.ResponseWriter, r *http.Request) {
	users := getUsers()
	json.NewEncoder(w).Encode(users)
}

func main() {
	// http.HandleFunc("/", handler)
	// http.HandleFunc("/users", usersPage)
	// log.Fatal(http.ListenAndServe(":8082", nil))

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{ // Use() を通してアタッチされたミドルウェアは、 すべてのリクエストのハンドラチェーンに含まれる
		AllowOrigins:     []string{"http://localhost:8083", "*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
	}))

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	r.GET("/", gin.WrapF(playground.Handler("GraphQL playground", "/query")))
	r.POST("/query", gin.WrapH(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	r.Run(":" + port) // デフォルトでPanicが発生するので、log.Fatalは不要
}
