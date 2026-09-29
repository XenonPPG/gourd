package counter

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const storagePath = "./storage/counter.txt"

type Service struct {
	mu    sync.RWMutex
	value int64
}

func New() (*Service, error) {
	dir := filepath.Dir(storagePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("mkdir error: %w", err)
	}

	s := &Service{}

	data, err := os.ReadFile(storagePath)
	if os.IsNotExist(err) {
		if err := s.saveToFile(0); err != nil {
			return nil, fmt.Errorf("initial save error: %w", err)
		}
		return s, nil
	} else if err != nil {
		return nil, fmt.Errorf("read file error: %w", err)
	}

	valStr := strings.TrimSpace(string(data))
	num, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse count error: %w", err)
	}

	s.value = num
	return s, nil
}

func (s *Service) GetValue() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

func (s *Service) Increment() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.value++
}

func (s *Service) Save() error {
	s.mu.RLock()
	val := s.value
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveToFile(val)
}

func (s *Service) saveToFile(value int64) error {
	data := []byte(strconv.FormatInt(value, 10))
	return os.WriteFile(storagePath, data, 0644)
}
