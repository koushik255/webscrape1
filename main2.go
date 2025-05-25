package main

import (
	"fmt"
	"log"
	"net/http"

	soccer "chill/soccerthing/crawler"

	"github.com/gorilla/mux"
)

func main() {
	// Init the database
	dbPath := "./players.db" // SQLite database file will be created in current directory
	err := soccer.InitDatabase(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/{player}", soccer.GetPlayer)

	fmt.Println("Server is running on port localhost:8000/{player}")
	fmt.Printf("Database initialized at: %s\n", dbPath)

	log.Fatal(http.ListenAndServe(":8000", router))
}
