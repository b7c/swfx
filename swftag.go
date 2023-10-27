package swfx

import (
	"github.com/b7c/swfx/tagcode"
)

type SwfTag interface {
	tagDataReader
	Code() tagcode.TagCode
}

type CharacterTag interface {
	SwfTag
	CharacterId() int
	SetCharacterId(id int)
}

type SwfTagStruct struct {
	Type int
}

type ArbitraryTag struct {
	tagCode tagcode.TagCode
	Data    []byte
}

func (tag *ArbitraryTag) Code() tagcode.TagCode {
	return tag.tagCode
}

func NewArbitraryTag(code tagcode.TagCode, data []byte) *ArbitraryTag {
	return &ArbitraryTag{
		tagCode: code,
		Data:    data,
	}
}

type tagDataReader interface {
	readData(r SwfReader, length int)
}

func (tag *ArbitraryTag) readData(r SwfReader, length int) {
	tag.Data = make([]byte, length)
	r.MustRead(tag.Data)
}

type characterTag struct {
	characterId int
}

func (tag *characterTag) CharacterId() int {
	return tag.characterId
}

func (tag *characterTag) SetCharacterId(id int) {
	tag.characterId = id
}

func (tag *characterTag) readData(r SwfReader, length int) {
	tag.characterId = int(r.ReadUint16())
}

func MakeTag(code tagcode.TagCode) SwfTag {
	switch code {
	case tagcode.Metadata:
		return &Metadata{}
	case tagcode.DefineBitsJpeg2:
		return &DefineBitsJpeg2{}
	case tagcode.DefineBitsLossless2:
		return &DefineBitsLossless2{}
	case tagcode.DefineBinaryData:
		return &DefineBinaryData{}
	case tagcode.SymbolClass:
		return &SymbolClass{}
	default:
		return &ArbitraryTag{tagCode: code}
	}
}
