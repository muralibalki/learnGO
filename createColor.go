package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())

	for i := 1; i < 10; i++ {
		outFile := fmt.Sprint("/Users/muralib/go/output/output", strconv.Itoa(i), ".jpg")
		out, err := os.Create(outFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// generate some QR code look a like image
		c := color.Gray{uint8(128)}
		imgRect := image.Rect(0, 0, 100, 100)
		img := image.NewGray(imgRect)
		draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
		for y := 0; y < 100; y += 10 {
			for x := 0; x < 100; x += 10 {
				fill := &image.Uniform{c}
				if rand.Intn(10)%2 == 0 {
					fill = &image.Uniform{color.White}
				}
				draw.Draw(img, image.Rect(x, y, x+10, y+10), fill, image.ZP, draw.Src)
			}
		}

		var opt jpeg.Options

		opt.Quality = 80
		// ok, write out the data into the new JPEG file

		err = jpeg.Encode(out, img, &opt) // put quality to 80%
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Generated image to output.jpg \n")
	}
}
