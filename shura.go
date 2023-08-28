package shura

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
)

func Run(repeats int) {

	for i := 0; i < repeats; i++ {
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

	if err := takeScreenshot(address); err != nil {
		return content, err
	}

	return content, nil
}

func takeScreenshot(address string) error {

	const saveDir = "screenshots"

	if err := createScreenshotsDirectory(saveDir); err != nil {
		// logging
		return err
	}

	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	var buf []byte

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(address, 90, &buf)); err != nil {
		// logging
		return err
	}

	// Get the domain name excluding TLD (.onion)
	parsedURL, _ := url.Parse(address)
	name := strings.Split(parsedURL.Host, ".")[0]

	if err := os.WriteFile(filepath.Join(saveDir, name+".png"), buf, 0o644); err != nil {
		// logging
		return err
	}

	return nil
}

func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
	}
}

func createScreenshotsDirectory(dirName string) error {
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		err := os.Mkdir(dirName, os.ModePerm)
		if err != nil {
			// logging
			return err
		}
	}

	return nil
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
