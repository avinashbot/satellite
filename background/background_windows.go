package background

import (
	"image"
	"os/user"
	"path/filepath"
	"syscall"
	"unsafe"
)

const (
	spiSETDESKWALLPAPER = 0x14
	spifUPDATEINIFILE   = 0x2
)

var (
	user32 = syscall.MustLoadDLL("user32.dll")
	proc   = user32.MustFindProc("SystemParametersInfoW")
)

// PlatformDownload dowloads the image to the preferred location for the
// platform and returns the path it downloaded to.
func PlatformDownload(img image.Image) (string, error) {
	// Get the absolute path of the directory.
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	absPath := filepath.Join(usr.HomeDir, "AppData", "Roaming", "Satellite", "background.png")

	// Create the file.
	return absPath, DownloadOnly(img, absPath)
}

// Set the background on windows.
func Set(absPath string) error {
	// Set the background, hoping it worked.
	proc.Call(
		spiSETDESKWALLPAPER,
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(absPath))),
		spifUPDATEINIFILE,
	)
	return nil
}
