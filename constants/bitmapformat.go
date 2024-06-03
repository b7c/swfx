//go:generate stringer -type=BitmapFormat

package constants

type BitmapFormat int

const (
	Colormap8 BitmapFormat = 3
	Argb32    BitmapFormat = 5
)
