package main

import (
	"log"
	"os"

	"github.com/icholy/draw"
)

func main() {
	cv := draw.NewCanvas(80, 40)

	a := draw.Point{40, 20}
	b := draw.Point{20, 0}
	c := draw.Circle{a, 10}
	l := draw.Line{a, b}
	t := draw.Text{draw.Point{10, 10}, "Hello World!"}

	cv.Draw(c, '%')
	cv.Draw(l, '*')
	cv.Draw(a, 'A')
	cv.Draw(b, 'B')
	cv.Draw(t, 0)

	if err := cv.WriteTo(os.Stdout); err != nil {
		log.Fatal(err)
	}
}
