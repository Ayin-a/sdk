package http

// 改index死妈
func (s *Server) initRouter() {
	s.router.Use(clientIPMiddleware())

	//改index死妈
	s.router.Any("/", s.handleDefault)
	s.router.Any("/index.html", s.handleDefault)
	s.router.GET("/mihoyo/common/accountSystemSandboxFE/index.html", s.handleViewAccount)
	s.router.GET("/sdkFacebookLogin.html", s.handleViewAccount)         //第三方登录
	s.router.GET("/sdkTwitterLogin.html", s.handleViewAccount)          //第三方登录
	s.router.GET("/sandbox/sdkFacebookLogin.html", s.handleViewAccount) //第三方登录
	s.router.GET("/sandbox/sdkTwitterLogin.html", s.handleViewAccount)  //第三方登录
	s.router.GET("/sandbox/index.html", s.handleViewAccount)
	s.router.Any("/:game_biz/combo/panda/qrcode/fetch", s.handleSDKScanLoginConfig) ////扫码登录配置
	s.router.Any("/:game_biz/combo/panda/qrcode/query", s.handleSDKScanStatusCheck) //扫码登录轮询
	s.router.Any("/:game_biz/combo/panda/qrcode/login", s.handleSDKQRLoginCheck)    //扫码登录验证
	s.router.POST("/:game_biz/mdk/shield/api/loginMobile", s.handleSDKLoginMobile)
	s.router.POST("/:game_biz/mdk/shield/api/loginCaptcha", s.handleSDKSendSMS) //短信登录
	s.router.Any("/:game_biz/mdk/guest/guest/v2/login", s.handleSDKGuestLogin)  // 游客登录
	s.router.POST("/account/risky/api/check", s.handleSDKRiskyCheck)
	s.router.POST("/account/device/api/listNewerDevices")
	s.router.POST("/account/auth/api/bindRealname", s.handleSDKRealName)         //实名认证
	s.router.POST("/:game_biz/mdk/shield/api/bindRealname", s.handleSDKRealName) //实名认证
	s.router.Any("/combo/postman/device/setUserTags", s.handleSDKSetUserTags)
	s.router.Any("/bat/game/gameLoginNotify", s.handleSDKGameLoginNotify)           //登录通知
	s.router.Any("/bat/game/gameLogoutNotify", s.handleSDKGameLogoutNotify)         //登出通知
	s.router.Any("/bat/game/gameHeartBeatNotify", s.handleSDKGameHeartBeatNotify)   //游戏心跳包通知
	s.router.Any("/outer_api/Outer/Call", s.handleSDKNicknameCall)                  //昵称审核
	s.router.Any("/hk4e_homeland/OuterApi/Record", s.handleSDKRecord)               //HOME验证
	s.router.Any("/hk4e_homeland/OuterApi/DungeonRecord", s.handleSDKDungeonRecord) //地牢记录
	s.router.Any("/hk4e/monitor", s.handleSDKMonitor)                               //输出控制监测
	s.router.Any("/event", s.handleSDKEvent)                                        //事件调整跟踪配置
	s.router.POST("/2g/dataUpload", s.handleLogUpload)
	s.router.Any("/perf_report_config/config/verify", s.handleSdkEmptySDKRspJson)
	s.router.POST("/perf/dataUpload", s.handleLogUpload)
	s.router.POST("/sdk/dataUpload", s.handleSdkDataUpload) //
	s.router.POST("/client/event/dataUpload", s.handleSdkDataUpload)
	s.router.POST("/crash/dataUpload", s.handleCrashDataUpload) //
	s.router.POST("/crashdump/dataUpload", s.handleCrashDataUpload)
	s.router.POST("/apm/dataUpload", s.handleLogUpload)
	s.router.POST("/log/sdk/upload", s.handleLogUpload)
	s.router.POST("/perf/config/verify", s.handleLogUpload)
	s.router.GET("/combo/box/api/config/sdk/combo", s.handleSDKConfigCombo)
	s.router.Any("/ys/event/e20210830cloud/index.html", s.handleDefault) //云原神
	s.router.GET("/admin/mi18n/plat_oversea/*any", s.handleWebStaticJSON)
	s.router.GET("/admin/mi18n/plat_os/*any", s.handleWebStaticJSON)
	s.router.GET("/hk4e/announcement/index.html", s.handleViewAnnouncement) //公告模板地址
	s.router.GET("/combo/box/api/config/sw/precache", s.handleSdkEmptySDKRspJson)
	s.router.Any("/common/h5log/log/batch", s.handleSDKh5Log)
	s.router.Any("/privacy/policy/authorization/status", s.handleSdkAuthStatus)
	s.router.Any("/dsign", s.handleSDKDSIGN)
	s.router.Any("/dinfo", s.handleSDKDINFO)
	s.router.Any("/v5/gcl", s.handleSDKV5GCL)
	s.router.Any("/v5/gcf", s.handleSDKV5GCF)
	s.router.POST("/:game_biz/mdk/shield/api/emailCaptcha")
	s.router.POST("/:game_biz/mdk/shield/api/login", s.handleSDKShieldLogin)
	s.router.POST("/:game_biz/mdk/shield/api/verify", s.handleSDKShieldVerify) //token account_id 验证
	s.router.POST("/:game_biz/mdk/shield/api/loginByThirdparty", s.handleSDKLoginByThirdparty)
	s.router.POST("/:game_biz/combo/granter/login/beforeVerify", s.handleSDKBeforeVerify) //在之前验证
	s.router.POST("/:game_biz/combo/granter/login/v2/login", s.handleSDKComboLogin)
	s.router.GET("/query_region_list", s.handleQueryRegionList())
	s.router.POST("/log", s.handleNLogUpload) //
	s.router.GET("/common/:game_biz/announcement/api/getAlertPic", s.handleSDKGetAlertPic)
	s.router.GET("/common/:game_biz/announcement/api/getAlertAnn", s.handleSDKGetAlertAnn)
	s.router.GET("/common/:game_biz/announcement/api/getAnnList", s.handleFileRequest("./data/json/ann_list.json")) //公告列表
	s.router.GET("/common/:game_biz/announcement/api/getAnnContent", s.handleFileRequest("./data/json/ann_content.json"))
	s.router.GET("/:game_biz/mdk/agreement/api/getAgreementInfos", s.handleSDKGetAgreementInfos)
	s.router.Any("/:game_biz/mdk/shopwindow/shopwindow/listPriceTier", s.handleSDKGetShopPriceTier)
	s.router.Any("/:game_biz/mdk/shopwindow/shopwindow/listPriceTierV2", s.handleSDKGetShopPriceTier)
	s.router.Any("/:game_biz/combo/granter/api/compareProtocolVersion", s.handleSDKCompareProtocolVersion)
	s.router.POST("/:game_biz/mdk/shield/api/actionTicket", s.handleSDKActionTicket)      //实名认证
	s.router.POST("/:game_biz/mdk/luckycat/luckycat/createOrder", s.handleSDKCreateOrder) //创建订单
	s.router.POST("/:game_biz/mdk/luckycat/luckycat/checkOrder", s.handleSDKCheckOrder)   //支付状态检测
	s.router.Any("/:game_biz/combo/granter/api/getFont", s.handleSdkEmptySDKRspJson)
	s.router.Any("/:game_biz/combo/red_dot", s.handleSDKRedDotList)
	s.router.GET("/:game_biz/combo/granter/api/getConfig", s.handleSDKGetConfig)
	s.router.Any("/:game_biz/mdk/shield/api/loadConfig", s.handleSDKLoadConfig)
	s.router.POST("/account/auth/api/getConfig", s.handleFileRequest("./data/json/auth_config.json"))
	s.router.Any("/:game_biz/combo/postman/device/setAlias", s.handleSdkEmptySDKRspJson)
	s.router.POST("/data_abtest_api/config/experiment/list", s.handleSDKABTest)
	s.router.Any("/:game_biz/combo/red_dot/list", s.handleSDKRedDotList)
	s.router.Any("/receive/tkio/payment", s.handleSDKPayment)
	s.router.Any("/device-fp/api/getExtList", s.handleFp)
	s.router.Any("/device-fp/api/getFp", s.handleSDKGetFp)
	//改index死妈
	view := s.router.Group("/view")
	{
		view.Any("/qr_code_pay", s.handleViewQrCodePay) //支付回调
		view.Any("/qr_code_login", s.handleViewQrCodeLogin)

	}
	//改index死妈
	api := s.router.Group("/api")
	{
		api.Any("/reportQRScanned", s.handleSDKScanned)
		api.Any("/sdkUploadLoginToken", s.handleSDKUploadLoginToken)
		api.Any("/createAndFinishOrder", s.handleAPICreateAndFinishOrder)
		api.Any("/getLoginInfoByUid", s.handleSDKGetLoginInfoByUid)
		api.Any("/status", s.handleAPIStatus)
		api.Any("/restart", s.handleAPIRestart)

	}
}
