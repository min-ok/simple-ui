package ui


import (
	"github.com/hajimehoshi/ebiten/v2"
)


type rect struct {
	minX, minY float32
	maxX, maxY float32
}


type Node interface {
	draw(screen *ebiten.Image, showDebugInfo bool, depth int)
	getInitSize() (float32, float32)

	base() *baseNode
}


type baseNode struct {
	weight float32
	border float32
	spacing float32

	selfAlign float32
	contentAlign float32

	// If fixX != 0, Weight has no effect on the X axis.
	// If fixY != 0, Weight has no effect on the Y axis.
	fixX, fixY float32

	occupiedX, occupiedY float32
	bounds rect
	children []Node

	flags uint8

	rootNode *RootNode
}

const (
	flagVisible uint8 = 1 << 0
	flagEnabled uint8 = 1 << 1
	flagVertical uint8 = 1 << 2
)


func (b *baseNode) setRoot(r *RootNode) {
	b.rootNode = r
	for _, child := range b.children {
		child.base().setRoot(r)
	}
}



// MakeDynamic makes the Node dynamic if it was fixed.
func MakeDynamic(v float32) CommonOpt {
	return func(n Node) {
		n.base().weight = v
		n.base().fixX = 0
		n.base().fixY = 0
	}
}

//
func Weight(v float32) CommonOpt {
	return func(n Node) { n.base().weight = v }
}

func (bn *baseNode) SetWeight(v float32) {
	bn.weight = v
	bn.rootNode.changed = true
}

func (bn *baseNode) GetWeight() float32 {
	return bn.weight
}



// Border sets the border of the content along the edges of the main diagonal.
// Measured in pixels.
func Border(v float32) CommonOpt {
	return func(n Node) { n.base().border = v }
}


func (bn *baseNode) SetBorder(v float32) {
	bn.border = v
	bn.rootNode.changed = true
}

func (bn *baseNode) GetBorder() float32 {
	return bn.border
}



// Direction

func Vertical() CommonOpt {
	return func(n Node) {
		n.base().flags |= flagVertical
	}
}

func Horizontal() CommonOpt {
	return func(n Node) {
		n.base().flags &= ^flagVertical
	}
}

func (bn *baseNode) SetVertical() {
	bn.flags |= flagVertical
	bn.rootNode.changed = true
}


func (bn *baseNode) SetHorizontal() {
	bn.flags &= ^flagVertical
	bn.rootNode.changed = true
}


func (bn *baseNode) IsVertical() bool {
	return bn.flags & flagVertical != 0
}



// Spacing sets the spacing of the content along the edges of the main diagonal.
// Measured in pixels.
func Spacing(v float32) CommonOpt {
	return func(n Node) { n.base().spacing = v }
}


func (bn *baseNode) SetSpacing(v float32) {
	bn.spacing = v
	bn.rootNode.changed = true
}

func (bn *baseNode) GetSpacing() float32 {
	return bn.spacing
}



// SelfAlign sets the aligning along the secondary axis.
// Works if the Node is fixed on the secondary axis.
// Measured from 0 to 1, in relation to the Node size.
func SelfAlign(v float32) CommonOpt  {
	return func(n Node) {
		n.base().selfAlign = v
	}
}

func (bn *baseNode) SetSelfAlign(v float32) {
	bn.selfAlign = v
	bn.rootNode.changed = true
}

func (bn *baseNode) GetSelfAlign() float32 {
	return bn.selfAlign
}

// ContentAlign sets the content aligning.
// Works if there are no dynamic objects in the box.
// Measured from 0 to 1, in relation to the Node size.
func ContentAlign(v float32) CommonOpt  {
	return func(n Node) {
		n.base().contentAlign = v
	}
}


func (bn *baseNode) SetContentAlign(v float32) {
	bn.contentAlign = v
	bn.rootNode.changed = true
}

func (bn *baseNode) GetContentAlign() float32 {
	return bn.contentAlign
}


// FixH fixes the node along the X axis.
func FixX(x float32) CommonOpt {
	return func(n Node) {
		n.base().fixX = x
	}
}

func (bn *baseNode) SetFixX(v float32) {
	bn.fixX = v
	bn.rootNode.changed = true
}

func (bn *baseNode) GetFixX() float32 {
	return bn.fixX
}

// FixV fixes the node along the Y axis.
func FixY(y float32) CommonOpt {
	return func(n Node) {
		n.base().fixY = y
	}
}

func (bn *baseNode) SetFixY(v float32) {
	bn.fixY = v
	bn.rootNode.changed = true
}

func (bn *baseNode) GetFixY() float32 {
	return bn.fixY
}

// FixSize fixes the node on both axes.
func FixSize(x, y float32) CommonOpt {
	return func(n Node) {
		n.base().fixX = x
		n.base().fixY = y
	}
}


func (bn *baseNode) SetFixSize(x, y float32) {
	bn.fixX = x
	bn.fixY = y
	bn.rootNode.changed = true
}

func (bn *baseNode) GetFixSize() (float32, float32) {
	return bn.fixX, bn.fixY
}





func Visible() CommonOpt {
	return func(n Node) {
		n.base().flags |= flagVisible
	}
}

func Invisible() CommonOpt {
	return func(n Node) {
		n.base().flags &= ^flagVisible
	}
}

func (bn *baseNode) Show() {
	bn.flags |= flagVisible
	bn.rootNode.changed = true
}


func (bn *baseNode) Hide() {
	bn.flags &= ^flagVisible
	bn.rootNode.changed = true
}


func (bn *baseNode) IsVisible() bool {
	return bn.flags & flagVisible != 0
}





func (bn *baseNode) enableBranch() {
	bn.flags |= flagEnabled
	for _, c := range bn.children {
		c.base().enableBranch()
	}
}


func (bn *baseNode) disableBranch() {
	bn.flags &= ^flagEnabled
	for _, c := range bn.children {
		c.base().disableBranch()
	}
}


func Enabled() CommonOpt {
	return func(n Node) {
		n.base().enableBranch()
	}
}

func Disabled() CommonOpt {
	return func(n Node) {
		n.base().disableBranch()
	}
}

func (bn *baseNode) Enable() {
	bn.enableBranch()
	bn.rootNode.changed = true
}


func (bn *baseNode) Disable() {
	bn.disableBranch()
	bn.rootNode.changed = true
}


func (bn *baseNode) IsEnabled() bool {
	return bn.flags & flagEnabled != 0
}







// DefaultSize Должен быть вызван после всех настроек, которые могут повлиять

func DefaultSize() CommonOpt {
	return func(n Node) {
		base := n.base()
		base.fixX, base.fixY = n.getInitSize()
	}
}

func (bn *baseNode) GetBounds() (float32, float32, float32, float32) {
	bounds := bn.bounds
	return bounds.minX, bounds.minY, bounds.maxX, bounds.maxY
}

// func setRootNode(rn *RootNode) CommonOpt {
// 	return func(n Node) {
// 		n.base().rootNode = rn
// 	}
// }
