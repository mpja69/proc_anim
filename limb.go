package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Limb struct {
	chain       *Chain
	anchorJoint *Joint
	footPos     Point
	rightSide   bool
	frontSide   bool
	maxLength   float64
	vertices    []ebiten.Vertex
	indices     []uint16
}

func LimbNew(anchorJoint *Joint, distance int, rightSide, frontSide bool) *Limb {
	joints := []int{20, 40, 30}
	chain := ChainNew(joints, 0, 0, distance)
	maxLength := float64(len(joints) * distance)
	return &Limb{chain: chain, anchorJoint: anchorJoint, rightSide: rightSide, frontSide: frontSide, maxLength: maxLength}
}

func (l *Limb) totalLength() float64 {
	return distance(l.chain.first().pos, l.chain.last().pos)
}

// Update the limbs position and shape, and take step if it's time to move
func (l *Limb) update() (didMove bool) {
	didMove = false
	var anchorPos Point
	var moveAngle float64
	var newFootPos Point
	var length float64

	if l.rightSide {
		if l.frontSide {
			moveAngle = l.anchorJoint.angle + math.Pi/8 // 6
			length = l.maxLength
		} else {
			moveAngle = l.anchorJoint.angle + math.Pi/4
			length = l.maxLength * 0.5
		}
		newFootPos.x = length*math.Cos(moveAngle) + l.anchorJoint.Right().x
		newFootPos.y = length*math.Sin(moveAngle) + l.anchorJoint.Right().y
		anchorPos = l.anchorJoint.getAdjustedPos(math.Pi/2, -14)
	} else {
		if l.frontSide {
			moveAngle = l.anchorJoint.angle - math.Pi/8 //6
			length = l.maxLength
		} else {
			moveAngle = l.anchorJoint.angle - math.Pi/4
			length = l.maxLength * 0.5
		}
		newFootPos.x = length*math.Cos(moveAngle) + l.anchorJoint.Left().x
		newFootPos.y = length*math.Sin(moveAngle) + l.anchorJoint.Left().y
		anchorPos = l.anchorJoint.getAdjustedPos(-math.Pi/2, -14)
	}

	delta := distance(l.anchorJoint.pos, l.footPos) * 0.7
	if delta > l.maxLength {
		l.footPos = newFootPos
		didMove = true
	}

	l.chain.FABRIK(l.footPos, anchorPos)
	return didMove
}

func (b *Limb) draw(screen *ebiten.Image) {

	path := b.createPath()

	// Triangle options
	top := &ebiten.DrawTrianglesOptions{}

	// Stroke options
	sop := &vector.StrokeOptions{}
	sop.Width = 40
	sop.LineJoin = vector.LineJoinRound
	sop.LineCap = vector.LineCapRound
	b.vertices, b.indices = path.AppendVerticesAndIndicesForStroke(b.vertices[:0], b.indices[:0], sop)
	screen.DrawTriangles(b.vertices, b.indices, outlineSubImage, top)

	sop.Width = 32
	sop.LineJoin = vector.LineJoinRound
	b.vertices, b.indices = path.AppendVerticesAndIndicesForStroke(b.vertices[:0], b.indices[:0], sop)
	screen.DrawTriangles(b.vertices, b.indices, fillSubImage, top)
}
func (l *Limb) createPath() *vector.Path {
	path := vector.Path{}

	shoulder := l.chain.joints[0].pos
	elbow := l.chain.joints[1].pos
	foot := l.chain.joints[2].pos

	para := foot.Sub(shoulder)
	perp := Point{-para.y, para.x}.SetMag(30)
	if l.frontSide == false {
		if l.rightSide {
			elbow = elbow.Sub(perp)
		} else {
			elbow = elbow.Add(perp)
		}
	}

	path.MoveTo(float32(shoulder.x), float32(shoulder.y))
	path.CubicTo(
		float32(elbow.x), float32(elbow.y),
		float32(elbow.x), float32(elbow.y),
		float32(foot.x), float32(foot.y),
	)
	return &path
}
