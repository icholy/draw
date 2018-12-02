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

func (cv Canvas) Bounds() Rect {
	return Rect{
		Min: Point{0, 0},
		Max: Point{
			X: float64(cv.Width()) - 1,
			Y: float64(cv.Height()) - 1,
		},
	}
}

func (cv Canvas) Center() Point {
	return Point{
		X: (float64(cv.Width()) - 1) / 2,
		Y: (float64(cv.Height()) - 1) / 2,
	}
}

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
		for x := 0; x < cv.Width(); x++ {
			b := cv.At(x, y)
			if b == 0 {
				b = ' '
			}
			ww.WriteByte(b)
		}
	}
	return ww.Flush()
}

type Point struct{ X, Y float64 }

func (p Point) Bounds() Rect {
	return Rect{p, p}
}

func (p Point) Add(other Point) Point {
	return Point{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

func (p Point) Sub(other Point) Point {
	return Point{
		X: p.X - other.X,
		Y: p.Y - other.Y,
	}
}

func (p Point) Min(other Point) Point {
	return Point{
		X: math.Min(p.X, other.X),
		Y: math.Min(p.Y, other.Y),
	}
}

func (p Point) Max(other Point) Point {
	return Point{
		X: math.Max(p.X, other.X),
		Y: math.Max(p.Y, other.Y),
	}
}

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

func (l Line) String() string {
	return fmt.Sprintf("Line(%s, %s)", l.A, l.B)
}

func (l Line) Bounds() Rect {
	return Rect{
		Min: l.A.Min(l.B),
		Max: l.A.Max(l.B),
	}
}

func (l Line) Mid() Point {
	return l.A.Between(l.B, 0.5)
}

type Orientation int

const (
	Horizonal Orientation = iota
	Verical
	Angled
)

func (l Line) Orientation() Orientation {
	const delta = 0.00001
	switch {
	case math.Abs(l.A.X-l.B.X) < delta:
		return Verical
	case math.Abs(l.A.Y-l.B.Y) < delta:
		return Horizonal
	default:
		return Angled
	}
}

func (l Line) Draw(cv Canvas, b byte) {
	switch l.Orientation() {
	case Verical:
		min := l.A.Min(l.B)
		max := l.A.Max(l.B)
		for y := min.Y; y <= max.Y; y++ {
			cv.Draw(Point{l.A.X, y}, b)
		}
	case Horizonal:
		min := l.A.Min(l.B)
		max := l.A.Max(l.B)
		for x := min.X; x <= max.X; x++ {
			cv.Draw(Point{x, l.A.Y}, b)
		}
	case Angled:
		var factor float64
		step := 1 / l.A.Distance(l.B)
		for factor < 1 {
			cv.Draw(l.A.Between(l.B, factor), b)
			factor += step
		}
		cv.Draw(l.B, b)
	default:
		panic("invalid orientation")
	}
}

type Circle struct {
	Center Point
	Radius float64
}

func (c Circle) Bounds() Rect {
	p := Point{c.Radius * 2, c.Radius}
	return Rect{
		Min: c.Center.Sub(p),
		Max: c.Center.Add(p),
	}
}

func (c Circle) Circumference() float64 {
	return 2 * math.Pi * c.Radius
}

func (c Circle) Point(t float64) Point {
	return Point{
		X: c.Center.X + math.Sin(t)*(c.Radius*2),
		Y: c.Center.Y + math.Cos(t)*c.Radius,
	}
}

func (c Circle) Draw(cv Canvas, b byte) {
	t := -math.Pi
	step := 1 / c.Circumference()
	for t <= math.Pi {
		cv.Draw(c.Point(t), b)
		t += step
	}
}

type Text struct {
	Origin Point
	Text   string
}

func (t Text) Bounds() Rect {
	width := float64(len([]byte(t.Text)))
	return Rect{
		Min: t.Origin,
		Max: t.Origin.Add(Point{width - 1, 0}),
	}
}

func (t Text) Draw(cv Canvas, _ byte) {
	x, y := t.Origin.Round()
	for i, b := range []byte(t.Text) {
		cv.SetAt(x+i, y, b)
	}
}

type Rect struct {
	Min, Max Point
}

func (r Rect) String() string {
	return fmt.Sprintf("Rect(%s, %s)", r.Min, r.Max)
}

func (r Rect) Bounds() Rect {
	return r
}

func (r Rect) Pad(n float64) Rect {
	padding := Point{n, n}
	return Rect{
		Min: r.Min.Sub(padding),
		Max: r.Max.Add(padding),
	}
}

func (r Rect) TopLeft() Point     { return r.Min }
func (r Rect) TopRight() Point    { return Point{r.Max.X, r.Min.Y} }
func (r Rect) BottomLeft() Point  { return Point{r.Min.X, r.Max.Y} }
func (r Rect) BottomRight() Point { return r.Max }
func (r Rect) Top() Line          { return Line{r.TopLeft(), r.TopRight()} }
func (r Rect) Bottom() Line       { return Line{r.BottomLeft(), r.BottomRight()} }
func (r Rect) Left() Line         { return Line{r.TopLeft(), r.BottomLeft()} }
func (r Rect) Right() Line        { return Line{r.TopRight(), r.BottomRight()} }

func (r Rect) Draw(cv Canvas, b byte) {
	cv.Draw(r.Top(), b)
	cv.Draw(r.Bottom(), b)
	cv.Draw(r.Left(), b)
	cv.Draw(r.Right(), b)
}

type Bounder interface {
	Bounds() Rect
}

func BoxAround(b Bounder) Box {
	return b.Bounds().Pad(1).Box()
}

func (r Rect) Box() Box {
	return Box{r}
}

type Box struct {
	Rect
}

func (b Box) Draw(cv Canvas, _ byte) {
	cv.Draw(b.Top(), '-')
	cv.Draw(b.Bottom(), '-')
	cv.Draw(b.Left(), '|')
	cv.Draw(b.Right(), '|')
}
