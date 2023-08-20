package http

import (
	"crypto/rsa"
	"errors"
	"hk4e_sdk/pkg/config"
	"hk4e_sdk/pkg/database"
	"hk4e_sdk/pkg/logger"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// ************************sdk***********************************
type ID int64

// risky
const RISKY_ACTION_NONE = "ACTION_NONE"

// account
//const ACCOUNT_NORMAL = 0
//const ACCOUNT_GUEST = 1

// channel
//const CHANNEL_ID_MIHOYO = 1
//const CHANNEL_ID_BILIBILI = 14

// scenes
// const SCENE_NORMAL = "S_NORMAL"   // mobile + account; default to mobile
const SCENE_ACCOUNT = "S_ACCOUNT" // mobile + account; default to account
//const SCENE_USER = "S_USER"       // account only
//const SCENE_TEMPLE = "S_TEMPLE"   // account only; no registration

// platform
var PLATFORM_TYPE_STR = map[int32]string{
	0:  "EDITOR",
	1:  "IOS",
	2:  "ANDROID",
	3:  "PC",
	4:  "PS4",
	5:  "SERVER",
	6:  "CLOUD_ANDROID",
	7:  "CLOUD_IOS",
	8:  "PS5",
	9:  "CLOUD_WEB",
	10: "CLOUD_TV",
	11: "CLOUD_MAC",
	12: "CLOUD_PC",
	13: "CLOUD_THIRD_PARTY_MOBILE",
	14: "CLOUD_THIRD_PARTY_PC",
}

var (
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidLoginToken = errors.New("invalid login token")
	//ErrInvalidComboToken = errors.New("invalid combo token")
)

type Response struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
	Retcode int32  `json:"retcode"`
}

var ResponseMessage = map[int32]string{
	0:    "OK",
	-101: "系统错误",
	-102: "密码格式错误，密码格式为8-30位，并且由数字、大小写字母、英文特殊符号两种以上组合",
	-103: "参数错误",
	-104: "缺少配置",
	-106: "协议加载失败",
	-107: "渠道错误",
	-111: "???",
	-115: "请前往官网/商店下载最新版本",
	-202: "账号或密码错误",
	-203: "用户名异常",
	-204: "登录异常",
	-205: "用户名为空",
	-206: "您的ip或ip段在黑名单",
	-207: "账号内含有非法字符",
	-208: "检查ip失败",
	-210: "为了您的账号安全，请重新登录。",
	-301: "SDK数据库读取获取用户信息失败",
	-701: "GM请求参数不能为空",
	-702: "错误的请求",
	-703: "GM分区MUIP未配置或不存在",
	-713: "PAY分区OA不存在或未配置",

	-2001: "创建游客账号失败",
	-2002: "创建账号用户失败",
	-2003: "修改密码失败",
	-2100: "创建hash失败",
	-3000: "验证码邮件发送错误",
	-3001: "邮箱验证码错误",
	-3101: "未知的ACTION",
}

// 改index死妈
func NewResponse(retcode int32, data any) *Response {
	return &Response{
		Data:    data,
		Message: ResponseMessage[retcode],
		Retcode: retcode,
	}
}

type SimpleResponse struct {
	Message string `json:"message"`
	RetCode int32  `json:"retcode"`
}

// 改index死妈
func NewSimpleResponse(retCode int32, message string) *SimpleResponse {
	return &SimpleResponse{
		Message: message,
		RetCode: retCode,
	}
}

type GuestLoginRequestData struct {
	// "client": 3,
	// "device": "9e50ff97dedbfbaba48fa7958e9a5aca90af2caf1661473377723",
	// "g_version": "OSRELWin3.2.0",
	// "game_key": "hk4e_global",
	// "sign": "ab7388665dfcb88e5b8db742028550991f4fca9214c0dc801aa8f99ad9f85267"

	Client   int32  `json:"client"`
	Device   string `json:"device"`
	GVersion string `json:"g_version"`
	GameKey  string `json:"game_key"`
	Sign     string `json:"sign"`
}
type ShopwindowListReq struct {
	Currency string `json:"currency"`
}

//	type H5Log struct {
//		Data string `json:"data"`
//	}
type LoadConfigReq struct {
	Client  int    `json:"client"`
	GameKey string `json:"game_key"`
}
type CompareProtocolVersionRequestData struct {
	ID        string `json:"id"`
	AppID     string `json:"app_id"`
	ChannelID string `json:"channel_id"`
	Language  string `json:"language"`
	Major     string `json:"major"`
	Minimum   string `json:"minimum"`
}
type ProtocolVersion struct {
	ID            int64  `json:"id"`
	AppID         int64  `json:"app_id"`
	Language      string `json:"language"`
	UserProto     string `json:"user_proto"`
	PrivProto     string `json:"priv_proto"`
	Major         int32  `json:"major"`
	Minimum       int32  `json:"minimum"`
	CreateTime    string `json:"create_time"`
	TeenagerProto string `json:"teenager_proto"`
	ThirdProto    string `json:"third_proto"`
}
type CompareProtocolVersionResponseData struct {
	Modified bool             `json:"modified"`
	Protocol *ProtocolVersion `json:"protocol"`
}
type InnerAccountVerifyRequestData struct {
	AppID      uint32 `json:"app_id"`
	ChannelID  uint32 `json:"channel_id"`
	OpenID     string `json:"open_id"`
	ComboToken string `json:"combo_token"`
	Sign       string `json:"sign"`
	Region     string `json:"region"`
}

//	type TokenCheckResponseData struct {
//		AccountType int32 `json:"account_type"`
//		IPInfo      struct {
//			CountryCode string `json:"country_code"`
//		} `json:"ip_info"`
//	}
type AccountBaseData struct {
	Action          string `json:"action"`
	Account         string `json:"account"`
	Password        string `json:"password"`
	EMail           string `json:"email"`
	EMailVerifyCode string `json:"email_verify_code"`
	IsCrypto        bool   `json:"is_crypto"`
}
type MobileLoginData struct {
	// {
	//     "action": "Login",
	//     "area": "+86",
	//     "captcha": "652202",
	//     "mobile": "13764100025"
	// }
	Action  string `json:"action"`
	Area    string `json:"area"`
	Captcha string `json:"captcha"`
	Mobile  string `json:"mobile"`
}
type ShieldLoginRequestData struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	IsCrypto bool   `json:"is_crypto"`
}
type GameLoginLogoutNotify struct {
	Uid         int64  `json:"uid"`
	AccountType int    `json:"account_type"`
	AccountId   string `json:"account"`
	Platform    int    `json:"platform"`
	Region      string `json:"region"`
	BizGame     string `json:"biz_game"`
}
type Account struct {
	UID               ID     `json:"uid"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	Mobile            string `json:"mobile"`
	IsEmailVerify     string `json:"is_email_verify"`
	RealName          string `json:"realname"`
	IdentityCard      string `json:"identity_card"`
	Token             string `json:"token"`
	SafeMobile        string `json:"safe_mobile"`
	FacebookName      string `json:"facebook_name"`
	GoogleName        string `json:"google_name"`
	TwitterName       string `json:"twitter_name"`
	GameCenterName    string `json:"game_center_name"`
	AppleName         string `json:"apple_name"`
	SonyName          string `json:"sony_name"`
	TapName           string `json:"tap_name"`
	Country           string `json:"country"`
	ReactivateTicket  string `json:"reactivate_ticket"`
	AreaCode          string `json:"area_code"`
	DeviceGrantTicket string `json:"device_grant_ticket"`
	SteamName         string `json:"steam_name"`
}
type ShieldLoginResponseData struct {
	Account             *Account `json:"account"`
	DeviceGrantRequired bool     `json:"device_grant_required"`
	SafeMobileRequired  bool     `json:"safe_mobile_required"`
	RealPersonRequired  bool     `json:"realperson_required"`
	ReactivateRequired  bool     `json:"reactivate_required"`
	RealNameOperation   string   `json:"realname_operation"`
}

// const GuestIdStart = 100000000
type ShieldVerifyRequestData struct {
	UID   ID     `json:"uid"`
	Token string `json:"token"`
}
type ThirdpartyLoginRequestData struct {
	AccessToken string `json:"access_token"`
	CbUrl       string `json:"cb_url"`
	NoRegist    bool   `json:"no_regist"`
	Thirdparty  string `json:"thirdparty"`
}
type ShieldVerifyResponseData struct {
	Account             *Account `json:"account"`
	DeviceGrantRequired bool     `json:"device_grant_required"`
	SafeMobileRequired  bool     `json:"safe_mobile_required"`
	RealPersonRequired  bool     `json:"realperson_required"`
	RealNameOperation   string   `json:"realname_operation"`
}
type ComboLoginRequestData struct {
	AppID     any    `json:"app_id"`
	ChannelID any    `json:"channel_id"`
	Data      string `json:"data"`
	Device    string `json:"device"`
	Sign      string `json:"sign"`
}

type ComboLoginData struct {
	UID   ID     `json:"uid"`
	Guest bool   `json:"guest"`
	Token string `json:"token"`
}

type ComboLoginResponseData struct {
	ComboID       string `json:"combo_id"`
	OpenID        ID     `json:"open_id"`
	AccountType   int32  `json:"account_type"`
	ComboToken    string `json:"combo_token"`
	Data          string `json:"data"`
	FatigueRemind any    `json:"fatigue_remind"`
	Heartbeat     bool   `json:"heartbeat"`
}

type GetFpReq struct {
	AppName   string `json:"app_name"`
	DeviceFp  string `json:"device_fp"`
	DeviceId  string `json:"device_id"`
	ExtFields string `json:"ext_fields"`
	Platform  string `json:"platform"`
	SeedId    string `json:"seed_id"`
	SeedTime  string `json:"seed_time"`
}
type ScanLoginData struct {
	AccountUid string `json:"account_uid"`
	GameToken  string `json:"game_token"`
	IsScanned  bool   `json:"is_scanned"`
}
type ScanStatusCheckData struct {
	AppId  string `json:"app_id"`
	Device string `json:"device"`
	Ticket string `json:"ticket"`
}

/*支付创建*/
type CreateOrderDataReq struct {
	DoNotNoticeAgain bool                      `json:"do_not_notice_again"`
	Order            *CreateOrderDataReq_Order `json:"order"`
	Sign             string                    `json:"sign"`
	Who              *CreateOrderDataReq_Who   `json:"who"`
}
type CreateOrderDataReq_Order struct {
	Account     string `json:"account"`
	Amount      int    `json:"amount"`
	ChannelId   any    `json:"channel_id"`
	ClientType  int    `json:"client_type"`
	Country     string `json:"country"`
	Currency    string `json:"currency"`
	DeliveryUrl string `json:"delivery_url"`
	Device      string `json:"device"`
	Game        string `json:"game"`
	GoodsId     string `json:"goods_id"`
	GoodsNum    string `json:"goods_num"`
	GoodsTitle  string `json:"goods_title"`
	PayPlat     string `json:"pay_plat"`
	PriceTier   string `json:"price_tier"`
	Region      string `json:"region"`
	Uid         any    `json:"uid"`
}
type CreateOrderDataReq_Who struct {
	Account string `json:"account"`
	Token   string `json:"token"`
}

type CreateOrderDataRsp struct {
	Account              string `json:"account"`
	Action               string `json:"action"`
	Amount               string `json:"amount"`
	Blance               int    `json:"blance"`
	CreateTime           string `json:"create_time"`
	Currency             string `json:"currency"`
	DisplayDontShowAgain bool   `json:"display_dont_show_again"`
	EncodeOrder          string `json:"encode_order"`
	ExtInfo              string `json:"ext_info"`
	ForeignSerial        string `json:"foreign_serial"`
	GoodsId              string `json:"goods_id"`
	Method               string `json:"method"`
	NoticeAmount         int    `json:"notice_amount"`
	OrderNo              string `json:"order_no"`
	RedirectUrl          string `json:"redirect_url"`
	SessionToKen         string `json:"session_token"`
}

type CheckOrderDataReq struct {
	Game    string                  `json:"game"`
	OrderNo string                  `json:"order_no"`
	Region  string                  `json:"region"`
	Uid     string                  `json:"uid"`
	Who     *CreateOrderDataReq_Who `json:"who"`
}

// 检查订单状态返回
type CheckOrderDataRsp struct {
	Amount     string `json:"amount"`
	GoodsNum   string `json:"goods_num"`
	GoodsTitle string `json:"goods_title"`
	OrderNo    string `json:"order_no"`
	PayPlat    string `json:"pay_plat"`
	Status     int    `json:"status"` //1 未支付 900 支付成功
}
type OrderFullData struct {
	CreateOrderDataReq *CreateOrderDataReq `json:"create_order_data_req"`
	CreateOrderDataRsp *CreateOrderDataRsp `json:"create_order_data_rsp"`
	CheckOrderDataReq  *CheckOrderDataReq  `json:"check_order_req"`
	CheckOrderDataRsp  *CheckOrderDataRsp  `json:"check_order_data_rsp"`
}

//	type GmTalkData struct {
//		Region string `json:"region"`
//		Uid    string `json:"uid"`
//		Msg    string `json:"msg"`
//	}
//
//	type SendMailData struct {
//		///cfg *config.Config, region string, uid string, sender string, title string, content string, item_list string, ticket string
//		Region   string `json:"region"`
//		Uid      string `json:"uid"`
//		Sender   string `json:"sender"`
//		Title    string `json:"title"`
//		Content  string `json:"content"`
//		ItemList string `json:"item_list"`
//	}
type RealNameDataRsp struct {
	IdentityCard      string `json:"identity_card"`
	RealName          string `json:"realname"`
	RealNameOperation string `json:"realname_operation"`
}

type SdkData struct {
	ApplicationId   int                   `json:"applicationId"`
	ApplicationName string                `json:"applicationName"`
	EventId         int                   `json:"eventId"`
	EventName       string                `json:"eventName"`
	MsgId           string                `json:"msgId"`
	UploadContent   *SdkUploadContentData `json:"uploadContent"`
}

type SdkUploadContentData struct {
	DeviceInfo  *SdkDeviceInfo  `json:"deviceInfo"`
	LogInfo     *SdkLogInfo     `json:"logInfo"`
	UserInfo    *SdkUserInfo    `json:"userInfo"`
	VersionInfo *SdkVersionInfo `json:"versionInfo"`
}

type SdkDeviceInfo struct {
	AddressMac         string  `json:"addressMac" `
	BundleId           string  `json:"bundleId"`
	DeviceId           string  `json:"deviceId"`
	DeviceModel        string  `json:"deviceModel"`
	DeviceName         string  `json:"deviceName"`
	DeviceFp           string  `json:"device_fp"`
	DeviceSciX         int     `json:"device_sciX"`
	DeviceSciY         int     `json:"device_sciY"`
	GpuMemSize         float32 `json:"gpuMemSize"`
	GpuName            string  `json:"gpuName"`
	Platform           int     `json:"platform"`
	ProcessorCount     int     `json:"processorCount"`
	ProcessorFrequency float32 `json:"processorFrequency"`
	ProcessorType      string  `json:"processorType"`
	RamCapacity        float32 `json:"ramCapacity"`
	RamRemain          float32 `json:"ramRemain"`
	RomCapacity        float32 `json:"romCapacity"`
	RomRemain          float32 `json:"romRemain"`
	SoftSciX           int     `json:"soft_sciX"`
	SoftSciY           int     `json:"soft_sciY"`
	SystemInfo         string  `json:"systemInfo"`
}

type SdkLogInfo struct {
	ActionID   int    `json:"actionId"`
	ActionName string `json:"actionType"`
	CBody      string `json:"cBody"`
	LogTime    string `json:"logTime"`
}

type SdkUserInfo struct {
	AccountId   string `json:"accountId"`
	AccountType string `json:"accountType"`
	ChannelId   string `json:"channelId"`
	UserId      string `json:"userId"`
}

type SdkVersionInfo struct {
	ClientVersion string `json:"clientVersion"`
	LogVersion    string `json:"logVersion"`
}

type CrashData struct {
	ApplicationId   int                     `json:"applicationId"`
	ApplicationName string                  `json:"applicationName"`
	EventId         int                     `json:"eventId"`
	EventName       string                  `json:"eventName"`
	MsgId           string                  `json:"msgId"`
	UploadContent   *CrashUploadContentData `json:"uploadContent"`
}

type CrashUploadContentData struct {
	Auid               string `json:"auid"`
	ClientIp           string `json:"clientIp"`
	CpuInfo            string `json:"cpuInfo"`
	DeviceModel        string `json:"deviceModel"`
	DeviceName         string `json:"deviceName"`
	ErrorCategory      string `json:"errorCategory"`
	SErrorCode         string `json:"errorCode"`
	ErrorLevel         string `json:"errorLevel"`
	ErrorCode          int    `json:"error_code"`
	ExceptionSerialNum int    `json:"exceptionSerialNum"`
	Frame              string `json:"frame"`
	GpuInfo            string `json:"gpuInfo"`
	Guid               string `json:"guid"`
	IsRelease          bool   `json:"isRelease"`
	LogType            string `json:"logType"`
	MemoryInfo         string `json:"memoryInfo"`
	Message            string `json:"message"`
	NotifyUser         string `json:"notifyUser"`
	OperatingSystem    string `json:"operatingSystem"`
	ProjectNick        string `json:"projectNick"`
	ServerName         string `json:"serverName"`
	StackTrace         string `json:"stackTrace"`
	SubEorrorCode      string `json:"subErrorCode"`
	Time               int    `json:"time"`
	UserName           string `json:"userName"`
	UserNick           string `json:"userNick"`
	UserId             int    `json:"user_id"`
	Version            string `json:"version"`
}

type RedDotReqData struct {
	GameBiz     string `json:"game_biz"`
	PlayerLevel uint32 `json:"player_level"`
	Region      string `json:"region"`
	Uid         any    `json:"uid"`
}
type LogData struct {
	Auid                string `json:"auid"`
	BuildUrl            string `json:"buildUrl"`
	ClientIp            string `json:"clientIp"`
	CpuInfo             string `json:"cpuInfo"`
	DeviceModel         string `json:"deviceModel"`
	DeviceName          string `json:"deviceName"`
	ErrorCategory       string `json:"errorCategory"`
	ErrorCode           string `json:"errorCode"`
	ErrorCodeToPlatform int    `json:"errorCodeToPlatform"`
	ErrorLevel          string `json:"errorLevel"`
	ExceptionSerialNum  string `json:"exceptionSerialNum"`
	Frame               string `json:"frame"`
	GpuInfo             string `json:"gpuInfo"`
	Guid                string `json:"guid"`
	LogStr              string `json:"logStr"`
	LogType             string `json:"logType"`
	MemoryInfo          string `json:"memoryInfo"`
	NotifyUser          string `json:"notifyUser"`
	OperatingSystem     string `json:"operatingSystem"`
	Pos                 string `json:"pos"`
	ServerName          string `json:"serverName"`
	StackTrace          string `json:"stackTrace"`
	SubEorrorCode       int    `json:"subErrorCode"`
	Time                string `json:"time"`
	Uid                 int    `json:"uid"`
	UserName            string `json:"userName"`
	Version             string `json:"version"`
}

// ***********************************Key*********************************************//
type PrivateKey struct {
	*rsa.PrivateKey
	PrivateKeyPEM *PrivateKeyPEM
}
type PublicKey struct {
	*rsa.PublicKey
	PublicKeyPEM *PublicKeyPEM
}
type PrivateKeyPEM struct {
	PrivateKeyPKCS1 string
	PrivateKeyPKCS8 string
}
type PublicKeyPEM struct {
	PublicKeyPKCS1 string
	PublicKeyPKCS8 string
}
type Secret struct {
	mutex               sync.Mutex
	ServerPrivateKeyMap map[uint32]*PrivateKey
	ServerPublicKeyMap  map[uint32]*PublicKey
	ClientPublicKeyMap  map[uint32]*PublicKey
	ClientPrivateKeyMap map[uint32]*PrivateKey
	PasswordPrivateKey  *PrivateKey
	PasswordPublicKey   *PublicKey
	PayPrivateKey       *PrivateKey
	PayPublicKey        *PublicKey
}

// ***********************************HttpServer*********************************************//
type Server struct {
	config       *config.Config
	configLock   sync.Mutex
	logger       *logger.Logger
	secret       *Secret
	router       *gin.Engine
	server       *http.Server
	store        *database.Store
	clientLogger *zerolog.Logger
	staticFiles  map[string][]byte
	jsonFiles    map[string][]byte
}
