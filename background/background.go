package background

import (
	"errors"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

var (
	// ErrDEUnsupported means that I haven't gotten around to implementing
	// the desktop environment you're using. Submit an issue!
	ErrDEUnsupported = errors.New("desktop environment not supported")
)

// Create the background at a given path.
func createFile(img image.Image, absPath string) error {
	// Create the directory if it doesn't exist.
	if err := os.MkdirAll(filepath.Dir(absPath), 0777); err != nil {
		return err
	}

	// Get the handle to the file.
	file, err := os.Create(absPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the image to the file.
	return png.Encode(file, img)
}
