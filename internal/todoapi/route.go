package todoapi

import (
	"net/http"
	"todo/internal/session"
	"todo/internal/storage"
)

var todoDB storage.TodoDatabase

func Init(db storage.TodoDatabase) {
	todoDB = db
	http.HandleFunc("/api/auth", routeAuth)
	http.HandleFunc("/api/add", session.Authenticate(routeAdd))
	http.HandleFunc("/api/delete", session.Authenticate(routeDelete))
	http.HandleFunc("/api/user", session.Authenticate(routeUser))
	http.HandleFunc("/api/list", session.Authenticate(routeList))
	http.HandleFunc("/api/ping", routePing)
	http.HandleFunc("/api/register", routeRegister)
	http.HandleFunc("/api/update", session.Authenticate(routeUpdate))
}
