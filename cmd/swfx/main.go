/*
Copyright © 2023 b7c

*/
package main

import (
	"github.com/b7c/swfx/cmd/swfx/cmd"
	_ "github.com/b7c/swfx/cmd/swfx/cmd/ls"
	_ "github.com/b7c/swfx/cmd/swfx/cmd/extract"
)

func main() {
	cmd.Execute()
}
