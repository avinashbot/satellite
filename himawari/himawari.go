package himawari

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"time"

	_ "image/png" // Decode PNG images from imagePath
)

const (
	updatePath   = "http://himawari8-dl.nict.go.jp/himawari8/img/D531106/latest.json"
	updateFormat = "2006-01-02 15:04:05"
	imagePath    = "http://himawari8.nict.go.jp/img/D531106/%dd/550/%s_%d_%d.png"
	imageFormat  = "2006/01/02/150405"
)

// Latest image time.
func Latest() (*time.Time, error) {
	// Get the response from the url.
	res, err := http.Get(updatePath)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Convert the response to bytes.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Convert to JSON and get the date string.
	ts := &struct {
		Date string `json:"date"`
	}{}
	if err := json.Unmarshal(body, ts); err != nil {
		return nil, err
	}

	// Parse the date string.
	latestTime, err := time.Parse(updateFormat, ts.Date)
	return &latestTime, err
}

// GridAt returns the image at a certain row and column.
func GridAt(t *time.Time, depth, row, col int) (image.Image, error) {
	// Construct the url.
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, imagePath, depth, t.Format(imageFormat), row, col)

	// Get the image from the url.
	res, err := http.Get(buf.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Create an image from the response.
	img, _, err := image.Decode(res.Body)
	return img, err
}
