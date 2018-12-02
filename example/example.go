package main

import (
	"fmt"
	"log"
	"os"

	"github.com/icholy/draw"
)

func main() {
	cv := draw.NewCanvas(80, 40)

	cv.Draw(cv.Bounds().Border(), 0)

	c := draw.Circle{
		Center: cv.Center(),
		Radius: 15,
	}
	cv.Draw(c, '*')

	cv.Draw(c.Bounds().Border(), '!')

	if err := cv.WriteTo(os.Stdout); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}
