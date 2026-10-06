package ui

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type LabelNode struct {
	baseNode
	labelElement
}

type labelElement struct {
	text           string
	goTextFace     *text.GoTextFace
	color          color.RGBA
	primaryAlign   text.Align
	secondaryAlign text.Align
	lineSpacing    float64
	wordWrap       bool
}

func (ln *LabelNode) draw(screen *ebiten.Image, showDebugInfo bool, depth int) {
	bounds := ln.bounds
	cx, cy := 0.0, 0.0

	switch ln.primaryAlign {
	case text.AlignStart:
		cx = float64(bounds.minX)
	case text.AlignEnd:
		cx = float64(bounds.maxX)
	default:
		cx = float64(bounds.minX + (bounds.maxX-bounds.minX)/2)
	}
	switch ln.secondaryAlign {
	case text.AlignStart:
		cy = float64(bounds.minY)
	case text.AlignEnd:
		cy = float64(bounds.maxY)
	default:
		cy = float64(bounds.minY + (bounds.maxY-bounds.minY)/2)
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(cx, cy)
	op.PrimaryAlign = ln.primaryAlign
	op.SecondaryAlign = ln.secondaryAlign

	lineSpacing := ln.lineSpacing
	if lineSpacing == 0 {
		lineSpacing = ln.goTextFace.Size
	}
	op.LineSpacing = lineSpacing

	op.ColorScale.ScaleWithColor(ln.color)

	t := ln.text
	if ln.wordWrap {
		t = wrapText(ln.text, ln.goTextFace, float64(bounds.maxX-bounds.minX))
	}

	text.Draw(screen, t, ln.goTextFace, op)
}

func wrapText(s string, face *text.GoTextFace, maxWidth float64) string {
	var result strings.Builder
	lines := strings.Split(s, "\n")
	space := text.Advance(" ", face)
	for i, line := range lines {
		if i > 0 {
			result.WriteByte('\n')
		}
		words := strings.Fields(line)
		if len(words) == 0 {
			continue
		}
		var lineWidth float64
		for j, word := range words {
			w := text.Advance(word, face)
			if j == 0 {
				result.WriteString(word)
				lineWidth = w
			} else if lineWidth+space+w <= maxWidth {
				result.WriteByte(' ')
				result.WriteString(word)
				lineWidth += space + w
			} else {
				result.WriteByte('\n')
				result.WriteString(word)
				lineWidth = w
			}
		}
	}
	return result.String()
}

func (ln *LabelNode) getInitSize() (float32, float32) {
	lineSpacing := ln.lineSpacing
	if lineSpacing == 0 {
		lineSpacing = ln.goTextFace.Size
	}
	sizeX, sizeY := text.Measure(ln.text, ln.goTextFace, lineSpacing)
	return float32(sizeX), float32(sizeY)
}

func (ln *LabelNode) base() *baseNode {
	return &ln.baseNode
}

func Label(opts ...LabelOption) *LabelNode {
	ln := &LabelNode{}
	ln.flags |= flagVisible | flagEnabled

	ln.labelElement = labelElement{
		color:          color.RGBA{255, 255, 255, 255},
		primaryAlign:   text.AlignCenter,
		secondaryAlign: text.AlignCenter,
	}

	for _, opt := range opts {
		opt.applyToNode(ln)
		opt.applyToLabel(ln)
	}

	return ln
}

func Text(t string) LabelOpt {
	return func(ln *LabelNode) {
		ln.text = t
	}
}
func (ln *LabelNode) SetText(t string) {
	ln.text = t
	if ln.IsAutoSize() && ln.rootNode != nil {
		ln.rootNode.changed = true
	}
}
func (ln *LabelNode) GetText() string {
	return ln.text
}

func Color(c color.RGBA) LabelOpt {
	return func(ln *LabelNode) {
		ln.color = c
	}
}
func (ln *LabelNode) SetColor(c color.RGBA) {
	ln.color = c
}
func (ln *LabelNode) GetColor() color.RGBA {
	return ln.color
}

func PrimaryAlign(a text.Align) LabelOpt {
	return func(ln *LabelNode) {
		ln.primaryAlign = a
	}
}
func (ln *LabelNode) SetPrimaryAlign(a text.Align) {
	ln.primaryAlign = a
}
func (ln *LabelNode) GetPrimaryAlign() text.Align {
	return ln.primaryAlign
}

func SecondaryAlign(a text.Align) LabelOpt {
	return func(ln *LabelNode) {
		ln.secondaryAlign = a
	}
}
func (ln *LabelNode) SetSecondaryAlign(a text.Align) {
	ln.secondaryAlign = a
}
func (ln *LabelNode) GetSecondaryAlign() text.Align {
	return ln.secondaryAlign
}

func LineSpacing(s float64) LabelOpt {
	return func(ln *LabelNode) {
		ln.lineSpacing = s
	}
}
func (ln *LabelNode) SetLineSpacing(s float64) {
	ln.lineSpacing = s
	if ln.IsAutoSize() && ln.rootNode != nil {
		ln.rootNode.changed = true
	}
}
func (ln *LabelNode) GetLineSpacing() float64 {
	return ln.lineSpacing
}

func GoTextFace(goTextFace *text.GoTextFace) LabelOpt {
	return func(ln *LabelNode) {
		ln.goTextFace = goTextFace
	}
}

func (ln *LabelNode) SetGoTextFace(goTextFace *text.GoTextFace) {
	ln.goTextFace = goTextFace
	if ln.IsAutoSize() && ln.rootNode != nil {
		ln.rootNode.changed = true
	}
}

func (ln *LabelNode) GetGoTextFace() *text.GoTextFace {
	return ln.goTextFace
}

func WordWrap() LabelOpt {
	return func(ln *LabelNode) {
		ln.wordWrap = true
	}
}

func NoWordWrap() LabelOpt {
	return func(ln *LabelNode) {
		ln.wordWrap = false
	}
}

func (ln *LabelNode) EnableWordWrap() {
	ln.wordWrap = true
}
func (ln *LabelNode) DisableWordWrap() {
	ln.wordWrap = false
}

func (ln *LabelNode) IsWordWrap() bool {
	return ln.wordWrap
}
