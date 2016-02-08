package background

import (
	"image"
	"os"
	"os/exec"
	"path/filepath"
)

func setGnome3(imgPath string) error {
	// Darken background area
	exec.Command(
		"gsettings", "set", "org.gnome.desktop.background",
		"primary-color", "#000000",
	).Run()

	// Set background mode (again, testing on gnome3)
	exec.Command(
		"gsettings", "set", "org.gnome.desktop.background",
		"picture-options", "scaled",
	).Run()

	// Set the background (gnome3 only atm)
	return exec.Command(
		"gsettings", "set", "org.gnome.desktop.background",
		"picture-uri", "file://"+imgPath,
	).Run()
}

// Set the background on windows.
func Set(img image.Image) error {
	// Get the absolute path of the directory.
	imgPath := filepath.Join(os.Getenv("HOME"), ".local", "share", "himawari", "background.png")

	// Create the file.
	if err := createFile(img, imgPath); err != nil {
		return err
	}

	return setGnome3(imgPath)
}
