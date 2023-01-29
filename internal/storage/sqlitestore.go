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
		{"GetTodo", "SELECT id, title, description, completed FROM todo WHERE user_id=? AND id=?"},
		{"GetTodos", "SELECT id, title, description, completed FROM todo WHERE user_id=?"},
		{"AddTodo", "INSERT INTO todo (user_id, title, description, completed) VALUES (?,?,?,?) RETURNING id, title, description, completed"},
		{"UpdateTodo", "UPDATE todo SET title=?, description=?, completed=? WHERE id=? AND user_id=?"},
		{"DeleteTodo", "DELETE FROM todo WHERE user_id=? AND id=?"},

		{"GetUser", "SELECT id, username, password, is_admin, fullname FROM user WHERE id=?"},
		{"FindUser", "SELECT id, username, password, is_admin, fullname FROM user WHERE username=?"},
		{"AddUser", "INSERT INTO user (username,fullname,password,is_admin) VALUES (?,?,?,?)"},
		{"UpdateUser", "UPDATE user SET username=?, fullname=?, password=?, is_admin=? WHERE id=?"},
		{"DeleteUser", "DELETE FROM user WHERE id=?"},
	}

	for _, stmt := range statements {
		db.prep[stmt.k], err = db.DB.Prepare(stmt.v)
		if err != nil {
			return nil, fmt.Errorf("error preparing %q: %v", stmt.k, err)
		}
	}

	return db, nil
}

func (db SqliteDB) GetDB() *sql.DB {
	return db.DB
}

/****************** TODO ******************/
/****************** TODO ******************/
/****************** TODO ******************/

func (db SqliteDB) GetTodo(userId int64, todoId int64) (Todo, error) {
	row := db.prep["GetTodo"].QueryRow(userId, todoId)

	todo := Todo{}
	err := row.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed) 
	if err != nil {
		return todo, err
	}

	return todo, nil
}

func (db SqliteDB) GetTodos(userId int64) ([]Todo, error) {
	rows, err := db.prep["GetTodos"].Query(userId)
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

func (db SqliteDB) AddTodo(userId int64, todo Todo) (*Todo, error) {
	newTodo := Todo{}
	row := db.prep["AddTodo"].QueryRow(userId, todo.Title, todo.Description, todo.Completed)
	if err := row.Scan(&newTodo.Id, &newTodo.Title, &newTodo.Description, &newTodo.Completed); err != nil {
		return nil, err
	}
	return &newTodo, nil
}

func (db SqliteDB) UpdateTodo(userId int64, todo Todo) error {
	result, err := db.prep["UpdateTodo"].Exec(todo.Title, todo.Description, todo.Completed, todo.Id, userId)
	if rows, _ := result.RowsAffected(); err == nil && rows == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (db SqliteDB) DeleteTodo(userId, itemId int64) error {
	result, err := db.prep["DeleteTodo"].Exec(userId, itemId)
	if rows, _ := result.RowsAffected(); err == nil && rows == 0 {
		return sql.ErrNoRows
	}
	return err
}

/****************** USER ******************/
/****************** USER ******************/
/****************** USER ******************/

func (db SqliteDB) GetUser(id int64) (*User, error) {
	row := db.prep["GetUser"].QueryRow(id)
	var u User
	var isAdmin int64
	if err := row.Scan(&u.Id, &u.UserName, &u.Password, &isAdmin, &u.FullName); err != nil {
		return nil, err
	}
	u.IsAdmin = isAdmin != 0
	return &u, nil
}

// Will return sql.ErrNoRows if user not found
func (db SqliteDB) FindUser(username string) (*User, error) {
	row := db.prep["FindUser"].QueryRow(username)
	var u User
	var isAdmin int64
	if err := row.Scan(&u.Id, &u.UserName, &u.Password, &isAdmin, &u.FullName); err != nil {
		return nil, err
	}
	u.IsAdmin = isAdmin != 0
	return &u, nil
}

func (db SqliteDB) AddUser(u *User) error {
	_, err := db.prep["AddUser"].Exec(u.UserName, u.FullName, u.Password, u.IsAdmin)
	return err
}

func (db SqliteDB) UpdateUser(u *User) error {
	result, err := db.prep["UpdateUser"].Exec(u.UserName, u.FullName, u.Password, u.IsAdmin, u.Id)
	if rows, _ := result.RowsAffected(); err == nil && rows == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (db SqliteDB) DeleteUser(id int64) error {
	result, err := db.prep["DeleteUser"].Exec(id)
	if rows, _ := result.RowsAffected(); err == nil && rows == 0 {
		return sql.ErrNoRows
	}
	return err
}
