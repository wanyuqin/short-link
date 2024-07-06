package netx

import (
	"fmt"
	"net"
)

// IPToInt 将IP地址转换为整数表示.
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

	IPInt := uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3])
	return IPInt, nil
}

// IntToIP 将整数表示的IP地址转换为字符串形式.
func IntToIP(IPInt uint32) string {
	IP := make(net.IP, 4)
	IP[0] = byte(IPInt >> 24 & 0xFF)
	IP[1] = byte(IPInt >> 16 & 0xFF)
	IP[2] = byte(IPInt >> 8 & 0xFF)
	IP[3] = byte(IPInt & 0xFF)
	return IP.String()
}
