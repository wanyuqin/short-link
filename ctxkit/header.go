package ctxkit

import "context"

func GetUserId(ctx context.Context) uint64 {
	value, ok := ctx.Value("userId").(float64)
	if !ok {
		return 0
	}
	return uint64(value)
}

func GetIp(ctx context.Context) string {
	ip, ok := ctx.Value("ip").(string)
	if !ok {
		return ""
	}
	return ip
}
