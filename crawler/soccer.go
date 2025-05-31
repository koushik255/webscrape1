package soccer

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Player struct is used to store the player's information in the database
type Player struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Goals     int       `json:"goals"`
	Assists   int       `json:"assists"`
	Photo     string     `json:"photo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var db *sql.DB

// InitDatabase initializes the SQLite database and creates the players table
func InitDatabase(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Testing connectino by pinging the database, if it fails then the database is not connected
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// create the players table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS players (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		goals INTEGER NOT NULL DEFAULT 0,
		assists INTERGER NOT NULL DEFAULT 0,
		photo TEXT UNIQUE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	// name is the player's name and is unique so we can't have two players with the same name

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// SavePlayer saves or updates a player's goal count in the database
func SavePlayer(name string, goals int, assists int,photo string) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	// If the player already exists in the database we update their goals and updated_at time so we can keep track of when they last scored
	// Problem is we are scraping the data again because we need to check if its updated so need to find a way to do this without scraping all the data again if its in database
	
	// updateSQL := `UPDATE players SET goals = ?, updated_at = CURRENT_TIMESTAMP WHERE name = ?,UPDATE players SET assists = ?`
	// result, err := db.Exec(updateSQL, goals,assists,name)
	// if err != nil {
	// 	return fmt.Errorf("failed to update player: %v", err)
	// }
	
	updateSQL := `UPDATE players SET goals = ?, assists = ?, photo = ?, updated_at = CURRENT_TIMESTAMP WHERE name = ?`
	result, err := db.Exec(updateSQL, goals, assists,photo, name)
	if err != nil {
    	return fmt.Errorf("failed to update player: %v", err)
	}


	// checking if rows were affected, if not then the player is not in the database already so we need to INSERT them into players
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	// If 0 rows were affected, that means the player is not in the database already so we need to INSERT them into players
	if rowsAffected == 0 {
		insertSQL := `INSERT INTO players (name, goals,assists,photo) VALUES (?, ?,?,?)`
		_, err = db.Exec(insertSQL, name, goals,assists,photo)
		if err != nil {
			return fmt.Errorf("failed to insert player: %v", err)
		}
		log.Printf("Inserted new player: %s with %d goals and assists %d photolink %s", name, goals,assists,photo)
	} else {
		log.Printf("Updated player: %s with %d goals and %d assists photolink %s", name, goals, assists,photo)
	}

	return nil
}




func GetPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var player string
	player = vars["player"]

	fmt.Fprintf(w, "Searching for %s...\n", player)
	// Construct the search URL
    searchURL := "https://fbref.com/en/search/search.fcgi?search=" + player

    // ---- CALL DEBUGHTML HERE ----
    log.Println("----------- DEBUGGING HTML FOR:", searchURL, "-----------")
    DebugHTML(searchURL) // Call your debug function with the URL
    log.Println("----------- FINISHED DEBUGGING HTML -----------")
    // ---- END DEBUG CALL ----

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	goals, err := GetPlayerGoals(searchURL)
	if err != nil {
		fmt.Fprintf(w, "Error getting player goals: %v\n", err)
		return
	}

	assists, err := GetPlayerAssists(searchURL)
	if err != nil {
		fmt.Fprintf(w,"Error getting player assists: %v\n",err)
		return
	}

	photo,err := getPlayerPhoto(searchURL, player)
		if err != nil {
			fmt.Fprintf(w,"Error receiveing player photo! %v\n",err)
		}
	

	fmt.Println("the photo is", photo)

	fmt.Printf("%s has scored %d goals\n", player, goals)
	fmt.Fprintf(w, "%s has scored %d goals\n", player, goals)


	fmt.Printf("%s has  %d assists\n", player, assists)
	fmt.Fprintf(w,"%s has assists  %d times\n", player,assists)


	// Saving player to database
	err = SavePlayer(player, goals,assists,photo)
	if err != nil {
		fmt.Printf("Error saving player to database: %v\n", err)
		fmt.Fprintf(w, "Error saving to database: %v\n", err)
	} else {
		fmt.Fprintf(w, "Player data saved to database successfully!\n")
	}

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Wait()
}

func GetPlayerGoals(url string) (int, error) {
	var goals int
	var parseErr error
	var found bool

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.OnHTML(`td[data-stat="team"]`, func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)

		// Check if this cell contains "Clubs" or club because if they 1 team player its club
		if strings.Contains(strings.ToLower(text), "clubs") || strings.Contains(strings.ToLower(text), "club") {
			fmt.Printf("Found clubs summary row: %s\n", text)

			// gets the parrent row
			// tr is the parent row, and we are looking for the goals cell within this row
			// and then with that row we loop through the cells and find the goals cell
			row := e.DOM.Parent()

			// find the goals cell within this row
			goalsCell := row.Find(`td[data-stat="goals"]`)
			if goalsCell.Length() > 0 {
				goalsText := strings.TrimSpace(goalsCell.Text())
				if goalsText != "" && goalsText != "-" {
					goals, parseErr = strconv.Atoi(goalsText)
					if parseErr != nil {
						fmt.Printf("Error converting '%s' to int: %v\n", goalsText, parseErr)
					} else {
						found = true
						fmt.Printf("Found total goals across all clubs: %s -> %d\n", goalsText, goals)

					}
				}
			} else {
				fmt.Println("No goals cell found in clubs summary row")
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
		parseErr = err
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit(url)
	c.Wait()

	if !found {
		return 0, fmt.Errorf("no goals found")
	}

	return goals, parseErr
}


func GetPlayerAssists(url string) (int, error) {
	var assists int
	var parseErr error
	var found bool

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.OnHTML(`td[data-stat="team"]`, func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)

		// Check if this cell contains "Clubs" or club because if they 1 team player its club
		if strings.Contains(strings.ToLower(text), "clubs") || strings.Contains(strings.ToLower(text), "club") {
			fmt.Printf("Found clubs summary row: %s\n", text)

			// gets the parrent row
			// tr is the parent row, and we are looking for the goals cell within this row
			// and then with that row we loop through the cells and find the goals cell
			row := e.DOM.Parent()

			// find the goals cell within this row
			assistsCell := row.Find(`td[data-stat="assists"]`)
			if assistsCell.Length() > 0 {
				assistsText := strings.TrimSpace(assistsCell.Text())
				if assistsText != "" && assistsText != "-" {
					assists, parseErr = strconv.Atoi(assistsText)
					if parseErr != nil {
						fmt.Printf("Error converting '%s' to int: %v\n", assistsText, parseErr)
					} else {
						found = true
						fmt.Printf("Found total goals across all clubs: %s -> %d\n", assistsText, assists)

					}
				}
			} else {
				fmt.Println("No assists cell found in clubs summary row")
			}
		}

	})
	///---//


		//--
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
		parseErr = err
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit(url)
	c.Wait()

	if !found {
		return 0, fmt.Errorf("no assists found")
	}

	return assists, parseErr
}



func getPlayerPhoto(url string, name string) (string,error) {

		player := name
		var src string
		altText := fmt.Sprintf("%s headshot",player)
		selector := fmt.Sprintf(`img[alt="%s"]`,altText)

		c := colly.NewCollector(
		colly.AllowURLRevisit(),

	)

		c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
		// Replace "Bukayo Saka headshot" with the actual alt text you want to match
   		c.OnHTML(selector, func(e *colly.HTMLElement) {
        src = e.Attr("src")
        fmt.Printf("Image src for %s: %s\n", player, src)
    	})

    // Replace with the actual page URL you want to scrape
    c.Visit(url)
    
    err := c.Visit(url)
    if err != nil {
        return "", err
    }

    if src == "" {
        return "", fmt.Errorf("image not found for player: %s", player)
    }
    return src, nil
}




func DebugHTML(url string) {
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"

	// c.OnHTML("body", func(e *colly.HTMLElement) {
	// 	fmt.Println("=== PAGE CONTENT ===")
	// 	fmt.Println(e.Text[:500]) // First 500 characters
	// 	fmt.Println("=== END CONTENT ===")
	// })

	c.OnHTML("td", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "assists") || e.Attr("data-stat") == "assists" {
			fmt.Printf("Found TD: %s, data-stat: %s, class: %s\n",
				e.Text, e.Attr("data-stat"), e.Attr("class"))
		}
	})

	c.Visit(url)
	c.Wait()    


}
