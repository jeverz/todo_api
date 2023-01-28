package main

import (
	"log"
	"os"
	"todo/internal/session"
	"todo/internal/storage"
	"todo/internal/todoapi"
	"todo/pkg/httplisten"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":8080"
	}

	secretKey := []byte(os.Getenv("SECRET_KEY"))

	db, err := storage.SqliteOpen("todo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.GetDB().Close()

	todoapi.Init(db)
	session.Init(secretKey)
	session.InitAuth(db)

	httplisten.Serve(port, nil)
}
