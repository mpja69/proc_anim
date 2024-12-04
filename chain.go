package main

import (
	"image/color"
	"math"
)

type Chain struct {
	joints   []*Joint
	distance float64
	x, y     float64
}

func lerp(t, lo, hi float64) float64 {
	return float64((1-t)*lo + t*hi)
}

// Loop to create segments. Head first, tail last
func ChainNew(bodyShape []int, x, y, distance int) *Chain {
	nbrSegments := len(bodyShape)
	segments := make([]*Joint, nbrSegments)

	// Create all segments, incl head and tail
	for i := 0; i < nbrSegments; i++ {

		//color := HSVtoRGBA(HSVtoRGBNorm)(hue, 1, 1)
		segments[i] = JointNew(
			float64(x-i*distance), float64(y),
			float64(distance), float64(bodyShape[i]),
			color.RGBA{255, 255, 255, 255},
		)
	}
	return &Chain{joints: segments, distance: float64(distance), x: float64(x), y: float64(y)}
}

func (c *Chain) first() *Joint {
	return c.joints[0]
}

func (c *Chain) last() *Joint {
	return c.joints[len(c.joints)-1]
}
func (c *Chain) getAdjustedPosX(i int, angleOffset, lengthOffset float64) float64 {
	s := c.joints[i]
	return s.pos.x + math.Cos(s.angle+angleOffset)*(s.radius+lengthOffset)
}
func (c *Chain) getAdjustedPosY(i int, angleOffset, lengthOffset float64) float64 {
	s := c.joints[i]
	return s.pos.y + math.Sin(s.angle+angleOffset)*(s.radius+lengthOffset)
}
func (c *Chain) getAdjustedPos(i int, angleOffset, lengthOffset float64) Point {
	s := c.joints[i]
	return s.getAdjustedPos(angleOffset, lengthOffset)
}

func (c *Chain) SetAnchorPos(p Point) {
	c.joints[0].pos = p
}
func (c *Chain) easyFollow(target Point) {
	// The first joint follow the target (set position and angle)
	head := c.first()

	head.pos = target
	// update the other segments
	for i := 1; i < len(c.joints); i++ {
		prev := c.joints[i-1]
		curr := c.joints[i]
		curr.DirectlyFollow(prev.pos)
	}
}

// Update all segments of the body
func (c *Chain) DIRECT(target Point, speed float64) {
	// The first joint follow the target (set position and angle)
	head := c.first()
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
		head.pos.x += math.Cos(head.angle) * speed
		head.pos.y += math.Sin(head.angle) * speed
	}

	// update the other segments
	for i := 1; i < len(c.joints); i++ {
		prev := c.joints[i-1]
		curr := c.joints[i]
		curr.DirectlyFollow(prev.pos)
	}
}

// Update all segments of the limb with respect to the target point
// Using Cyclic Coordinate Descent Inverse Kinematics
func (c *Chain) CCDIK(target Point) {
	lastIdx := len(c.joints) - 1
	lastSegment := c.joints[lastIdx]

	// Outer loop: Iterate from the last segment to the first, and update each.
	for i := lastIdx; i >= 0; i-- {
		end := lastSegment.End()
		c.joints[i].CCDIKUpdateAngle(end, target)

		// Inner loop: Update all the following segments start positions, based on the previous segments
		for j := i; j < lastIdx; j++ {
			curr := c.joints[j]
			next := c.joints[j+1]
			next.pos = curr.End()
		}
	}
}

func (c *Chain) FABRIK(target, anchor Point) {
	// Backward loop: Iterate from the last segment to the first, and update each.
	c.joints[len(c.joints)-1].pos = target
	for i := len(c.joints) - 2; i >= 0; i-- {
		curr := c.joints[i]
		next := c.joints[i+1]
		curr.pos = SetConstraint(curr.pos, next.pos, curr.distance)
	}

	// Forward loop: Update all the following segments start positions, based on the previous segments
	c.joints[0].pos = anchor
	for i := 1; i < len(c.joints); i++ {
		curr := c.joints[i]
		prev := c.joints[i-1]
		curr.pos = SetConstraint(curr.pos, prev.pos, curr.distance)
		curr.angle = prev.pos.Angle(curr.pos)
	}
	c.joints[0].angle = c.joints[0].pos.Angle(c.joints[1].pos)

}
