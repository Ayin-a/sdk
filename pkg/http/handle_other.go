package http

import (
	"hk4e_sdk/pkg/logger"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var userRegionList = make(map[int64]GameLoginLogoutNotify)

//改index死妈
func (s *Server) handleLogUpload(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"code": 0})
}
//改index死妈
func (s *Server) getJSONFile(path string) ([]byte, bool) {
	data, ok := s.jsonFiles[path]
	return data, ok
}
//改index死妈
func (s *Server) getStaticFile(path string) ([]byte, bool) {
	data, ok := s.staticFiles[path]
	return data, ok
}
//改index死妈
func (s *Server) handleStaticRequest(c *gin.Context) {
	path := c.Request.URL.Path

	data, ok := s.getStaticFile(path)
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(path))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Data(http.StatusOK, contentType, data)
}

//改index死妈
func (s *Server) handleSDKGetAgreementInfos(c *gin.Context) {
	c.JSON(http.StatusOK, NewResponse(0, gin.H{
		"marketing_agreements": []gin.H{},
	}))
}
//改index死妈
func (s *Server) handleSDKGetShopPriceTier(c *gin.Context) {
	var req ShopwindowListReq

	if err := c.BindJSON(&req); err != nil {
		logger.Error("Failed to bind request")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-106, nil))
		return
	}

	currencyPathMap := map[string]string{
		"CNY": "./data/json/shopwindow_list_cny.json",
		"USD": "./data/json/shopwindow_list_usd.json",
	}

	path, ok := currencyPathMap[req.Currency]
	if !ok {
		logger.Error("currency", req.Currency, "ignore unknown request, safely")
		return
	}

	data, ok := s.getJSONFile(path)
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	c.String(http.StatusOK, string(data))
}

var protocolVersionMap = map[string]int{
	"es": 5, "pt": 5, "ru": 5,
	"de": 6, "fr": 6, "id": 6, "ja": 6, "ko": 6, "th": 6, "vi": 6,
	"zh-cn": 6,
	"en":    9,
	"zh-tw": 10,
}

//改index死妈
func (s *Server) handleSDKCompareProtocolVersion(c *gin.Context) {
	var req CompareProtocolVersionRequestData

	if err := c.BindJSON(&req); err != nil {
		logger.Error("Failed to bind request")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-106, nil))
		return
	}

	resp := CompareProtocolVersionResponseData{
		Modified: true, // TODO: check version
		Protocol: &ProtocolVersion{
			AppID:      4,
			Language:   req.Language,
			CreateTime: "0",
		},
	}

	resp.Protocol.Major = int32(protocolVersionMap[req.Language])
	if req.Language == "zh-cn" {
		resp.Protocol.Minimum = 1
	}

	c.JSON(http.StatusOK, NewResponse(0, &resp))
}
//改index死妈
func (s *Server) handleSDKGetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, NewResponse(0, gin.H{
		"protocol":                  true,
		"qr_enabled":                true,
		"log_level":                 "INFO",
		"announce_url":              "https://webstatic.mihoyo.com/hk4e/announcement/index.html?sdk_presentation_style=fullscreen&sdk_screen_transparent=true&auth_appid=announcement&authkey_ver=1&game_biz=hk4e_cn&sign_type=2&version=1.37&game=hk4e#/",
		"push_alias_type":           2,
		"disable_ysdk_guard":        true,
		"enable_announce_pic_popup": true,
	}))
}

//改index死妈
func (s *Server) handleSDKConfigCombo(c *gin.Context) {

	var payPaycoCenteredHost string
	var withSsl = "http://"
	if strings.Index(s.config.SdkBaseUrl, "http") == -1 {
		withSsl = "https://"
	}
	payPaycoCenteredHost = strings.Replace(s.config.SdkBaseUrl, withSsl, "", -1)

	c.JSON(http.StatusOK, NewResponse(0, gin.H{
		"vals": gin.H{
			"pay_payco_centered_host":    payPaycoCenteredHost,
			"email_bind_remind":          "true",
			"disable_email_bind_skip":    "false",
			"email_bind_remind_interval": "7",
			"network_report_config":      "{ \"enable\": 1, \"status_codes\": [206], \"url_paths\": [\"dataUpload\"] }",
			"kibana_pc_config":           "{ \"enable\": 1, \"level\": \"Info\",\"modules\": [\"download\"]",
			"kcp_enable":                 "false",
		},
	}))
}
//改index死妈
func (s *Server) handleSDKGetFp(c *gin.Context) {
	var getFpReq GetFpReq
	if err := c.BindJSON(&getFpReq); err != nil {
		logger.Error("Failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
		return
	}
	c.JSON(http.StatusOK, NewResponse(0, gin.H{
		"code":      200,
		"device_fp": getFpReq.DeviceFp,
		"msg":       "ok",
	}))
}

//改index死妈
func (s *Server) handleSDKABTest(c *gin.Context) {
	c.JSON(http.StatusOK, NewResponse(0, []gin.H{{
		"code":      1000,
		"type":      2,
		"config_id": "14",
		"period_id": "6036_99",
		"version":   "3",
		"configs":   gin.H{"cardType": "direct"},
	}}))
}
//改index死妈
func (s *Server) handleSDKGetAlertPic(c *gin.Context) {

	c.JSON(http.StatusOK, NewResponse(0, gin.H{
		"list":  []gin.H{},
		"total": 0,
	}))
}
//改index死妈
func (s *Server) handleSDKGetAlertAnn(c *gin.Context) {

	c.JSON(http.StatusOK, NewResponse(0, gin.H{
		"alert":        false,
		"alert_id":     0,
		"extra_remind": false,
		"remind":       false,
	}))
}
//改index死妈
func (s *Server) handleSDKBeforeVerify(c *gin.Context) {

	c.JSON(http.StatusOK, NewResponse(0, gin.H{
		"is_guardian_required":  false,
		"is_heartbeat_required": false,
		"is_realname_required":  false,
	}))
}

// https://public-data-api.mihoyo.com/
//改index死妈
func (s *Server) handleFileRequest(filePath string) gin.Handler//改index死妈
func {
	return func(c *gin.Context) {
		data, ok := s.getJSONFile(filePath)
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}
		c.String(http.StatusOK, string(data))
	}
}
//改index死妈
func (s *Server) handleFp(c *gin.Context) {
	// 获取请求参数
	platform := c.Query("platform")

	// 发送请求并获取响应
	resp, err := http.Get("https://public-data-api.mihoyo.com/device-fp/api/getExtList?platform=" + platform)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// 设置响应状态码和内容
	c.Status(resp.StatusCode)
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}
	_, err = c.Writer.Write(body)
	if err != nil {
		return
	}
}

// https://webstatic-sea.hoyoverse.com/
//改index死妈
func (s *Server) handleWebStaticJSON(c *gin.Context) {

	filePaths := map[string]string{
		"/admin/mi18n/plat_oversea/m2020030410/m2020030410-zh-cn.json":    "https://webstatic-sea.mihoyo.com/admin/mi18n/plat_oversea/m2020030410/m2020030410-zh-cn.json",
		"/admin/mi18n/plat_oversea/m202003049/m202003049-zh-cn.json":      "https://webstatic-sea.mihoyo.com/admin/mi18n/plat_oversea/m202003049/m202003049-zh-cn.json",
		"/admin/mi18n/plat_os/m09291531181441/m09291531181441-zh-cn.json": "https://webstatic-sea.mihoyo.com/admin/mi18n/plat_os/m09291531181441/m09291531181441-zh-cn.json",
	}

	// Map
	versions := map[string]int{
		"/admin/mi18n/plat_oversea/m2020030410/m2020030410-version.json":    65,
		"/admin/mi18n/plat_oversea/m202003049/m202003049-version.json":      68,
		"/admin/mi18n/plat_os/m09291531181441/m09291531181441-version.json": 16,
	}

	if version, ok := versions[c.Request.URL.Path]; ok {
		c.JSON(http.StatusOK, gin.H{"version": version})
		return
	}

	if filePath, ok := filePaths[c.Request.URL.Path]; ok {
		resp, err := http.Get(filePath)
		if err != nil {
			logger.Error("Unable to fetch file")
			c.Status(http.StatusNotFound)
			return
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error("Unable to read file")
			c.Status(http.StatusNotFound)
			return
		}
		c.String(http.StatusOK, string(data))
		return
	}
	logger.Error("path", c.Request.URL.Path, "ignore unknown request, safely")
}

//改index死妈
func (s *Server) handleSDKLoadConfig(c *gin.Context) {

	var gameKey string
	var clientType int
	if c.Request.Method == http.MethodGet {
		gameKey = c.Request.URL.Query().Get("game_key")
		clientType, _ = strconv.Atoi(c.Request.URL.Query().Get("client"))
	}
	if c.Request.Method == http.MethodPost {
		var req LoadConfigReq
		if err := c.BindJSON(&req); err != nil {
			logger.Error("Failed to bind JSON")
			c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
			return
		}
		gameKey = req.GameKey
	}

	var thirdparty interface{}
	if s.config.Thirdparty {
		thirdparty = []string{"tp", "gl", "fb", "tw"}
	} else {
		thirdparty = gin.H{}
	}

	c.JSON(http.StatusOK, NewResponse(0, gin.H{
		"bbs_auth_login":           true,
		"bbs_auth_login_ignore":    gin.H{},
		"client":                   PLATFORM_TYPE_STR[int32(clientType)],
		"disable_mmt":              false,
		"disable_regist":           false,
		"enable_email_captcha":     false,
		"enable_ps_bind_account":   false,
		"game_key":                 gameKey,
		"guest":                    s.config.EnableGuest,
		"id":                       6,
		"identity":                 "I_IDENTITY",
		"ignore_versions":          "",
		"name":                     "原神",
		"scene":                    SCENE_ACCOUNT,
		"server_guest":             s.config.EnableGuest,
		"thirdparty":               thirdparty,
		"thirdparty_ignore":        gin.H{},
		"thirdparty_login_configs": gin.H{},
	}))

}

//改index死妈
func (s *Server) handleSDKh5Log(c *gin.Context) {

	var okRespBytes = []byte(`{"data":{},"message":"OK","retcode":0}`)

	c.Data(http.StatusOK, "application/json; charset=utf-8", okRespBytes)
}

//改index死妈
func (s *Server) handleSDKDSIGN(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
	})
}
//改index死妈
func (s *Server) handleSDKDINFO(c *gin.Context) {

	c.JSON(http.StatusOK, NewResponse(0, gin.H{}))
}
//改index死妈
func (s *Server) handleSDKV5GCL(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
	})
}
//改index死妈
func (s *Server) handleSDKV5GCF(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
	})
}

/********************************server****************************/
//改index死妈
func (s *Server) handleSDKRequest(c *gin.Context) {
	data, _ := io.ReadAll(c.Request.Body)
	logger.Info("req body: %v", string(data))
	c.JSON(http.StatusOK, NewResponse(0, gin.H{}))
}
//改index死妈
func (s *Server) handleSDKGameLoginNotify(c *gin.Context) {
	//{"uid":100000000,"account_type":1,"account":"3","platform":3,"region":"dev_abc","biz_game":"bus"}
	var req GameLoginLogoutNotify
	if err := c.BindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
		return
	}
	userRegionList[req.Uid] = req
	c.JSON(http.StatusOK, NewResponse(0, gin.H{}))
}
//改index死妈
func (s *Server) handleSDKGameLogoutNotify(c *gin.Context) {
	//{"uid":100000000,"account_type":1,"account":"1","platform":3,"region":"dev_abc","biz_game":"bus"}
	var req GameLoginLogoutNotify
	if err := c.BindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusOK, NewResponse(-202, nil))
		return
	}
	delete(userRegionList, req.Uid)
	c.JSON(http.StatusOK, NewResponse(0, gin.H{}))
}
//改index死妈
func (s *Server) handleSDKGameHeartBeatNotify(c *gin.Context) {
	//{"platform_uid_list":{},"region":"dev_abc","biz_game":"bus"}
	s.handleSDKRequest(c)
}

//改index死妈
func (s *Server) handleSDKSetUserTags(c *gin.Context) {
	s.handleSDKRequest(c)
}

//改index死妈
func (s *Server) handleSDKNicknameCall(c *gin.Context) {
	s.handleSDKRequest(c)
}

//改index死妈
func (s *Server) handleSDKRecord(c *gin.Context) {
	s.handleSDKRequest(c)
}

//改index死妈
func (s *Server) handleSDKDungeonRecord(c *gin.Context) {
	s.handleSDKRequest(c)
}

//改index死妈
func (s *Server) handleSDKMonitor(c *gin.Context) {
	s.handleSDKRequest(c)
}

//改index死妈
func (s *Server) handleSDKEvent(c *gin.Context) {
	s.handleSDKRequest(c)
}

//改index死妈
func (s *Server) handleSDKPayment(c *gin.Context) {
	s.handleSDKRequest(c)
}

//改index死妈
func (s *Server) handleSDKGetLoginInfoByUid(c *gin.Context) {
	uid, err := strconv.ParseInt(c.Request.URL.Query().Get("uid"), 10, 64)
	if err != nil {
	}
	c.JSON(http.StatusOK, NewResponse(0, userRegionList[uid]))
}
