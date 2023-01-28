package session

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"todo/pkg/hsjwt"
)

var Timeout int64 = 60 * 60 * 24 // This should be something more reasonable like 5 mins
var RefreshTimeout int64 = 60

var secretKey = []byte("")

func Clear(w http.ResponseWriter) {
	w.Header().Del("Session")
}

func Get(w http.ResponseWriter, r *http.Request) (*Session, error) {
	signedJWT := r.Header.Get("Session")
	if signedJWT == "" {
		return nil, fmt.Errorf("no session")
	}

	s := Session{}
	if err := hsjwt.Decode(signedJWT, &s, secretKey); err != nil {
		return nil, err
	}

	t := time.Now().Unix()
	if t > s.Timeout {
		log.Println("Session timeout", s.Timeout, time.Now().Unix())
		return nil, fmt.Errorf("session expired")
	}

	return &s, nil
}

func Init(key []byte) {
	if len(key) == 0 {
		var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()-_=+,./<>?;:'[]{}|~")

		rand.Seed(time.Now().UnixMilli())
		b := make([]byte, 128)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		secretKey = b
	} else {
		secretKey = key
	}
}

func Set(w http.ResponseWriter, sesh *Session) error {
	var err error

	t := time.Now().Unix()
	sesh.Timeout = t + Timeout
	sesh.RefreshTimeout = t + RefreshTimeout

	signedJWT, err := hsjwt.Encode(&sesh, secretKey)
	if err != nil {
		return err
	}

	w.Header().Del("Session")
	w.Header().Add("Session", signedJWT)

	return nil
}
