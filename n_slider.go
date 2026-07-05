package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var _ Node = &ScrollNode{}

type ScrollNode struct {
	baseNode
	scrollElement
}


type scrollElement struct {
	scrollX, scrollY float32 // 0 - 1
	buffer *ebiten.Image
}

func (sn *ScrollNode) draw(screen *ebiten.Image, showDebugInfo bool, depth int) {

}

func (sn *ScrollNode) base() *baseNode {
	return &sn.baseNode
}

func (sn *ScrollNode) getInitSize() (float32, float32) {
	return float32(sn.buffer.Bounds().Dx()), float32(sn.buffer.Bounds().Dy())
}
