package model

import (
	"database/sql"
)

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
		err := results.Scan(&t.ID, &t.UserID, &t.Title, &t.Content, &t.IsDone, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			panic(err)
		}
		todos = append(todos, &t)
	}
	return todos
}
