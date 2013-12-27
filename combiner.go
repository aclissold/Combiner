package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
)

// Converts a collection of 32x32 .jpg sprites into a single png spritesheet.
// The .jpg files must be named 0.jpg, 1.jpg, etc.
func main() {
	setup()
	read()
	write()
}

var sprites, cols int
var dst *image.RGBA
var square image.Rectangle

func setup() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: combiner <dirname>")
		os.Exit(1)
	}
	if err := os.Chdir(os.Args[1]); err != nil {
		log.Fatal(err)
	}
	// Find the number of sprites and create the spritesheet
	for i := 0; ; i++ {
		filename := fmt.Sprintf("%d.jpg", i)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			break
		}
		sprites = i + 1
	}
	fmt.Print("Enter the number of sprites per row: ")
	_, err := fmt.Scanf("%d", &cols)
	if err != nil {
		log.Fatal(err)
	}
	x1 := 32 * cols
	y1 := int(32 * math.Ceil(float64(sprites)/float64(cols)))
	dst = image.NewRGBA(image.Rect(0, 0, x1, y1))

	// Create a square to use in specifying the size of the sprite to write
	square = image.Rect(0, 0, 32, 32)
}

func read() {
	fmt.Println("Reading sprites into memory...")
	for i := 0; i < sprites; i += cols {
		square.Min.X = 0
		square.Max.X = 32
		for j := 0; j < cols; j++ {
			// Avoid attempting to open files that don't exist
			if i+j == sprites {
				return
			}

			// Get a sprite name
			filename := fmt.Sprintf("%d.jpg", i+j)

			// Open its corresponding file
			file, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}

			// Decode it
			src, err := jpeg.Decode(file)
			if err != nil {
				log.Fatal(err)
			}
			file.Close()

			// Draw the sprite onto the spritesheet
			draw.Draw(dst, square, src, image.ZP, draw.Src)
			square.Min.X += 32
			square.Max.X += 32
		}
		square.Min.Y += 32
		square.Max.Y += 32
	}
}

func write() {
	// Open a png for writing and write the spritesheet to it
	filename := "spritesheet.png"
	fmt.Println("Writing to", filename+"...")
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(file, dst)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")
}
