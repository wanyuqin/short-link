package netx

import (
	"fmt"
	"net"
)

// IPToInt 将IP地址转换为整数表示
func IPToInt(ipStr string) (uint32, error) {
	if ipStr == "" {
		return 0, nil
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, fmt.Errorf("无效的IP地址：%s", ipStr)
	}

	// IPv4地址转换为整数
	ip = ip.To4()
	if ip == nil {
		return 0, fmt.Errorf("不是有效的IPv4地址：%s", ipStr)
	}

	ipInt := uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3])
	return ipInt, nil
}

// IntToIP 将整数表示的IP地址转换为字符串形式
func IntToIP(ipInt uint32) string {
	ip := make(net.IP, 4)
	ip[0] = byte(ipInt >> 24 & 0xFF)
	ip[1] = byte(ipInt >> 16 & 0xFF)
	ip[2] = byte(ipInt >> 8 & 0xFF)
	ip[3] = byte(ipInt & 0xFF)
	return ip.String()
}
