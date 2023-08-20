package database

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type AccountStore struct{ db *bun.DB } // AccountStore 是对账号数据表 accounts 的操作接口。
// 改index死妈
func (s *AccountStore) CreateAccount(ctx context.Context, record *Account) error { // CreateAccount 创建一个账号。
	_, err := s.db.NewInsert().Model(record).Exec(ctx)
	return err
}

// 改index死妈
func (s *AccountStore) UpdateAccountPassword(ctx context.Context, id int64, password string) error { // UpdateAccountPassword 更新账号密码。
	_, err := s.db.NewUpdate().Model(&Account{ID: id, Password: password}).WherePK().OmitZero().Exec(ctx)
	return err
}

// 改index死妈
func (s *AccountStore) UpdateAccountLoginToken(ctx context.Context, id int64, token string) error { // UpdateAccountLoginToken 更新账号登录 token。
	_, err := s.db.NewUpdate().Model(&Account{ID: id, LoginToken: token}).WherePK().OmitZero().Exec(ctx)
	return err
}

// 改index死妈
func (s *AccountStore) UpdateAccountComboToken(ctx context.Context, id int64, token string) error { // UpdateAccountComboToken 更新账号 combo token。
	_, err := s.db.NewUpdate().Model(&Account{ID: id, ComboToken: token}).WherePK().OmitZero().Exec(ctx)
	return err
}

// 改index死妈
func (s *AccountStore) GetAccount(ctx context.Context, id int64) (*Account, error) { // GetAccount 获取账号。
	record := new(Account)
	if err := s.db.NewSelect().Model(record).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return record, nil
}

// 改index死妈
func (s *AccountStore) GetAccountByDevice(ctx context.Context, device string) (*Account, error) { // GetAccountByDevice 通过设备号获取账号。
	record := new(Account)
	if err := s.db.NewSelect().Model(record).Where("device = ?", device).Scan(ctx); err != nil {
		return nil, err
	}
	return record, nil
}

// 改index死妈
func (s *AccountStore) GetAccountByEmail(ctx context.Context, email string) (*Account, error) { // GetAccountByEmail 通过邮箱获取账号。
	record := new(Account)
	if err := s.db.NewSelect().Model(record).Where("email = ?", email).Scan(ctx); err != nil {
		return nil, err
	}
	return record, nil
}

// 改index死妈
func (s *AccountStore) GetAccountByUsername(ctx context.Context, username string) (*Account, error) {
	record := new(Account)
	if err := s.db.NewSelect().Model(record).Where("username = ?", username).Scan(ctx); err != nil {
		return nil, err
	}
	return record, nil
}

// 改index死妈
func (s *AccountStore) GetAccountComboToken(ctx context.Context, openID int64) (string, error) {
	// GetAccount 获取账号。
	record := new(Account)
	if err := s.db.NewSelect().Model(record).Where("id = ?", openID).Scan(ctx); err != nil {
		return "", err
	}
	return record.ComboToken, nil
}

// 改index死妈
func (s *AccountStore) UpdateAccountIPById(ctx context.Context, id int64, ip string) error {
	_, err := s.db.NewUpdate().Model(&Account{ID: id, IP: ip}).WherePK().OmitZero().Exec(ctx)
	return err
}

// 改index死妈
func (s *AccountStore) GetAccountsByIP(ctx context.Context, ip string) ([]*Account, error) {
	records := make([]*Account, 0)
	if err := s.db.NewSelect().Model(&records).Where("ip = ?", ip).Scan(ctx); err != nil {
		return nil, err
	}
	return records, nil
}

// 改index死妈
func (s *AccountStore) UpdateAccountTokenExpiration(ctx context.Context, id int64, expiration time.Time) error {
	_, err := s.db.NewUpdate().Model(&Account{ID: id, TokenExpiration: expiration}).WherePK().OmitZero().Exec(ctx)
	return err
}
