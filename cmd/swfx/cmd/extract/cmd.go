package extract

import (
	"fmt"
	"image/png"
	"os"
	"path"
	"strconv"

	"github.com/gabriel-vasile/mimetype"

	"github.com/spf13/cobra"

	"github.com/b7c/swfx"
	"github.com/b7c/swfx/cmd/swfx/cmd"
)

var (
	outputDir     string
	quiet         bool
	extractAll    bool
	extractImages bool
	extractBinary bool

	extractCount int
)

var extractCmd = &cobra.Command{
	Use:     "extract",
	Aliases: []string{"x"},
	Short:   "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

		extractAll = !(extractBinary || extractImages)

		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "failed to create output directory: %s\n", err)
			return
		}

		for _, fileName := range args {
			f, err := os.Open(fileName)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "%s: failed to open file\n", fileName)
				continue
			}
			defer f.Close()

			swf, err := swfx.ReadSwf(f)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "%s: failed to read swf - %s\n", fileName, err)
				continue
			}

			for _, tag := range swf.Tags {
				switch tag := tag.(type) {
				case *swfx.DefineBitsLossless2:
					if extractAll || extractImages {
						extractBitsLossless2(swf, tag)
					}
				case *swfx.DefineBitsJpeg2:
					if extractAll || extractImages {
						extractBitsJpeg2(swf, tag)
					}
				case *swfx.DefineBinaryData:
					if extractAll || extractBinary {
						extractBinaryData(swf, tag)
					}
				}
			}
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(extractCmd)

	extractCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "The output directory.")
	extractCmd.Flags().BoolVarP(&extractAll, "all", "a", false, "Extract all resources.")
	extractCmd.Flags().BoolVarP(&extractImages, "images", "i", false, "Extract images.")
	extractCmd.Flags().BoolVarP(&extractBinary, "binary", "b", false, "Extract binary data.")
	extractCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Don't output filenames.")
}

func extractBitsLossless2(swf *swfx.Swf, tag *swfx.DefineBitsLossless2) error {
	var names []string
	var ok bool
	if names, ok = swf.ReverseSymbols[tag.CharacterId()]; !ok {
		names = []string{strconv.Itoa(int(tag.CharacterId()))}
	}

	for _, name := range names {
		outputFile := path.Join(outputDir, name+".png")
		img, err := tag.Decode()
		if err != nil {
			return fmt.Errorf("failed to decode tag (%d): %s",
				tag.CharacterId(), err)
		}
		f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("failed to open file: %q", outputFile)
		}
		defer f.Close()
		err = png.Encode(f, img)
		if err != nil {
			return fmt.Errorf("failed to encode image: %s", err)
		}
		extractCount++
		if !quiet {
			fmt.Println(outputFile)
		}
	}
	return nil
}

func extractBitsJpeg2(swf *swfx.Swf, tag *swfx.DefineBitsJpeg2) error {
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
		outputFile := path.Join(outputDir, name+ext)
		f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("failed to open file: %q", outputFile)
		}
		defer f.Close()
		n, err := f.Write(tag.ImageData)
		if n < len(tag.ImageData) {
			return fmt.Errorf("failed to write all data")
		}
		extractCount++
		if !quiet {
			fmt.Println(outputFile)
		}
	}
	return nil
}

func extractBinaryData(swf *swfx.Swf, tag *swfx.DefineBinaryData) error {
	var names []string
	var ok bool
	if names, ok = swf.ReverseSymbols[tag.CharacterId()]; !ok {
		names = []string{strconv.Itoa(int(tag.CharacterId()))}
	}

	mtype := mimetype.Detect(tag.Data)
	ext := mtype.Extension()

	for _, name := range names {
		outputFile := path.Join(outputDir, name+ext)

		f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("failed to open file: %q", outputFile)
		}
		defer f.Close()

		n, err := f.Write(tag.Data)
		if n < len(tag.Data) {
			return fmt.Errorf("failed to write all data")
		}
		extractCount++
		if !quiet {
			fmt.Println(outputFile)
		}
	}
	return nil
}
