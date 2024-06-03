package main

import (
	"b7c.io/swfx/cmd/swfx/cmd"
	_ "b7c.io/swfx/cmd/swfx/cmd/extract"
	_ "b7c.io/swfx/cmd/swfx/cmd/ls"
)

func main() {
	cmd.Execute()
}
