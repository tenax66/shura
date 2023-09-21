package shura

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

func init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar = logger.Sugar()
}

func Collect(startUrls []string) {

	db, err := initDB("links.db")
	if err != nil {
		sugar.Error("Error initializing table:", err)
		return
	}

	for _, url := range startUrls {
		links, err := Extract(url)
		if err != nil {
			sugar.Error("Error creating table:", err)
			return
		}

		for _, link := range links {
			_, err = db.Exec("INSERT INTO links (link) VALUES (?)", link)
			if err != nil {
				//TODO: ignore UNIQUE constraint errors
				fmt.Println("Error inserting link:", err)
				return
			}
		}
	}
}

func Extract(url string) ([]string, error) {

	regex := regexp.MustCompile(`(http|https)://[a-z2-7]{56}\.onion`)
	resp, err := http.Get(url)
	if err != nil {
		sugar.Error("Error fetching the URL: ", err)
		return []string{""}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		sugar.Error("Error reading response body: ", err)
		return []string{""}, err
	}

	return extractLinks(string(body), regex), nil

}

func extractLinks(html string, regex *regexp.Regexp) []string {
	matches := regex.FindAllStringSubmatch(html, -1)

	var links []string
	for _, match := range matches {
		links = append(links, match[0])
	}

	return links
}

func initDB(name string) (*sql.DB, error){
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		sugar.Info("Error connecting to the database:", err)
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS links (link TEXT PRIMARY KEY)")
	if err != nil {
		sugar.Error("Error creating table:", err)
		return nil, err
	}

	return db, err
}
