package main

import (
	"fmt"
	"net/http"

	soccer "chill/soccerthing/crawler"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/{player}", soccer.GetPlayer)
	fmt.Println("Server is running on port localhost:8000/{player}")
	http.ListenAndServe(":8000", router)
}
