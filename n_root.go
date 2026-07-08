package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)


type RootNode struct {
	lastW, lastH int
	changed bool

	animations map[tweenKey]Animation

	dragging interactable
	child Node
}



type tweenKey struct {
	node Node
	id string
}


type Animation struct {
	ID string
	Setter func(float32)
	From, To float32
	Duration float32
	elapsed float32
	Curve func(float32) float32
	Callback func()
}


func NewAnimation(id string, setter func(float32),  from, to, duration float32, curve func(float32) float32, callback func()) Animation {
	return Animation{
		id,
		setter,
		from, to,
		duration,
		0,
		curve,
		callback,
	}
}


func Animate(n Node, a Animation) {
	rn := n.base().rootNode
	if rn != nil {rn.animations[tweenKey{n, a.ID}] = a}
}

func AnimateTo(n Node, a Animation, getter func() float32) {
	a.From = getter()
	rn := n.base().rootNode
	if rn != nil {rn.animations[tweenKey{n, a.ID}] = a}
}

func AnimationSequence(n Node, animations ...Animation) {
	if len(animations) == 0 {
		return
	}

	first := animations[0]
	rest := animations[1:]
	original := first.Callback
	first.Callback = func() {
		if original != nil {
			original()
		}
		AnimationSequence(n, rest...)
	}
	Animate(n, first)
}

func (rn *RootNode) updateTween(dt float32) bool {
	if len(rn.animations) == 0 {
		return false
	}

	for id, a := range rn.animations {
		a.elapsed += dt
		time := a.elapsed / a.Duration

		if time >= 1 {
			time = 1
			a.Setter(a.From + (a.To - a.From) * a.Curve(time))
			delete(rn.animations, id)
			if a.Callback != nil {
				a.Callback()
			}

		} else {
			a.Setter(a.From + (a.To - a.From) * a.Curve(time))
			rn.animations[id] = a
		}
	}

	return true
}






func Root(firstNode Node) *RootNode {
	rn := &RootNode{
		lastW: -1,
		lastH: -1,
		changed: false,
		animations: make(map[tweenKey]Animation),
		child: firstNode,
	}

	firstNode.base().setRoot(rn)

	return rn
}


func (rn *RootNode) UpdateLayout(dt float32, outsideWidth, outsideHeight int) {
	if rn.updateTween(dt) { rn.changed = true }

	if outsideWidth != rn.lastW || outsideHeight != rn.lastH || rn.changed {

		rn.changed = false
		rn.lastW = outsideWidth
		rn.lastH = outsideHeight

		firstNode := rn.child

		if firstNode.base().IsVisible() {
			firstNode.base().bounds = rect{0, 0, float32(outsideWidth), float32(outsideHeight)}
		}

		arrange(firstNode)
		calculateBounds(firstNode)
	}
}



func (rn *RootNode) UpdateCursorLogic() {
	imX, imY := ebiten.CursorPosition()
	mX, mY := float32(imX), float32(imY)

	prepareFrame(rn.child)

	if rn.dragging != nil {
		rn.dragging.mouseOn(mX, mY)
	} else {
		findCursorInChildren(rn.child, mX, mY)
	}
}


func (rn *RootNode) Draw(screen *ebiten.Image, showDebugInfo bool) {
	drawTree(rn.child, screen, showDebugInfo, 0)
}
