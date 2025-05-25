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
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}



func getPlayerName(w http.ResponseWriter, r *http.Request) {
	player := r.PathValue("player")
	fmt.Fprintf(w,"Received request for player: %s",player)
	fmt.Println(player)
}


func getAllPlayersFromDb(db *sql.DB) ([]Player,error){
	query := "SELECT id, name, goals, created_at, updated_at FROM players "
	rows,err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Error : %w",err)
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		var p Player 
		if err := rows.Scan(&p.ID,&p.Name,&p.Goals,&p.CreatedAt,&p.UpdatedAt); err != nil {
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


func main() {

	db, err := sql.Open("sqlite3", "/home/koushikk/go/src/chill/soccerthing/players.db")
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
		log.Printf("Error getting all players : %v\n,err")
		return
	} else {
		for _,p := range allPlayers {
			fmt.Printf("ID: %d, Name: %s, Total Goals: %d\n",p.ID,p.Name,p.Goals)
		}
	}


	router := http.NewServeMux()
	getAllPlayersHandler:= makegetAllPlayersHandler(db)
	//routes
	router.HandleFunc("GET /players",getAllPlayersHandler)
	router.HandleFunc("GET /players/{player}", getPlayerName)
server := http.Server {
	Addr: ":3000",
	Handler: router,
}


fmt.Println("Server is up and running on Port:3000!")
server.ListenAndServe()
}