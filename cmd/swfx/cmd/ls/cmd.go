package ls

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/b7c/swfx"
	"github.com/b7c/swfx/cmd/swfx/cmd"
	"github.com/b7c/swfx/tagcode"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List items in a SWF file",
	RunE:  runLs,
}

var (
	listSymbols bool
)

func init() {
	cmd.RootCmd.AddCommand(lsCmd)

	lsCmd.Flags().BoolVarP(&listSymbols, "symbols", "s", false, "List symbols")
}

type symbolItem struct {
	Name string
	Id   swfx.CharacterId
}

func runLs(cmd *cobra.Command, args []string) (err error) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	cmd.SilenceUsage = true

	f, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "failed to open input file")
		return err
	}
	defer f.Close()

	if listSymbols {
		var swf *swfx.Swf
		swf, err = swfx.ReadSwf(f)
		if err != nil {
			return
		}

		symbols := []symbolItem{}
		for symbol, ch := range swf.Symbols {
			symbols = append(symbols, symbolItem{symbol, ch})
		}

		slices.SortFunc(symbols, func(a, b symbolItem) int {
			return int(a.Id) - int(b.Id)
		})

		for _, item := range symbols {
			symbol, ch := item.Name, item.Id
			tag, ok := swf.Characters[ch]
			if !ok {
				if ch == 0 {
					fmt.Printf("%5d %s (root class)\n", ch, symbol)
				} else {
					fmt.Printf("%5d %s: not found\n", ch, symbol)
				}
				continue
			}
			tagName := tag.Code().String()
			if strings.ContainsRune(tagName, '(') {
				fmt.Printf("%5d %s %s\n", ch, symbol, tagName)
			} else {
				fmt.Printf("%5d %s %s (%d)\n", ch, symbol, tagName, tag.Code())
			}
		}

	} else {

		reader := swfx.NewReader(f)
		_, err = swfx.ReadHeader(reader)
		if err != nil {
			return err
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
	}

	return
}
