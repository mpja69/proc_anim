package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Point struct {
	x, y float64
}

func distance(p, q Point) float64 {
	return math.Sqrt(math.Pow(p.x-q.x, 2) + math.Pow(p.y-q.y, 2))
}
func angle(p Point) float64 {
	return math.Atan2(p.y, p.x)
}
func (p Point) Angle(q Point) float64 {
	return angle(p.Sub(q))
}
func SetConstraint(pos, anchor Point, constraint float64) Point {
	delta := pos.Sub(anchor)
	deltaConstrained := delta.SetMag(constraint)
	return anchor.Add(deltaConstrained)
}
func (p Point) Add(q Point) Point {
	p.x += q.x
	p.y += q.y
	return p
}

func (p Point) Sub(q Point) Point {
	x := p.x - q.x
	y := p.y - q.y
	return Point{x, y}
}

func (p Point) Mag() float64 {
	return math.Sqrt(math.Pow(p.x, 2) + math.Pow(p.y, 2))
}

func (p Point) SetMag(m float64) Point {
	angle := math.Atan2(p.y, p.x)
	x := m * math.Cos(angle)
	y := m * math.Sin(angle)
	return Point{x, y}
}

type Joint struct {
	pos        Point
	distance   float64
	radius     float64
	color      color.RGBA
	angle      float64
	adjustment float64
}

func JointNew(x, y, distance, radius float64, color color.RGBA) *Joint {
	return &Joint{
		distance:   distance,
		pos:        Point{x, y},
		radius:     radius,
		color:      color,
		adjustment: 1.0,
	}
}

// Rotates the segment towards the segment it follows. And translates it to the distance from the segment it follows
func (s *Joint) DirectlyFollow(prev Point) {
	// Turn
	targetAngle := math.Atan2(prev.y-s.pos.y, prev.x-s.pos.x)
	// delta := targetAngle - s.angle
	// if delta < -math.Pi*0.5 {
	// 	targetAngle = -math.Pi * 0.5
	// }
	// if delta > math.Pi*0.5 {
	// 	targetAngle = math.Pi * 0.5
	// }
	s.angle = targetAngle

	// Move
	dist := math.Sqrt(math.Pow(prev.x-s.pos.x, 2) + math.Pow(prev.y-s.pos.y, 2))
	if dist > s.distance {
		delta := dist - s.distance
		s.pos.x += delta * math.Cos(s.angle)
		s.pos.y += delta * math.Sin(s.angle)

	}
}

// Updates the angle of the segment, (to make the endpoint align with the target position)
func (s *Joint) CCDIKUpdateAngle(end, target Point) {
	angle := math.Atan2(end.y-s.pos.y, end.x-s.pos.x)
	targetAngle := math.Atan2(target.y-s.pos.y, target.x-s.pos.x)
	delta := targetAngle - angle
	for delta < -math.Pi {
		delta += 2 * math.Pi
	}
	for delta > math.Pi {
		delta -= 2 * math.Pi
	}
	s.angle += delta * s.adjustment * 0.1 // Tip: Don't change if "end arm" is close to target
}

func (s *Joint) Left() Point {
	angle := s.angle - math.Pi*0.5

	x := s.pos.x + s.radius*math.Cos(angle)
	y := s.pos.y + s.radius*math.Sin(angle)

	return Point{x, y}
}

func (s *Joint) Right() Point {
	angle := s.angle + math.Pi*0.5

	x := s.pos.x + s.radius*math.Cos(angle)
	y := s.pos.y + s.radius*math.Sin(angle)

	return Point{x, y}
}

func (s *Joint) End() Point {
	x1 := s.pos.x + math.Cos(s.angle)*s.distance
	y1 := s.pos.y + math.Sin(s.angle)*s.distance
	return Point{x1, y1}
}

func (s *Joint) getAdjustedPos(angleOffset, lengthOffset float64) Point {
	x := s.pos.x + math.Cos(s.angle+angleOffset)*(s.radius+lengthOffset)
	y := s.pos.y + math.Sin(s.angle+angleOffset)*(s.radius+lengthOffset)
	return Point{x, y}
}

func (s *Joint) DrawCircle(screen *ebiten.Image) {
	vector.StrokeCircle(screen,
		float32(s.pos.x), float32(s.pos.y),
		float32(s.radius), 2,
		color.RGBA{255, 255, 255, 255}, true)

}
