package main

import (
	"os"

	"github.com/tenax66/shura"
)

func main() {
	os.Setenv("HTTP_PROXY", "socks5://localhost:9050")

	const REPEATS = 10000
	shura.Run(REPEATS)

}
