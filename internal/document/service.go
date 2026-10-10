package document

import (
	"fmt"
	"gourd/internal/counter"
	"gourd/internal/document/document_instance"
	"gourd/internal/document/rule"
	"gourd/internal/domain"
	"io"
	"strings"
)

type Service struct {
	counterService *counter.Service
}

func New(counterService *counter.Service) *Service {
	return &Service{counterService: counterService}
}

func getFileType(filename string) domain.FileType {
	parts := strings.Split(filename, ".")
	return domain.FileType(parts[len(parts)-1])
}

func NewDocument(filename string) (Document, error) {
	fileType := getFileType(filename)

	switch fileType {
	case domain.Text:
		return &document_instance.CommonTextDocument{EscapeMarkdown: false}, nil
	case domain.Markdown:
		return &document_instance.CommonTextDocument{EscapeMarkdown: true}, nil
	case domain.Docx:
		return &document_instance.DocxDocument{}, nil
	default:
		return nil, fmt.Errorf("unsupported file type: %v", fileType)
	}
}

func (s *Service) Process(in io.Reader, out io.Writer, ruleIDs []int, filename string) error {
	doc, err := NewDocument(filename)
	if err != nil {
		return err
	}
	rules := make([]rule.Rule, 0, len(ruleIDs))
	for _, id := range ruleIDs {
		r, ok := rule.QuickRegistry[id]
		if !ok {
			return fmt.Errorf("rule with id %d not found", id)
		}
		rules = append(rules, r)
	}

	if err = doc.Read(in); err != nil {
		return fmt.Errorf("read: %w", err)
	}

	if err = doc.Apply(rules); err != nil {
		return fmt.Errorf("apply: %w", err)
	}

	if err = doc.Write(out); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	s.counterService.Increment()

	return nil
}

func (s *Service) ProcessRaw(text string, ruleIDs []int, fileType domain.FileType) (string, error) {
	if fileType != domain.Text && fileType != domain.Markdown {
		return "", fmt.Errorf("unsupported file type: %v", fileType)
	}

	var buf strings.Builder
	err := s.Process(strings.NewReader(text), &buf, ruleIDs, string(fileType))
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
