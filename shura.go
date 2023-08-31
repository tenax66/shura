package shura

import (
	"io"
	"net/http"
	"regexp"

	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

func init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar = logger.Sugar()
}

func Run(url string) {
	regex := regexp.MustCompile(`http.*[a-z2-7]{56}\.onion`)
	resp, err := http.Get(url)
	if err != nil {
		sugar.Error("Error fetching the URL: ", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		sugar.Error("Error reading response body: ", err)
		return
	}

	links := extractLinks(string(body), regex)

	for _, link := range links {
		// TODO: save url
		sugar.Info("URL: ", link)
	}
}

func extractLinks(html string, regex *regexp.Regexp) []string {
	matches := regex.FindAllStringSubmatch(html, -1)

	var links []string
	for _, match := range matches {
		links = append(links, match[0])
	}

	return links
}
