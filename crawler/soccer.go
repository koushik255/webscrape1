package soccer

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
)

func GetPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	player := vars["player"]

	fmt.Fprintf(w, "You searched for %s", player)

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	goals, err := GetPlayerGoals("https://fbref.com/en/search/search.fcgi?search=" + player)
	if err != nil {
		fmt.Println("Error getting player goals:", err)
	}

	fmt.Println("Goals:", goals)

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	DebugHTML("https://fbref.com/en/search/search.fcgi?search= " + player)

	// c.Visit("https://fbref.com/en/search/search.fcgi?search= " + player)
	// c.PostMultipart("https://fbref.com/en/search/search.fcgi?search= "+player, map[string][]byte{
	// 	"search": []byte(player),
	// })
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

func DebugHTML(url string) {
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"

	// c.OnHTML("body", func(e *colly.HTMLElement) {
	// 	fmt.Println("=== PAGE CONTENT ===")
	// 	fmt.Println(e.Text[:500]) // First 500 characters
	// 	fmt.Println("=== END CONTENT ===")
	// })

	c.OnHTML("td", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "goal") || e.Attr("data-stat") == "goals" {
			fmt.Printf("Found TD: %s, data-stat: %s, class: %s\n",
				e.Text, e.Attr("data-stat"), e.Attr("class"))
		}
	})

	c.Visit(url)
	c.Wait()

}
