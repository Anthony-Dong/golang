package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

type SplitFile struct {
	file     *os.File
	filename string

	start int
	end   int
	index int
}

func NewSplitFile(filename string, start, end int) *SplitFile {
	return &SplitFile{
		filename: filename,
		start:    start,
		end:      end,
	}
}

func (s *SplitFile) Index() (int, int) {
	return s.start, s.end
}

func (s *SplitFile) Size() int64 {
	return int64(s.end - s.start)
}

func (s *SplitFile) FileName() string {
	return s.filename
}

func (s *SplitFile) Init() error {
	file, err := os.OpenFile(s.filename, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	if _, err := file.Seek(int64(s.start), io.SeekStart); err != nil {
		file.Close()
		return err
	}
	s.file = file
	s.index = s.start
	return nil
}

func (s *SplitFile) Read(p []byte) (n int, err error) {
	size := s.malloc(p)
	if size <= 0 {
		return 0, io.EOF
	}
	n, err = s.file.Read(p[:size])
	s.update(n)
	return
}

func (s *SplitFile) malloc(p []byte) int {
	size := len(p)
	if s.index+size > s.end {
		size = s.end - s.index
	}
	return size
}

func (s *SplitFile) update(size int) {
	if size <= 0 {
		return
	}
	s.index = s.index + size
}

var ErrLimitExceeded = errors.New("write limit exceeded")

func (s *SplitFile) Write(p []byte) (n int, err error) {
	size := s.malloc(p)
	if size <= 0 {
		return 0, ErrLimitExceeded
	}
	n, err = s.file.Write(p[:size])
	s.update(n)
	return
}

func (s *SplitFile) Close() error {
	if s.file == nil {
		return nil
	}
	return s.file.Close()
}

func NewMultiSplitFile(filename string, split int) ([]*SplitFile, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	result := make([]*SplitFile, 0)
	if err := Split(int(stat.Size()), split, func(start, end int) error {
		result = append(result, &SplitFile{
			filename: filename,
			start:    start,
			end:      end,
		})
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func Split(size int, split int, handle func(start, end int) error) error {
	if handle == nil {
		return fmt.Errorf(`split handle is nil`)
	}
	start := 0
	end := 0
	for {
		if start >= size {
			return nil
		}
		end = start + split
		if end >= size {
			end = size
		}
		if err := handle(start, end); err != nil {
			return err
		}
		start = end
	}
}
