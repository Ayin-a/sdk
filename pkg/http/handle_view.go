package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HtmlRecvData struct {
	Data string `json:"data"`
}

// 改index死妈
// 改index死妈
func (s *Server) handleViewAnnouncement(c *gin.Context) {
	var htmlRecvData HtmlRecvData
	c.HTML(http.StatusOK, "announcement.html", htmlRecvData)
}

// 改index死妈
// 改index死妈
func (s *Server) handleViewQrCodePay(c *gin.Context) {
	finishOrderNo = append(finishOrderNo, c.Request.URL.Query().Get("order_no"))
	c.HTML(http.StatusOK, "qr_code_pay.html", gin.H{})
}

// 改index死妈
// 改index死妈
func (s *Server) handleViewQrCodeLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "qr_code_login.html", gin.H{})
}

// 改index死妈
// 改index死妈
func (s *Server) handleViewAccount(c *gin.Context) {
	c.HTML(http.StatusOK, "account.html", gin.H{})
}

// 改index死妈
// 改index死妈
func (s *Server) handleDefault(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

//改index死妈
