package swfx

import (
	"errors"
	"fmt"
	"io"

	"b7c.io/swfx/tagcode"
)

type Swf struct {
	Header         SwfHeader
	Tags           []SwfTag
	Characters     map[CharacterId]CharacterTag
	Symbols        map[string]CharacterId
	ReverseSymbols map[CharacterId][]string
}

func NewSwf() *Swf {
	return &Swf{}
}

func ReadSwf(reader io.Reader) (swf *Swf, err error) {
	defer func() {
		if e := recover(); e != nil {
			switch v := e.(type) {
			case string:
				err = errors.New(v)
			case error:
				err = v
			default:
				err = errors.New(fmt.Sprint(v))
			}
		}
	}()

	swfReader := NewReader(reader)

	var header SwfHeader
	header, err = ReadHeader(swfReader)
	if err != nil {
		return
	}

	swf = &Swf{
		Header:         header,
		Tags:           []SwfTag{},
		Characters:     map[CharacterId]CharacterTag{},
		Symbols:        map[string]CharacterId{},
		ReverseSymbols: map[CharacterId][]string{},
	}

	for {
		tagCode, length := swfReader.ReadTagCodeAndLength()
		tag := swfReader.ReadTag(tagCode, length)
		if tag, ok := tag.(CharacterTag); ok {
			swf.Characters[tag.CharacterId()] = tag
		}
		if tag, ok := tag.(*SymbolClass); ok {
			for name, id := range tag.Symbols {
				swf.Symbols[name] = id
				swf.ReverseSymbols[id] = append(swf.ReverseSymbols[id], name)
			}
		}
		swf.Tags = append(swf.Tags, tag)
		if tagCode == tagcode.End {
			break
		}
	}

	err = nil
	return
}
