package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/min-ok/simple-ui"
	"github.com/min-ok/simple-ui/kit"
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
	g.root.UpdateLayout(1.0/60.0, outsideWidth, outsideHeight)
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
