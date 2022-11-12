package router

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", Status)
}

func Status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Continuous is running...")
}
