# swfx

A command-line SWF disassembly tool & Go library.

# installation

## go install

Requires the Go compiler. To install the command-line toolkit:

```sh
go install b7c.io/swfx/cmd/swfx@latest
```

# usage

## listing resources

The `ls` command will list each SWF tag's offset, size, name and tag code. If the tag has an associated character ID, it will be shown in square brackets.

```
$ swfx ls duck.swf
0x00000015        4 FileAttributes (69)
0x0000001b      459 Metadata (77)
0x000001ec        4 ScriptLimits (65)
0x000001f2        3 SetBackgroundColor (9)
0x000001f7       26 TagCode(41)
0x00000213        5 FrameLabel (43)
0x0000021a      129 DefineBinaryData (87) [1]
0x000002a1      766 DefineBinaryData (87) [2]
0x000005a5     1145 DefineBinaryData (87) [3]
0x00000a24       26 DefineBitsLossless2 (36) [4]
0x00000a44      117 DefineBitsLossless2 (36) [5]
0x00000abf      233 DefineBitsLossless2 (36) [6]
0x00000bae      766 DefineBinaryData (87) [7]
0x00000eb2      292 DefineBinaryData (87) [8]
0x00000fdc       91 DefineBitsLossless2 (36) [9]
0x0000103d       38 DefineBitsLossless2 (36) [10]
0x00001069      175 DefineBitsLossless2 (36) [11]
0x0000111e      238 DefineBitsLossless2 (36) [12]
0x00001212     5169 TagCode(82)
0x00002649      248 SymbolClass (76)
0x00002747        0 ShowFrame (1)
0x00002749        0 End (0)
```

Use the `--symbols` flag to extract symbol information from SymbolClass tags. This will list each symbol's character ID, name and its associated tag.

```
$ swfx ls --symbols duck.swf
    0 duck (root class)
    1 duck_index DefineBinaryData (87)
    2 duck_duck_visualization DefineBinaryData (87)
    3 duck_duck_assets DefineBinaryData (87)
    4 duck_duck_32_sd_4_0 DefineBitsLossless2 (36)
    5 duck_duck_32_a_4_0 DefineBitsLossless2 (36)
    6 duck_duck_64_a_4_0 DefineBitsLossless2 (36)
    7 duck_manifest DefineBinaryData (87)
    8 duck_duck_logic DefineBinaryData (87)
    9 duck_duck_32_a_6_0 DefineBitsLossless2 (36)
   10 duck_duck_64_sd_4_0 DefineBitsLossless2 (36)
   11 duck_duck_64_a_6_0 DefineBitsLossless2 (36)
   12 duck_duck_icon_a DefineBitsLossless2 (36)
```

## extracting resources

Use the `extract` / `x` command to extract resources from the SWF file. \
This can currently extract binary data from `DefineBinaryData` tags and images from `DefineBitsLossless2` and `DefineBitsJpeg2` tags. \
Files will be named by their symbol defined in the `SymbolClass` tags, or by their character ID if no symbol exists. \
Note that multiple symbols can point to the same tag, in which case multiple files with the same content will be created. \
The file extension for binary data is chosen via MIME sniffing.

```
$ swfx x duck.swf
duck.swf
duck\duck_index.xml
duck\duck_duck_visualization.xml
duck\duck_duck_assets.xml
duck\duck_duck_32_sd_4_0.png
duck\duck_duck_32_a_4_0.png
duck\duck_duck_64_a_4_0.png
duck\duck_manifest.xml
duck\duck_duck_logic.xml
duck\duck_duck_32_a_6_0.png
duck\duck_duck_64_sd_4_0.png
duck\duck_duck_64_a_6_0.png
duck\duck_duck_icon_a.png

Extracted 12 files.
* DefineBinaryData: 5
* DefineBitsLossless2: 7
```