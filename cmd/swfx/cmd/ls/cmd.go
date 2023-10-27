package ls

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/b7c/swfx"
	"github.com/b7c/swfx/cmd/swfx/cmd"
	"github.com/b7c/swfx/tagcode"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List tags from a SWF file.",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Fprintf(cmd.ErrOrStderr(), "no file specified")
			return
		}

		f, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "failed to open input file")
			return
		}
		defer f.Close()

		reader := swfx.NewReader(f)
		_, err = swfx.ReadHeader(reader)
		if err != nil {
			panic(err)
		}

		for {
			offset := reader.Position()
			tagCode, length := reader.ReadTagCodeAndLength()
			tagName := tagCode.String()
			if strings.ContainsRune(tagName, '(') {
				fmt.Printf("0x%08x %8d %s\n", offset, length, tagCode)
			} else {
				fmt.Printf("0x%08x %8d %s (%d)\n", offset, length, tagCode, tagCode)
			}
			reader.ReadTag(tagCode, length)
			if tagCode == tagcode.End {
				break
			}
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(lsCmd)
}
