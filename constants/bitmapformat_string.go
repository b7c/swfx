// Code generated by "stringer -type=BitmapFormat"; DO NOT EDIT.

package constants

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Colormap8-3]
	_ = x[Argb32-5]
}

const (
	_BitmapFormat_name_0 = "Colormap8"
	_BitmapFormat_name_1 = "Argb32"
)

func (i BitmapFormat) String() string {
	switch {
	case i == 3:
		return _BitmapFormat_name_0
	case i == 5:
		return _BitmapFormat_name_1
	default:
		return "BitmapFormat(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
