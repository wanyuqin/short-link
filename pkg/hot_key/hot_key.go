package hot_key

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

//  client 作为计算传递来key 并且接受热key的消息，进行存储

// worker 用于进行阈值计算

type HotKey struct {
	store IHotKeyStore
}

var (
	once sync.Once
	tool *HotKey
)

func GetHotKeyTool() *HotKey {
	once.Do(func() {
		tool = &HotKey{
			store: NewHotKeyStore(),
		}
	})
	return tool
}

func (h *HotKey) Fetch(ctx context.Context, key string, res any, fetchFn func(ctx context.Context, key any) (any, error)) error {
	// 先从本地拿
	value := h.store.GetValue(key)
	if len(value) > 0 {
		return json.Unmarshal(value, res)
	}
	// 本地没有 上报 执行自带方法
	h.store.Push(key)
	ret, err := fetchFn(ctx, key)
	if err != nil {
		return err
	}
	err = h.handleFetchResult(ret, res)
	if err != nil {
		return err
	}
	// 进行本地缓存
	return h.store.SmartSet(key, ret)
}

// handleFetchResult 处理 fetchFn 返回的值并解码到 res 中
func (h *HotKey) handleFetchResult(ret any, res any) error {
	switch v := ret.(type) {
	case []byte:
		return json.Unmarshal(v, res)
	case string:
		return json.Unmarshal([]byte(v), res)
	case int, int8, int16, int32, int64:
		return json.Unmarshal([]byte(fmt.Sprintf("%d", v)), res)
	case uint, uint8, uint16, uint32, uint64:
		return json.Unmarshal([]byte(fmt.Sprintf("%d", v)), res)
	case float32, float64:
		return json.Unmarshal([]byte(fmt.Sprintf("%f", v)), res)
	case bool:
		return json.Unmarshal([]byte(fmt.Sprintf("%t", v)), res)
	case map[string]interface{}, []interface{}:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		return json.Unmarshal(jsonBytes, res)
	case struct{}, *struct{}:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		return json.Unmarshal(jsonBytes, res)
	default:
		return errors.New("unsupported return type from fetchFn")
	}
}
