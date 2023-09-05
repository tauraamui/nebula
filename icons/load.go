package icons

import (
	_ "embed"

	"github.com/inkeliz/giosvg"
)

//go:embed  mouse-pointer.svg
var mousePointer []byte

func MousePointer() (*giosvg.Icon, error) {
	vec, err := giosvg.NewVector(mousePointer)
	if err != nil {
		return nil, err
	}

	return giosvg.NewIcon(vec), nil
}
