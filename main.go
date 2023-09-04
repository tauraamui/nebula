package main

import (
	"log"
	"os"

	"gioui.org/app"
	"github.com/tauraamui/nebula/nebula"
)

func main() {
	appx := nebula.New()
	go func() {
		if err := appx.Run(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
