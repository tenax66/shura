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

const LINKS = "links.db"

func init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar = logger.Sugar()
}

func Collect(startUrls []string) {

	db, err := initDB(LINKS)
	if err != nil {
		sugar.Error("Failed to initialize database:", err)
		return
	}

	defer db.Close()

	for _, url := range startUrls {
		links, err := Extract(url)
		if err != nil {
			sugar.Warn("Error occured while extracting .onion urls:", err)
		}
		for _, link := range links {
			sugar.Infof("Extracted: %s", link)
		}

		for _, link := range links {
			_, err = db.Exec("INSERT INTO links (link) VALUES (?)", link)
			if err != nil {
				//TODO: ignore UNIQUE constraint errors
				sugar.Warn("This url is already known:", err)
			}
		}
	}
}

func Extract(url string) ([]string, error) {

	sugar.Infof("Extracting from %s", url)

	regex := regexp.MustCompile(`(http|https)://[a-z2-7]{56}\.onion`)
	resp, err := http.Get(url)
	if err != nil {
		sugar.Error("Error occured while fetching the content: ", err)
		return []string{""}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		sugar.Error("Error occured while reading response body: ", err)
		return []string{""}, err
	}

	return extractLinks(string(body), regex), nil

}

func LoadAllSavedLinks() ([]string, error) {
	db, err := initDB(LINKS)
	if err != nil {
		sugar.Error("Failed to initialize database:", err)
		return nil, err
	}

	defer db.Close()

	query := "SELECT link FROM links"

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Failed to execute query:", err)
		return nil, err
	}
	defer rows.Close()

	var links []string

	for rows.Next() {
		var link string

		err := rows.Scan(&link)
		if err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func extractLinks(html string, regex *regexp.Regexp) []string {
	matches := regex.FindAllStringSubmatch(html, -1)

	var links []string
	for _, match := range matches {
		links = append(links, match[0])
	}

	return links
}

func initDB(name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		sugar.Info("Error occured while connecting to the database:", err)
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS links (link TEXT PRIMARY KEY)")
	if err != nil {
		sugar.Error("Failed to create a table:", err)
		return nil, err
	}

	return db, err
}
