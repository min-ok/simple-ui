package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type buttonState int

type ButtonNode struct {
	baseNode
	buttonElement
}

type buttonElement struct {
	spriteElement

	state     buttonState
	prevState buttonState

	mouseWasOn     bool
	prevMouseWasOn bool

	eventsSource *ButtonEventsSource
}

type ButtonEventsSource struct {
	onJustPressedFunc  func(bn *ButtonNode)
	onJustReleasedFunc func(bn *ButtonNode)
	onPressedFunc      func(bn *ButtonNode)
	onEnterFunc        func(bn *ButtonNode)
	onLeaveFunc        func(bn *ButtonNode)
}

var _ Node = &ButtonNode{}
var _ interactable = &ButtonNode{}

// for element interface
func (bn *ButtonNode) draw(screen *ebiten.Image, showDebugInfo bool, depth int) {
	bn.drawSpriteElement(bn.baseNode.bounds, screen)
}

func (bn *ButtonNode) base() *baseNode {
	return &bn.baseNode
}

func (bn *ButtonNode) mouseOn(mX, mY float32) {
	bn.mouseWasOn = true

	if !bn.IsEnabled() {
		bn.mouseWasOn = false
		bn.state = normal
		return
	}

	switch {
	case inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft):
		bn.state = justPressed
		if bn.eventsSource == nil {
			return
		}
		if bn.eventsSource.onJustPressedFunc != nil {
			bn.eventsSource.onJustPressedFunc(bn)
		}
	case inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft):
		bn.state = justReleased
		if bn.eventsSource == nil {
			return
		}
		if bn.eventsSource.onJustReleasedFunc != nil {
			bn.eventsSource.onJustReleasedFunc(bn)
		}
	case ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft):
		bn.state = pressed
		if bn.eventsSource == nil {
			return
		}
		if bn.eventsSource.onPressedFunc != nil {
			bn.eventsSource.onPressedFunc(bn)
		}
	default:
		bn.state = hovered

		if bn.prevState != hovered {
			if bn.eventsSource == nil {
				return
			}
			if bn.eventsSource.onEnterFunc != nil {
				bn.eventsSource.onEnterFunc(bn)
			}
		}
	}
}

func (bn *ButtonNode) prepareFrame() {
	if !bn.mouseWasOn && bn.prevMouseWasOn {
		if bn.eventsSource != nil {
			if bn.eventsSource.onLeaveFunc != nil {
				bn.eventsSource.onLeaveFunc(bn)
			}
		}
	}

	bn.prevMouseWasOn = bn.mouseWasOn
	bn.mouseWasOn = false

	bn.prevState = bn.state
	bn.state = normal
}

// Button creates a Button Node.
// By default, the Button gets the size of its normal Sprite. To avoid this, use MakeDynamic.
func Button(opts ...ButtonOption) *ButtonNode {
	bn := &ButtonNode{}
	bn.flags |= flagVisible | flagEnabled

	for _, opt := range opts {
		opt.applyToNode(bn)
		opt.applyToButton(bn)
	}
	return bn
}

// ButtonEvents
func NewButtonEventsSource(onJustPressedFunc, onJustReleasedFunc, onPressedFunc, onEnterFunc, onLeaveFunc func(bn *ButtonNode)) *ButtonEventsSource {
	return &ButtonEventsSource{
		onJustPressedFunc:  onJustPressedFunc,
		onJustReleasedFunc: onJustReleasedFunc,
		onPressedFunc:      onPressedFunc,
		onEnterFunc:        onEnterFunc,
		onLeaveFunc:        onLeaveFunc,
	}
}

func ButtonEvents(bes *ButtonEventsSource) ButtonOpt {
	return func(be *ButtonNode) {
		be.eventsSource = bes
	}
}
func (bn *ButtonNode) SetButtonEvents(bes *ButtonEventsSource) {
	bn.eventsSource = bes
}
func (bn *ButtonNode) GetButtonEvents() *ButtonEventsSource {
	return bn.eventsSource
}

func OnJustPressedFunc(onJustPressedFunc func(*ButtonNode)) ButtonOpt {
	return func(bn *ButtonNode) {
		bn.ensureEventsSource()
		bn.eventsSource.onJustPressedFunc = onJustPressedFunc
	}
}
func (bn *ButtonNode) SetOnJustPressedFunc(onJustPressedFunc func(*ButtonNode)) {
	bn.ensureEventsSource()
	bn.eventsSource.onJustPressedFunc = onJustPressedFunc
}
func (bn *ButtonNode) GetOnJustPressedFunc() func(*ButtonNode) {
	if bn.eventsSource == nil {
		return nil
	}
	return bn.eventsSource.onJustPressedFunc
}

func OnJustReleasedFunc(onJustReleasedFunc func(*ButtonNode)) ButtonOpt {
	return func(bn *ButtonNode) {
		bn.ensureEventsSource()
		bn.eventsSource.onJustReleasedFunc = onJustReleasedFunc
	}
}
func (bn *ButtonNode) SetOnJustReleasedFunc(onJustReleasedFunc func(*ButtonNode)) {
	bn.ensureEventsSource()
	bn.eventsSource.onJustReleasedFunc = onJustReleasedFunc
}
func (bn *ButtonNode) GetOnJustReleasedFunc() func(*ButtonNode) {
	if bn.eventsSource == nil {
		return nil
	}
	return bn.eventsSource.onJustReleasedFunc
}

func OnPressedFunc(onPressedFunc func(*ButtonNode)) ButtonOpt {
	return func(bn *ButtonNode) {
		bn.ensureEventsSource()
		bn.eventsSource.onPressedFunc = onPressedFunc
	}
}
func (bn *ButtonNode) SetOnPressedFunc(onPressedFunc func(*ButtonNode)) {
	bn.ensureEventsSource()
	bn.eventsSource.onPressedFunc = onPressedFunc
}
func (bn *ButtonNode) GetOnPressedFunc() func(*ButtonNode) {
	if bn.eventsSource == nil {
		return nil
	}
	return bn.eventsSource.onPressedFunc
}

func OnEnterFunc(onEnterFunc func(*ButtonNode)) ButtonOpt {
	return func(bn *ButtonNode) {
		bn.ensureEventsSource()
		bn.eventsSource.onEnterFunc = onEnterFunc
	}
}
func (bn *ButtonNode) SetOnEnterFunc(onEnterFunc func(*ButtonNode)) {
	bn.ensureEventsSource()
	bn.eventsSource.onEnterFunc = onEnterFunc
}
func (bn *ButtonNode) GetOnEnterFunc() func(*ButtonNode) {
	if bn.eventsSource == nil {
		return nil
	}
	return bn.eventsSource.onEnterFunc
}

func OnLeaveFunc(onLeaveFunc func(*ButtonNode)) ButtonOpt {
	return func(bn *ButtonNode) {
		bn.ensureEventsSource()
		bn.eventsSource.onLeaveFunc = onLeaveFunc
	}
}
func (bn *ButtonNode) SetOnLeaveFunc(onLeaveFunc func(*ButtonNode)) {
	bn.ensureEventsSource()
	bn.eventsSource.onLeaveFunc = onLeaveFunc
}
func (bn *ButtonNode) GetOnLeaveFunc() func(*ButtonNode) {
	if bn.eventsSource == nil {
		return nil
	}
	return bn.eventsSource.onLeaveFunc
}
