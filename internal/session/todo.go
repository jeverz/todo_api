package session

import (
	"net/http"
	"todo/internal/storage"
)

func SetUser(w http.ResponseWriter, u *storage.User) (*Session, error) {
	sesh := Session{u.Id, 0, 0, u.IsAdmin}

	if err := Set(w, &sesh); err != nil {
		return nil, err
	}

	return &sesh, nil
}
