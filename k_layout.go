package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func arrange(node Node) {
	nodeBase := node.base()
	nodeBase.occupiedX = 0
	nodeBase.occupiedY = 0

	var occupiedSumX, occupiedSumY float32
	var maxX, maxY float32
	var visibleChildren = 0

	for _, childNode := range nodeBase.children {
		childNodeBase := childNode.base()

		if !childNodeBase.IsVisible() {
			continue
		}

		visibleChildren += 1

		arrange(childNode)

		occupiedSumX += childNodeBase.fixX
		occupiedSumY += childNodeBase.fixY

		maxX = max(maxX, childNodeBase.fixX)
		maxY = max(maxY, childNodeBase.fixY)
	}

	if visibleChildren == 0 {
		if nodeBase.IsAutoSize() {
			nodeBase.fixX, nodeBase.fixY = node.getInitSize()
		}
		return
	}

	spacing := float32(visibleChildren-1)*nodeBase.spacing + 2*nodeBase.border

	if nodeBase.IsVertical() {
		nodeBase.occupiedX = max(nodeBase.occupiedX, maxX)
		nodeBase.occupiedY = max(nodeBase.occupiedY, occupiedSumY+spacing)
	} else {
		nodeBase.occupiedX = max(nodeBase.occupiedX, occupiedSumX+spacing)
		nodeBase.occupiedY = max(nodeBase.occupiedY, maxY)
	}

	if nodeBase.IsAutoSize() {
		nodeBase.fixX, nodeBase.fixY = nodeBase.occupiedX, nodeBase.occupiedY
	}
}

func calculateBounds(node Node) {
	nodeBase := node.base()

	if len(nodeBase.children) == 0 {
		return
	}

	var sumWeight float32

	for _, childNode := range nodeBase.children {
		childNodeBase := childNode.base()

		if !childNodeBase.IsVisible() {
			continue
		}

		if nodeBase.IsVertical() && childNodeBase.fixY == 0 {
			sumWeight += effectiveWeight(childNodeBase.weight)
		} else if !nodeBase.IsVertical() && childNodeBase.fixX == 0 {
			sumWeight += effectiveWeight(childNodeBase.weight)
		}
	}

	nodeSizeX := nodeBase.bounds.maxX - nodeBase.bounds.minX
	nodeSizeY := nodeBase.bounds.maxY - nodeBase.bounds.minY

	remainingX := nodeSizeX - nodeBase.occupiedX
	remainingY := nodeSizeY - nodeBase.occupiedY

	shift := nodeBase.border
	if sumWeight == 0 {
		if nodeBase.IsVertical() {
			shift += remainingY * float32(nodeBase.contentAlign)
		} else {
			shift += remainingX * float32(nodeBase.contentAlign)
		}
	}

	for _, childNode := range nodeBase.children {
		childNodeBase := childNode.base()
		var x1, x2, y1, y2 float32

		if !childNodeBase.IsVisible() {
			continue
		}

		if nodeBase.IsVertical() {
			childSizeX := calculateSize(childNodeBase.fixX, 1, 1, nodeSizeX)
			x1 = nodeBase.bounds.minX + (nodeSizeX-childSizeX)*float32(childNodeBase.selfAlign)
			x2 = x1 + childSizeX

			childSizeY := calculateSize(childNodeBase.fixY, effectiveWeight(childNodeBase.weight), sumWeight, remainingY)
			y1 = nodeBase.bounds.minY + shift
			y2 = y1 + childSizeY
			shift += childSizeY
		} else {
			childSizeX := calculateSize(childNodeBase.fixX, effectiveWeight(childNodeBase.weight), sumWeight, remainingX)
			x1 = nodeBase.bounds.minX + shift
			x2 = x1 + childSizeX
			shift += childSizeX

			childSizeY := calculateSize(childNodeBase.fixY, 1, 1, nodeSizeY)
			y1 = nodeBase.bounds.minY + (nodeSizeY-childSizeY)*float32(childNodeBase.selfAlign)
			y2 = y1 + childSizeY
		}

		childNodeBase.bounds = rect{x1, y1, x2, y2}

		shift += nodeBase.spacing
		calculateBounds(childNode)
	}
}

func effectiveWeight(w float32) float32 {
	if w == 0 {
		return 1
	}
	return w
}

func calculateSize(init, weight, sumWeight, remaining float32) float32 {
	if init != 0 {
		return init
	}

	return (weight / sumWeight) * remaining
}

func drawTree(n Node, screen *ebiten.Image, showDebugInfo bool, depth int) {
	nodeBase := n.base()

	if !nodeBase.IsVisible() {
		return
	}

	n.draw(screen, showDebugInfo, depth)

	for _, v := range nodeBase.children {
		drawTree(v, screen, showDebugInfo, depth+1)
	}
}
