package dscovr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"time"

	_ "image/jpeg" // Decode JPEG images from imagePath
)

const (
	updatePath   = "http://epic.gsfc.nasa.gov/api/images.php"
	updateFormat = "2006-01-02 15:04:05"
	imagePath    = "http://epic.gsfc.nasa.gov/epic-archive/jpg/%s.jpg"
)

// Latest image name.
func Latest() (string, error) {
	// Get the response from the url.
	res, err := http.Get(updatePath)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Convert the response to bytes.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	// Convert from JSON.
	var ts []struct {
		Image string `json:"image"`
		Date  string `json:"date"`
	}
	if err := json.Unmarshal(body, &ts); err != nil {
		return "", err
	}

	// Now we iterate through the slice and find the newest image.
	latestImage := ""
	latestDate := time.Unix(0, 0)
	for _, obj := range ts {
		objDate, err := time.Parse(updateFormat, obj.Date)
		if err != nil {
			return "", err
		}
		if objDate.After(latestDate) {
			latestImage = obj.Image
			latestDate = objDate
		}
	}

	return latestImage, nil
}

// GetImage returns the image with the given name. No stitching this time!
func GetImage(name string) (image.Image, error) {
	// Construct the url.
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, imagePath, name)

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
