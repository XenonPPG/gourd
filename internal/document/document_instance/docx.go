package document_instance

import (
	"archive/zip"
	"bytes"
	"gourd/internal/constants"
	"gourd/internal/document/rule"
	"gourd/internal/domain"
	"gourd/internal/utils"
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
	if err := utils.ValidateZipEntries(zr.File); err != nil {
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

func (d *DocxDocument) Apply(rules []rule.Rule) error {
	for _, r := range rules {
		if err := r.Apply(d.document, domain.Docx); err != nil {
			return err
		}
	}
	return nil
}
