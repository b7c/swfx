package swfx

import (
	"bytes"

	"github.com/b7c/swfx/tagcode"
)

type ImageType int

const (
	Jpeg ImageType = iota
	Png
	Gif
)

var jpegHeader = []byte{0xFF, 0xD8}
var pngHeader = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
var gifHeader = []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61}

type DefineBitsJpeg2 struct {
	characterTag
	ImageData []byte
}

func (tag *DefineBitsJpeg2) Code() tagcode.TagCode {
	return tagcode.DefineBitsJpeg2
}

func (tag *DefineBitsJpeg2) readData(r SwfReader, length int) {
	tag.characterTag.readData(r, length)
	tag.ImageData = make([]byte, length-2)
	r.MustRead(tag.ImageData)
}

func (tag *DefineBitsJpeg2) ImageType() ImageType {
	if bytes.Equal(tag.ImageData[:len(jpegHeader)], jpegHeader) {
		return Jpeg
	} else if bytes.Equal(tag.ImageData[:len(pngHeader)], pngHeader) {
		return Png
	} else if bytes.Equal(tag.ImageData[:len(gifHeader)], gifHeader) {
		return Gif
	} else {
		return -1
	}
}
