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

func routeUser(w http.ResponseWriter, r *http.Request) {
	sesh, err := session.Cached(r)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal server error")
		return
	}

	idRequest := getId("/api/user/", r.URL.Path)

	switch r.Method {
	case "GET":
		var user *storage.User
		var err error

		if idRequest == nil {
			user, err = todoDB.GetUser(sesh.Id)
		} else {
			if !sesh.IsAdmin {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("unauthorized"))
				return
			}
			user, err = todoDB.GetUser(*idRequest)
		}
		if err != nil { // No error
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			return
		}
		user.Password = ""
		jstr, err := json.Marshal(&user)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "internal server error")
			return
		}
		fmt.Fprint(w, string(jstr))
	case "POST":
		if !sesh.IsAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "unauthorized")
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "internal server error")
			return
		}
		user := storage.User{}
		if err := json.Unmarshal(body, &user); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "error parsing json")
			return
		}
		_, err = todoDB.FindUser(user.UserName)
		if err != sql.ErrNoRows {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprintln(w, "username taken")
			return
		}
		if user.Password == "" {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprintln(w, "password not complex enough")
			return
		}
		user.Password = *utils.ShaEncode(&user.Password)

		todoDB.AddUser(&user)

		if _, err = session.SetUser(w, &user); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "internal server error")
			return
		}

		jstr, err := json.Marshal(&user)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "internal server error, but user created")
			return
		}
		fmt.Fprint(w, string(jstr))
	case "PUT":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "internal server error")
			return
		} 
		user := storage.User{}
		if err := json.Unmarshal(body, &user); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "error parsing json")
			return
		}
		oldUser, err := todoDB.GetUser(user.Id)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprintln(w, "cannot find user")
			return
		}

		if user.UserName == "" {
			user.UserName = oldUser.UserName
		}
		if user.FullName == "" {
			user.FullName = oldUser.FullName
		}
		if user.Password == "" {
			user.Password = oldUser.Password
		} else {
			user.Password = *utils.ShaEncode(&user.Password)
		}

		todoDB.UpdateUser(&user)

		fmt.Fprintln(w, "user modified")
	case "DELETE":
		if !sesh.IsAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}
		if idRequest == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("no id supplied"))
			return
		}
		if err := todoDB.DeleteUser(*idRequest); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error deleting user"))
			return
		}
		w.Write([]byte("user deleted"))
	}
}
