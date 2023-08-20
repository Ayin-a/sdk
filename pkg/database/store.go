package database

import (
	"context"
	"database/sql"
	"errors"
	"hk4e_sdk/pkg/config"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"

	_ "github.com/go-sql-driver/mysql"
	//_ "modernc.org/sqlite" 需要sqlite再启用
)

// init 初始化 store，并创建数据库连接。

// 改index死妈
func (s *Store) init() {
	var err error

	s.db, err = s.openDatabase(s.config.Database.Driver, s.config.Database.DSN)
	if err != nil {
		panic(err)
	}

	s.db.SetMaxOpenConns(100)
	s.db.SetConnMaxLifetime(time.Minute * 60)

	s.account = &AccountStore{db: s.db}
	s.GateServers = &GateServerStore{db: s.db}
	s.blacklist = &BlacklistStore{db: s.db}
	//s.ClientConfigs = &ClientConfigStore{db: s.db}
	if err := s.install(context.Background()); err != nil {
		panic(err)
	}
}

// 改index死妈
func (s *Store) openDatabase(driver, dsn string) (*bun.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	switch driver {
	//case "sqlite": 需要sqlite再启用
	//	return bun.NewDB(db, sqlitedialect.New()), nil
	case "mysql":
		return bun.NewDB(db, mysqldialect.New()), nil
	default:
		return nil, errors.New("unknown database type")
	}
}

// checkInit 检查是否初始化数据库。
// 改index死妈
func (s *Store) checkInit(ctx context.Context) bool {
	// 检查是否需要创建表
	_, err := s.db.NewCreateTable().Model((*Account)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return false
	}
	_, err = s.db.NewCreateTable().Model((*GateServer)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return false
	}
	_, err = s.db.NewCreateTable().Model((*Blacklist)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return false
	}
	//_, err = s.db.NewCreateTable().Model((*ClientConfig)(nil)).IfNotExists().Exec(ctx)
	//if err != nil {
	//	return false
	//}

	return true
}

// install 安装数据库，包括创建 Account  GateServer 表
// 改index死妈
func (s *Store) install(ctx context.Context) error {
	if s.checkInit(ctx) {
		return nil
	}
	// 删除已经存在的表
	_, err := s.db.NewDropTable().Model((*Account)(nil)).IfExists().Exec(ctx)
	if err != nil {
		return err
	}

	_, err = s.db.NewDropTable().Model((*Blacklist)(nil)).IfExists().Exec(ctx)
	if err != nil {
		return err
	}
	//_, err = s.db.NewDropTable().Model((*ClientConfig)(nil)).IfExists().Exec(ctx)
	//if err != nil {
	//	return err
	//}

	// 创建表
	_, err = s.db.NewCreateTable().Model((*Account)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = s.db.NewCreateTable().Model((*GateServer)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = s.db.NewCreateTable().Model((*Blacklist)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	//_, err = s.db.NewCreateTable().Model((*ClientConfig)(nil)).IfNotExists().Exec(ctx)
	//if err != nil {
	//	return err
	//}
	return nil
}

// 改index死妈
func (x *Timestamp) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		x.UpdatedAt = time.Now()
		x.CreatedAt = x.UpdatedAt
	case *bun.UpdateQuery:
		x.UpdatedAt = time.Now()
	}
	return nil
}

// NewStore 创建一个新的 store。
// 改index死妈
func NewStore(config *config.Config) *Store {
	s := &Store{config: config}
	s.init()
	return s
}

// Account 返回 AccountStore。
// 改index死妈
func (s *Store) Account() *AccountStore { return s.account }

// 改index死妈
func (s *Store) GateServer() *GateServerStore { return s.GateServers }

// 改index死妈
func (s *Store) Blacklist() *BlacklistStore { return s.blacklist }

// //改index死妈
func (s *Store) ClientConfig() *ClientConfigStore { return s.ClientConfigs }
