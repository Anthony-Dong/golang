package utils

import "time"

type StopWatch struct {
	start   time.Time // global start
	records map[string]time.Duration
	starts  map[string]time.Time
}

// NewStopWatch 创建一个秒表计时器
func NewStopWatch() StopWatch {
	return StopWatch{start: time.Now()}
}

// Record Record
func (s *StopWatch) Record(name string) {
	// 如果不存在则取全局的
	start, isExist := s.starts[name]
	if !isExist {
		start = s.start // 不存在使用全局的start
		if s.starts == nil {
			s.starts = map[string]time.Time{}
		}
		s.starts[name] = start
	}
	if s.records == nil {
		s.records = map[string]time.Duration{}
	}
	s.records[name] = time.Since(start)
}

func (s *StopWatch) Start(name string) {
	if s.starts == nil {
		s.starts = map[string]time.Time{}
	}
	s.starts[name] = time.Now()
}

func (s *StopWatch) GetStarts() map[string]time.Time {
	return s.starts
}

func (s *StopWatch) GetRecords() map[string]time.Duration {
	return s.records
}

func (s *StopWatch) GetStart() time.Time {
	return s.start
}
