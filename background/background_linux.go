package background

import (
	"image"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

func setXfce4(absPath string) error {
	const (
		xfconfCmd     = "xfconf-query"
		lastImageProp = "last-image"
	)

	channelArg := []string{"--channel", "xfce4-desktop"}

	re := regexp.MustCompile(`/backdrop.*/` + lastImageProp)

	args := append(channelArg, "--list")

	// Retrieve the Xfce desktop properties
	propList, err := exec.Command(xfconfCmd, args...).Output()
	if err != nil {
		log.Fatal(err)
	}

	// Keep the properties related to each monitor and workspace
	backdropProps := re.FindAllString(string(propList), -1)

	if len(backdropProps) == 0 {
		log.Fatal("No", lastImageProp, "property found")
	}

	// For each workspace
	for _, property := range backdropProps {
		// Set the background image
		args = append(
			channelArg,
			"--property", property,
			"--set", absPath,
		)
		err := exec.Command(xfconfCmd, args...).Run()
		if err != nil {
			log.Fatal(err)
		}

		// Set the image style to 4 (Streched), so that the image fits in the
		// center of the screen
		imageStyle := strings.Replace(property, lastImageProp, "image-style", 1)

		args = append(
			channelArg,
			"--property", imageStyle,
			"--set", "4",
		)
		err = exec.Command(xfconfCmd, args...).Run()
		if err != nil {
			// Modifying the property didn't work, try to create it instead
			args = append(
				channelArg,
				"--property", imageStyle,
				"--create",
				"--type", "int",
				"--set", "4",
			)
			err = exec.Command(xfconfCmd, args...).Run()
			if err != nil {
				log.Fatal(err)
			}
		}

		// Set the primary background color to black
		rgba1 := strings.Replace(property, lastImageProp, "rgba1", 1)

		args = append(
			channelArg,
			"--property", rgba1,
			"--set", "0", // Red
			"--set", "0", // Green
			"--set", "0", // Blue
			"--set", "1", // Alpha
		)
		err = exec.Command(xfconfCmd, args...).Run()
		if err != nil {
			// Modifying the property didn't work, try to create it instead
			args = append(
				channelArg,
				"--property", rgba1,
				"--create",
				"--type", "uint16",
				"--type", "uint16",
				"--type", "uint16",
				"--type", "uint16",
				"--set", "0", // Red
				"--set", "0", // Green
				"--set", "0", // Blue
				"--set", "1", // Alpha
			)
			err = exec.Command(xfconfCmd, args...).Run()
			if err != nil {
				log.Fatal(err)
			}
		}

		// Set the color style to 0 (Solid color)
		colorStyle := strings.Replace(property, lastImageProp, "color-style", 1)

		args = append(
			channelArg,
			"--property", colorStyle,
			"--set", "0",
		)
		err = exec.Command(xfconfCmd, args...).Run()
		if err != nil {
			// Modifying the property didn't work, try to create it instead
			args = append(
				channelArg,
				"--property", colorStyle,
				"--create",
				"--type", "int",
				"--set", "0",
			)
			err = exec.Command(xfconfCmd, args...).Run()
			if err != nil {
				log.Fatal(err)
			}
		}

		// Make sure "backdrop-cycle-enable" is set to false
		backdropCycleEnable := strings.Replace(property, lastImageProp, "backdrop-cycle-enable", 1)

		args = append(
			channelArg,
			"--property", backdropCycleEnable,
			"--set", "false",
		)
		err = exec.Command(xfconfCmd, args...).Run()
		if err != nil {
			// Modifying the property didn't work, try to create it instead
			args = append(
				channelArg,
				"--property", backdropCycleEnable,
				"--create",
				"--type", "bool",
				"--set", "false",
			)
			err = exec.Command(xfconfCmd, args...).Run()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
}

func setFeh(absPath string) error {
	// Set the background
	return exec.Command(
		"feh", "--bg-max", absPath,
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
	var useDesktop string
	if CustomDesktop == "" {
		useDesktop = os.Getenv("XDG_CURRENT_DESKTOP")
	} else {
		useDesktop = CustomDesktop
	}

	// Check for desktop-specific methods.
	switch strings.ToLower(useDesktop) {
	case "gnome", "x-cinnamon":
		return setGnome3(absPath)
	case "mate":
		return setMate(absPath)
	case "xfce":
		return setXfce4(absPath)
	}

	// None found, now let's try to check if feh is installed.
	if _, err := exec.LookPath("feh"); err != nil {
		return ErrDEUnsupported
	}
	return setFeh(absPath)
}
