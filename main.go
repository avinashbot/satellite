package main

import (
	"image"
	"log"
	"sync"

	"github.com/avinashbot/himawari/background"
	"github.com/avinashbot/himawari/download"
)

func makeGrid(size int) [][]image.Image {
	a := make([][]image.Image, size)
	for i := range a {
		a[i] = make([]image.Image, size)
	}
	return a
}

func handle(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func main() {
	m := makeGrid(4)
	t, err := download.Latest()
	handle(err)

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				part, err := download.GridAt(t, 4, i, j)
				handle(err)
				log.Printf("Part [%d:%d] completed.\n", i, j)
				m[i][j] = part
			}(i, j)
		}
	}
	wg.Wait()
	log.Println("Done!")

	var img image.Image
	img = background.Join(m, 550*4, 550*4)
	img = background.Expand(img, 16/9)
	err = background.Set(img)
	handle(err)
}
