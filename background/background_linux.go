package background

import (
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func setGnome3(absPath string) error {
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
		"picture-uri", "file://"+absPath,
	).Run()
}

func setMate(absPath string) error {
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
		"picture-filename", absPath,
	).Run()
}

// PlatformDownload dowloads the image to the preferred location for the
// platform and returns the path it downloaded to.
func PlatformDownload(img image.Image) (string, error) {
	// Get the absolute path of the directory.
	homePath := os.Getenv("HOME")
	absPath := filepath.Join(homePath, ".local", "share", "satellite", "background.png")

	// Create the file.
	return absPath, DownloadOnly(img, absPath)
}

// Set the background on linux.
func Set(absPath string) error {
	if Desktop == "" {
		Desktop = os.Getenv("XDG_CURRENT_DESKTOP")
	}
	switch strings.ToLower(Desktop) {
	case "gnome", "x-cinnamon":
		return setGnome3(absPath)
	case "mate":
		return setMate(absPath)
	}

	return ErrDEUnsupported
}
