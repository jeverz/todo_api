# Todo API

## Features
This is a Todo API I created in order to learn Go. It has one import for Sqlite3. It implements the following.

* Fully Go API
* REST API
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
type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id       int64  `json:"id"`
	UserName string `json:"username"`
	FullName string `json:"fullname"`
	IsAdmin  bool   `json:"isadmin"`
}

type RegisterRequest struct {
	UserName string `json:"username"`
	FullName string `json:"fullname"`
	Password string `json:"password"`
}

type Todo struct {
	Id          int64   `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Completed   *string `json:"completed"`
}
```

## /api/auth
Login with **LoginRequest**. Returns status 200 and **UserResponse** if successful.

## /api/add
Add with **Todo**. **id** is ignored. Returns status 200 and the newly added **Todo** if successful.

## /api/delete
Delete with **Todo**. **id** is to only field required. Returns status 200 if successful.

## /api/user
Get current user details from session. Returns status 200 and **UserResponse** if successful.

## /api/list
Get list of all Todos. Returns status 200 and **[]Todo** if successful.

## /api/ping
Returns "pong"

## /api/register
Register new user using **RegisterRequest**. Returns status 200 and **UserResponse** if successful.

## /api/update
Update a todo with **Todo**. Returns status 200 if successful.
