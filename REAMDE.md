# Basic 2d drawing primitives for the terminal

[![GoDoc](https://godoc.org/github.com/icholy/draw?status.svg)](https://godoc.org/github.com/icholy/draw)

``` go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/icholy/draw"
)

func main() {
	cv := draw.NewCanvas(80, 40)

	// draw a point
	p := draw.Point{
		X: 10,
		Y: 10,
	}
	cv.Draw(p, '$')

	// draw a line
	l := draw.Line{
		A: cv.Bounds().TopRight(),
		B: draw.Point{
			X: 40,
			Y: 5,
		},
	}
	cv.Draw(l, '*')

	// add a border
	cv.Draw(cv.Bounds().Box(), 0)

	// draw text
	t := draw.Text{
		Origin: cv.Bounds().TopLeft().AddXY(20, 10),
		Text:   "You can draw\nmulti-line text",
	}
	cv.Draw(t, 0)

	// add border to text
	cv.Draw(draw.BoxAround(t), 0)

	// draw circle
	c := draw.Circle{
		Center: cv.Bounds().Center().AddXY(5, 5),
		Radius: 10,
	}
	cv.Draw(c, '%')

	if err := cv.WriteTo(os.Stdout); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}
```