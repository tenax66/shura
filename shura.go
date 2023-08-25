package shura

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
)

func Run() {

	const TRY = 100

	for i := 0; i < TRY; i++ {
		address := generateOnionAddress()
		_, err := tryAccess(address)
		if err != nil {
			//TODO: logging
		}
	}
}

func tryAccess(address string) (string, error) {
	content, err := fetchContent(address)

	if err != nil {
		return content, err
	}

	takeScreenshot(address)

	return content, nil
}

func takeScreenshot(address string) {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	var buf []byte

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(address, 90, &buf)); err != nil {
		// logging
	}

	// Get the domain name excluding TLD (.onion)
	parsedURL, _ := url.Parse(address)
	name := strings.Split(parsedURL.Host, ".")[0]

	if err := os.WriteFile(name+".png", buf, 0o644); err != nil {
		// logging
	}
}

func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
	}
}

func fetchContent(address string) (string, error) {
	resp, err := http.Get(address)

	if err != nil {
		// TODO: error handling
		return "", err
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)

	if err != nil {
		// TODO error handling
		return "", err
	}

	body := string(bytes)

	if resp.StatusCode != http.StatusOK {
		// TODO: error handling
		return body, errors.New(fmt.Sprintf("unsuccessful request to %s: %s", address, resp.Status))
	}

	return body, nil
}

func generateOnionAddress() string {
	const addressLength = 56
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

	var name string
	for i := 0; i < addressLength; i++ {
		name += string(letters[rand.Intn(len(letters))])
	}

	return "http://" + name + ".onion"
}
