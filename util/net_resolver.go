package util

import (
	"context"
	"net"
)

// NewDialerResolver 使用 s 将 dialer.Resolver 设置为新的 net.Resolver。
//
// s：用于创建新 net.Resolver 的字符串。
func NewDialerResolver(s string) {
	dialer.Resolver = newNetResolver(s)
}

// newNetResolver 当 s 不为空时返回使用 s 的 Go 内置 DNS 解析器。
//
// s：net.Resolver 的 DNS 服务器地址。
func newNetResolver(s string) *net.Resolver {
	if s == "" {
		return net.DefaultResolver
	}

	return &net.Resolver{
		//两者都配置的情况下，会优先使用Dial里面的配置
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			//注意这里的实现
			return net.Dial("udp", s)
		},
	}
}
