package session

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"todo/internal/storage"
)

type cacheItem struct {
	session *Session
	timeout int64
}

// Authentication middleware session cache
var sessionCache = map[*http.Request]cacheItem{}
var cacheDuration int64 = 20
var cacheCleanTimeout int64 = 0
var cacheCleanDuration time.Duration = 120 * time.Second

var todoDB storage.TodoDatabase

// Authentication middleware
func Authenticate(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sesh, err := Get(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "session not authenticated")
			log.Println(err)
			return
		}

		t := time.Now().Unix()

		if t > sesh.RefreshTimeout {
			u, err := todoDB.GetUser(sesh.Id)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, "could not find user")
				return
			}

			sesh, err = SetUser(w, u) // Recreate JWT incase fields inside Session are added whos values may change in the database
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "internal server error")
			}
		}

		if t > cacheCleanTimeout {
			cleanCache()
		}
		sessionCache[r] = cacheItem{sesh, t + cacheDuration}

		fn.ServeHTTP(w, r)
	})
}

// Can be called only once to retrieve the cache of a previous call to Authenticate middleware
func Cached(r *http.Request) (*Session, error) {
	s, ok := sessionCache[r]
	if !ok {
		return nil, fmt.Errorf("not in cache")
	}
	delete(sessionCache, r)
	return s.session, nil
}

func cleanCache() {
	t := time.Now().Unix()
	for k, v := range sessionCache {
		if t > v.timeout {
			delete(sessionCache, k)
		}
	}
	cacheCleanTimeout = time.Now().Add(cacheCleanDuration).Unix()
}

func InitAuth(db storage.TodoDatabase) {
	todoDB = db
}
