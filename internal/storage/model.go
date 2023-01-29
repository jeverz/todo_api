package storage

import (
	"database/sql"
)

type User struct {
	Id       int64  `json:"id"`
	UserName string `json:"username"`
	FullName string `json:"fullname"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isadmin"`
}

type Todo struct {
	Id          int64   `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Completed   *string `json:"completed"`
}

type TodoDatabase interface {
	GetDB() *sql.DB

	GetTodo(userId int64, todoId int64) (Todo, error)
	GetTodos(userId int64) ([]Todo, error)
	AddTodo(userId int64, todo Todo) (*Todo, error)
	UpdateTodo(userId int64, todo Todo) error
	DeleteTodo(userId, itemId int64) error

	GetUser(id int64) (*User, error)
	FindUser(s string) (*User, error)
	AddUser(u *User) error
	UpdateUser(u *User) error
	DeleteUser(id int64) error
}
