package todoapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"todo/internal/session"
	"todo/internal/storage"
	"todo/internal/utils"
)

func routeAuth(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "internal server error")
		return
	}
	login := storage.User{}
	if err := json.Unmarshal(body, &login); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "error parsing json")
		return
	}
	user, err := todoDB.FindUser(login.UserName)
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
	jstr, err := json.Marshal(storage.User{Id: user.Id, UserName: user.UserName, FullName: user.FullName, IsAdmin: user.IsAdmin})
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal server error")
		return
	}
	fmt.Fprint(w, string(jstr))
}
