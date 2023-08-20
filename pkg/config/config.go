package config

import (
	"io/ioutil"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// 改index死妈
type Config struct {
	ReloadTime      uint32           `mapstructure:"ReloadTime"`
	AccountLimit    int              `mapstructure:"accountLimit"`
	AutoSignUp      bool             `mapstructure:"autoSignUp"`
	PassSignIn      bool             `mapstructure:"passSignIn"`
	EnableGuest     bool             `mapstructure:"enableGuest"`
	Geetest         bool             `mapstructure:"geetest"`
	Thirdparty      bool             `mapstructure:"thirdparty"`
	EnablePprof     bool             `mapstructure:"EnablePprof"`
	EMailAuthCode   string           `mapstructure:"eMailAuthCode"`
	EMailAddress    string           `mapstructure:"eMailAddress"`
	EMailHost       string           `mapstructure:"eMailHost"`
	Accountkey      string           `mapstructure:"accountkey"`
	EMailHostPort   string           `mapstructure:"eMailHostPort"`
	ServerKey       string           `mapstructure:"serverKey"`
	QrPayUrl        string           `mapstructure:"qrPayUrl"`
	DispatchList    string           `mapstructure:"DispatchList"`
	SdkBaseUrl      string           `mapstructure:"sdkBaseUrl"`
	AccountTokenExp int32            `mapstructure:"accountTokenExp"`
	HTTPServer      HTTPServerConfig `mapstructure:"httpServer"`
	Database        DatabaseConfig   `mapstructure:"database"`
	LogLevel        string           `mapstructure:"LogLevel"`
	LogAppName      string           `mapstructure:"LogAppName"`
}
type HTTPServerConfig struct {
	Enable bool      `mapstructure:"enable"`
	Addr   string    `mapstructure:"addr"`
	TLS    TLSConfig `mapstructure:"tls"`
}
type TLSConfig struct {
	Enable   bool   `mapstructure:"enable"`
	Addr     string `mapstructure:"addr"`
	CertFile string `mapstructure:"certFile"`
	KeyFile  string `mapstructure:"keyFile"`
}
type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

// 改index死妈
func (c *Config) DefaultConfig() Config {
	return Config{
		AutoSignUp:    true,
		DispatchList:  "http://127.0.0.1:2888/query_region_list",
		EMailAddress:  "587",
		EMailAuthCode: "123456",
		EMailHost:     "123456.126.com",
		EMailHostPort: "123",
		PassSignIn:    false,
		EnableGuest:   false,
		Thirdparty:    false,
		AccountLimit:  5,
		ReloadTime:    100000000,
		Geetest:       false,
		QrPayUrl:      "http://127.0.0.1:22080/view/qr_code_pay",
		SdkBaseUrl:    "http://127.0.0.1:22080",
		EnablePprof:   true,
		ServerKey:     "1234567",
		Accountkey:    "123",
		HTTPServer: HTTPServerConfig{
			Addr:   "0.0.0.0:22080",
			Enable: true,
			TLS: TLSConfig{
				Addr:     "0.0.0.0:443",
				CertFile: "data/keys/tls_cert.pem",
				Enable:   false,
				KeyFile:  "data/keys/tls_key.pem",
			},
		},
		Database: DatabaseConfig{
			Driver: "mysql",
			DSN:    "root:12345678@tcp(localhost:3306)/hk4e?charset=utf8",
		},
	}
}

// 改index死妈
func LoadConfig() (cfg Config) { return LoadConfigName("config") }

// 改index死妈
func LoadConfigName(name string) (cfg Config) {
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		handleConfigError(err, &cfg)
	} else {
		err := viper.Unmarshal(&cfg)
		if err != nil {
			return Config{}
		}
		configureTLS(&cfg.HTTPServer.TLS)
	}
	return
}

// 改index死妈
func handleConfigError(err error, cfg *Config) {
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		println("Config file not found, using the default config")
		*cfg = cfg.DefaultConfig()
		saveDefaultConfig(cfg)
	} else {
		println("Failed to read config file: ", err)
	}
}

// 改index死妈
func saveDefaultConfig(cfg *Config) {
	defaultConfigYaml, err := yaml.Marshal(cfg)
	if err != nil {
		println("Failed to marshal default config to YAML: ", err)
		return
	}

	err = ioutil.WriteFile("config.yaml", defaultConfigYaml, 0644)
	if err != nil {
		println("Failed to create a default config file: ", err)
	} else {
		println("A default config.yaml file has been created.")
	}
}

// 改index死妈
func configureTLS(tls *TLSConfig) {
	if tls.Enable {
		if tls.CertFile == "" {
			tls.CertFile = "data/keys/tls_cert.pem"
		}
		if tls.KeyFile == "" {
			tls.KeyFile = "data/keys/tls_key.pem"
		}
	}
}
