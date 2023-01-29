package todoapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"todo/internal/session"
	"todo/internal/storage"
)

func routeTodo(w http.ResponseWriter, r *http.Request) {
	sesh, err := session.Cached(r)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal server error")
		return
	}

	idRequest := getId("/api/todo/", r.URL.Path)

	switch r.Method {
	case "GET":
		if idRequest != nil {
			todo, err := todoDB.GetTodo(sesh.Id, *idRequest)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "error accessing record")
				return
			}
			jstr, err := json.Marshal(&todo)
			if err != nil {
				log.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "internal server error")
				return
			}
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintln(w, string(jstr))
			return
		}
		todos, err := todoDB.GetTodos(sesh.Id)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "internal server error")
			return
		}
		jstr, err := json.Marshal(todos)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "internal server error")
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(jstr))
	case "POST":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "internal server error")
			return
		}

		todo := storage.Todo{}
		if err := json.Unmarshal(body, &todo); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "error parsing json")
			return
		}

		newTodo, err := todoDB.AddTodo(sesh.Id, todo)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "internal server error")
			return
		}
		jstr, err := json.Marshal(newTodo)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "internal server error")
			return
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(jstr))
	case "PUT":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "internal server error")
			return
		}

		todo := storage.Todo{}
		if err := json.Unmarshal(body, &todo); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "error parsing json")
			return
		}

		err = todoDB.UpdateTodo(sesh.Id, todo)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "error modifying item")
			return
		}
		fmt.Fprintln(w, "updated")
	case "DELETE":
		if idRequest == nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "no id suplied")
			return
		}
		if err := todoDB.DeleteTodo(sesh.Id, *idRequest); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "error deleting record")
			return
		}
		fmt.Fprintln(w, "deleted")
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "unhandled method")
	}
}
