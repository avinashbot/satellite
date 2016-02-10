package main

import (
	"flag"
	"log"
	"path/filepath"
	"time"

	"github.com/avinashbot/satellite/background"
	"github.com/avinashbot/satellite/download"
)

var (
	satellite    string
	depth        int
	every        time.Duration
	downloadPath string
	dontSet      bool
)

func init() {
	flag.StringVar(&satellite, "use", "himawari", `The satellite to use: "himawari" or "dscovr".`)
	flag.IntVar(&depth, "depth", 4, "Resolution of the Himawari image. One of 4, 8, 16, 20.")
	flag.DurationVar(&every, "every", 0, "Time to wait between each rerun.")

	flag.StringVar(&downloadPath, "path", "", "The path to download to (default is platform-dependent).")
	flag.BoolVar(&dontSet, "dontset", false, "Just download the image.")
}

func main() {
	flag.Parse()

	// Set the satellite.
	var dl download.Downloader
	switch satellite {
	case "himawari":
		dl = download.Himawari{Depth: depth}
	case "dscovr":
		dl = download.Dscovr{}
	default:
		log.Fatalln("Satellite not recognized. Exiting.")
	}

	// If a path is supplied, absolute-ify it.
	imgPath := ""
	if downloadPath != "" {
		var err error
		if imgPath, err = filepath.Abs(downloadPath); err != nil {
			log.Fatalln(err)
		}
	}

	// Start off with a zero time.
	lastTime := time.Time{}

	for ; ; time.Sleep(every) {
		// Get the filename to the latest image.
		filename, err := dl.ModifiedSince(lastTime)
		if err != nil {
			log.Println(err)
			continue
		}
		if filename == "" {
			log.Println("No changes since last time. Trying again later...")
			continue
		}

		// There is new image out. Download it.
		log.Println("Starting download...")
		benchmarkTime := time.Now()
		img, err := dl.Download(filename)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("Done! Download took %s.\n", time.Now().Sub(benchmarkTime))

		// Download the image to the location preferred by the platform or user.
		backgroundPath := ""
		if imgPath != "" {
			backgroundPath = imgPath
			if err := background.DownloadOnly(img, imgPath); err != nil {
				log.Fatalln(err)
			}
		} else {
			var err error
			backgroundPath, err = background.PlatformDownload(img)
			if err != nil {
				log.Fatalln(err)
			}
		}

		// Set the image as the background (if needed).
		if !dontSet {
			log.Println("Setting image as background...")
			if err := background.Set(backgroundPath); err != nil {
				log.Fatalln(err)
			}
		}

		// Success. Replace lastTime with the current time.
		lastTime = time.Now()

		// If we're only doing this once, quit.
		if every == 0 {
			break
		}
	}
}
