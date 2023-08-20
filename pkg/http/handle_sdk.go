package http

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hk4e_sdk/pkg/database"
	"hk4e_sdk/pkg/logger"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var orderFullDataMap = make(map[string]OrderFullData)
var finishOrderNo []string
var scanLoginDataMap = make(map[string]ScanLoginData)
var loginTokenDataMap = make(map[string]ShieldVerifyRequestData)
var loginTokenPool = &sync.Pool{New: func() interface{} { return make([]byte, 24) }}

// 改index死妈
func (s *Server) handleSdkEmptySDKRspJson(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{}) //空响应
}

// handleSDKRiskyCheck 处理SDK风险检查 //滑块可以在此启用
// 改index死妈
func (s *Server) handleSDKRiskyCheck(c *gin.Context) {
	c.JSON(http.StatusOK, NewResponse(0, gin.H{
		"id":      "",
		"action":  RISKY_ACTION_NONE,
		"geetest": nil,
	}))
}

// handleSdkAuthStatus 处理SDK认证状态
// 改index死妈
func (s *Server) handleSdkAuthStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

// 改index死妈
func (s *Server) handleSDKGetGateAddress(c *gin.Context) {

	var gateserveraddr string

	if gateserveraddr == "" {
		c.JSON(http.StatusOK, NewResponse(0, gin.H{"address_list": []gin.H{}}))
		return
	}
	ip, portStr := strings.Split(gateserveraddr, ":")[0], strings.Split(gateserveraddr, ":")[1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		c.JSON(http.StatusOK, NewResponse(0, gin.H{"address_list": []gin.H{}}))
		return
	}
	c.JSON(http.StatusOK, NewResponse(0, gin.H{"address_list": []gin.H{{"ip": ip, "port": port}}}))
}

// 用户登录，以及验证滑块
// 改index死妈
func (s *Server) handleSDKShieldLogin(c *gin.Context) {
	clientIP := c.ClientIP() //取客户端ip

	isInBlacklist, comment, err := s.store.Blacklist().IsIPInBlacklist(c, clientIP) //是否在黑名单
	if err != nil {
		logger.Error("Failed to check IP blacklist") //检查黑名单失败，或许为ip格式错误
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-208, nil))
		return
	}
	if isInBlacklist { //在黑名单，返回原因
		message := fmt.Sprintf("您的ip在黑名单-原因:%s", comment)
		c.AbortWithStatusJSON(http.StatusOK, NewSimpleResponse(-206, message))
		return
	}

	var req ShieldLoginRequestData //解析请求的数据
	if err := c.BindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON") //解析失败
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-111, nil))
		return
	}

	account, err := s.serviceShieldLogin(c, req.Account, req.Password, req.IsCrypto) //取用户数据

	if err != nil && err == sql.ErrNoRows && s.config.AutoSignUp { //不在数据库内，如果开启自动注册，继续下一步
		account, err = s.serviceCreateAccountWithEmail(c, clientIP, req.Account, req.Password, req.Account+"@a.com", req.IsCrypto) //创建用户
		if err != nil {
			logger.Error("Failed to create account")
			c.AbortWithStatusJSON(http.StatusOK, NewResponse(-207, nil))
			return
		}
	} else if err != nil {
		logger.Error("Failed to shield login")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-204, nil))
		return
	}

	err = s.store.Account().UpdateAccountIPById(c, account.ID, clientIP) //更新客户端ip
	if err != nil {
		logger.Error("Failed to update account IP")
	}

	var resp ShieldLoginResponseData //返回对应数据
	resp.Account = &Account{
		UID:           ID(account.ID),
		Name:          account.Username,
		Email:         account.Email,
		IsEmailVerify: "0",
		Token:         account.LoginToken,
		Country:       "US",
		AreaCode:      "**",
	}
	resp.RealNameOperation = "None"
	c.JSON(http.StatusOK, NewResponse(0, &resp))
	logger.Info("[ account:%s, login succ IP: %s ]", resp.Account.Name, clientIP)

}

// 用户短信登录
// 改index死妈
func (s *Server) handleSDKLoginMobile(c *gin.Context) {
	var req MobileLoginData
	//获取client ip
	clientIP := c.ClientIP()
	// 判断IP是否在黑名单中
	isInBlacklist, comment, err := s.store.Blacklist().IsIPInBlacklist(c, clientIP)
	if err != nil {
		logger.Error("Failed to check IP blacklist")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-206, nil))
		return
	}
	if isInBlacklist {
		message := fmt.Sprintf("您的ip在黑名单-原因:%s", comment)
		c.AbortWithStatusJSON(http.StatusOK, NewSimpleResponse(-206, message))
		return
	}

	if err := c.BindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
		return
	}
}

// 改index死妈
func (s *Server) generateAndAssignLoginToken(record *database.Account, ctx context.Context) error {
	loginToken := loginTokenPool.Get().([]byte)
	defer loginTokenPool.Put(loginToken)

	_, err := rand.Read(loginToken)
	if err != nil {
		return err
	}

	record.LoginToken = base64.RawStdEncoding.EncodeToString(loginToken)
	if err := s.store.Account().UpdateAccountLoginToken(ctx, record.ID, record.LoginToken); err != nil {
		return err
	}

	// 更新token的过期时间，
	expiration := time.Now().Add(time.Hour * time.Duration(s.config.AccountTokenExp))
	if err := s.store.Account().UpdateAccountTokenExpiration(ctx, record.ID, expiration); err != nil {
		return err
	}

	return nil
}

// 改index死妈
func (s *Server) serviceShieldLogin(ctx context.Context, username, password string, isCrypto bool) (*database.Account, error) {
	var record *database.Account
	var err error

	if !strings.Contains(username, "@") {
		record, err = s.store.Account().GetAccountByUsername(ctx, username)
	} else {
		record, err = s.store.Account().GetAccountByEmail(ctx, username)
	}
	if err != nil {
		return nil, err
	}

	if s.config.PassSignIn {
		p, err := s.secret.PasswordPrivateKey.DecryptBase64(password)
		if err != nil {
			return nil, err
		}
		password = string(p)

		concatenated := password + s.config.Accountkey

		hasher := sha256.New()
		hasher.Write([]byte(concatenated))
		computedHash := hex.EncodeToString(hasher.Sum(nil))

		if computedHash != record.Password {
			return nil, ErrInvalidPassword
		}
	}

	// 如果 LoginToken 已经过期，则重新生成并更新到数据库中
	if time.Now().After(record.TokenExpiration) {
		err = s.generateAndAssignLoginToken(record, ctx)
		if err != nil {
			return nil, err
		}
	}

	return record, nil
}

// 改index死妈
func (s *Server) serviceCreateGuestAccount(ctx context.Context, device string) (record *database.Account, err error) {
	record = &database.Account{IsGuest: true, Device: device}
	if err = s.store.Account().CreateAccount(ctx, record); err != nil {
		return nil, err
	}

	return record, s.generateAndAssignLoginToken(record, ctx)
}

// 改index死妈
func (s *Server) serviceCreateAccountWithEmail(ctx context.Context, clientIP, username, password, eMailAddress string, isCrypto bool) (record *database.Account, err error) {
	// Check username
	reg := regexp.MustCompile("^[\u4e00-\u9fa5a-zA-Z0-9.@#&_-]*$")
	if !reg.MatchString(username) {
		return nil, errors.New("-206")
	}

	// Check IP
	accountsFromSameIP, err := s.store.Account().GetAccountsByIP(ctx, clientIP)
	if err != nil {
		logger.Error("Failed to fetch accounts from the same IP")
		return nil, err
	}

	if len(accountsFromSameIP) >= s.config.AccountLimit {
		err = s.store.Blacklist().AddIPToBlacklist(ctx, clientIP, "系统拉黑")
		if err != nil {
			logger.Error("Failed to add IP to the blacklist")
			return nil, err
		}
	}

	if !isValidEmail(username) {
		record = &database.Account{Username: username, Email: eMailAddress, IsGuest: false}
		if err = s.store.Account().CreateAccount(ctx, record); err != nil {
			return nil, err
		}
	} else {
		record = &database.Account{Username: username, Email: username}
		if err = s.store.Account().CreateAccount(ctx, record); err != nil {
			return nil, err
		}
	}

	if s.config.PassSignIn {
		p, err := s.secret.PasswordPrivateKey.DecryptBase64(password)
		if err != nil {
			return nil, err
		}
		password = string(p)
		concatenated := password + s.config.Accountkey

		hasher := sha256.New()
		hasher.Write([]byte(concatenated))
		computedHash := hex.EncodeToString(hasher.Sum(nil))

		if err := s.store.Account().UpdateAccountPassword(ctx, record.ID, computedHash); err != nil {
			return nil, err
		}
	}

	// 设置登录token的过期时间
	expiration := time.Now().Add(time.Hour * time.Duration(s.config.AccountTokenExp)) // 可以根据实际需要调整
	if err := s.store.Account().UpdateAccountTokenExpiration(ctx, record.ID, expiration); err != nil {
		return nil, err
	}

	logger.Info("Account created successfully. Username: %s, Login Token: %s, Client IP: %s", username, record.LoginToken, clientIP)
	return record, s.generateAndAssignLoginToken(record, ctx)
}

// handleSDKShieldVerify 登录token验证，成功进入选择服务器页面
// 改index死妈
func (s *Server) handleSDKShieldVerify(c *gin.Context) {

	clientIP := c.ClientIP()
	// 判断IP是否在黑名单中
	isInBlacklist, comment, err := s.store.Blacklist().IsIPInBlacklist(c, clientIP)
	if err != nil {
		logger.Error("Failed to check IP blacklist")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-208, nil))
		return
	}
	if isInBlacklist {
		message := fmt.Sprintf("您的ip在黑名单-原因:%s", comment)
		c.AbortWithStatusJSON(http.StatusOK, NewSimpleResponse(-206, message))
		return
	}
	var req ShieldVerifyRequestData
	if err := c.BindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-210, nil))
		return
	}

	account, err := s.serviceShieldVerify(c, int64(req.UID), req.Token)
	if err != nil {
		logger.Error("验证失败, UID: %s Token: %s", req.UID, req.Token)
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-210, nil))
		return
	}
	// 检查token是否过期
	if time.Now().After(account.TokenExpiration) {
		logger.Error("Token已过期, UID: %s", req.UID)
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-210, nil)) // -210是token过期的错误代码，
		return
	}

	// 更新ip到数据库内
	err = s.store.Account().UpdateAccountIPById(c, account.ID, clientIP)
	if err != nil {
		logger.Error("Failed to update account IP")
	}

	var resp ShieldVerifyResponseData
	resp.Account = &Account{
		UID:               req.UID,
		Email:             account.Email,
		Name:              account.Username,
		IsEmailVerify:     "0",
		Token:             account.LoginToken,
		Country:           "US",
		AreaCode:          "",
		IdentityCard:      "320************110",
		RealName:          "*游",
		SafeMobile:        "188****8888",
		ReactivateTicket:  "",
		SonyName:          "",
		TwitterName:       "",
		GoogleName:        "",
		TapName:           "",
		FacebookName:      "",
		GameCenterName:    "",
		DeviceGrantTicket: "",
		AppleName:         "",
	}
	resp.RealNameOperation = "None"
	c.JSON(http.StatusOK, NewResponse(0, &resp))
	logger.Info("[ VerifySucc, account: %s, loginToken: %s, clientIp: %s ]", account.Username, account.LoginToken, clientIP)
}

// handleSDKActionTicket 处理SDK动作票据？？
// 改index死妈
func (s *Server) handleSDKActionTicket(c *gin.Context) {
	c.JSON(http.StatusOK, NewResponse(0, gin.H{
		"ticket": "",
	}))
}

// handleSDKRealName 处理SDK实名认证
// 改index死妈
func (s *Server) handleSDKRealName(c *gin.Context) {
	var realNameDataRsp RealNameDataRsp
	realNameDataRsp.IdentityCard = "320************110"
	realNameDataRsp.RealName = "*游"
	realNameDataRsp.RealNameOperation = "None"
	c.JSON(http.StatusOK, NewResponse(0, realNameDataRsp))
}

// serviceShieldVerify 验证登录token
// 改index死妈
func (s *Server) serviceShieldVerify(ctx context.Context, id int64, token string) (record *database.Account, err error) {
	record, err = s.store.Account().GetAccount(ctx, id)
	if err != nil {
		return nil, err
	}
	if record.LoginToken == "" || record.LoginToken != token {
		return nil, ErrInvalidLoginToken
	}
	return record, nil
}

// handleSDKComboLogin 处理SDK组合登录
// 改index死妈
func (s *Server) handleSDKComboLogin(c *gin.Context) {
	clientIP := c.ClientIP()

	isInBlacklist, comment, err := s.store.Blacklist().IsIPInBlacklist(c, clientIP)
	if err != nil {
		logger.Error("Failed to check IP blacklist")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-208, nil))
		return
	}
	if isInBlacklist {
		message := fmt.Sprintf("您的ip在黑名单-原因:%s", comment)
		c.AbortWithStatusJSON(http.StatusOK, NewSimpleResponse(-206, message))
		return
	}

	var req ComboLoginRequestData
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON")
		c.JSON(http.StatusOK, NewResponse(-202, nil))
		return
	}

	var data ComboLoginData
	if err := json.Unmarshal([]byte(req.Data), &data); err != nil {
		logger.Error("Failed to unmarshal data")
		c.JSON(http.StatusOK, NewResponse(-202, nil))
		return
	}

	if data.Guest { //Guest login verification
		tAccount, _ := s.store.Account().GetAccount(c, int64(data.UID))
		if tAccount.Device != req.Device {
			c.JSON(http.StatusOK, ErrInvalidLoginToken)
			return
		}
		data.Token = tAccount.LoginToken //Set login Token
	}

	account, err := s.serviceComboLogin(c, int64(data.UID), data.Token)
	if err != nil {
		logger.Error("Failed to combo login")
		c.JSON(http.StatusOK, NewResponse(-202, nil))
		return
	}

	// Update client IP
	err = s.store.Account().UpdateAccountIPById(c, account.ID, clientIP)
	if err != nil {
		logger.Error("Failed to update account IP")
	}

	respData := map[string]bool{"guest": data.Guest}
	respDataJSON, _ := json.Marshal(respData)

	c.JSON(http.StatusOK, NewResponse(0, ComboLoginResponseData{
		ComboID:       "0",
		OpenID:        ID(account.ID),
		ComboToken:    account.ComboToken,
		Data:          string(respDataJSON),
		Heartbeat:     false,
		AccountType:   1,
		FatigueRemind: nil,
	}))
	logger.Info("[ comboToken login successful, Account ID: %d, ComboToken: %s, Login token: %s, Client IP: %s ]", account.ID, account.ComboToken, account.LoginToken, clientIP)

}

// serviceComboLogin 组合登录服务
// 改index死妈
func (s *Server) serviceComboLogin(ctx context.Context, id int64, token string) (record *database.Account, err error) {
	if record, err = s.serviceShieldVerify(ctx, id, token); err != nil {
		return nil, err
	}
	comboToken := make([]byte, 20)
	_, err = rand.Read(comboToken)
	if err != nil {
		return nil, err
	}
	record.ComboToken = hex.EncodeToString(comboToken)
	err = s.store.Account().UpdateAccountComboToken(ctx, record.ID, record.ComboToken)
	if err != nil {
		return nil, err
	}
	return record, nil
}

var LogDataPool = &sync.Pool{
	New: func() interface{} {
		return &LogData{}
	},
}

// handleNLogUpload handleNLogUpload 处理客户端日志上传
// 改index死妈
func (s *Server) handleNLogUpload(c *gin.Context) {
	nlog := LogDataPool.Get().(*LogData) // 从池中获取对象
	defer LogDataPool.Put(nlog)          // 使用完后放回池中

	if err := c.BindJSON(nlog); err != nil {
		logger.Error("Failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
		return
	}

	switch nlog.LogType {
	case "Warning", "Error":
		s.clientLogger.Info().Str("UID", strconv.Itoa(nlog.Uid)).Str("Message", nlog.LogStr).Str("StackTrace", nlog.StackTrace).Msg("Client log")
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})
}

// handleSDKRedDotList 处理SDK红点列表
// 改index死妈
func (s *Server) handleSDKRedDotList(c *gin.Context) {
	var req RedDotReqData
	if err := c.BindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
		return
	}

	c.JSON(http.StatusOK, NewResponse(0, gin.H{
		"infos": []string{
			"Hot Update",
			"Reload",
		},
	}))
}

// handleCrashDataUpload 处理崩溃数据上传
// 改index死妈
func (s *Server) handleCrashDataUpload(c *gin.Context) {
	var crashlog []CrashData
	if err := c.BindJSON(&crashlog); err != nil {
		logger.Error("Failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
		return
	}
	data, _ := io.ReadAll(c.Request.Body)

	s.clientLogger.Info().Msgf("req body: %v", string(data))
	c.JSON(http.StatusOK, gin.H{"code": 0})
}

// handleSdkDataUpload 处理SDK数据上传
// 改index死妈
func (s *Server) handleSdkDataUpload(c *gin.Context) {
	//var sdkData []SdkData
	//if err := c.BindJSON(&sdkData); err != nil {
	//	logger.Error("Failed to bind JSON")
	///	c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
	//	return
	////}
	//dataLen := len(sdkData)
	//for i := 0; i < dataLen; i++ {
	//	logger.Info("ClientLog::Sdk:" + sdkData[i].UploadContent.UserInfo.AccountId)
	///	logger.Info("ClientLog::Sdk::UID:" + sdkData[i].UploadContent.UserInfo.UserId)
	//	logger.Info("ClientLog::Sdk::CBody:" + sdkData[i].UploadContent.LogInfo.CBody)
	//}
	//data, _ := io.ReadAll(c.Request.Body)
	//logger.Info("req body: %v", string(data))
	c.JSON(http.StatusOK, gin.H{"code": 0})
}

/******************************创建订单*******************************/
//改index死妈
func (s *Server) handleSDKCreateOrder(c *gin.Context) { //创建订单
	gameBiz := c.Param("game_biz")
	go logger.Info("game_biz:%s", gameBiz)

	var ts = time.Now().Unix()
	var orderData CreateOrderDataReq
	if err := c.BindJSON(&orderData); err != nil {
		logger.Error("Failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
		return
	}

	var orderNo = RandNumStr(19)
	qrPayUrl := s.config.QrPayUrl + "?order_no=" + orderNo
	var rsp Response
	var createOrderData CreateOrderDataRsp
	createOrderData.Account = orderData.Order.Account
	createOrderData.Action = ""
	createOrderData.Amount = strconv.Itoa(orderData.Order.Amount)
	createOrderData.CreateTime = strconv.Itoa(int(ts))
	createOrderData.DisplayDontShowAgain = orderData.DoNotNoticeAgain
	createOrderData.EncodeOrder = qrPayUrl
	createOrderData.ExtInfo = ""
	createOrderData.ForeignSerial = ""
	createOrderData.GoodsId = orderData.Order.GoodsId
	createOrderData.Method = ""
	createOrderData.NoticeAmount = 0
	createOrderData.OrderNo = orderNo
	createOrderData.RedirectUrl = ""
	createOrderData.SessionToKen = ""
	createOrderData.Currency = orderData.Order.Currency
	rsp.Message = "OK"
	rsp.Retcode = 0
	rsp.Data = createOrderData
	var sOrderData CheckOrderDataRsp
	sOrderData.Amount = strconv.Itoa(orderData.Order.Amount)
	sOrderData.GoodsNum = orderData.Order.GoodsNum
	sOrderData.GoodsTitle = orderData.Order.GoodsTitle
	sOrderData.OrderNo = createOrderData.OrderNo
	sOrderData.PayPlat = orderData.Order.PayPlat
	sOrderData.Status = 1

	var orderFullDataMapTmp OrderFullData
	orderFullDataMapTmp.CreateOrderDataReq = &orderData
	orderFullDataMapTmp.CreateOrderDataRsp = &createOrderData
	orderFullDataMapTmp.CheckOrderDataRsp = &sOrderData
	orderFullDataMap[createOrderData.OrderNo] = orderFullDataMapTmp
	go logger.Info("\n扫码支付:\n账号UID:%s\n游戏UID:%s\n商品名称:%s\n扫码支付url:%s\n", orderData.Order.Account, orderData.Order.Uid, orderData.Order.GoodsTitle, qrPayUrl)
	c.JSON(http.StatusOK, &rsp)
}

// 改index死妈
func (s *Server) handleSDKCheckOrder(c *gin.Context) { //检查订单
	gameBiz := c.Param("game_biz")
	go logger.Info("game_biz:%s", gameBiz)
	var req CheckOrderDataReq
	if err := c.BindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
		return
	}
	var orderFullDataMapTmp = orderFullDataMap[req.OrderNo]
	orderFullDataMapTmp.CheckOrderDataReq = &req
	var rsp Response
	var checkOrderDataRsp CheckOrderDataRsp
	if orderFullDataMap[req.OrderNo].CheckOrderDataRsp == nil {
		checkOrderDataRsp.Status = 0
		rsp.Message = "OK"
		rsp.Retcode = 0
		rsp.Data = checkOrderDataRsp

		c.JSON(http.StatusOK, &rsp)
		return
	}
	checkOrderDataRsp = *orderFullDataMap[req.OrderNo].CheckOrderDataRsp
	for i := 0; i < len(finishOrderNo); i++ {
		if req.OrderNo == finishOrderNo[i] { //判断支付信息
			checkOrderDataRsp.Status = 900

			defer s.PaySuccReq(orderFullDataMapTmp)

		}
	}

	rsp.Message = "OK"
	rsp.Retcode = 0
	rsp.Data = checkOrderDataRsp

	c.JSON(http.StatusOK, &rsp)

}

/********************************第三方验证*****************************************/
//改index死妈
func (s *Server) handleSDKLoginByThirdparty(c *gin.Context) {

	var req ThirdpartyLoginRequestData
	if err := c.BindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
		return
	}

	shieldLoginRequestData := loginTokenDataMap[req.AccessToken]
	account, err := s.serviceShieldVerify(c, int64(shieldLoginRequestData.UID), shieldLoginRequestData.Token)

	if err != nil {
		logger.Error("Failed to shield verify")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-210, nil))
		return
	}

	resp := ShieldVerifyResponseData{
		Account: &Account{
			UID:               ID(account.ID),
			Email:             account.Email,
			Name:              account.Username,
			IsEmailVerify:     "0",
			Token:             account.LoginToken,
			Country:           "US",
			AreaCode:          "",
			IdentityCard:      "320************110",
			RealName:          "*游",
			SafeMobile:        "188****8888",
			ReactivateTicket:  "",
			SonyName:          "",
			TwitterName:       "",
			GoogleName:        "",
			TapName:           "",
			FacebookName:      "",
			GameCenterName:    "",
			DeviceGrantTicket: "",
			AppleName:         "",
		},
		RealNameOperation: "None",
	}

	c.JSON(http.StatusOK, NewResponse(0, &resp))
	delete(loginTokenDataMap, req.AccessToken)
}

// 登录信息上传
// 改index死妈
func (s *Server) handleSDKUploadLoginToken(c *gin.Context) {
	var req ShieldVerifyRequestData

	if err := c.BindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
		return
	}
	var ticket = c.Request.URL.Query().Get("ticket")
	loginTokenDataMap[ticket] = req
	c.JSON(http.StatusOK, NewResponse(0, gin.H{}))
}

// ********************************游客登录*********************************//
// 改index死妈
func (s *Server) handleSDKGuestLogin(c *gin.Context) {
	//gameBiz := c.Param("game_biz")
	//logger.Info("game_biz:%s", gameBiz)
	var req GuestLoginRequestData
	var isNew = false
	if err := c.BindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
		return
	}

	account, err := s.store.Account().GetAccountByDevice(c, req.Device)
	if err != nil {
		// c.JSON(http.StatusOK, NewResponse(-101, gin.H{}))
		// return
		isNew = true
		account, err = s.serviceCreateGuestAccount(c, req.Device)
		if err != nil {
			c.JSON(http.StatusOK, NewResponse(-2001, gin.H{}))
			return
		}

	}
	logger.Info("游客登录账号UID:", account.ID)
	c.JSON(http.StatusOK, NewResponse(0, gin.H{
		"guest_id": account.ID,
		"newly":    isNew,
	}))

}

// **********************************扫码登录******************************//
// 扫码登录配置
// 改index死妈
func (s *Server) handleSDKScanLoginConfig(c *gin.Context) {
	var scanLoginData ScanLoginData
	var ticket = RandC16Data(24)
	scanLoginData.IsScanned = false
	scanLoginData.AccountUid = ""
	scanLoginData.GameToken = ""
	scanLoginDataMap[ticket] = scanLoginData
	logger.Info("ticket:" + ticket)
	scanLoginUrl := s.config.SdkBaseUrl + "/view/qr_code_login?app_id=" + c.Request.URL.Query().Get("app_id") + "&app_name=%E5%8E%9F%E7%A5%9E&bbs=true&biz_key=hk4e_cn&expire=1673612993&ticket=" + ticket
	logger.Info("\n账号UID:%s\n扫描登录地址:%s", scanLoginData.AccountUid, scanLoginUrl)
	c.JSON(http.StatusOK, NewResponse(0, gin.H{
		//"url": s.config.SdkBaseUrl + "/view/qr_login",
		"url": scanLoginUrl,
	}))

}

// 改index死妈
func (s *Server) handleSDKScanned(c *gin.Context) { //开始扫码
	var ticket = c.Request.URL.Query().Get("ticket") //扫码的标识
	var tmpScanLoginData = scanLoginDataMap[ticket]
	tmpScanLoginData.IsScanned = true
	scanLoginDataMap[ticket] = tmpScanLoginData
	c.JSON(http.StatusOK, NewResponse(0, gin.H{}))
}

// 扫码登录的按钮
// 改index死妈
func (s *Server) handleSDKQRLoginCheck(c *gin.Context) {
	var req ShieldLoginRequestData

	var ticket = c.Request.URL.Query().Get("ticket") //扫码的标识
	if ticket == "undefined" {
		ticket = RandC16Data(24)
	}
	if err := c.BindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
		return
	}
	if req.Account == "" {
		logger.Error("account is empty", "Bad request")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
		return
	}
	var record *database.Account
	var err error
	if !strings.Contains(req.Account, "@") { //是否是邮箱
		record, err = s.store.Account().GetAccountByUsername(c, req.Account) //获取数据库中的数据
	} else {
		record, err = s.store.Account().GetAccountByEmail(c, req.Account)
	}
	if err != nil {
		c.JSON(http.StatusOK, NewResponse(-301, gin.H{})) //用户信息获取失败
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(record.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusOK, ErrInvalidPassword)
		return
	}
	loginToken := make([]byte, 24) //
	_, err = rand.Read(loginToken)
	if err != nil {
		return
	}
	record.LoginToken = base64.RawStdEncoding.EncodeToString(loginToken)
	err = s.store.Account().UpdateAccountLoginToken(c, record.ID, record.LoginToken)
	if err != nil {
		return
	}
	var tmpScanLoginData = scanLoginDataMap[ticket]
	tmpScanLoginData.AccountUid = strconv.FormatInt(record.ID, 10) //int64 转10进制数据
	tmpScanLoginData.GameToken = record.LoginToken
	scanLoginDataMap[ticket] = tmpScanLoginData

	c.JSON(http.StatusOK, NewResponse(0, gin.H{
		// "account_uid": scanLoginDataMap[ticket].AccountUid,
		//"game_token":  scanLoginDataMap[ticket].GameToken,
		// "ticket":      ticket,
	}))
}

// 改index死妈
func (s *Server) handleSDKScanStatusCheck(c *gin.Context) {

	var req ScanStatusCheckData
	if err := c.BindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
		return
	}

	ticket := req.Ticket
	scanData := scanLoginDataMap[ticket]

	if scanData.AccountUid != "" && scanData.GameToken != "" {
		c.JSON(http.StatusOK, NewResponse(0, gin.H{
			"payload": gin.H{"ext": "", "proto": "Account", "raw": "{\"uid\":\"" + scanData.AccountUid + "\",\"token\":\"" + scanData.GameToken + "\"}"},
			"stat":    "Confirmed",
		}))
		delete(scanLoginDataMap, ticket)
	} else {
		stat := "Init"
		if scanData.IsScanned {
			stat = "Scanned"
		}
		c.JSON(http.StatusOK, NewResponse(0, gin.H{
			"payload": gin.H{"ext": "", "proto": "Raw", "raw": ""},
			"stat":    stat,
		}))
	}
}

// 改index死妈
func (s *Server) handleSDKSendSMS(c *gin.Context) {
	c.JSON(http.StatusOK, NewResponse(0, gin.H{
		"action": "Login",
	}))
}
