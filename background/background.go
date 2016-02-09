package background

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
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
