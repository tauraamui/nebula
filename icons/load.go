package icons

import (
	_ "embed"

	"github.com/inkeliz/giosvg"
)

//go:embed  mouse-pointer.svg
var mousePointer []byte

//go:embed  square.svg
var square []byte

type IconResolver func() (*giosvg.Icon, error)

func MousePointer() (*giosvg.Icon, error) {
	return makeIcon(mousePointer)
}

func Square() (*giosvg.Icon, error) {
	return makeIcon(square)
}

func makeIcon(b []byte) (*giosvg.Icon, error) {
	vec, err := giosvg.NewVector(b)
	if err != nil {
		return nil, err
	}

	return giosvg.NewIcon(vec), nil
}
