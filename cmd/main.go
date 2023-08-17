package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	result, err := run()

	if err != nil {
		sugar.Warn(err)
		return
	}

	fmt.Print(result)
}

func run() (string, error) {
	os.Setenv("HTTP_PROXY", "socks5://localhost:9050")

	address := generateOnionAddress()

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
