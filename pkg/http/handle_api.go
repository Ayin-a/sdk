package http

import (
	"crypto/sha256"
	"encoding/hex"
	"hk4e_sdk/pkg/logger"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/process"
)

var (
	// mutex用于确保同时只有一个请求可以触发重启
	mutex            sync.Mutex
	emailVerifyCodes []string
)

// 改index死妈
func (s *Server) handleAPICreateAndFinishOrder(c *gin.Context) {
	var ts = time.Now().Unix()
	var data OrderFullData
	data.CheckOrderDataRsp = &CheckOrderDataRsp{
		OrderNo: RandNumStr(19),
		PayPlat: "alipay",
	}

	data.CheckOrderDataReq = &CheckOrderDataReq{
		Region: "dev_client",
	}
	data.CreateOrderDataReq = &CreateOrderDataReq{

		Order: &CreateOrderDataReq_Order{
			Uid:        137641255,
			PriceTier:  "Tier_60",
			GoodsId:    "ys_chn_primogem6thstall_tier60",
			GoodsTitle: "60创世结晶",
			GoodsNum:   "1",
			Currency:   "CNY",
			ChannelId:  1,
		},
	}
	data.CreateOrderDataRsp = &CreateOrderDataRsp{
		CreateTime: strconv.Itoa(int(ts)),
	}

	if c.Request.Method == http.MethodGet {
		uidStr := c.Request.URL.Query().Get("uid")
		uid, err := strconv.ParseUint(uidStr, 10, 32)
		if err != nil {
			logger.Error("ParseUint Uid Fail")
		}
		data.CreateOrderDataReq.Order.Uid = uint32(uid)
		data.CheckOrderDataReq.Region = c.Request.URL.Query().Get("region")
		data.CreateOrderDataReq.Order.PriceTier = c.Request.URL.Query().Get("price_tier")
		data.CreateOrderDataReq.Order.GoodsId = c.Request.URL.Query().Get("goods_id")
		data.CreateOrderDataReq.Order.GoodsTitle = c.Request.URL.Query().Get("goods_title")
	}
	if c.Request.Method == http.MethodPost {

	}
	x := s.PaySuccReq(data)
	c.String(http.StatusOK, x)

}

// 改index死妈
func (s *Server) handleAPIChangePassword(c *gin.Context) {
	var req AccountBaseData

	var emailVerifyStatus = false
	if err := c.BindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
		return
	}
	switch req.Action {
	case "email_check":
		accountData, err := s.store.Account().GetAccountByUsername(c, req.Account)
		if err != nil {

			c.JSON(http.StatusOK, NewResponse(-301, gin.H{}))
			return
		}
		var emailVerifyCode = RandNumStr(6)
		emailVerifyCodes = append(emailVerifyCodes, emailVerifyCode)
		emailErr := SendMail(s.config.EMailAddress, s.config.EMailAuthCode, s.config.EMailHost, s.config.EMailHostPort, accountData.Email, "SYSTEM_MANAGER", "注册验证码", "<h1>你的验证码是:"+emailVerifyCode+" </h1>")

		if emailErr != nil {
			c.JSON(http.StatusOK, NewResponse(-3000, gin.H{}))
			return
		}

		c.JSON(http.StatusOK, NewResponse(0, gin.H{}))
		return
	case "change_password":
		for code := range emailVerifyCodes {
			if req.EMailVerifyCode == emailVerifyCodes[code] {
				emailVerifyStatus = true
				FromStringsRemoveString(emailVerifyCodes, req.EMailVerifyCode)
				break
			}

		}
		if emailVerifyStatus {
			accountData, err := s.store.Account().GetAccountByUsername(c, req.Account)
			if err != nil {

				c.JSON(http.StatusOK, NewResponse(-301, gin.H{}))
				return
			}

			concatenated := req.Password + s.config.Accountkey
			hasher := sha256.New()
			hasher.Write([]byte(concatenated))
			computedHash := hex.EncodeToString(hasher.Sum(nil))
			if err := s.store.Account().UpdateAccountPassword(c, accountData.ID, computedHash); err != nil { //保存哈希到数据库 用于密码验证
				c.JSON(http.StatusOK, NewResponse(-2003, gin.H{}))
				return
			}

			c.JSON(http.StatusOK, NewResponse(0, gin.H{}))
			return

		} else {

			c.JSON(http.StatusOK, NewResponse(-3001, gin.H{}))
			return
		}
	default:
		c.JSON(http.StatusOK, NewResponse(-3101, gin.H{}))

	}
}

// 改index死妈
func (s *Server) handleAPILogin(c *gin.Context) {
	s.handleSDKShieldLogin(c)
}

// 改index死妈
func (s *Server) handleAPIStatus(c *gin.Context) {
	pid := int32(os.Getpid())
	p, err := process.NewProcess(pid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cpuPercent, err := p.Percent(time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	memInfo, err := p.MemoryInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	memUsageMB := float64(memInfo.RSS) / 1024.0 / 1024.0
	virtualMemUsageMB := float64(memInfo.VMS) / 1024.0 / 1024.0

	c.JSON(http.StatusOK, gin.H{
		"Physical Memory (MB)": memUsageMB,
		"Virtual Memory (MB)":  virtualMemUsageMB,
		"CPU (%)":              cpuPercent,
	})
}

// 改index死妈
func (s *Server) handleAPIRestart(c *gin.Context) {
	// 检查 key 参数
	auth := c.Query("key")
	if auth != s.config.ServerKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// 使用新的 goroutine 在1秒后执行 restart()
	go func() {
		time.Sleep(1 * time.Second)
		err := Restart()
		if err != nil {
			logger.Error("Error restarting application: %v\n", err)
		}
	}()

	c.String(http.StatusOK, "succ")
}
