// 参考 https://zenn.dev/shimpo/articles/setup-go-mysql-with-docker-compose
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID            int64          `db:"id" json:"id"`
	Email         string         `db:"email" json:"email"`
	Password      string         `db:"password" json:"password"`
	Name          string         `db:"name" json:"name"`
	ProfileImage  sql.NullString `db:"profile_image" json:"profileImage"`
	EmailVerified bool           `db:"email_verified" json:"emailVerified"`
	Via           string         `db:"via" json:"via"`
	CreatedAt     *time.Time     `db:"created_at" json:"createdAt"`
	UpdatedAt     *time.Time     `db:"updated_at" json:"updatedAt"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func getUsers() []*User {
	db, err := sql.Open("mysql", "tester:password@tcp(db:3306)/test")
	fmt.Println("⭐️", db)
	if err != nil {
				fmt.Println("🌞", err)
		panic(err)
	}
	defer db.Close()

	results, err := db.Query("SELECT * FROM users")
		fmt.Println("🐎", results)
	if err != nil {
		panic(err)
	}

	var users []*User
	for results.Next() {
		var u User
		err := results.Scan(&u.ID, &u.Name)
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
	http.HandleFunc("/", handler)
	http.HandleFunc("/users", usersPage)
	log.Fatal(http.ListenAndServe(":8082", nil))
}