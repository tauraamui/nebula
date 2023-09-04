package main

import (
	"log"
	"os"

	"gioui.org/app"
	"github.com/tauraamui/nebula/widgets"
)

func main() {
	canvas := widgets.Canvas{}
	go func() {
		w := app.NewWindow(
			app.Title("github.com/tauraamui/nebula"),
		)
		err := canvas.Run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
