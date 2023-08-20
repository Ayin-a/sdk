package database

import (
	"github.com/uptrace/bun"
	"hk4e_sdk/pkg/config"
	"time"
)

// ***********************model*************************
// Timestamp 是用于自动添加 CreatedAt 和 UpdatedAt 时间戳的 struct。
type Timestamp struct {
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"` // 记录创建时间
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp"`          // 记录更新时间
	DeletedAt time.Time `bun:",soft_delete,nullzero"`                       // 记录软删除时间
}

// Account  数据库中的 accounts 表的结构
type Account struct {
	ID         int64  `bun:",pk,autoincrement"`
	Email      string `bun:",nullzero"`
	Username   string `bun:",nullzero"`
	Password   string `bun:",nullzero"`
	LoginToken string `bun:",nullzero"`
	ComboToken string `bun:",nullzero"`
	Device     string `bun:",nullzero"`
	IsGuest    bool   `bun:",nullzero"`
	IP         string `bun:",nullzero"`
	Timestamp
	TokenExpiration time.Time `bun:",nullzero"`
}

// GateServerConfig 数据库中的 gate_servers 表的结构
type GateServer struct {
	ID               int64  `bun:",pk,autoincrement" json:"ID"`
	Name             string `bun:",nullzero" json:"Name"`
	Title            string `bun:",nullzero" json:"Title"`
	Addr             string `bun:",nullzero" json:"Addr"`
	DispatchUrl      string `bun:",nullzero" json:"DispatchUrl"`
	ProxyDispatchUrl string `bun:",nullzero" json:"ProxyDispatchUrl"`
	Platform         string `bun:",nullzero" json:"Platform"`
	MuipUrl          string `bun:",nullzero" json:"MuipUrl"`
	PayCallbackUrl   string `bun:",nullzero" json:"PayCallbackUrl"`
	MuipSign         string `bun:",nullzero" json:"MuipSign"`
	PaySign          string `bun:",nullzero" json:"PaySign"`
}

// Store 是 hk4e 数据库的 store。
type Store struct {
	config      *config.Config
	db          *bun.DB
	account     *AccountStore
	GateServers *GateServerStore
	blacklist   *BlacklistStore
	//ClientConfigs *ClientConfigStore
}
type Blacklist struct {
	ID      int64  `bun:",pk,autoincrement"` // 主键，自增长
	IP      string `bun:",nullzero"`         // IP或IP段
	Comment string `bun:",nullzero"`         // 注释，可以用来描述为何被加入黑名单
}

/*
// ClientConfig 对应数据库中的 client_configs 表
type ClientConfig struct {
	ID                          int
	Retcode                     int32 //强更:20  弹窗:11
	Version                     string
	GateID                      string
	Msg                         string //顶部提示
	StopServerConfigStr         string //弹窗
	ForceUpdateConfigStr        string //强更
	ClientCustomConfigStr       string //list 的
	ClientRegionCustomConfigStr string //CUR 的region_custom_config_encrypted和region_custom_config_encrypted相似。
}*/
