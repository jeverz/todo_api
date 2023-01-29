# Todo API

## Features
This is a Todo API I created in order to learn Go. It has one import for Sqlite3. It implements the following.

* Fully Go API
* CRUD
* Session Management
* Authentication
* Simple HMAC JWT
* Sqlite3 database
* Graceful server shutdown

## Starting
To setup the database execute (this will drop tables) or copy the sample to /cmd
```
cat schema | sqlite3 cmd/todo.db
```
To run the example
```
cd cmd
go run main.go
```

# API
The server will provide a JWT in the header of an /api/auth request and occasionally provide a new token when it is due to be refreshed during any api call. The default port is 8080 and runs on plain HTTP.

## Models
```golang
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
```

## General

### GET /api/auth
Login with json "username" and "password". Returns status 200 and **UserResponse** if successful.

### GET /api/ping
Returns "pong"

## Todo API

### GET /api/todo/[id]
Returns status 200 and current user **[]Todo** with no params or a **Todo** when *id* is supplied.

### POST /api/todo/
Create new todo with json **Todo**. "id" is ignored. Returns 200 on success.

### PUT /api/todo/
Send json **Todo** to update a todo. Returns 200 on success.

### DELETE /api/todo/*id*
Deletes todo with *id*. Returns 200 on success.

## User API

### GET /api/user/[id]
No params returns current **User** information and 200 OK. If user is admin can get **User** information by supplying *id*.

### POST /api/user/
Create a new user sending **User**. Must be admin. Returns **User** and 200 OK on success.

### PUT /api/user/
Modify a user by sending **User**. An empty or missing password will retain the password. Must be admin. Returns 200 OK on success.

### DELETE /api/user/*id*
Must be admin. Returns 200 OK on success.
