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
	AddItem(userId int64, todo Todo) (*Todo, error)
	DeleteItem(userId, itemId int64) error
	GetDB() *sql.DB
	GetUser(s string) (*User, error)
	GetUserById(id int64) (*User, error)
	AddUser(u *User) error
	ListItems(userId int64) ([]Todo, error)
	UpdateItem(userId int64, todo Todo) error
}
