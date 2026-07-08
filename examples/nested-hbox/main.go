package main

import (
	"log"

	ui "github.com/min-ok/simple-ui"
	"github.com/min-ok/simple-ui/kit"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	root *ui.RootNode
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.root.Draw(screen, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.root.UpdateLayout(1.0/60.0, outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

func main() {
	game := &Game{}

	ebiten.SetWindowSize(680, 360)
	ebiten.SetWindowTitle("Nested HBox + VBox")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

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

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
