package swfx

import "b7c.io/swfx/tagcode"

type DefineBinaryData struct {
	characterTag
	Data []byte
}

func (tag *DefineBinaryData) Code() tagcode.TagCode {
	return tagcode.DefineBinaryData
}

func (tag *DefineBinaryData) readData(r SwfReader, length int) {
	tag.characterTag.readData(r, length)
	r.ReadUint32() // reserved
	tag.Data = make([]byte, length-6)
	r.MustRead(tag.Data)
}
