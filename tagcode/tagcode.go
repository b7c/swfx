//go:generate stringer -type=TagCode

package tagcode

type TagCode int

const (
	End                          TagCode = 0
	ShowFrame                    TagCode = 1
	PlaceObject                  TagCode = 4
	RemoveObject                 TagCode = 5
	DefineBits                   TagCode = 6
	JpegTables                   TagCode = 8 // JPEGTables
	SetBackgroundColor           TagCode = 9
	DefineBitsLossless           TagCode = 20
	DefineBitsJpeg2              TagCode = 21 // DefineBitsJPEG2
	Protect                      TagCode = 24
	PlaceObject2                 TagCode = 26
	RemoveObject2                TagCode = 28
	DefineBitsJpeg3              TagCode = 35 // DefineBitsJPEG3
	DefineBitsLossless2          TagCode = 36
	FrameLabel                   TagCode = 43
	NamedAnchor                  TagCode = 43 // extension of FrameLabel
	ExportAssets                 TagCode = 56
	ImportAssets                 TagCode = 57
	EnableDebugger               TagCode = 58
	EnableDebugger2              TagCode = 64
	ScriptLimits                 TagCode = 65
	SetTabIndex                  TagCode = 66
	FileAttributes               TagCode = 69
	PlaceObject3                 TagCode = 70
	ImportAssets2                TagCode = 71
	SymbolClass                  TagCode = 76
	Metadata                     TagCode = 77
	DefineScalingGrid            TagCode = 78
	DefineSceneAndFrameLabelData TagCode = 86
	DefineBinaryData             TagCode = 87
	DefineBitsJpeg4              TagCode = 90 // DefineBitsJPEG4
	EnableTelemetry              TagCode = 93

	// SWF 3 actions

	DoAction TagCode = 12
)
