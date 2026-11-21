package document_instance

import (
	"archive/zip"
	"bytes"
	"fmt"
	"gourd/internal/common"
	"gourd/internal/constants"
	"gourd/internal/document/result"
	"gourd/internal/document/rule"
	"io"

	"github.com/fumiama/go-docx"
)

type DocxDocument struct {
	document *docx.Docx
}

func (d *DocxDocument) Read(in io.Reader) error {
	limited := io.LimitReader(in, constants.MaxSingleFileSize)

	data, err := io.ReadAll(limited)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(data)
	zr, err := zip.NewReader(reader, int64(len(data)))
	if err != nil {
		return err
	}
	if err := validateZipEntries(zr.File); err != nil {
		return err
	}

	doc, err := docx.Parse(reader, int64(len(data)))
	if err != nil {
		return err
	}
	d.document = doc
	return nil
}

func (d *DocxDocument) Write(out io.Writer) error {
	_, err := d.document.WriteTo(out)
	return err
}

func (d *DocxDocument) Apply(rules []rule.Rule) ([]result.Result, error) {
	results := make([]result.Result, 0, len(rules))
	for _, r := range rules {
		res, err := r.Apply(d.document, common.Docx)
		if err != nil {
			return nil, err
		}
		if res != nil {
			results = append(results, result.Result{
				AnyResult: *res,
				RuleID:    r.ID,
			})
		}
	}
	return results, nil
}

func validateZipEntries(files []*zip.File) error {
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
