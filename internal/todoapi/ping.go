package todoapi

import (
	"fmt"
	"net/http"
)

func routePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pong")
}
