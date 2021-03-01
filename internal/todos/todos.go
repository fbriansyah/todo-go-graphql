package todos

import (
	"log"

	database "github.com/fbriansyah/todo-go-graphql/internal/pkg/db/mysql"
	"github.com/fbriansyah/todo-go-graphql/internal/users"
)

// Todo struct
type Todo struct {
	ID   string
	Text string
	Done bool
	User *users.User
}

// Save todo ke database
func (todo Todo) Save() int64 {
	//#3
	stmt, err := database.Db.Prepare("INSERT INTO Todos(Text, Done, UserID) VALUES(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	//#4
	// res, err := stmt.Exec(todo.Text, todo.Done)
	res, err := stmt.Exec(todo.Text, todo.Done, todo.User.ID)
	if err != nil {
		log.Fatal(err)
	}
	//#5
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}
	log.Print("Row inserted!")
	return id
}

// GetAll mengambil semua data Todo pada database, return berupa array Todo
func GetAll() []Todo {
	// stmt, err := database.Db.Prepare("SELECT id, Text, Done from Todos")
	stmt, err := database.Db.Prepare("select T.id, T.Test, T.Done, T.UserID, U.Username from Todos T inner join Users U on T.UserID = U.ID")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var todos []Todo
	var username string
	var id string
	for rows.Next() {
		var todo Todo
		// err := rows.Scan(&todo.ID, &todo.Text, &todo.Done)
		err := rows.Scan(&todo.ID, &todo.Text, &todo.Done, &id, &username)
		if err != nil {
			log.Fatal(err)
		}
		todo.User = &users.User{
			ID:       id,
			Username: username,
		} // changed
		todos = append(todos, todo)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return todos
}
