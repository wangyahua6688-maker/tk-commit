package reqx

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// ParseIntOrDefault 解析正整数；非法值时回退默认值。
func ParseIntOrDefault(raw string, fallback int) int {
	v, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}

// ParsePathID 从路径里按 prefix 读取下一个段作为正整数 ID。
func ParsePathID(path string, prefix string) (uint64, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for idx := range parts {
		if parts[idx] == prefix && idx+1 < len(parts) {
			id, err := strconv.ParseUint(parts[idx+1], 10, 64)
			if err != nil || id == 0 {
				return 0, fmt.Errorf("invalid id")
			}
			return id, nil
		}
	}
	return 0, fmt.Errorf("invalid id")
}

// DeviceID 从请求中提取设备标识：优先 Header，再回退 Query。
func DeviceID(r *http.Request) string {
	deviceID := strings.TrimSpace(r.Header.Get("X-Device-ID"))
	if deviceID != "" {
		return deviceID
	}
	return strings.TrimSpace(r.URL.Query().Get("device_id"))
}

// ClientIP 从代理头中提取客户端 IP，最后回退 RemoteAddr。
func ClientIP(r *http.Request) string {
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 && strings.TrimSpace(parts[0]) != "" {
			return strings.TrimSpace(parts[0])
		}
	}
	if rip := strings.TrimSpace(r.Header.Get("X-Real-IP")); rip != "" {
		return rip
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}
