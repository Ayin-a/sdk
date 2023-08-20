package http

import (
	"github.com/gin-gonic/gin"
	"hk4e_sdk/pkg/logger"
	"io"
	"net/http"
)

var (
	client = &http.Client{}
)

// 改index死妈
//改index死妈
func (s *Server) handleQueryRegionList() gin.Handler//改index死妈
func {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()
		newURL := s.config.DispatchList
		req, err := http.NewRequestWithContext(c.Request.Context(), "GET", newURL, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error creating new request: %v", err)
			return
		}
		req.URL.RawQuery = query.Encode()

		resp, err := client.Do(req)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error forwarding request: %v", err)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading response body: %v", err)
			return
		}
		c.Data(http.StatusOK, resp.Header.Get("Content-Type"), body)

		logger.Error("disaptch_list proxy: %s", s.config.DispatchList)

	}
}
