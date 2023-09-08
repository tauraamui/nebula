package main

import (
	"log"
	"os"
	"runtime/trace"

	"gioui.org/app"
	"github.com/tauraamui/nebula/nebula"
)

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}

	appx := nebula.New()
	go func() {
		if err := appx.Run(); err != nil {
			log.Fatal(err)
		}
		trace.Stop()
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
		os.Exit(0)
	}()
	app.Main()
}
