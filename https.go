package main

import (
	"crypto/tls"
	"fmt"
	"time"
)

// 获取 HTTPS 证书信息
func getHttps(domain string, days_due int) (map[string]int, error) {
	var domain_status = make(map[string]int)

	// 建立 TCP 连接
	conn, err := tls.Dial("tcp", domain+":443", &tls.Config{
		InsecureSkipVerify: true, // 忽略证书验证
	})
	if err != nil {
		return domain_status, err
	}
	defer conn.Close()

	// 获取证书链
	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return domain_status, fmt.Errorf("找不到证书")
	}

	// 获取第一个证书（通常是服务器证书）
	cert := state.PeerCertificates[0]

	// 检查证书是否有效
	now := time.Now()
	if now.Before(cert.NotBefore) {
		return domain_status, fmt.Errorf("证书尚未有效")
	} else if now.After(cert.NotAfter) {
		return domain_status, fmt.Errorf("证书已过期")
	}

	// 计算剩余有效期
	remaining := cert.NotAfter.Sub(now)
	daysRemaining := int(remaining.Hours() / 24)
	if daysRemaining <= days_due {
		domain_status[domain] = daysRemaining
	}
	return domain_status, nil
}
