# Todo API

This is a Todo API I created in order to learn Go. It has one import for Sqlite3. It implements the following.

* Fully Go API
* REST API
* Session Management
* Authentication
* Simple HMAC JWT
* Sqlite3 database
* Graceful server shutdown

## Usage
To setup the database execute (this will drop tables) or copy the sample to /cmd
```
cat schema | sqlite3 cmd/todo.db
```
To run the example
```
cd cmd
go run main.go
```
