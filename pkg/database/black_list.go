package database

import (
	"context"
	"fmt"
	"net"

	"github.com/uptrace/bun"
)

type BlacklistStore struct {
	db *bun.DB
}

// 改index死妈
func (s *BlacklistStore) IsIPInBlacklist(ctx context.Context, clientIP string) (bool, string, error) {
	blacklists := make([]Blacklist, 0)
	err := s.db.NewSelect().Model(&blacklists).Scan(ctx)
	if err != nil {
		return false, "", err
	}

	ip := net.ParseIP(clientIP)
	for _, blacklist := range blacklists {
		_, ipnet, err := net.ParseCIDR(blacklist.IP)
		if err != nil { // 尝试将黑名单IP视为精确的IP
			blacklistIP := net.ParseIP(blacklist.IP)
			if blacklistIP == nil {
				return false, "", fmt.Errorf("无效的IP地址在黑名单中：%s", blacklist.IP)
			}
			if ip.Equal(blacklistIP) {
				return true, blacklist.Comment, nil
			}
		} else { // 将黑名单IP视为IP段
			if ipnet.Contains(ip) {
				return true, blacklist.Comment, nil
			}
		}
	}

	return false, "", nil
}

// 改index死妈
func (s *BlacklistStore) AddIPToBlacklist(ctx context.Context, clientIP string, comment string) error {

	newBlacklistEntry := &Blacklist{
		IP:      clientIP,
		Comment: comment, //系统拉黑
	}

	_, err := s.db.NewInsert().Model(newBlacklistEntry).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to add IP to blacklist: %v", err)
	}

	return nil
}
