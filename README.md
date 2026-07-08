# Simple-ui
## What is this
- This is little ui system based on [ebiten](https://github.com/hajimehoshi/ebiten). Using flex elements, almost like CSS [flexbox](https://css-tricks.com/snippets/css/a-guide-to-flexbox/).
## For whom
- For those who want to add ui to their projects without having to learn large frameworks and maintain development flexibility.
## Examples
### [Minimal start](examples/minimal/main.go)
```go
package main

import (
	"log"
	ui "simple-ui"
	"simple-ui/kit"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	root *ui.RootNode
}

func (g *Game) Update() error {
	// g.root.UpdateCursorLogic() // if you have interactive nodes like button
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.root.Draw(screen, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.root.UpdateLayout(1.0 / 60.0, outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

func main() {
	game := &Game{}

	ebiten.SetWindowSize(512, 512)
	ebiten.SetWindowTitle("Hello simple-ui")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game.root = ui.Root(
		kit.VBox(1),
	)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
```
### [Nested Vbox](examples/nested-vbox/main.go)
![](docs/layout/03-nested-vbox.svg)
```go
game.root = ui.Root(
	kit.VBox(1, ui.Border(20), ui.Spacing(20),
		ui.Box(ui.FixY(50)),

		kit.HBox(1, ui.Spacing(20),
			ui.Box(ui.FixX(140)),
			ui.Box(ui.Weight(1)),
			ui.Box(ui.Weight(2)),
		),

		ui.Box(ui.FixY(50)),
	),
)
```
### [Nested Hbox](examples/nested-hbox/main.go)
![](docs/layout/04-nested-hbox.svg)
```go
game.root = ui.Root(
	kit.HBox(1, ui.Border(20), ui.Spacing(20),
		ui.Box(ui.FixX(140)),

		kit.VBox(1, ui.Spacing(15),
			ui.Box(ui.FixY(80)),
			ui.Box(ui.Weight(1)),
		),

		ui.Box(ui.Weight(1)),
	),
)
```
## A few words about align
![](docs/layout/01-align.svg)
## A few words about spacing and border
![](docs/layout/02-spacing-border.svg)
