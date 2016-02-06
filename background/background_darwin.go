package background

import (
	"image"
	"os/exec"
	"os/user"
	"path/filepath"
)

// Set the background on OSX.
func Set(img image.Image) error {
	// Get the absolute path of the directory.
	usr, err := user.Current()
	if err != nil {
		return err
	}
	imgPath := filepath.Join(usr.HomeDir, "Pictures", "himawari.png")

	// Create the file.
	if err := createFile(img, imgPath); err != nil {
		return err
	}

	// Only tested on El Capitan
	err = exec.Command("osascript", "-e", `tell application "Finder" to set desktop picture to POSIX file "`+imgPath+`"`).Run()

	return err
}
