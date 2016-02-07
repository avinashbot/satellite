package download

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"net/http"
	"sync"
	"time"

	_ "image/png" // To parse PNG images.
)

const (
	// The width and height of the returned grid part.
	gridSize = 550

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

// Return the image at a certain row and column.
func (h Himawari) gridAt(filename string, row, col int) (image.Image, error) {
	// Get the image from the url.
	res, err := http.Get(fmt.Sprintf(hwImagePath, h.Depth, filename, row, col))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Create an image from the response.
	img, _, err := image.Decode(res.Body)
	return img, err
}

// Join the matrix of images into a single image.
func (h Himawari) joinGrid(grid [][]image.Image) image.Image {
	dst := image.NewNRGBA(image.Rect(0, 0, gridSize*h.Depth, gridSize*h.Depth))

	var wg sync.WaitGroup
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == nil {
				break
			}
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				// Modification is performed pixel-by-pixel, so async
				// should be fine, theoretically.
				src := grid[i][j]
				sdx, sdy := src.Bounds().Dx(), src.Bounds().Dy()
				rect := image.Rect(
					i*sdx, j*sdy, i*sdx+sdx, j*sdy+sdy,
				)
				draw.Draw(dst, rect, src, image.ZP, draw.Src)
			}(i, j)
		}
	}
	wg.Wait()

	return dst
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
	// Make the image grid.
	m := make([][]image.Image, h.Depth)
	for i := range m {
		m[i] = make([]image.Image, h.Depth)
	}

	// Download images in parallel. TODO: handle potential errors somehow?
	var wg sync.WaitGroup
	for i := 0; i < h.Depth; i++ {
		for j := 0; j < h.Depth; j++ {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				m[i][j], _ = h.gridAt(filename, i, j)
			}(i, j)
		}
	}
	wg.Wait()

	// Join the pieces.
	// A depth of 20 crashed my VM around this part, so watch out.
	return h.joinGrid(m), nil
}
