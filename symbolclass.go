package swfx

import (
	"fmt"

	"github.com/b7c/swfx/tagcode"
)

type SymbolClass struct {
	Symbols map[string]int
}

func (tag *SymbolClass) Code() tagcode.TagCode {
	return tagcode.SymbolClass
}

func (tag *SymbolClass) readData(r SwfReader, length int) {
	end := r.Position() + length
	n := int(r.ReadUint16())
	tag.Symbols = make(map[string]int, n)
	for i := 0; i < n; i++ {
		id := int(r.ReadUint16())
		name := r.ReadString(end - r.Position())
		if _, exist := tag.Symbols[name]; exist {
			panic(fmt.Errorf("duplicate symbol name: %q", name))
		}
		tag.Symbols[name] = id
	}
}
