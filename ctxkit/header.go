package ctxkit

import "context"

func GetUserId(ctx context.Context) uint64 {
	value := ctx.Value("userId")
	if value == nil {
		return 0
	}
	if u, ok := value.(uint64); ok {
		return u
	}
	return 0
}
