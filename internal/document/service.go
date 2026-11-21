package document

import (
	"context"
	"fmt"
	"gourd/internal/common"
	"gourd/internal/counter"
	"gourd/internal/document/document_instance"
	"gourd/internal/document/result"
	"gourd/internal/document/rule"
	"gourd/internal/workerpool"
	"io"
	"strings"
	"time"
)

type Service struct {
	counterService *counter.Service
	workerPool     *workerpool.Pool[[]result.Result]
}

func New(workerCount int, counterService *counter.Service) *Service {
	return &Service{
		counterService: counterService,
		workerPool:     workerpool.New[[]result.Result](workerCount),
	}
}

func (s *Service) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.workerPool.Stop(ctx)
}

func NewDocument(fileType common.FileType) (Document, error) {
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

func (s *Service) Process(ctx context.Context, in io.Reader, out io.Writer, fileType common.FileType, ruleIDs []string) ([]result.Result, error) {
	results, err := s.workerPool.Do(ctx, func(fnCtx context.Context) ([]result.Result, error) {
		doc, err := NewDocument(fileType)
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

		res, err := doc.Apply(rules)
		if err != nil {
			return nil, fmt.Errorf("apply: %w", err)
		}

		if err = doc.Write(out); err != nil {
			return nil, fmt.Errorf("write: %w", err)
		}

		s.counterService.Increment()

		return res, nil
	})

	return results, err
}

func (s *Service) ProcessRaw(ctx context.Context, text string, ruleIDs []string, fileType common.FileType) (string, []result.Result, error) {
	if fileType != common.Text && fileType != common.Markdown {
		return "", nil, fmt.Errorf("unsupported file type: %v", fileType)
	}

	var buf strings.Builder
	results, err := s.Process(ctx, strings.NewReader(text), &buf, fileType, ruleIDs)
	if err != nil {
		return "", nil, err
	}

	return buf.String(), results, nil
}
