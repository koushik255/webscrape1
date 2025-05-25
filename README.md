# Soccer Player Database

A simple web scraper that fetches soccer player goal statistics and stores them in a local SQLite database.

## Features

- Scrapes player goal data from fbref.com
- Stores player names and goal counts in SQLite database
- Automatically creates or updates player records
- Simple HTTP API for searching players

## Setup

1. Install dependencies:
```bash
go mod tidy
```

2. Run the server:
```bash
go run main2.go
```

The server will start on `localhost:8000` and create a `players.db` SQLite file in the current directory.

## Usage

### Search and Store Player Data

Make a GET request to search for a player and store their goals in the database:

```bash
curl http://localhost:8000/messi
curl http://localhost:8000/ronaldo
curl http://localhost:8000/"harry kane"
```

The application will:
1. Search for the player on fbref.com
2. Scrape their total career goals
3. Store or update the data in the SQLite database
4. Return the results

## Database Schema

The `players` table has the following structure:

```sql
CREATE TABLE players (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    goals INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Dependencies

- `github.com/gocolly/colly` - Web scraping
- `github.com/gorilla/mux` - HTTP routing
- `github.com/mattn/go-sqlite3` - SQLite database driver 