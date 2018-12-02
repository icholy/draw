package draw

import (
	"bufio"
	"fmt"
	"io"
	"math"
)

// Canvas is 2D array x = 0, y = 0 is top left
type Canvas [][]byte

func NewCanvas(width, height int) Canvas {
	cv := make(Canvas, width)
	for i := range cv {
		cv[i] = make([]byte, height)
	}
	return cv
}

func (cv Canvas) At(x, y int) byte       { return cv[x][y] }
func (cv Canvas) SetAt(x, y int, b byte) { cv[x][y] = b }
func (cv Canvas) Width() int             { return len(cv) }
func (cv Canvas) Height() int            { return len(cv[0]) }

type Drawer interface {
	Draw(Canvas, byte)
}

func (cv Canvas) Draw(d Drawer, b byte) { d.Draw(cv, b) }

func (cv Canvas) WriteTo(w io.Writer) error {
	ww := bufio.NewWriter(w)
	for y := 0; y < cv.Height(); y++ {
		if y > 0 {
			ww.WriteByte('\n')
		}
		ww.WriteByte('|')
		for x := 0; x < cv.Width(); x++ {
			b := cv.At(x, y)
			if b == 0 {
				b = ' '
			}
			ww.WriteByte(b)
		}
		ww.WriteByte('|')
	}
	return ww.Flush()
}

type Point struct{ X, Y float64 }

func (p Point) Round() (x, y int) {
	x = int(math.Round(p.X))
	y = int(math.Round(p.Y))
	return x, y
}

func (p Point) String() string {
	return fmt.Sprintf("Point(%f, %f)", p.X, p.Y)
}

func (p Point) Distance(other Point) float64 {
	xDelta := p.X - other.X
	yDelta := p.Y - other.Y
	return math.Sqrt(xDelta*xDelta + yDelta*yDelta)
}

func (p Point) Between(other Point, factor float64) Point {
	xDelta := other.X - p.X
	yDelta := other.Y - p.Y
	return Point{
		X: p.X + (xDelta * factor),
		Y: p.Y + (yDelta * factor),
	}
}

func (p Point) Draw(cv Canvas, b byte) {
	x, y := p.Round()
	cv.SetAt(x, y, b)
}

type Line struct{ A, B Point }

func (l Line) Draw(cv Canvas, b byte) {
	var factor float64
	step := 1 / l.A.Distance(l.B)
	for factor < 1 {
		cv.Draw(l.A.Between(l.B, factor), b)
		factor += step
	}
	cv.Draw(l.B, b)
}

type Circle struct {
	P Point
	R float64
}

func (c Circle) Circumference() float64 {
	return 2 * math.Pi * c.R
}

func (c Circle) Point(factor float64) Point {
	phase := (2 * math.Pi * factor) - math.Pi
	return Point{
		X: c.P.X + math.Sin(phase)*(c.R*2),
		Y: c.P.Y + math.Cos(phase)*c.R,
	}
}

func (c Circle) Draw(cv Canvas, b byte) {
	var factor float64
	step := 1 / c.Circumference()
	for factor < 1 {
		cv.Draw(c.Point(factor), b)
		factor += step
	}
}

type Spiral struct {
	P      Point
	R      float64
	DeltaR float64
}

func (s Spiral) Draw(cv Canvas, b byte) {
	var factor float64
	radius := s.R
	deltaR := s.DeltaR
	if deltaR <= 0 {
		deltaR = 1
	}
	for radius > 0 {
		c := Circle{P: s.P, R: radius}
		cv.Draw(c.Point(factor), b)
		factor += 1 / c.Circumference()
		if factor > 1 {
			factor = 0
		}
		radius -= deltaR
	}
}

type Text struct {
	P    Point
	Text string
}

func (t Text) Draw(cv Canvas, _ byte) {
	x, y := t.P.Round()
	for i, b := range []byte(t.Text) {
		cv.SetAt(x+i, y, b)
	}
}
