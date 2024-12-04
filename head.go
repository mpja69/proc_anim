package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ------------------------ Head --------------------------
type Head struct {
	*Joint
	speed float64
	// vertices []ebiten.Vertex
	// indices  []uint16
}

func HeadNew(segment *Joint, speed float64) Head {
	head := Head{segment, speed} //, []ebiten.Vertex{}, []uint16{}}
	return head
}

func (head Head) follow(target Point) { //, nextAngle float64) {
	// Update the angle
	targetAngle := math.Atan2(target.y-head.pos.y, target.x-head.pos.x)
	delta := targetAngle - head.angle
	for delta < -math.Pi {
		delta += 2 * math.Pi
	}
	for delta > math.Pi {
		delta -= 2 * math.Pi
	}
	head.angle += delta * 0.01

	// Testing another method
	// delta := targetAngle - nextAngle
	// if delta < -math.Pi*0.1 {
	// 	targetAngle = -math.Pi * 0.1
	// }
	// if delta > math.Pi*0.1 {
	// 	targetAngle = math.Pi * 0.1
	// }
	// head.angle += targetAngle * 0.1

	// Update position
	dist := math.Sqrt(math.Pow(target.x-head.pos.x, 2) + math.Pow(target.y-head.pos.y, 2))
	if dist > head.distance {
		head.pos.x += math.Cos(head.angle) * head.speed
		head.pos.y += math.Sin(head.angle) * head.speed
	}
}

func (h Head) draw(path *vector.Path) {
	path.Arc(
		float32(h.pos.x), float32(h.pos.y), float32(h.radius),
		float32(h.angle-math.Pi*0.5), float32(h.angle+math.Pi*0.5), vector.Clockwise,
	)
}

func (h *Head) rightEye() Point {
	angle := h.angle + 3*math.Pi/5
	radius := h.radius - 7
	x := h.pos.x + math.Cos(angle)*(radius)
	y := h.pos.y + math.Sin(angle)*(radius)

	return Point{x, y}
}
func (h *Head) leftEye() Point {
	angle := h.angle - 3*math.Pi/5
	radius := h.radius - 7
	x := h.pos.x + math.Cos(angle)*(radius)
	y := h.pos.y + math.Sin(angle)*(radius)

	return Point{x, y}
}
func (h Head) drawEyes(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen,
		float32(h.rightEye().x), float32(h.rightEye().y),
		float32(10),
		color.White, true)

	vector.DrawFilledCircle(screen,
		float32(h.leftEye().x), float32(h.leftEye().y),
		float32(10),
		color.White, true)
}

// ------------------------ Tail --------------------------
type Tail struct {
	*Joint
	// vertices []ebiten.Vertex
	// indices  []uint16
}

func TailNew(segment *Joint) Tail {
	tail := Tail{segment} //, []ebiten.Vertex{}, []uint16{}}
	return tail
}

func (t Tail) draw(path *vector.Path) {
	path.Arc(
		float32(t.pos.x), float32(t.pos.y), float32(t.radius),
		float32(t.angle+math.Pi*0.5), float32(t.angle-math.Pi*0.5), vector.Clockwise,
	)
}
