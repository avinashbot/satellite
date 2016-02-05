package background

import (
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"sync"
)

// Join the matrix of image pointers into one massive image.
// TODO: center the image in a 1080p image.
func Join(mat [][]image.Image, width, height int) image.Image {
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))

	var wg sync.WaitGroup
	for i := range mat {
		for j := range mat[i] {
			if mat[i][j] == nil {
				break
			}

			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				// Modification is performed pixel-by-pixel, so async
				// should be fine, theoretically.
				src := mat[i][j]
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

// Expand the square-y image into a more widescreen-friendly version.
// TODO(avinash): Do it.
func Expand(img image.Image, ratio float64) image.Image {
	return img
}

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
