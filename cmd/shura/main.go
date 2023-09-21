package main

import (
	"os"

	"github.com/tenax66/shura"
)

func main() {
	os.Setenv("HTTP_PROXY", "socks5://localhost:9050")

	links, _ := shura.LoadAllSavedLinks()

	shura.Collect(links)
}
