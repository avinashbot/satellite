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

func setMate(imgPath string) error {
	// Darken background area
	exec.Command(
		"gsettings", "set", "org.mate.background",
		"primary-color", "#000000",
	).Run()

	// Enable solid background color
	exec.Command(
		"gsettings", "set", "org.mate.background",
		"color-shading-type", "solid",
	).Run()

	// Set background mode
	exec.Command(
		"gsettings", "set", "org.mate.background",
		"picture-options", "scaled",
	).Run()

	// Set the background
	return exec.Command(
		"gsettings", "set", "org.mate.background",
		"picture-filename", imgPath,
	).Run()
}

// Set the background on linux.
func Set(img image.Image) error {
	// Get the absolute path of the directory.
	imgPath := filepath.Join(os.Getenv("HOME"), ".local", "share", "himawari", "background.png")

	// Create the file.
	if err := createFile(img, imgPath); err != nil {
		return err
	}

	switch os.Getenv("XDG_CURRENT_DESKTOP") {
	case "GNOME":
		return setGnome3(imgPath)
	case "MATE":
		return setMate(imgPath)
	}

	return ErrDEUnsupported
}
