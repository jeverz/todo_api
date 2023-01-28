package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteDB struct {
	DB   *sql.DB
	prep map[string]*sql.Stmt
}

func (db SqliteDB) AddItem(userId int64, todo Todo) (*Todo, error) {
	newTodo := Todo{}
	row := db.prep["AddItem"].QueryRow(userId, todo.Title, todo.Description, todo.Completed)
	if err := row.Scan(&newTodo.Id, &newTodo.Title, &newTodo.Description, &newTodo.Completed); err != nil {
		return nil, err
	}
	return &newTodo, nil
}

func (db SqliteDB) AddUser(u *User) error {
	_, err := db.prep["AddUser"].Exec(u.UserName, u.FullName, u.Password, u.IsAdmin)
	return err
}

func (db SqliteDB) DeleteItem(userId, itemId int64) error {
	result, err := db.prep["DeleteItem"].Exec(userId, itemId)
	if rows, _ := result.RowsAffected(); err == nil && rows == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (db SqliteDB) GetDB() *sql.DB {
	return db.DB
}

// Will return sql.ErrNoRows if user not found
func (db SqliteDB) GetUser(username string) (*User, error) {
	row := db.prep["GetUser"].QueryRow(username)
	var u User
	var isAdmin int64
	if err := row.Scan(&u.Id, &u.UserName, &u.Password, &isAdmin, &u.FullName); err != nil {
		return nil, err
	}
	u.IsAdmin = isAdmin != 0
	return &u, nil
}

func (db SqliteDB) GetUserById(id int64) (*User, error) {
	row := db.prep["GetUserById"].QueryRow(id)
	var u User
	var isAdmin int64
	if err := row.Scan(&u.Id, &u.UserName, &u.Password, &isAdmin, &u.FullName); err != nil {
		return nil, err
	}
	u.IsAdmin = isAdmin != 0
	return &u, nil
}

func (db SqliteDB) ListItems(userId int64) ([]Todo, error) {
	rows, err := db.prep["ListItems"].Query(userId)
	if err != nil {
		return nil, err
	}

	todos := []Todo{}
	for rows.Next() {
		var t Todo
		err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.Completed)
		if err != nil {
			log.Println(err)
			break
		}
		todos = append(todos, t)
	}

	return todos, nil
}

func (db SqliteDB) UpdateItem(userId int64, todo Todo) error {
	result, err := db.prep["UpdateItem"].Exec(todo.Title, todo.Description, todo.Completed, todo.Id, userId)
	if rows, _ := result.RowsAffected(); err == nil && rows == 0 {
		return sql.ErrNoRows
	}
	return err
}

func SqliteOpen(cs string) (TodoDatabase, error) {
	var err error
	db := SqliteDB{}
	db.DB, err = sql.Open("sqlite3", cs)
	if err != nil {
		return nil, err
	}

	db.prep = make(map[string]*sql.Stmt)

	statements := []struct {
		k string
		v string
	}{
		{"AddItem", "INSERT INTO todo (user_id, title, description, completed) VALUES (?,?,?,?) RETURNING id, title, description, completed"},
		{"AddUser", "INSERT INTO user (username,fullname,password,is_admin) VALUES (?,?,?,?)"},
		{"DeleteItem", "DELETE FROM todo WHERE user_id=? AND id=?"},
		{"GetUser", "SELECT id, username, password, is_admin, fullname FROM user WHERE username=?"},
		{"GetUserById", "SELECT id, username, password, is_admin, fullname FROM user WHERE id=?"},
		{"ListItems", "SELECT id, title, description, completed FROM todo WHERE user_id=?"},
		{"UpdateItem", "UPDATE todo SET title=?, description=?, completed=? WHERE id=? AND user_id=?"},
	}

	for _, stmt := range statements {
		db.prep[stmt.k], err = db.DB.Prepare(stmt.v)
		if err != nil {
			return nil, fmt.Errorf("error preparing %q: %v", stmt.k, err)
		}
	}

	return db, nil
}
