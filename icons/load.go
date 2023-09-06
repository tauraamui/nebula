package icons

import (
	_ "embed"

	"github.com/inkeliz/giosvg"
)

//go:embed  mouse-pointer.svg
var mousePointer []byte

//go:embed  mouse-pointer-outline.svg
var mousePointerOutline []byte

//go:embed  square.svg
var square []byte

//go:embed  square-border.svg
var squareBorder []byte

type IconResolver func() (*giosvg.Icon, error)

func MousePointer() (*giosvg.Icon, error) {
	return makeIcon(mousePointer)
}

func MousePointerOutline() (*giosvg.Icon, error) {
	return makeIcon(mousePointerOutline)
}

func Square() (*giosvg.Icon, error) {
	return makeIcon(square)
}

func SquareBorder() (*giosvg.Icon, error) {
	return makeIcon(squareBorder)
}

func makeIcon(b []byte) (*giosvg.Icon, error) {
	vec, err := giosvg.NewVector(b)
	if err != nil {
		return nil, err
	}

	return giosvg.NewIcon(vec), nil
}
