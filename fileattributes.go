package swfx

import "github.com/b7c/swfx/tagcode"

type FileAttributes struct {
	UseDirectBlit      bool
	UseGPU             bool
	HasMetaData        bool
	ActionScript3      bool
	NoCrossDomainCache bool
	UseNetwork         bool
}

func (tag *FileAttributes) Code() tagcode.TagCode {
	return tagcode.FileAttributes
}

func (tag *FileAttributes) readData(r SwfReader, length int) {
	r.ReadUbits(1) // reserved
	tag.UseDirectBlit = r.ReadBool()
	tag.UseGPU = r.ReadBool()
	tag.HasMetaData = r.ReadBool()
	tag.ActionScript3 = r.ReadBool()
	tag.NoCrossDomainCache = r.ReadBool()
	r.ReadUbits(1) // reserved
	tag.UseNetwork = r.ReadBool()
	r.ReadUbits(24) // reserved
}

func (tag *FileAttributes) String() string {
	return "FileAttributes"
}
