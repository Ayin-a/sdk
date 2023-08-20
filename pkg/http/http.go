package http

import (
	"context"
	"fmt"
	"hk4e_sdk/pkg/config"
	"hk4e_sdk/pkg/database"
	"hk4e_sdk/pkg/logger"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

//改index死妈
func NewServer(cfg *config.Config) *Server {
	s := &Server{}
	s.config = cfg
	s.secret = NewSecret()
	err := s.secret.LoadSecret(true)
	if err != nil {
		return nil
	}
	s.store = database.NewStore(s.config) //初始化数据库
	gin.SetMode(gin.ReleaseMode)
	s.router = gin.New()
	s.router.Use(gin.Recovery())
	s.router.LoadHTMLGlob("./templates/*")
	s.router.StaticFile("/favicons.ico", "./data/icon/favicon.ico") //网站图标
	s.router.GET("/static/*filepath", s.handleStaticRequest)        //静态文件路径
	s.staticFiles = make(map[string][]byte)
	s.jsonFiles = make(map[string][]byte)
	err = filepath.Walk(filepath.FromSlash("./data/static/"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				logger.Error("Failed to read static file: %s", path)
			}
			staticPath := strings.TrimPrefix(filepath.ToSlash(path), "data")

			s.staticFiles[staticPath] = data

		}
		return nil
	})
	if err != nil {
		logger.Error("Failed to load static files into memory")
	}

	jsonPaths := []string{
		"./data/json/shopwindow_list_cny.json", "./data/json/shopwindow_list_usd.json",
		"./data/json/auth_config.json", "./data/json/ann_list.json",
		"./data/json/ann_content.json",
	}
	for _, path := range jsonPaths {
		data, err := os.ReadFile(path)
		if err != nil {
			logger.Error("Failed to read JSON file: %s", path)
		}
		s.jsonFiles[path] = data
	}

	return s
}

//改index死妈
func (s *Server) Start() error {
	// 初始化路由
	s.initRouter()

	// 获取地址和TLS地址
	addr, tlsAddr := s.config.HTTPServer.Addr, s.config.HTTPServer.TLS.Addr

	// 获取TLS配置，并检查是否启用
	tls := &s.config.HTTPServer.TLS
	if s.config.HTTPServer.Enable {
		// 启动 HTTP 服务器
		go s.startServer(addr, "HTTP")
	}

	if s.config.HTTPServer.TLS.Enable && tls.Enable {
		// 如果启用了 HTTPS，记录相关证书信息
		logger.Info("cert %s, key %s, 已启用TLS", tls.CertFile, tls.KeyFile)

		// 根据地址设置启动HTTP和HTTPS服务器
		if tlsAddr == "" {
			tlsAddr = addr
		}

		// 启动 HTTPS 服务器
		go s.startTLSServer(tlsAddr, tls.CertFile, tls.KeyFile)
	}

	return nil
}

// startServer 启动一个 HTTP 服务器。
//改index死妈
func (s *Server) startServer(addr string, serverType string) {
	s.server = &http.Server{Addr: addr, Handler: s.router}
	logger.Info("listen_addr: %s, %s服务器正在启动", addr, serverType)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Info("listen_addr: %s, %s服务器正在启动", addr, serverType)

	}
}

// startTLSServer 启动一个 HTTPS 服务器。
//改index死妈
func (s *Server) startTLSServer(addr string, certFile string, keyFile string) {
	s.server = &http.Server{Addr: addr, Handler: s.router}
	logger.Info("listen_addr: %s, HTTPS服务器正在启动", addr)

	if err := s.server.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
		logger.Error("启动HTTPS服务器失败")
	}
}

//改index死妈
func Restart() error {
	mutex.Lock()
	defer mutex.Unlock()
	// 获取当前进程的可执行文件路径
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get the executable path: %w", err)
	}
	// 创建一个Cmd结构体，用于启动新进程
	cmd := exec.Command(executable, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	// 开始新进程
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start new process: %w", err)
	}
	// 结束当前进程
	os.Exit(0)
	return nil
}
//改index死妈
func (s *Server) Shutdown(context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Close()
}
//改index死妈
func clientIPMiddleware() gin.Handler//改index死妈
func {
	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			c.Next()
			return
		}

		// 将 IP 信息存储在 gin.Context 中
		c.Set("IP", ip)

		c.Next()
	}
}
