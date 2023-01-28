package todoapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"todo/internal/session"
	"todo/internal/storage"
	"todo/internal/utils"
)

func routeAdd(w http.ResponseWriter, r *http.Request) {
	sesh, err := session.Cached(r)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal server error")
		return
	}
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

	//fmt.Fprintln(w, todo, *todo.Description, *todo.Completed)
	newTodo, err := todoDB.AddItem(sesh.Id, todo)
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
	fmt.Fprint(w, string(jstr))
}

func routeDelete(w http.ResponseWriter, r *http.Request) {
	sesh, err := session.Cached(r)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal server error")
		return
	}
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
	if err := todoDB.DeleteItem(sesh.Id, todo.Id); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "error deleting record")
		return
	}
	fmt.Fprintln(w, "deleted")
}

func routeAuth(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "internal server error")
		return
	}
	login := LoginRequest{}
	if err := json.Unmarshal(body, &login); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "error parsing json")
		return
	}
	user, err := todoDB.GetUser(login.UserName)
	if err != nil || *utils.ShaEncode(&login.Password) != user.Password { // No error
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "invalid username or password")
		if err != nil {
			log.Println(err)
		}
		return
	}
	if _, err = session.SetUser(w, user); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "internal server error")
		return
	}
	jstr, err := json.Marshal(UserResponse{user.Id, user.UserName, user.FullName, user.IsAdmin })
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal server error")
		return
	}
	fmt.Fprint(w, string(jstr))
}

func routeList(w http.ResponseWriter, r *http.Request) {
	sesh, err := session.Cached(r)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal server error")
		return
	}
	todos, err := todoDB.ListItems(sesh.Id)
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
	fmt.Fprint(w, string(jstr))
}

func routePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pong")
}

func routeRegister(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "internal server error")
		return
	}
	register := RegisterRequest{}
	if err := json.Unmarshal(body, &register); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "error parsing json")
		return
	}
	_, err = todoDB.GetUser(register.UserName)
	if err != sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "username taken")
		return
	}
	if register.Password == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "password not complex enough")
		return
	}
	user := storage.User{
		Id:       0,
		UserName: register.UserName,
		FullName: register.FullName,
		Password: *utils.ShaEncode(&register.Password),
	}

	todoDB.AddUser(&user)

	if _, err = session.SetUser(w, &user); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "internal server error")
		return
	}

	fmt.Fprintln(w, "user added")
}

func routeUpdate(w http.ResponseWriter, r *http.Request) {
	sesh, err := session.Cached(r)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal server error")
		return
	}
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

	//fmt.Fprintln(w, todo, *todo.Description, *todo.Completed)
	err = todoDB.UpdateItem(sesh.Id, todo)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "error modifying item")
		return
	}
	fmt.Fprintln(w, "updated")
}

func routeUser(w http.ResponseWriter, r *http.Request) {
	sesh, err := session.Cached(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	user, err := todoDB.GetUserById(sesh.Id)
	if err != nil { // No error
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	jstr, err := json.Marshal(UserResponse{user.Id, user.UserName, user.FullName, user.IsAdmin })
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal server error")
		return
	}
	fmt.Fprint(w, string(jstr))
}

