package swfx

import (
	"b7c.io/swfx/tagcode"
)

type SwfTag interface {
	tagDataReader
	Code() tagcode.TagCode
}

type CharacterId int

type CharacterTag interface {
	SwfTag
	CharacterId() CharacterId
	SetCharacterId(id CharacterId)
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
	characterId CharacterId
}

func (tag *characterTag) CharacterId() CharacterId {
	return tag.characterId
}

func (tag *characterTag) SetCharacterId(id CharacterId) {
	tag.characterId = id
}

func (tag *characterTag) readData(r SwfReader, length int) {
	tag.characterId = CharacterId(r.ReadUint16())
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
