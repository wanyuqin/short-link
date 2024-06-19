package ctxkit

import "context"

func GetUserId(ctx context.Context) uint64 {
	value := ctx.Value("userId")
	if value == nil {
		return 0
	}
	if u, ok := value.(float64); ok {
		return uint64(u)
	}
	return 0
}
