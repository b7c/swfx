package swfx

import "github.com/b7c/swfx/tagcode"

type SymbolClass struct {
	Names map[int]string
}

func (tag *SymbolClass) Code() tagcode.TagCode {
	return tagcode.SymbolClass
}

func (tag *SymbolClass) readData(r SwfReader, length int) {
	tag.Names = map[int]string{}

	end := r.Position() + length
	n := int(r.ReadUint16())
	for i := 0; i < n; i++ {
		id := int(r.ReadUint16())
		name := r.ReadString(end - r.Position())
		tag.Names[id] = name
	}
}
