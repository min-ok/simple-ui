package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)


type SpriteNode struct {
	baseNode
	spriteElement
}


var _ Node = &SpriteNode{}


// for element interface
func (sn *SpriteNode) draw(screen *ebiten.Image, showDebugInfo bool, depth int) {
	sn.drawSpriteElement(sn.baseNode.bounds, screen)
}


func (sn *SpriteNode) base() *baseNode {
	return &sn.baseNode
}


// Sprite creates a Sprite Node.
// By default, the Sprite gets the size of its Image. To avoid this, use MakeDynamic.
func Sprite(opts ...SpriteOption) *SpriteNode {
	sn := &SpriteNode{}
	sn.flags |= flagVisible | flagEnabled

	for _, opt := range opts {
		opt.applyToNode(sn)
		opt.applyToSprite(sn)
	}
	return sn
}
