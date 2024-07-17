package hot_key

import (
	"short-link/config"
	"sync"
)

type SlidingWindow struct {
	timestamps []int64
	mutex      sync.Mutex
}

type SlidingWindowManager struct {
	windows map[string]*SlidingWindow
	mutex   sync.Mutex
}

func NewSlidingWindowManager() *SlidingWindowManager {
	return &SlidingWindowManager{
		windows: make(map[string]*SlidingWindow),
	}
}

func (m *SlidingWindowManager) AddAccess(key string, timestamp int64) {
	m.mutex.Lock()
	window, ok := m.windows[key]
	if !ok {
		window = &SlidingWindow{}
		m.windows[key] = window
	}
	m.mutex.Unlock()

	window.mutex.Lock()
	defer window.mutex.Unlock()

	window.timestamps = append(window.timestamps, timestamp)
	cutoff := timestamp - int64(config.GetConfig().HotKey.Interval)
	for len(window.timestamps) > 0 && window.timestamps[0] < cutoff {
		window.timestamps = window.timestamps[1:]
	}
}

func (m *SlidingWindowManager) ExceedsThreshold(key string) bool {
	m.mutex.Lock()
	window, exists := m.windows[key]
	m.mutex.Unlock()

	if !exists {
		return false
	}

	window.mutex.Lock()
	defer window.mutex.Unlock()

	return len(window.timestamps) > config.GetConfig().HotKey.Threshold
}
