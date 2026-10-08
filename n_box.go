package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type BoxNode struct {
	baseNode
	boxElement
}

type boxElement struct {
	image *ebiten.Image
}

func (bn *BoxNode) draw(screen *ebiten.Image, showDebugInfo bool, depth int) {
	if showDebugInfo {
		op := &ebiten.DrawImageOptions{}

		sizeX := bn.bounds.maxX - bn.bounds.minX
		sizeY := bn.bounds.maxY - bn.bounds.minY

		colors := [][]color.Color{
			{color.RGBA{70, 160, 180, 180}, color.RGBA{40, 90, 105, 180}},
			{color.RGBA{120, 190, 100, 180}, color.RGBA{70, 110, 60, 180}},
			{color.RGBA{250, 160, 110, 180}, color.RGBA{180, 100, 60, 180}},
			{color.RGBA{210, 80, 100, 180}, color.RGBA{120, 30, 50, 180}},
			{color.RGBA{150, 110, 190, 180}, color.RGBA{80, 50, 110, 180}},
			{color.RGBA{240, 210, 110, 180}, color.RGBA{170, 140, 60, 180}},
			{color.RGBA{110, 130, 150, 180}, color.RGBA{60, 70, 80, 180}},
		}

		clr := colors[depth%len(colors)]

		strokeWidth := float32(4)

		if bn.image == nil {
			bn.image = ebiten.NewImage(1, 1)
			bn.image.Fill(clr[0])
		}

		op.GeoM.Scale(float64(sizeX-2*strokeWidth), float64(sizeY-2*strokeWidth))
		op.GeoM.Translate(float64(bn.bounds.minX+strokeWidth), float64(bn.bounds.minY+strokeWidth))
		screen.DrawImage(bn.image, op)

		vector.StrokeRect(
			screen,
			bn.bounds.minX+strokeWidth/2,
			bn.bounds.minY+strokeWidth/2,
			sizeX-strokeWidth,
			sizeY-strokeWidth,
			strokeWidth,
			clr[1],
			false,
		)
	}
}

func (bn *BoxNode) getInitSize() (float32, float32) {
	return 0, 0
}

func (bn *BoxNode) base() *baseNode {
	return &bn.baseNode
}

func Box(opts ...BoxOption) *BoxNode {
	bn := &BoxNode{}
	bn.flags |= flagVisible | flagEnabled | flagVertical

	for _, opt := range opts {
		opt.applyToNode(bn)
		opt.applyToBox(bn)
	}
	return bn
}
