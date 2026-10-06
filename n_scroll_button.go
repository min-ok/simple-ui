package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ScrollButtonNode struct {
	baseNode
	scrollButtonElement
}

type scrollButtonElement struct {
	stripe spriteElement
	knob   spriteElement

	value float32

	knobAlign  float32
	trackStart float32
	trackEnd   float32

	dragStartPos   float32
	dragStartValue float32

	state     buttonState
	prevState buttonState

	stripeMouseWasOn     bool
	prevStripeMouseWasOn bool

	knobMouseWasOn     bool
	prevKnobMouseWasOn bool

	eventsSource *ScrollButtonEventsSource
}

type ScrollButtonEventsSource struct {
	onDragStartFunc func(sbn *ScrollButtonNode)
	onDragEndFunc   func(sbn *ScrollButtonNode)
	onDragingFunc   func(sbn *ScrollButtonNode)

	onStripeEnterFunc func(sbn *ScrollButtonNode)
	onStripeLeaveFunc func(sbn *ScrollButtonNode)
	onKnobEnterFunc   func(sbn *ScrollButtonNode)
	onKnobLeaveFunc   func(sbn *ScrollButtonNode)
}

var _ Node = &ScrollButtonNode{}
var _ interactable = &ScrollButtonNode{}

func clamp(v float32) float32 {
	return min(max(v, 0), 1)
}

func (sbn *ScrollButtonNode) knobRect() rect {
	knobSizeX, knobSizeY := sbn.knob.getInitSize()
	sbnBounds := sbn.baseNode.bounds
	barX := sbnBounds.maxX - sbnBounds.minX
	barY := sbnBounds.maxY - sbnBounds.minY

	var knobPosX, knobPosY float32
	if sbn.IsVertical() {
		trackSize := barY - sbn.trackStart - sbn.trackEnd
		knobPosX = sbnBounds.minX + (barX-knobSizeX)*sbn.knobAlign
		knobPosY = sbnBounds.minY + sbn.trackStart + sbn.value*(trackSize-knobSizeY)
	} else {
		trackSize := barX - sbn.trackStart - sbn.trackEnd
		knobPosX = sbnBounds.minX + sbn.trackStart + sbn.value*(trackSize-knobSizeX)
		knobPosY = sbnBounds.minY + (barY-knobSizeY)*sbn.knobAlign
	}

	return rect{knobPosX, knobPosY, knobPosX + knobSizeX, knobPosY + knobSizeY}
}

func (sbn *ScrollButtonNode) draw(screen *ebiten.Image, showDebugInfo bool, depth int) {
	sbn.stripe.drawSpriteElement(sbn.baseNode.bounds, screen)

	if sbn.knob.image == nil {
		return
	}

	sbn.knob.drawSpriteElement(sbn.knobRect(), screen)
}

func (sbn *ScrollButtonNode) getInitSize() (float32, float32) {
	return sbn.stripe.getInitSize()
}

func (sbn *ScrollButtonNode) base() *baseNode {
	return &sbn.baseNode
}

func (sbn *ScrollButtonNode) mouseOn(mX, mY float32) {
	sbn.stripeMouseWasOn = true

	if !sbn.IsEnabled() {
		sbn.stripeMouseWasOn = false
		sbn.state = normal
		return
	}

	switch {
	case inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft):
		r := sbn.knobRect()

		if mX < r.minX || mX > r.maxX || mY < r.minY || mY > r.maxY {
			sbn.state = hovered
			return
		}

		sbn.state = draging
		if sbn.rootNode != nil {
			sbn.rootNode.dragging = sbn
		}
		sbn.knobMouseWasOn = true

		if sbn.IsVertical() {
			sbn.dragStartPos = mY
		} else {
			sbn.dragStartPos = mX
		}
		sbn.dragStartValue = sbn.value

		if sbn.eventsSource == nil {
			return
		}
		if sbn.eventsSource.onDragStartFunc != nil {
			sbn.eventsSource.onDragStartFunc(sbn)
		}
	case inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft):
		sbn.state = normal
		if sbn.rootNode != nil {
			sbn.rootNode.dragging = nil
		}

		if sbn.prevState != draging {
			return
		}
		if sbn.eventsSource == nil {
			return
		}
		if sbn.eventsSource.onDragEndFunc != nil {
			sbn.eventsSource.onDragEndFunc(sbn)
		}
	case ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft):
		if sbn.prevState != draging {
			return
		}

		sbn.state = draging

		if sbn.IsVertical() {
			sbnBounds := sbn.baseNode.bounds
			_, knobSizeY := sbn.knob.getInitSize()
			barY := sbnBounds.maxY - sbnBounds.minY
			trackSize := barY - sbn.trackStart - sbn.trackEnd

			delta := mY - sbn.dragStartPos
			sbn.value = clamp(sbn.dragStartValue + delta/(trackSize-knobSizeY))
		} else {
			sbnBounds := sbn.baseNode.bounds
			knobSizeX, _ := sbn.knob.getInitSize()
			barX := sbnBounds.maxX - sbnBounds.minX
			trackSize := barX - sbn.trackStart - sbn.trackEnd

			delta := mX - sbn.dragStartPos
			sbn.value = clamp(sbn.dragStartValue + delta/(trackSize-knobSizeX))
		}

		if sbn.eventsSource == nil {
			return
		}
		if sbn.eventsSource.onDragingFunc != nil {
			sbn.eventsSource.onDragingFunc(sbn)
		}
	default:
		sbn.state = hovered

		if sbn.prevState != hovered {
			if sbn.eventsSource == nil {
				return
			}
			if sbn.eventsSource.onStripeEnterFunc != nil {
				sbn.eventsSource.onStripeEnterFunc(sbn)
			}
		}

		r := sbn.knobRect()

		if mX >= r.minX && mX <= r.maxX && mY >= r.minY && mY <= r.maxY {
			sbn.knobMouseWasOn = true
			if !sbn.prevKnobMouseWasOn {
				if sbn.eventsSource == nil {
					return
				}
				if sbn.eventsSource.onKnobEnterFunc != nil {
					sbn.eventsSource.onKnobEnterFunc(sbn)
				}
			}
		}
	}
}

func (sbn *ScrollButtonNode) prepareFrame() {
	if !sbn.stripeMouseWasOn && sbn.prevStripeMouseWasOn {
		if sbn.eventsSource != nil {
			if sbn.eventsSource.onStripeLeaveFunc != nil {
				sbn.eventsSource.onStripeLeaveFunc(sbn)
			}
		}
	}

	if !sbn.knobMouseWasOn && sbn.prevKnobMouseWasOn {
		if sbn.eventsSource != nil {
			if sbn.eventsSource.onKnobLeaveFunc != nil {
				sbn.eventsSource.onKnobLeaveFunc(sbn)
			}
		}
	}

	sbn.prevStripeMouseWasOn = sbn.stripeMouseWasOn
	sbn.prevKnobMouseWasOn = sbn.knobMouseWasOn
	sbn.stripeMouseWasOn = false
	sbn.knobMouseWasOn = false

	sbn.prevState = sbn.state
	sbn.state = normal
}

func ScrollButton(opts ...ScrollButtonOption) *ScrollButtonNode {
	sbn := &ScrollButtonNode{}
	sbn.flags |= flagVisible | flagEnabled

	for _, opt := range opts {
		opt.applyToNode(sbn)
		opt.applyToScrollButton(sbn)
	}
	return sbn
}

func Value(v float32) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.value = v
	}
}
func (sbn *ScrollButtonNode) GetValue() float32 {
	return sbn.value
}
func (sbn *ScrollButtonNode) SetValue(v float32) {
	sbn.value = v
}

func KnobAlign(ka float32) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.knobAlign = ka
	}
}
func (sbn *ScrollButtonNode) GetKnobAlign() float32 {
	return sbn.knobAlign
}
func (sbn *ScrollButtonNode) SetKnobAlign(ka float32) {
	sbn.knobAlign = ka
}

func TrackStart(ts float32) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.trackStart = ts
	}
}
func (sbn *ScrollButtonNode) GetTrackStart() float32 {
	return sbn.trackStart
}
func (sbn *ScrollButtonNode) SetTrackStart(ts float32) {
	sbn.trackStart = ts
}

func TrackEnd(te float32) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.trackEnd = te
	}
}
func (sbn *ScrollButtonNode) GetTrackEnd() float32 {
	return sbn.trackEnd
}
func (sbn *ScrollButtonNode) SetTrackEnd(te float32) {
	sbn.trackEnd = te
}

// ScrollButtonEvents
func NewScrollButtonEventsSource(onDragStartFunc, onDragEndFunc, onDragingFunc, onStripeEnterFunc, onStripeLeaveFunc, onKnobEnterFunc, onKnobLeaveFunc func(bn *ScrollButtonNode)) *ScrollButtonEventsSource {
	return &ScrollButtonEventsSource{
		onDragStartFunc: onDragStartFunc,
		onDragEndFunc:   onDragEndFunc,
		onDragingFunc:   onDragingFunc,

		onStripeEnterFunc: onStripeEnterFunc,
		onStripeLeaveFunc: onStripeLeaveFunc,
		onKnobEnterFunc:   onKnobEnterFunc,
		onKnobLeaveFunc:   onKnobLeaveFunc,
	}
}

func ScrollButtonEvents(sbes *ScrollButtonEventsSource) ScrollButtonOpt {
	return func(sbe *ScrollButtonNode) {
		sbe.eventsSource = sbes
	}
}
func (sbn *ScrollButtonNode) SetButtonEvents(sbes *ScrollButtonEventsSource) {
	sbn.eventsSource = sbes
}
func (sbn *ScrollButtonNode) GetButtonEvents() *ScrollButtonEventsSource {
	return sbn.eventsSource
}

func OnDragStartFunc(onDragStartFunc func(*ScrollButtonNode)) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.ensureEventsSource()
		sbn.eventsSource.onDragStartFunc = onDragStartFunc
	}
}
func (sbn *ScrollButtonNode) SetOnDragStartFunc(onDragStartFunc func(*ScrollButtonNode)) {
	sbn.ensureEventsSource()
	sbn.eventsSource.onDragStartFunc = onDragStartFunc
}
func (sbn *ScrollButtonNode) GetOnDragStartFunc() func(*ScrollButtonNode) {
	if sbn.eventsSource == nil {
		return nil
	}
	return sbn.eventsSource.onDragStartFunc
}

func OnDragEndFunc(onDragEndFunc func(*ScrollButtonNode)) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.ensureEventsSource()
		sbn.eventsSource.onDragEndFunc = onDragEndFunc
	}
}
func (sbn *ScrollButtonNode) SetOnDragEndFunc(onDragEndFunc func(*ScrollButtonNode)) {
	sbn.ensureEventsSource()
	sbn.eventsSource.onDragEndFunc = onDragEndFunc
}
func (sbn *ScrollButtonNode) GetOnDragEndFunc() func(*ScrollButtonNode) {
	if sbn.eventsSource == nil {
		return nil
	}
	return sbn.eventsSource.onDragEndFunc
}

func OnDragingFunc(onDragingFunc func(*ScrollButtonNode)) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.ensureEventsSource()
		sbn.eventsSource.onDragingFunc = onDragingFunc
	}
}
func (sbn *ScrollButtonNode) SetOnDragingFunc(onDragingFunc func(*ScrollButtonNode)) {
	sbn.ensureEventsSource()
	sbn.eventsSource.onDragingFunc = onDragingFunc
}
func (sbn *ScrollButtonNode) GetOnDragingFunc() func(*ScrollButtonNode) {
	if sbn.eventsSource == nil {
		return nil
	}
	return sbn.eventsSource.onDragingFunc
}

func OnStripeEnterFunc(onStripeEnterFunc func(*ScrollButtonNode)) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.ensureEventsSource()
		sbn.eventsSource.onStripeEnterFunc = onStripeEnterFunc
	}
}
func (sbn *ScrollButtonNode) SetOnStripeEnterFunc(onStripeEnterFunc func(*ScrollButtonNode)) {
	sbn.ensureEventsSource()
	sbn.eventsSource.onStripeEnterFunc = onStripeEnterFunc
}
func (sbn *ScrollButtonNode) GetOnStripeEnterFunc() func(*ScrollButtonNode) {
	if sbn.eventsSource == nil {
		return nil
	}
	return sbn.eventsSource.onStripeEnterFunc
}

func OnStripeLeaveFunc(onStripeLeaveFunc func(*ScrollButtonNode)) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.ensureEventsSource()
		sbn.eventsSource.onStripeLeaveFunc = onStripeLeaveFunc
	}
}
func (sbn *ScrollButtonNode) SetOnStripeLeaveFunc(onStripeLeaveFunc func(*ScrollButtonNode)) {
	sbn.ensureEventsSource()
	sbn.eventsSource.onStripeLeaveFunc = onStripeLeaveFunc
}
func (sbn *ScrollButtonNode) GetOnStripeLeaveFunc() func(*ScrollButtonNode) {
	if sbn.eventsSource == nil {
		return nil
	}
	return sbn.eventsSource.onStripeLeaveFunc
}

func OnKnobEnterFunc(onKnobEnterFunc func(*ScrollButtonNode)) ScrollButtonOpt {
	return func(bn *ScrollButtonNode) {
		bn.ensureEventsSource()
		bn.eventsSource.onKnobEnterFunc = onKnobEnterFunc
	}
}
func (bn *ScrollButtonNode) SetOnKnobEnterFunc(onKnobEnterFunc func(*ScrollButtonNode)) {
	bn.ensureEventsSource()
	bn.eventsSource.onKnobEnterFunc = onKnobEnterFunc
}
func (bn *ScrollButtonNode) GetOnKnobEnterFunc() func(*ScrollButtonNode) {
	if bn.eventsSource == nil {
		return nil
	}
	return bn.eventsSource.onKnobEnterFunc
}

func OnKnobLeaveFunc(onKnobLeaveFunc func(*ScrollButtonNode)) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.ensureEventsSource()
		sbn.eventsSource.onKnobLeaveFunc = onKnobLeaveFunc
	}
}
func (sbn *ScrollButtonNode) SetOnKnobLeaveFunc(onKnobLeaveFunc func(*ScrollButtonNode)) {
	sbn.ensureEventsSource()
	sbn.eventsSource.onKnobLeaveFunc = onKnobLeaveFunc
}
func (sbn *ScrollButtonNode) GetOnKnobLeaveFunc() func(*ScrollButtonNode) {
	if sbn.eventsSource == nil {
		return nil
	}
	return sbn.eventsSource.onKnobLeaveFunc
}

// Image
func StripeImage(img *ebiten.Image) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.stripe.image = img
	}
}
func (sbn *ScrollButtonNode) SetStripeImage(img *ebiten.Image) {
	sbn.stripe.image = img
	if sbn.IsAutoSize() && sbn.rootNode != nil {
		sbn.rootNode.changed = true
	}
}
func (sbn *ScrollButtonNode) GetStripeImage() *ebiten.Image {
	return sbn.stripe.image
}

func StripeNineSlice(nss *NineSliceSource) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.stripe.nineSliceSource = nss
	}
}
func (sbn *ScrollButtonNode) SetStripeNineSlice(nss *NineSliceSource) {
	sbn.stripe.nineSliceSource = nss
}
func (sbn *ScrollButtonNode) GetStripeNineSlice() *NineSliceSource {
	return sbn.stripe.nineSliceSource
}

func StripeColorScale(cs ebiten.ColorScale) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.stripe.colorScale = cs
	}
}
func (sbn *ScrollButtonNode) SetStripeColorScale(cs ebiten.ColorScale) {
	sbn.stripe.colorScale = cs
}
func (sbn *ScrollButtonNode) GetStripeColorScale() ebiten.ColorScale {
	return sbn.stripe.colorScale
}

func StripeFilter(f ebiten.Filter) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.stripe.filter = f
	}
}
func (sbn *ScrollButtonNode) SetStripeFilter(f ebiten.Filter) {
	sbn.stripe.filter = f
}
func (sbn *ScrollButtonNode) GetStripeFilter() ebiten.Filter {
	return sbn.stripe.filter
}

func StripeBlend(b ebiten.Blend) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.stripe.blend = b
	}
}
func (sbn *ScrollButtonNode) SetStripeBlend(b ebiten.Blend) {
	sbn.stripe.blend = b
}
func (sbn *ScrollButtonNode) GetStripeBlend() ebiten.Blend {
	return sbn.stripe.blend
}

func KnobImage(img *ebiten.Image) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.knob.image = img
	}
}
func (sbn *ScrollButtonNode) SetKnobImage(img *ebiten.Image) {
	sbn.knob.image = img
}
func (sbn *ScrollButtonNode) GetKnobImage() *ebiten.Image {
	return sbn.knob.image
}

func KnobNineSlice(nss *NineSliceSource) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.knob.nineSliceSource = nss
	}
}
func (sbn *ScrollButtonNode) SetKnobNineSlice(nss *NineSliceSource) {
	sbn.knob.nineSliceSource = nss
}
func (sbn *ScrollButtonNode) GetKnobNineSlice() *NineSliceSource {
	return sbn.knob.nineSliceSource
}

func KnobColorScale(cs ebiten.ColorScale) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) { sbn.knob.colorScale = cs }
}
func (sbn *ScrollButtonNode) SetKnobColorScale(cs ebiten.ColorScale) {
	sbn.knob.colorScale = cs
}
func (sbn *ScrollButtonNode) GetKnobColorScale() ebiten.ColorScale {
	return sbn.knob.colorScale
}

func KnobFilter(f ebiten.Filter) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.knob.filter = f
	}
}
func (sbn *ScrollButtonNode) SetKnobFilter(f ebiten.Filter) {
	sbn.knob.filter = f
}
func (sbn *ScrollButtonNode) GetKnobFilter() ebiten.Filter {
	return sbn.knob.filter
}

func KnobBlend(b ebiten.Blend) ScrollButtonOpt {
	return func(sbn *ScrollButtonNode) {
		sbn.knob.blend = b
	}
}
func (sbn *ScrollButtonNode) SetKnobBlend(b ebiten.Blend) {
	sbn.knob.blend = b
}
func (sbn *ScrollButtonNode) GetKnobBlend() ebiten.Blend {
	return sbn.knob.blend
}
