package main

import (
	"flag"
	"log"
	"time"

	"github.com/avinashbot/himawari/background"
	"github.com/avinashbot/himawari/download"
)

const (
	// The width and height of the returned grid part.
	gridSize = 550
)

var (
	satellite string
	depth     int
	every     time.Duration
)

func init() {
	flag.StringVar(&satellite, "satellite", "himawari", `The satellite to use: "himawari" or "dscovr".`)
	flag.IntVar(&depth, "depth", 4, "Resolution of the Himawari image. One of 4, 8, 16, 20.")
	flag.DurationVar(&every, "every", 0, "Time to wait between each rerun.")
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

	// Start off with a zero time.
	for lastTime := (time.Time{}); ; time.Sleep(every) {
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

		// Set the image as the background.
		// This one's a serious error, so break if it happens.
		log.Println("Setting image as background...")
		if err := background.Set(img); err != nil {
			log.Println(err)
			break
		}

		// Success. Replace lastTime with the current time.
		lastTime = time.Now()

		// If we're only doing this once, quit.
		if every == 0 {
			break
		}
	}
}
