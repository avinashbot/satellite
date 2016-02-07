package download

import (
	"encoding/json"
	"image"
	"net/http"
	"time"

	_ "image/png" // To parse PNG images.
)

const (
	hwUpdatePath   = "http://himawari8-dl.nict.go.jp/himawari8/img/D531106/latest.json"
	hwUpdateFormat = "2006-01-02 15:04:05"
	hwImagePath    = "http://himawari8.nict.go.jp/img/D531106/%dd/550/%s_%d_%d.png"
	hwImageFormat  = "2006/01/02/150405"
)

// Himawari is a japanese weather satellite whose images are a pain in the
// ass to parse.
// See https://en.wikipedia.org/wiki/Himawari_8
type Himawari struct {
	// Depth is the resolution depth. One of 4, 8, 16, 20.
	Depth int
}

// Join the matrix of images into a single image.
// TODO: do.
func (h Himawari) joinGrid(grid [][]image.Image) image.Image {
	return nil
}

// ModifiedSince returns a filename if the string was modified since 'after'.
func (h Himawari) ModifiedSince(after time.Time) (string, error) {
	// Get the response from the url.
	res, err := http.Get(hwUpdatePath)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// JSON decoding...
	dec := json.NewDecoder(res.Body)
	var ts struct {
		Image string `json:"image"`
		Date  string `json:"date"`
	}
	if err := dec.Decode(&ts); err != nil {
		return "", err
	}

	// Parse the date string.
	latestDate, err := time.Parse(hwUpdateFormat, ts.Date)
	if err != nil {
		return "", err
	}

	// Check if the date is after the latest date.
	if latestDate.After(after) {
		return latestDate.Format(hwImageFormat), nil
	}
	return "", err
}

// Download the grid images and stitch them.
// TODO: do.
func (h Himawari) Download(filename string) (image.Image, error) {
	return nil, nil
}
