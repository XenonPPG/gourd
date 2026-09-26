package utils

import (
	"archive/zip"
	"fmt"
	"gourd/internal/constants"
	"strings"
)

func EscapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]",
		"(", "\\(", ")", "\\)", "~", "\\~", "`", "\\`",
		">", "\\>", "#", "\\#", "+", "\\+", "-", "\\-",
		"=", "\\=", "|", "\\|", "{", "\\{", "}", "\\}",
		".", "\\.", "!", "\\!",
	)
	return replacer.Replace(text)
}

func ValidateZipEntries(files []*zip.File) error {
	var total uint64
	for _, f := range files {
		if f.UncompressedSize64 > constants.MaxSingleFileSize {
			return fmt.Errorf("file %q declares size %d, exceeds limit", f.Name, f.UncompressedSize64)
		}
		total += f.UncompressedSize64
		if total > constants.MaxTotalSize {
			return fmt.Errorf("total uncompressed size exceeds limit")
		}
		if f.CompressedSize64 > 0 {
			ratio := f.UncompressedSize64 / f.CompressedSize64
			if ratio > constants.MaxCompressionRatio {
				return fmt.Errorf("file %q has suspicious compression ratio %d", f.Name, ratio)
			}
		}
	}
	return nil
}
