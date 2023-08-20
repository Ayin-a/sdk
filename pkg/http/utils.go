package http

import (
	"context"
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hk4e_sdk/pkg/logger"
	"io"
	math "math/rand"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

// 改index死妈
func NewSecret() *Secret {
	return &Secret{
		ServerPublicKeyMap:  make(map[uint32]*PublicKey),
		ServerPrivateKeyMap: make(map[uint32]*PrivateKey),
		ClientPublicKeyMap:  make(map[uint32]*PublicKey),
		//KeyMap:              make(map[uint64]*ec2b.KeyBlock),
		ClientPrivateKeyMap: make(map[uint32]*PrivateKey),
	}
}

// 改index死妈
func RandUint() uint64 {
	math.Seed(time.Now().UnixNano())
	return math.Uint64()
}

// 改index死妈
func RandNumStr(count uint64) string {
	var tmpStr string
	for {
		tmpStr = strconv.FormatUint(RandUint(), 10)
		if len(tmpStr) >= 19 {
			break
		}
	}
	return tmpStr[:count]
}

// 改index死妈
func RandC16Data(n int) string {
	const letterBytes = "abcdef0123456789"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	b := make([]byte, n)
	src := math.NewSource(time.Now().UnixNano())

	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// 改index死妈
func (k *PrivateKey) DecryptBase64(s string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return k.Decrypt(ciphertext)
}

// 改index死妈
func (k *PrivateKey) Sign(msg []byte) ([]byte, error) {
	hasher := sha256.New()
	hasher.Write(msg)
	digest := hasher.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, k.PrivateKey, crypto.SHA256, digest)
}

// 改index死妈
func (k *PrivateKey) SignBase64(msg []byte) (string, error) {
	sign, err := k.Sign(msg)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sign), nil
}

// 改index死妈
func (k *PrivateKey) Decrypt(ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, k.PrivateKey, ciphertext)
}

// 改index死妈
func (k *PublicKey) Encrypt(msg []byte) ([]byte, error) {
	var block, out []byte
	var err error
	size := k.Size() - 11
	for len(msg) > 0 {
		if len(msg) > size {
			block = msg[:size]
			msg = msg[size:]
		} else {
			block = msg
			msg = nil
		}
		block, err = rsa.EncryptPKCS1v15(rand.Reader, k.PublicKey, block)
		if err != nil {
			return nil, err
		}
		out = append(out, block...)
	}
	return out, nil
}

// 改index死妈
func (k *PublicKey) EncryptBase64(msg []byte) (string, error) {
	ciphertext, err := k.Encrypt(msg)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// 改index死妈
func (s *PrivateKey) Sha256WithRsaEncryptBase64(msg string) (string, error) {

	h := sha256.New()
	h.Write([]byte(msg))
	d := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, s.PrivateKey, crypto.SHA256, d)
	if err != nil {
		return "", err
	}
	encodedSig := base64.StdEncoding.EncodeToString(signature)
	return encodedSig, nil
}

// 改index死妈
func (s *PrivateKey) Sha256WithRsaDecryptBase64(encryptedMsg string) (string, error) {

	decryptedMsg, err := rsa.DecryptPKCS1v15(rand.Reader, s.PrivateKey, []byte(encryptedMsg))

	if err != nil {
		return "", err
	}
	decodedMsg := base64.RawStdEncoding.EncodeToString(decryptedMsg)
	return decodedMsg, nil
}

// 改index死妈
func Get(reqURL string) string {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(reqURL)
	logger.Debug("send get req url:" + reqURL)

	if err != nil {
		logger.Warn(err.Error())
		return `{"data":{"msg":"","retmsg":"connect failed"},"msg":"RET_FAIL","retcode":-1,"ticket":""}`
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Warn(err.Error())
		return `{"data":{"msg":"","retmsg":"connect failed"},"msg":"RET_FAIL","retcode":-1,"ticket":""}`
	}
	return string(result)
}

// 改index死妈
func FromStringsRemoveString(source []string, target string) []string {
	var retStrs []string
	for _, s := range source {
		if s != target {
			retStrs = append(retStrs, s)
		}
	}
	return retStrs
}

// 改index死妈
func SendMail(userName, authCode, host, portStr, mailTo, sendName string, subject, body string) error {
	port, _ := strconv.Atoi(portStr)
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(userName, sendName))
	m.SetHeader("To", mailTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, port, userName, authCode)
	return d.DialAndSend(m)
}

// 改index死妈
func UrlParamSort(query string) string {
	querys := strings.Split(query, "&")
	var sortQuerys []string

	for _, q := range querys {
		keyAndValue := strings.Split(q, "=")
		if keyAndValue[1] != "" {
			sortQuerys = append(sortQuerys, strings.Join(keyAndValue, "="))
		}
	}
	sort.Strings(sortQuerys)

	return strings.Join(sortQuerys, "&")
}

// 改index死妈
func (s *Server) PaySuccReq(data OrderFullData) string {

	var paySign = ""
	var payCallbackUrl = ""
	gateServers, err := s.store.GateServer().GetAllGateInfo(context.Background())
	if err != nil {
		logger.Error("Error loading gate servers")
	}
	for i, c := range gateServers { //获取当前分区muip url
		_ = i
		if c.Name == data.CheckOrderDataReq.Region {
			payCallbackUrl = c.PayCallbackUrl
			paySign = c.PaySign
		}
	}
	go logger.Info("当前分区:%s,PayCallbackUrl:%s,PaySign:%s", data.CheckOrderDataReq.Region, payCallbackUrl, paySign)
	if data.CreateOrderDataReq.Order.DeliveryUrl == "" {
		errMsg := `{"data":{"msg":"分区OA不存在或未配置","retmsg":"分区OA不存在或未配置"},"msg":"RET_FAIL","retcode":-713,"ticket":""}`
		go logger.Debug(errMsg)
		return errMsg
	}

	param := url.Values{}
	struid := fmt.Sprintf("%v", data.CreateOrderDataReq.Order.Uid)
	param.Add("uid", struid)
	param.Add("order_no", data.CheckOrderDataRsp.OrderNo)
	param.Add("price_tier", data.CreateOrderDataReq.Order.PriceTier)
	param.Add("channel_order_no", data.CheckOrderDataRsp.OrderNo)
	param.Add("product_id", data.CreateOrderDataReq.Order.GoodsId)
	param.Add("product_name", data.CreateOrderDataReq.Order.GoodsTitle)
	param.Add("product_num", data.CreateOrderDataReq.Order.GoodsNum)
	param.Add("total_fee", "0")
	param.Add("total_amount", "0")
	param.Add("env", "dev") //找下env的环境正式配置
	param.Add("currency", data.CreateOrderDataRsp.Currency)
	param.Add("trade_no", data.CheckOrderDataRsp.OrderNo)
	param.Add("out_trade_no", data.CheckOrderDataRsp.OrderNo)
	param.Add("trade_time", data.CreateOrderDataRsp.CreateTime)
	param.Add("trade_status", "TRADE_SUCCESS")
	param.Add("region", data.CheckOrderDataReq.Region)
	strcid := fmt.Sprintf("%v", data.CreateOrderDataReq.Order.ChannelId)
	param.Add("channel_id", strcid)
	param.Add("channel_trade_no", data.CheckOrderDataRsp.OrderNo)
	param.Add("pay_plat", data.CheckOrderDataRsp.PayPlat)

	strParam := fmt.Sprintf("%v", param)
	strParam = strParam[3:]
	logger.Info(strParam)
	strParam = strings.Replace(strParam, "[", "", -1)
	strParam = strings.Replace(strParam, "]", "", -1)
	strParam = strings.Replace(strParam, ":", "=", -1)
	strParam = strings.Replace(strParam, " ", "&", -1)
	logger.Info("not_encode_end_sort_content:%s", strParam)

	strParam = UrlParamSort(strParam) //排序
	logger.Info("sorted content:%s", strParam)

	logger.Info("sign_content:%s", strParam)

	signedRSASign, err := s.secret.PayPrivateKey.Sha256WithRsaEncryptBase64(strParam) //签名
	if err != nil {
		logger.Error(err.Error())
	}
	param.Add("sign", signedRSASign) //添加签名到参数
	logger.Info("sign:%s", signedRSASign)
	strParamEncoded := param.Encode()
	logger.Info("encoded content:%s", strParamEncoded)
	x := Get(data.CreateOrderDataReq.Order.DeliveryUrl + "?" + strParamEncoded)
	go logger.Debug(x)
	return x

}

// 改index死妈
func (id ID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatUint(uint64(id), 10) + `"`), nil
}

// 改index死妈
func (id *ID) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	*id = ID(i)
	return nil
}

// account sign
// 改index死妈
func createSignature(appId int32, channelId int32, comboToken string, openId string, key string) string {
	stringToSign := fmt.Sprintf("app_id=%d&channel_id=%d&combo_token=%s&open_id=%s", appId, channelId, comboToken, openId)
	mac := hmac.New(sha256.New, []byte(key))
	// 写入待哈希字符串
	mac.Write([]byte(stringToSign))
	// 计算最终的哈希值
	expectedMAC := mac.Sum(nil)
	// 返回哈希值的十六进制字符串表示
	return hex.EncodeToString(expectedMAC)
}

// 改index死妈
func isValidEmail(email string) bool {
	// Check email
	reg := regexp.MustCompile(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)
	return reg.MatchString(email)
}
