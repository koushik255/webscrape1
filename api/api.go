package main


import (
	"database/sql"
	"fmt"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
	"encoding/json"


)


type Player struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Goals     int       `json:"goals"`
	Assists   int       `json:"assists"`
	Photo     string    `json:"photo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


// CORS middleware function
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// handle "preflight" requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// call the next handler
		next.ServeHTTP(w, r)
	})
}



// take the path value then search the database for specfic player 

// func getPlayerName(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

// 	player := r.PathValue("player")
// 	fmt.Fprintf(w,"Received request for player: %s",player)

// 	fmt.Println(player)

// }


func getPlayerFromDb(db *sql.DB, name string) ([]Player,error){
	query := "SELECT id, name, goals,assists,photo, created_at, updated_at FROM players WHERE name LIKE ? "


	searchTerm := "%" + name + "%"

	rows,err := db.Query(query,searchTerm)
	if err != nil {
		return nil, fmt.Errorf("Error : %w",err)
	}
	defer rows.Close()

	var player []Player
	for rows.Next() {
		var p Player 
		if err := rows.Scan(&p.ID,&p.Name,&p.Goals,&p.Assists,&p.Photo,&p.CreatedAt,&p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("Error: %w",err)
		}
		player = append(player,p)
	 } 
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ERROR %w",err)
	}

	return player,nil
}




func getAllPlayersFromDb(db *sql.DB) ([]Player,error){
	query := "SELECT id, name, goals,assists,photo, created_at, updated_at FROM players "
	rows,err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Error : %w",err)
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		var p Player 
		if err := rows.Scan(&p.ID,&p.Name,&p.Goals,&p.Assists,&p.Photo,&p.CreatedAt,&p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("Error: %w",err)
		}
		players = append(players,p)
	 } 
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ERROR %w",err)
	}

	return players,nil
}


//create handler func since takes params
func makegetAllPlayersHandler(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers directly in the handler
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Set content type
		w.Header().Set("Content-Type", "application/json")

		allPlayers, err := getAllPlayersFromDb(db)
		if err != nil {
			log.Println("Error getting players from db",err)

			http.Error(w," faieled to retreive players",http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(allPlayers)
		if err != nil {
			fmt.Println("Error encoding players to JSON",err)
			return
		}
	}
}

func makegetPlayersHandler(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers directly in the handler
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Set content type
		w.Header().Set("Content-Type", "application/json")


	name := r.PathValue("player")

	fmt.Printf("Received request for player: %s",name)

	// fmt.Println(name)

		certainPlayer, err := getPlayerFromDb(db,name)
		if err != nil {
			log.Println("Error getting players from db",err)

			http.Error(w," faieled to retreive players",http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(certainPlayer)
		if err != nil {
			fmt.Println("Error encoding players to JSON",err)
			return
		}
	}
}


func main() {

	db, err := sql.Open("sqlite3", "/home/koushikk/go/src/chill/soccerthing/players3.db")
	if err != nil {
		log.Fatalf("Error opening database %v",err)
	}
	defer db.Close()

	//pinign to test connectinon
	err = db.Ping()
	if err != nil {
		log.Fatalf("Erroring pining database : %v",err)
	}
	fmt.Println("Connected to database!")

	//

	allPlayers, err := getAllPlayersFromDb(db)
	if err != nil {
		log.Printf("Error getting all players : %v\n",err)
		return
	} else {
		for _,p := range allPlayers {
			fmt.Printf("ID: %d, Name: %s, Total Goals: %d total assists :%d  photo link  %s\n",p.ID,p.Name,p.Goals,p.Assists,p.Photo)
		}
	}




 

	router := http.NewServeMux()
	getAllPlayersHandler:= makegetAllPlayersHandler(db)
	getPlayerHandler := makegetPlayersHandler(db)
	//routes
	router.HandleFunc("GET /players",getAllPlayersHandler)
	router.HandleFunc("GET /players/{player}", getPlayerHandler)

	corsRouter := corsMiddleware(router)

server := http.Server {
	Addr: ":3000",
	Handler: corsRouter,
}


fmt.Println("Server is up and running on Port:3000!")
server.ListenAndServe()
}