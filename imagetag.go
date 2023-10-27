package swfx

import "image"

type ImageTag interface {
	SwfTag
	Decode() (image.Image, error)
}
