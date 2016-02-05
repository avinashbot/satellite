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

// Set the background on windows.
func Set(img image.Image) error {
	// Get the absolute path of the directory.
	usr, err := user.Current()
	if err != nil {
		return err
	}
	imgPath := filepath.Join(usr.HomeDir, "AppData", "Roaming", "Himawari", "background.png")

	// Create the file.
	if err := createFile(img, imgPath); err != nil {
		return err
	}

	// Set the background, hoping it worked.
	proc.Call(
		spiSETDESKWALLPAPER,
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(imgPath))),
		spifUPDATEINIFILE,
	)
	return nil
}
