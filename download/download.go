package download

import (
	"image"
	"time"
)

// A Downloader is a thing that returns a downloadable image.
type Downloader interface {
	// Return the latest image name after a given time if found.
	// If no image was found, filename is "".
	ModifiedSince(after time.Time) (filename string, err error)

	// Download the image with the given filename.
	Download(filename string) (img image.Image, err error)
}
