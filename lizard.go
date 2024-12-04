package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Lizard struct {
	chain *Chain
	limbs []*Limb
	speed float64
	// To draw
	vertices []ebiten.Vertex
	indices  []uint16
}

// Loop to create segments. Head first, tail last
func LizardNew(x, y int) *Lizard {
	bodyShape := []int{52, 58, 40, 60, 68, 71, 65, 50, 28, 15, 11, 9, 7, 7, 7}
	chain := ChainNew(bodyShape, x, y, 64)

	limbs := make([]*Limb, 0, 4)
	limbs = append(limbs, LimbNew(chain.joints[3], 40, true, true))
	limbs = append(limbs, LimbNew(chain.joints[7], 40, false, false))
	limbs = append(limbs, LimbNew(chain.joints[3], 40, false, true))
	limbs = append(limbs, LimbNew(chain.joints[7], 40, true, false))

	return &Lizard{chain: chain, limbs: limbs, vertices: []ebiten.Vertex{}, indices: []uint16{}, speed: 2}
}

// Update all segments of the body
func (b *Lizard) update(target Point) {

	// Update the body (with each joint directly follow eachother)
	b.chain.DIRECT(target, b.speed)

	// Update each limb, (using FABRIK)
	for _, l := range b.limbs {
		l.update()
	}
}

func (b *Lizard) debugDraw(screen *ebiten.Image) {
	for _, j := range b.chain.joints {
		j.DrawCircle(screen)
	}
	for _, l := range b.limbs {
		for _, j := range l.chain.joints {
			j.DrawCircle(screen)
			x0 := j.pos.x
			y0 := j.pos.y
			x1 := j.pos.x + math.Cos(j.angle)*(j.radius)
			y1 := j.pos.y + math.Sin(j.angle)*(j.radius)
			ebitenutil.DrawLine(screen, x0, y0, x1, y1, color.White)

		}
	}
}

func (b *Lizard) createPath() *vector.Path {
	// Create the path clockwise around the whole body
	path := vector.Path{}

	chain := b.chain

	// Move to start point: First point
	path.MoveTo(float32(chain.first().Right().x), float32(chain.first().Right().y))

	// Draw the right side. p1 only used for getting control points, and move through p2 an p3. Take 2 steps in the loop
	for i := 0; i < len(chain.joints)-2; i += 1 {
		prev := chain.joints[i].Right()
		curr := chain.joints[i+1].Right()
		next := chain.joints[i+2].Right()
		interpolateBezierVertices(&path, prev, curr, next)
	}

	// Draw the tail
	tail := len(chain.joints) - 1
	prev := chain.joints[tail-1].Right()
	curr := chain.joints[tail].Right()
	next := Point{chain.getAdjustedPosX(tail, math.Pi, 20), chain.getAdjustedPosY(tail, math.Pi, 20)}
	interpolateBezierVertices(&path, prev, curr, next)

	prev = chain.joints[tail].Right()
	curr = Point{chain.getAdjustedPosX(tail, math.Pi, 20), chain.getAdjustedPosY(tail, math.Pi, 20)}
	next = chain.joints[tail].Left()
	interpolateBezierVertices(&path, prev, curr, next)

	prev = Point{chain.getAdjustedPosX(tail, math.Pi, 20), chain.getAdjustedPosY(tail, math.Pi, 20)}
	curr = chain.joints[tail].Left()
	next = chain.joints[tail-1].Left()
	interpolateBezierVertices(&path, prev, curr, next)

	// Draw the left side
	for i := len(chain.joints) - 1; i > 1; i -= 1 {
		prev := chain.joints[i].Left()
		curr := chain.joints[i-1].Left()
		next := chain.joints[i-2].Left()
		interpolateBezierVertices(&path, prev, curr, next)
	}

	// Draw the head
	prev = chain.joints[1].Left()
	curr = chain.joints[0].Left()
	next = Point{chain.getAdjustedPosX(0, -math.Pi/6, -8), chain.getAdjustedPosY(0, -math.Pi/6, -10)}
	interpolateBezierVertices(&path, prev, curr, next)

	// Top of the head (completes the loop)
	p1 := Point{chain.getAdjustedPosX(0, -math.Pi/6, -8), chain.getAdjustedPosY(0, -math.Pi/6, -10)}
	p2 := Point{chain.getAdjustedPosX(0, 0, -6), chain.getAdjustedPosY(0, 0, -4)}
	p3 := Point{chain.getAdjustedPosX(0, math.Pi/6, -8), chain.getAdjustedPosY(0, math.Pi/6, -10)}
	interpolate2BezierVertices(&path, p1, p2, p3)

	prev = Point{chain.getAdjustedPosX(0, math.Pi/6, -8), chain.getAdjustedPosY(0, math.Pi/6, -10)}
	curr = chain.joints[0].Right()
	next = chain.joints[1].Right()
	interpolateBezierVertices(&path, prev, curr, next)

	return &path
}

func (b *Lizard) draw(screen *ebiten.Image) {
	for _, l := range b.limbs {
		l.draw(screen)
	}

	path := b.createPath()

	// Render the filled area
	b.vertices, b.indices = path.AppendVerticesAndIndicesForFilling(b.vertices[:0], b.indices[:0])
	for i := range b.vertices {
		b.vertices[i].SrcX = 1
		b.vertices[i].SrcY = 1
		b.vertices[i].ColorR = 0x58 / float32(0xff)
		b.vertices[i].ColorG = 0x85 / float32(0xff)
		b.vertices[i].ColorB = 0x7A / float32(0xff)
		b.vertices[i].ColorA = 1
	}
	top := &ebiten.DrawTrianglesOptions{}
	top.AntiAlias = true
	top.FillRule = ebiten.FillRuleNonZero
	screen.DrawTriangles(b.vertices, b.indices, outlineSubImage, top)

	// Render the outline
	sop := &vector.StrokeOptions{}
	sop.Width = 3
	sop.LineJoin = vector.LineJoinRound
	b.vertices, b.indices = path.AppendVerticesAndIndicesForStroke(b.vertices[:0], b.indices[:0], sop)
	screen.DrawTriangles(b.vertices, b.indices, outlineSubImage, top)

	b.drawEyes(screen)
}

// Conclrete and detailed implementation of how to draw the lizard's eyes
// No need to generalize and regard DRY!
func (b *Lizard) drawEyes(screen *ebiten.Image) {
	p := b.chain.first()
	angle := p.angle + 3*math.Pi/5
	radius := p.radius - 7
	x := p.pos.x + math.Cos(angle)*(radius)
	y := p.pos.y + math.Sin(angle)*(radius)
	vector.DrawFilledCircle(screen,
		float32(x), float32(y),
		float32(10),
		color.White, true)

	angle = p.angle - 3*math.Pi/5
	radius = p.radius - 7
	x = p.pos.x + math.Cos(angle)*(radius)
	y = p.pos.y + math.Sin(angle)*(radius)
	vector.DrawFilledCircle(screen,
		float32(x), float32(y),
		float32(10),
		color.White, true)
}
