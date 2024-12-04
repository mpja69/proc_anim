package main

import "github.com/hajimehoshi/ebiten/v2/vector"

func getControlPoints(x1, y1, x2, y2, x3, y3 float64) (float64, float64, float64, float64, float64, float64, float64, float64) {
	cx1a := x1 + (x2-x1)/3
	cy1a := y1 + (y2-y1)/3
	cx1b := x2 - (x3-x1)/3
	cy1b := y2 - (y3-y1)/3
	cx2a := x2 + (x3-x1)/3
	cy2a := y2 + (y3-y1)/3
	cx2b := x3 - (x3-x2)/3
	cy2b := y3 - (y3-y2)/3
	return cx1a, cy1a, cx1b, cy1b, cx2a, cy2a, cx2b, cy2b
}

// Interpolates a line through the next 2 points: curr and next.
// The prev point is given to calculate the correct control points
func interpolate2BezierVertices(path *vector.Path, prev, curr, next Point) {
	// cx1a, cy1a, cx1b, cy1b, cx2a, cy2a, cx2b, cy2b := getControlPoints(p1.x, p1.y, p2.x, p2.y, p3.x, p3.y)
	cx1a := prev.x //+ (curr.x-prev.x)/3
	cy1a := prev.y //+ (curr.y-prev.y)/3
	cx1b := curr.x - (next.x-prev.x)/3
	cy1b := curr.y - (next.y-prev.y)/3
	cx2a := curr.x + (next.x-prev.x)/3
	cy2a := curr.y + (next.y-prev.y)/3
	cx2b := next.x //- (next.x-curr.x)/3
	cy2b := next.y //- (next.y-curr.y)/3
	path.CubicTo(float32(cx1a), float32(cy1a), float32(cx1b), float32(cy1b), float32(curr.x), float32(curr.y))
	path.CubicTo(float32(cx2a), float32(cy2a), float32(cx2b), float32(cy2b), float32(next.x), float32(next.y))
}

// Interpolates a line through the next 2 points: curr and next.
// The prev point is given to calculate the correct control points
func interpolateBezierVertices(path *vector.Path, prev, curr, next Point) {
	cx1a := prev.x //+ (curr.x-prev.x)/3
	cy1a := prev.y //+ (curr.y-prev.y)/3
	cx1b := curr.x - (next.x-prev.x)/3
	cy1b := curr.y - (next.y-prev.y)/3
	// cx2a := curr.x + (next.x-prev.x)/3
	// cy2a := curr.y + (next.y-prev.y)/3
	// cx2b := next.x //- (next.x-curr.x)/3
	// cy2b := next.y //- (next.y-curr.y)/3
	path.CubicTo(float32(cx1a), float32(cy1a), float32(cx1b), float32(cy1b), float32(curr.x), float32(curr.y))
	// path.CubicTo(float32(cx2a), float32(cy2a), float32(cx2b), float32(cy2b), float32(next.x), float32(next.y))
}
