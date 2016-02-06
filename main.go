package main

import (
	"flag"
	"image"
	"log"
	"sync"
	"time"

	"github.com/avinashbot/himawari/background"
	"github.com/avinashbot/himawari/himawari"
)

const (
	// The width and height of the returned grid part.
	gridSize = 550
)

var (
	depth int
	every int64
	once  bool
	dsc   bool
)

func init() {
	flag.IntVar(&depth, "depth", 4, "Resolution of the image. One of 4, 8, 16, 20.")
	flag.Int64Var(&every, "every", 600, "Re-run every x seconds.")
	flag.BoolVar(&once, "once", false, "Set the background and exit.")
	flag.BoolVar(&dsc, "dscovr", false, "Use DSCOVR imagery. It's not geostationary though.")
}

func run(t *time.Time) error {
	// Make the image grid.
	m := make([][]image.Image, depth)
	for i := range m {
		m[i] = make([]image.Image, depth)
	}

	// Download images in parallel. Woo go!
	log.Println("Starting download...")
	startTime := time.Now()
	var wg sync.WaitGroup
	var err error
	for i := 0; i < depth; i++ {
		for j := 0; j < depth; j++ {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				m[i][j], err = himawari.GridAt(t, depth, i, j)
			}(i, j)
		}
	}
	wg.Wait()
	log.Printf("Done! Downloading images took %s.\n", time.Now().Sub(startTime))

	// Join the pieces and set the background image.
	// A depth=20 crashed my VM around this part, so watch out.
	var img image.Image
	img = background.Join(m, gridSize*depth, gridSize*depth)
	img = background.Expand(img, 16/9) // FIXME

	log.Println("Setting image as background...")
	return background.Set(img)
}

func main() {
	flag.Parse()

	// Start off with a dummy time.
	t := time.Unix(0, 0)

	// Run the program.
	for ticker := time.NewTicker(time.Duration(every) * time.Second); ; <-ticker.C {
		// Get the latest timestamp.
		newt, err := himawari.Latest()
		if err != nil {
			continue // The update server threw an error. Try later.
		}
		log.Printf("The latest image is at %s.\n", newt)

		// Skip if the latest time hasn't changed.
		if newt.Equal(t) {
			log.Println("The image has not changed. Waiting...")
			continue
		}

		// If it has, run it.
		if err = run(newt); err != nil {
			log.Println(err)
			continue
		}

		// Everything has succeeded, set the time to the latest image time.
		t = *newt

		// If we're only running this once, exit.
		if once {
			break
		}
	}
}
