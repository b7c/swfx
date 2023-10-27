package swfx

import "github.com/b7c/swfx/tagcode"

type Metadata struct {
	Value string
}

func (tag *Metadata) Code() tagcode.TagCode {
	return tagcode.Metadata
}

func (tag *Metadata) readData(r SwfReader, length int) {
	tag.Value = r.ReadString(length)
}
