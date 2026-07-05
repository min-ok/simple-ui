package ui

const (
	normal buttonState = iota
	justPressed
	justReleased
	pressed
	hovered
	draging
)

type interactable interface {
	prepareFrame()
	mouseOn(mX, mY float32)
}

func prepareFrame(bn Node) {
	if inter, ok := bn.(interactable); ok {
		inter.prepareFrame()
	}

	for _, child := range bn.base().children {
		prepareFrame(child)
	}
}


func findCursorInChildren(node Node, mX, mY float32) {
	nodeBase := node.base()

	if mX < nodeBase.bounds.minX || mX > nodeBase.bounds.maxX || mY < nodeBase.bounds.minY || mY > nodeBase.bounds.maxY {
		return
	}

	if int, ok := node.(interactable); ok {
		int.mouseOn(mX, mY)
	}

	for _, childNode := range nodeBase.children {
		findCursorInChildren(childNode, mX, mY)
	}
}


func (bn *ButtonNode) ensureEventsSource() {
	if bn.eventsSource == nil {
		bn.eventsSource = &ButtonEventsSource{}
	}
}

func (sbn *ScrollButtonNode) ensureEventsSource() {
	if sbn.eventsSource == nil {
		sbn.eventsSource = &ScrollButtonEventsSource{}
	}
}
