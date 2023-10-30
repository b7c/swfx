package swfx

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"image"
	"image/color"
	"io"

	"github.com/b7c/swfx/constants"
	"github.com/b7c/swfx/tagcode"
)

type DefineBitsLossless2 struct {
	characterTag
	Format         constants.BitmapFormat
	Width          int
	Height         int
	ColorTableSize int
	ZlibBitmapData []byte
}

func (tag *DefineBitsLossless2) Code() tagcode.TagCode {
	return tagcode.DefineBitsLossless2
}

func (tag *DefineBitsLossless2) readData(r SwfReader, length int) {
	size := length - 7
	tag.characterTag.readData(r, length)
	tag.Format = constants.BitmapFormat(r.ReadUint8())
	tag.Width = int(r.ReadUint16())
	tag.Height = int(r.ReadUint16())
	if tag.Format == constants.Colormap8 {
		tag.ColorTableSize = int(r.ReadUint8())
		size--
	}
	tag.ZlibBitmapData = make([]byte, size)
	r.MustRead(tag.ZlibBitmapData)
}

func (tag *DefineBitsLossless2) Decode() (img image.Image, err error) {
	z, err := zlib.NewReader(bytes.NewReader(tag.ZlibBitmapData))
	if err != nil {
		return
	}
	pixels, err := io.ReadAll(z)
	if err != nil {
		return
	}
	switch tag.Format {
	case constants.Argb32:
		rgba := image.NewRGBA(image.Rect(0, 0, tag.Width, tag.Height))
		for i := 0; i < len(pixels); i += 4 {
			x, y := (i/4)%tag.Width, (i/4)/tag.Width
			rgba.SetRGBA(x, y, color.RGBA{
				A: pixels[i+0],
				R: pixels[i+1],
				G: pixels[i+2],
				B: pixels[i+3],
			})
		}
		img = rgba
	default:
		err = fmt.Errorf("unsupported image type")
		return
	}
	return
}

func (tag *DefineBitsLossless2) String() string {
	return fmt.Sprintf("DefineBitsLossless2[%dx%d:%s]", tag.Width, tag.Height, tag.Format)
}
