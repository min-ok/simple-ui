package kit


import (
	"github.com/min-ok/simple-ui"
)


const (
	Center float32 = 0.5
	Right = 1
	Left = 0
	Top = 0
	Bottom = 1
)


// Spacer creates a Box Node.
// Needed to arrange other boxes.
// By default, the box resizes based on the first weight parameter. To avoid this, use FixX FixY FixSize.
func Spacer(weight float32) *ui.BoxNode {
	return ui.Box(ui.Weight(weight))
}


// VBox creates a Vertical Node.
// The children of this box will be arranged vertically.
// By default, the box resizes based on the first weight parameter. To avoid this, use FixX FixY FixSize.
func VBox(weight float32, opts ...ui.BoxOption) *ui.BoxNode {
	all := []ui.BoxOption{ui.Weight(weight), ui.Vertical()}
	all = append(all, opts...)
	return ui.Box(all...)
}


// HBox creates a Horizontal Node.
// The children of this box will be arranged horizontally.
// By default, the box resizes based on the first weight parameter. To avoid this, use FixX FixY FixSize.
func HBox(weight float32, opts ...ui.BoxOption) *ui.BoxNode {
	all := []ui.BoxOption{ui.Weight(weight), ui.Horizontal()}
	all = append(all, opts...)
	return ui.Box(all...)
}

func FixedVBox(x, y float32, opts ...ui.BoxOption) *ui.BoxNode {
	all := []ui.BoxOption{ui.FixSize(x, y), ui.Vertical()}
	all = append(all, opts...)
	return ui.Box(all...)
}

func FixedHBox(x, y float32, opts ...ui.BoxOption) *ui.BoxNode {
	all := []ui.BoxOption{ui.FixSize(x, y), ui.Horizontal()}
	all = append(all, opts...)
	return ui.Box(all...)
}
