package background

import (
	"image"
	"os/exec"
	"os/user"
	"path/filepath"
)

// Set the background on windows.
func Set(img image.Image) error {
	// Get the absolute path of the directory.
	usr, err := user.Current()
	if err != nil {
		return err
	}
	imgPath := filepath.Join(usr.HomeDir, ".himawari", "background.png")

	// Create the file.
	if err := createFile(img, imgPath); err != nil {
		return err
	}

	// Set the background (gnome3 only atm)
	err = exec.Command(
		"gsettings", "set", "org.gnome.desktop.background",
		"picture-uri", "file://"+imgPath,
	).Run()

	// Set background mode (again, testing on gnome3)
	err = exec.Command(
		"gsettings", "set", "org.gnome.desktop.background",
		"picture-options", "scaled",
	).Run()

	return err
}
