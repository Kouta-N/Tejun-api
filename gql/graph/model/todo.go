package model

import (
	"database/sql"
	"net/http"

	// "tmhub/helpers"

	"github.com/mholt/binding"
)

func (n *Todo) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&n.ID:      "id",
		&n.UserID:  "userId",
		&n.Content: "content",
		&n.IsDone:  "isDone",
	}
}

func GetTodo() []*Todo {
	db, err := sql.Open("mysql", "tester:password@tcp(db:3306)/test?charset=utf8&parseTime=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	results, err := db.Query("SELECT * FROM todos;")
	if err != nil {
		panic(err)
	}

	var todos []*Todo
	for results.Next() {
		var t Todo
		err := results.Scan(&t.ID, &t.UserID, &t.Content, &t.IsDone, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			panic(err)
		}
		todos = append(todos, &t)
	}
	return todos
}
