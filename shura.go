package shura

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

func Run(url string) {
	regex := regexp.MustCompile(`http.*[a-z2-7]{56}\.onion`)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching the URL:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Print(string(body))
	links := extractLinks(string(body), regex)

	for _, link := range links {
		// TODO: save url
		fmt.Println(link)
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
