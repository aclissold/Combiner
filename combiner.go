package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

// Converts a collection of 32x32 .jpg sprites into a single png spritesheet.
// The .jpg files must be contained in the same directory as this program
// and named 0.jpg, 1.jpg, etc.
func main() {
	setup()
	read()
	write()
}

var sprites, cols int
var dst *image.RGBA
var square image.Rectangle

func setup() {
	// Find the number of sprites and create the spritesheet
	fmt.Print("Enter the number of sprites to combine: ") // largest#.jpg + 1
	_, err := fmt.Scanf("%d", &sprites)
	fmt.Print("Enter the number of sprites per row: ")
	_, err = fmt.Scanf("%d", &cols)
	if err != nil {
		log.Fatal(err)
	}
	dst = image.NewRGBA(image.Rect(0, 0, 32*cols, (32*sprites)/cols))

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
