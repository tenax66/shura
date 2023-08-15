package main

import (
	"fmt"
	"math/rand"
)

func main() {
	run()
}

func run() {
	address := generateOnionAddress()

	fmt.Println(address)
}

func generateOnionAddress() string {
	const addressLength = 56
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

	var url string
	for i := 0; i < addressLength; i++ {
		url += string(letters[rand.Intn(len(letters))])
	}

	return url + ".onion"
}
