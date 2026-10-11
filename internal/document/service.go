package document

import (
	"context"
	"fmt"
	"gourd/internal/common"
	"gourd/internal/counter"
	"gourd/internal/document/document_instance"
	"gourd/internal/document/rule"
	"gourd/internal/workerpool"
	"io"
	"strings"
	"time"
)

type Service struct {
	counterService *counter.Service
	workerPool     *workerpool.Pool[*string]
}

func New(counterService *counter.Service) *Service {
	// TODO: put the worker pool into env
	return &Service{
		counterService: counterService,
		workerPool:     workerpool.New[*string](10),
	}
}

func (s *Service) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.workerPool.Stop(ctx)
}

func getFileType(filename string) common.FileType {
	parts := strings.Split(filename, ".")
	return common.FileType(parts[len(parts)-1])
}

func NewDocument(filename string) (Document, error) {
	fileType := getFileType(filename)

	switch fileType {
	case common.Text:
		return &document_instance.CommonTextDocument{EscapeMarkdown: false}, nil
	case common.Markdown:
		return &document_instance.CommonTextDocument{EscapeMarkdown: true}, nil
	case common.Docx:
		return &document_instance.DocxDocument{}, nil
	default:
		return nil, fmt.Errorf("unsupported file type: %v", fileType)
	}
}

func (s *Service) Process(ctx context.Context, in io.Reader, out io.Writer, ruleIDs []string, filename string) error {
	_, err := s.workerPool.Do(ctx, func(fnCtx context.Context) (*string, error) {
		doc, err := NewDocument(filename)
		if err != nil {
			return nil, err
		}
		rules := make([]rule.Rule, 0, len(ruleIDs))
		for _, id := range ruleIDs {
			r, ok := rule.QuickRegistry[id]
			if !ok {
				return nil, fmt.Errorf("rule with id %s not found", id)
			}
			rules = append(rules, r)
		}

		if err = doc.Read(in); err != nil {
			return nil, fmt.Errorf("read: %w", err)
		}

		if err = doc.Apply(rules); err != nil {
			return nil, fmt.Errorf("apply: %w", err)
		}

		if err = doc.Write(out); err != nil {
			return nil, fmt.Errorf("write: %w", err)
		}

		s.counterService.Increment()

		return nil, nil
	})

	return err
}

func (s *Service) ProcessRaw(ctx context.Context, text string, ruleIDs []string, fileType common.FileType) (string, error) {
	if fileType != common.Text && fileType != common.Markdown {
		return "", fmt.Errorf("unsupported file type: %v", fileType)
	}

	var buf strings.Builder
	err := s.Process(ctx, strings.NewReader(text), &buf, ruleIDs, string(fileType))
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
