package background

import (
	"image"
	"os"
	"os/exec"
	"path/filepath"
)

// PlatformDownload dowloads the image to the preferred location for the
// platform and returns the path it downloaded to.
func PlatformDownload(img image.Image) (string, error) {
	// Get the absolute path of the directory.
	homePath := os.Getenv("HOME")
	absPath := filepath.Join(homePath, "Pictures", "satellite.png")

	// Create the file.
	return absPath, DownloadOnly(img, absPath)
}

// Set the background on OSX.
func Set(absPath string) error {
	// Only tested on El Capitan
	if err := exec.Command("osascript", "-e", `tell application "Finder" to set desktop picture to POSIX file "`+absPath+`"`).Run(); err != nil {
		return err
	}
	return exec.Command("killall", "Dock").Run()
}
