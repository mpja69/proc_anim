package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WIDTH  = 1200
	HEIGHT = 800
	FACTOR = 2.0
)

var (
	outlineImage    = ebiten.NewImage(3, 3)
	outlineSubImage = outlineImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	fillImage       = ebiten.NewImage(3, 3)
	fillSubImage    = fillImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	fillImage.Fill(color.RGBA{0x58, 0x85, 0x7A, 255})
	outlineImage.Fill(color.White)
}

// Use the mouse pointer as target
func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	tx := float64(x * FACTOR)
	ty := float64(y * FACTOR)

	g.lizard.update(Point{tx, ty})

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.backBuffer.Fill(color.RGBA{44, 51, 60, 255})
	g.lizard.draw(g.backBuffer)

	// g.lizard.debugDraw(g.backBuffer)
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Scale(1/FACTOR, 1/FACTOR)
	screen.DrawImage(g.backBuffer, &opts)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return WIDTH, HEIGHT
}

type Game struct {
	lizard     *Lizard
	backBuffer *ebiten.Image
}

func main() {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Inverse kinematics!")
	g := Game{
		lizard: LizardNew(WIDTH/2, HEIGHT/2),
	}
	// Create a bigger backbuffer
	g.backBuffer = ebiten.NewImage(WIDTH*FACTOR, HEIGHT*FACTOR)

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
