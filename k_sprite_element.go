package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type spriteElement struct {
	owner *baseNode

	image           *ebiten.Image
	nineSliceSource *NineSliceSource

	colorScale ebiten.ColorScale
	filter     ebiten.Filter
	blend      ebiten.Blend
}

type NineSliceSource struct {
	top    int
	bottom int
	right  int
	left   int
}

func (se *spriteElement) getInitSize() (float32, float32) {
	if se.image == nil {
		return 0, 0
	}

	bounds := se.image.Bounds()
	return float32(bounds.Dx()), float32(bounds.Dy())
}

func (se *spriteElement) drawSpriteElement(rect rect, screen *ebiten.Image) {
	if se.image == nil {
		return
	}

	if se.nineSliceSource == nil {
		se.drawScaled(rect, screen)
	} else {
		se.drawNineSlice(rect, screen)
	}
}

func (se *spriteElement) drawNineSlice(rect rect, screen *ebiten.Image) {
	bounds := se.image.Bounds()

	iSizeX, iSizeY := bounds.Dx(), bounds.Dy()
	fSizeX, fSizeY := float32(iSizeX), float32(iSizeY)

	offsetX := bounds.Min.X
	offsetY := bounds.Min.Y

	boxSizeX, boxSizeY := rect.maxX-rect.minX, rect.maxY-rect.minY

	top, bottom, right, left := se.nineSliceSource.top, se.nineSliceSource.bottom, se.nineSliceSource.right, se.nineSliceSource.left

	currSlicesX := []int{0, left, iSizeX - right, iSizeX}
	currSlicesY := []int{0, top, iSizeY - bottom, iSizeY}

	newSlicesX := []float32{0, float32(left), boxSizeX - float32(right), boxSizeX}
	newSlicesY := []float32{0, float32(top), boxSizeY - float32(bottom), boxSizeY}

	xFactor := (boxSizeX - float32(right+left)) / (fSizeX - float32(right+left))
	yFactor := (boxSizeY - float32(top+bottom)) / (fSizeY - float32(top+bottom))

	for i := 0; i < 3; i += 1 {
		for j := 0; j < 3; j += 1 {
			subRect := image.Rect(
				offsetX+currSlicesX[i],
				offsetY+currSlicesY[j],
				offsetX+currSlicesX[i+1],
				offsetY+currSlicesY[j+1],
			)
			res := se.image.SubImage(subRect).(*ebiten.Image)
			op := &ebiten.DrawImageOptions{}

			if i == 1 && j == 1 {
				op.GeoM.Scale(float64(xFactor), float64(yFactor))
			} else if i%2 == 1 {
				op.GeoM.Scale(float64(xFactor), 1)
			} else if j%2 == 1 {
				op.GeoM.Scale(1, float64(yFactor))
			}

			op.GeoM.Translate(float64(newSlicesX[i]+rect.minX), float64(newSlicesY[j]+rect.minY))

			op.ColorScale = se.colorScale
			op.Filter = se.filter
			op.Blend = se.blend

			screen.DrawImage(res, op)
		}
	}
}

func (se *spriteElement) drawScaled(rect rect, screen *ebiten.Image) {
	iSizeX, iSizeY := se.image.Bounds().Dx(), se.image.Bounds().Dy()
	fSizeX, fSizeY := float32(iSizeX), float32(iSizeY)

	boxSizeX, boxSizeY := rect.maxX-rect.minX, rect.maxY-rect.minY

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(float64(boxSizeX/fSizeX), float64(boxSizeY/fSizeY))
	op.GeoM.Translate(float64(rect.minX), float64(rect.minY))

	op.ColorScale = se.colorScale
	op.Filter = se.filter
	op.Blend = se.blend

	screen.DrawImage(se.image, op)
}

func (se *spriteElement) sizeChanged() {
	if se.owner != nil && se.owner.IsAutoSize() {
		se.owner.makeChanged()
	}
}

// NewNineSlice creates a NineSlice.
func NewNineSliceSource(top, bottom, right, left int) *NineSliceSource {
	return &NineSliceSource{
		top:    top,
		bottom: bottom,
		right:  right,
		left:   left,
	}
}

// Image

func Image(image *ebiten.Image) SpriteElementOpt {
	return func(se *spriteElement) {
		se.image = image
	}
}

func (se *spriteElement) SetImage(image *ebiten.Image) {
	se.image = image
	se.sizeChanged()
}

func (se *spriteElement) GetImage() *ebiten.Image {
	return se.image
}

// NineSlice

func NineSlice(nss *NineSliceSource) SpriteElementOpt {
	return func(se *spriteElement) {
		se.nineSliceSource = nss
	}
}

func (se *spriteElement) SetNineSlice(nss *NineSliceSource) {
	se.nineSliceSource = nss
	se.sizeChanged()
}

func (se *spriteElement) GetNineSlice() *NineSliceSource {
	return se.nineSliceSource
}

// ColorScale

func ColorScale(cs ebiten.ColorScale) SpriteElementOpt {
	return func(se *spriteElement) {
		se.colorScale = cs
	}
}

func (se *spriteElement) SetColorScale(cs ebiten.ColorScale) {
	se.colorScale = cs
}

func (se *spriteElement) GetColorScale() ebiten.ColorScale {
	return se.colorScale
}

// Filter

func Filter(f ebiten.Filter) SpriteElementOpt {
	return func(se *spriteElement) {
		se.filter = f
	}
}

func (se *spriteElement) SetFilter(f ebiten.Filter) {
	se.filter = f
}

func (se *spriteElement) GetFilter() ebiten.Filter {
	return se.filter
}

// Blend

func Blend(b ebiten.Blend) SpriteElementOpt {
	return func(se *spriteElement) {
		se.blend = b
	}
}

func (se *spriteElement) SetBlend(b ebiten.Blend) {
	se.blend = b
}

func (se *spriteElement) GetBlend() ebiten.Blend {
	return se.blend
}
