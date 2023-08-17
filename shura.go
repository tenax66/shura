package shura

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/pkg/errors"
)

func Run() (string, error) {
	address := generateOnionAddress()

	return fetchContent(address)
}

func fetchContent(address string) (string, error) {
	resp, err := http.Get(address)

	if err != nil {
		// TODO: error handling
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		// TODO: error handling
		return "", errors.New(fmt.Sprintf("unsuccessful request to %s: %s", address, resp.Status))
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)

	if err != nil {
		// TODO error handling
		return "", err
	}

	body := string(bytes)

	return body, nil
}

func generateOnionAddress() string {
	const addressLength = 56
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

	var url string
	for i := 0; i < addressLength; i++ {
		url += string(letters[rand.Intn(len(letters))])
	}

	return "http://" + url + ".onion"
}