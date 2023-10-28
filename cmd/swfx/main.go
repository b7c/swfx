/*
Copyright © 2023 b7c
*/
package main

import (
	"github.com/b7c/swfx/cmd/swfx/cmd"
	_ "github.com/b7c/swfx/cmd/swfx/cmd/extract"
	_ "github.com/b7c/swfx/cmd/swfx/cmd/ls"
)

func main() {
	cmd.Execute()
}
