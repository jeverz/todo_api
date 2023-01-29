package todoapi

import (
	"net/http"
	"todo/internal/session"
	"todo/internal/storage"
)

var todoDB storage.TodoDatabase

func Init(db storage.TodoDatabase) {
	todoDB = db
	http.HandleFunc("/api/todo/", session.Authenticate(routeTodo))
	http.HandleFunc("/api/user/", session.Authenticate(routeUser))
	http.HandleFunc("/api/auth", routeAuth)
	http.HandleFunc("/api/ping", routePing)
}
