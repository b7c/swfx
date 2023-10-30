package extract

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gabriel-vasile/mimetype"

	"github.com/spf13/cobra"

	"github.com/b7c/swfx"
	"github.com/b7c/swfx/cmd/swfx/cmd"
	"github.com/b7c/swfx/tagcode"
)

var (
	baseOutDir    string
	quiet         bool
	extractAll    bool
	extractImages bool
	extractBinary bool

	extractCount int
)

var extractCounts = make(map[tagcode.TagCode]int)

var extractCmd = &cobra.Command{
	Use:     "extract",
	Aliases: []string{"x"},
	Short:   "Extracts resources from tags.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		extractAll = !(extractBinary || extractImages)

		for _, fileName := range args {
			cmd.Println(fileName)
			err := extractSwf(cmd, fileName)
			if err != nil {
				cmd.Println(err)
			}
		}
		
		cmd.Printf("\nExtracted %d file%s.\n", extractCount, pluralize(extractCount))
		for code, count := range extractCounts {
			cmd.Printf("* %s: %d\n", code, count)
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(extractCmd)

	extractCmd.Flags().StringVarP(&baseOutDir, "output", "o", "", "The output directory. Creates a directory with the same name as the SWF file without the extension by default.")
	extractCmd.Flags().BoolVarP(&extractAll, "all", "a", false, "Extract all resources.")
	extractCmd.Flags().BoolVarP(&extractImages, "images", "i", false, "Extract images.")
	extractCmd.Flags().BoolVarP(&extractBinary, "binary", "b", false, "Extract binary data.")
	extractCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Don't output filenames.")
}

func pluralize(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func extractSwf(cmd *cobra.Command, fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	
	if stat.IsDir() {
		return fmt.Errorf("input must be a file")
	}

	outDir := baseOutDir
	if outDir == "" {
		index := strings.LastIndex(fileName, ".")
		if index != -1 {
			outDir = fileName[:index]
		} else {
			outDir = fileName + "_extracted"
		}
	}

	err = os.MkdirAll(outDir, 0755)
	if err != nil {
		return err
	}

	swf, err := swfx.ReadSwf(f)
	if err != nil {
		return err
	}

	for _, tag := range swf.Tags {
		switch tag := tag.(type) {
		case *swfx.DefineBitsLossless2:
			if extractAll || extractImages {
				err := extractBitsLossless2(outDir, swf, tag)
				if err != nil {
					cmd.PrintErrf("failed to extract image from character #%d: %s", tag.CharacterId(), err)
				}
			}
		case *swfx.DefineBitsJpeg2:
			if extractAll || extractImages {
				err := extractBitsJpeg2(outDir, swf, tag)
				if err != nil {
					cmd.PrintErrf("failed to extract image from character #%d: %s", tag.CharacterId(), err)
				}
			}
		case *swfx.DefineBinaryData:
			if extractAll || extractBinary {
				err := extractBinaryData(outDir, swf, tag)
				if err != nil {
					cmd.PrintErrf("failed to extract image from character #%d: %s", tag.CharacterId(), err)
				}
			}
		}
	}
	return nil
}

func extractBitsLossless2(outDir string, swf *swfx.Swf, tag *swfx.DefineBitsLossless2) error {
	var names []string
	var ok bool
	if names, ok = swf.ReverseSymbols[tag.CharacterId()]; !ok {
		names = []string{strconv.Itoa(int(tag.CharacterId()))}
	}

	for _, name := range names {
		outputFile := filepath.Join(outDir, name+".png")
		img, err := tag.Decode()
		if err != nil {
			return err
		}
		f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer f.Close()
		err = png.Encode(f, img)
		if err != nil {
			return err
		}
		extractCount++
		extractCounts[tag.Code()]++
		if !quiet {
			fmt.Println(outputFile)
		}
	}
	return nil
}

func extractBitsJpeg2(outDir string, swf *swfx.Swf, tag *swfx.DefineBitsJpeg2) error {
	var names []string
	var ok bool
	if names, ok = swf.ReverseSymbols[tag.CharacterId()]; !ok {
		names = []string{fmt.Sprintf("%d", tag.CharacterId())}
	}

	var ext string
	switch tag.ImageType() {
	case swfx.Jpeg:
		ext = ".jpg"
	case swfx.Png:
		ext = ".png"
	case swfx.Gif:
		ext = ".gif"
	default:
		return fmt.Errorf("unknown image type")
	}

	for _, name := range names {
		outputFile := filepath.Join(outDir, name+ext)
		f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(tag.ImageData)
		if err != nil {
			return err
		}
		extractCount++
		extractCounts[tag.Code()]++
		if !quiet {
			fmt.Println(outputFile)
		}
	}
	return nil
}

func extractBinaryData(outDir string, swf *swfx.Swf, tag *swfx.DefineBinaryData) error {
	var names []string
	var ok bool
	if names, ok = swf.ReverseSymbols[tag.CharacterId()]; !ok {
		names = []string{strconv.Itoa(int(tag.CharacterId()))}
	}

	mtype := mimetype.Detect(tag.Data)
	ext := mtype.Extension()

	for _, name := range names {
		outputFile := filepath.Join(outDir, name+ext)

		f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("failed to open file: %q", outputFile)
		}
		defer f.Close()

		_, err = f.Write(tag.Data)
		if err != nil {
			return err
		}
		extractCount++
		extractCounts[tag.Code()]++
		if !quiet {
			fmt.Println(outputFile)
		}
	}
	return nil
}
