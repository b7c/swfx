package swfx

type SwfHeader struct {
	Version    int
	FrameSize  Rect
	FrameRate  uint16
	FrameCount uint16
}
