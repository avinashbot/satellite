package main

import (
	"flag"
	"image"
	"log"
	"sync"
	"time"

	"github.com/avinashbot/himawari/background"
	"github.com/avinashbot/himawari/download"
)

var (
	depth int
	every int64
)

func init() {
	flag.IntVar(&depth, "depth", 4, "Resolution of the image. One of 4, 8, 16, 20.")
	flag.Int64Var(&every, "every", 0, "Optionally re-run every x seconds.")
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
				m[i][j], err = download.GridAt(t, depth, i, j)
			}(i, j)
		}
	}
	wg.Wait()
	log.Printf("Done! Downloading images took %s.", time.Now().Sub(startTime))

	// Join the pieces and set the background image.
	// A depth=20 crashed my computer around this part, so watch out.
	var img image.Image
	img = background.Join(m, 550*depth, 550*depth)
	img = background.Expand(img, 16/9)

	log.Println("Setting image as background...")
	return background.Set(img)
}

func main() {
	flag.Parse()

	// Get the latest timestamp.
	t, err := download.Latest()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("The latest image is at %s.", t)

	// Run the first time.
	if err := run(t); err != nil {
		log.Fatalln(err)
	}

	// Run this baby repeatedly.
	if every > 0 {
		ticker := time.NewTicker(time.Duration(every) * time.Second)
		for {
			// Get the latest timestamp for the second time.
			newt, err := download.Latest()
			if err != nil {
				log.Println(err)
				continue // Restart the download.
			}

			if !newt.Equal(*t) {
				t = newt
				log.Printf("The latest image is at %s.", newt)
				if err := run(t); err != nil {
					log.Print(err)
				}
			}

			// Wait until the next tick.
			<-ticker.C
		}
	}
}
