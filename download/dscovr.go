package download

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"time"

	_ "image/png" // To parse PNG images.
)

const (
	dsUpdatePath   = "http://epic.gsfc.nasa.gov/api/images.php"
	dsUpdateFormat = "2006-01-02 15:04:05"
	dsImagePath    = "http://epic.gsfc.nasa.gov/epic-archive/png/%s.png"
)

// Dscovr is a satellite launched by SpaceX.
// See https://en.wikipedia.org/wiki/Deep_Space_Climate_Observatory
type Dscovr struct{}

// ModifiedSince returns a filename if the string was modified since 'after'.
func (d Dscovr) ModifiedSince(after time.Time) (string, error) {
	// Get the response from the url.
	res, err := http.Get(dsUpdatePath)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// JSON decoding...
	dec := json.NewDecoder(res.Body)
	// Read opening bracket
	if _, err := dec.Token(); err != nil {
		return "", err
	}

	latestImage := ""
	latestDate := time.Time{}

	// Incrementally parse and update latestDate & latestImage.
	var ts struct {
		Image string `json:"image"`
		Date  string `json:"date"`
	}
	for dec.More() {
		if err := dec.Decode(&ts); err != nil {
			return "", err
		}
		tsDate, err := time.Parse(dsUpdateFormat, ts.Date)
		if err != nil {
			return "", err
		}
		if tsDate.After(latestDate) {
			latestImage = ts.Image
			latestDate = tsDate
		}
	}

	// Read closing bracket
	if _, err := dec.Token(); err != nil {
		return "", err
	}

	// Check if the date is after the latest date.
	if latestDate.After(after) {
		return latestImage, nil
	}
	return "", nil
}

// Download the image. Duh.
func (d Dscovr) Download(filename string) (image.Image, error) {
	// Construct the url.
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, dsImagePath, filename)

	// Get the response from the url.
	res, err := http.Get(buf.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Create an image from the response.
	img, _, err := image.Decode(res.Body)
	return img, err
}
