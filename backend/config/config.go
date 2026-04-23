package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Payment  PaymentConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type JWTConfig struct {
	Secret     string
	ExpireTime int64
}

type PaymentConfig struct {
	Alipay   AlipayConfig
	Wechat   WechatConfig
	PointsRate float64
}

type AlipayConfig struct {
	AppID      string
	PrivateKey string
	PublicKey  string
	NotifyURL  string
	ReturnURL  string
	Sandbox    bool
}

type WechatConfig struct {
	AppID      string
	MchID      string
	APIKey     string
	NotifyURL  string
	CertPath   string
	KeyPath    string
}

var AppConfig *Config

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using defaults")
	}

	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Driver: getEnv("DB_DRIVER", "sqlite"),
			DSN:    getEnv("DB_DSN", "token-hub.db"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "token-hub-secret-key-2024"),
			ExpireTime: getEnvAsInt64("JWT_EXPIRE_TIME", 86400),
		},
		Payment: PaymentConfig{
			Alipay: AlipayConfig{
				AppID:      getEnv("ALIPAY_APP_ID", ""),
				PrivateKey: getEnv("ALIPAY_PRIVATE_KEY", ""),
				PublicKey:  getEnv("ALIPAY_PUBLIC_KEY", ""),
				NotifyURL:  getEnv("ALIPAY_NOTIFY_URL", ""),
				ReturnURL:  getEnv("ALIPAY_RETURN_URL", ""),
				Sandbox:    getEnvAsBool("ALIPAY_SANDBOX", true),
			},
			Wechat: WechatConfig{
				AppID:     getEnv("WECHAT_APP_ID", ""),
				MchID:     getEnv("WECHAT_MCH_ID", ""),
				APIKey:    getEnv("WECHAT_API_KEY", ""),
				NotifyURL: getEnv("WECHAT_NOTIFY_URL", ""),
				CertPath:  getEnv("WECHAT_CERT_PATH", ""),
				KeyPath:   getEnv("WECHAT_KEY_PATH", ""),
			},
			PointsRate: getEnvAsFloat64("POINTS_RATE", 100.0),
		},
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsFloat64(key string, defaultValue float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}
