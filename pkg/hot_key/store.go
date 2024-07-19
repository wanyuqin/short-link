package hot_key

import (
	"encoding/json"
	"time"

	"github.com/go-redis/cache/v9"
)

type IHotKeyStore interface {
	// IsHotKey 判断是否是热key
	IsHotKey(key string) bool
	// SmartSet 给热key赋值value
	SmartSet(key string, value any) error
	// GetValue 集成热key判断。并且返回对应key的值
	GetValue(key string) []byte
	// Push 上报
	Push(key string)
}

type HotKeyStore struct {
	localCache           *cache.TinyLFU
	slidingWindowManager *SlidingWindowManager
}

func (h HotKeyStore) IsHotKey(key string) bool {
	return h.slidingWindowManager.ExceedsThreshold(key)
}

func (h HotKeyStore) Push(key string) {
	h.slidingWindowManager.AddAccess(key, time.Now().Unix())
}
func (h HotKeyStore) SmartSet(key string, value any) error {
	if !h.IsHotKey(key) {
		return nil
	}
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	h.localCache.Set(key, b)
	return nil
}

func (h HotKeyStore) GetValue(key string) []byte {
	if !h.IsHotKey(key) {
		return nil
	}
	res, _ := h.localCache.Get(key)
	return res
}

func NewHotKeyStore() *HotKeyStore {
	return &HotKeyStore{
		localCache:           cache.NewTinyLFU(DefaultSize, DefaultDuration),
		slidingWindowManager: NewSlidingWindowManager(),
	}
}
