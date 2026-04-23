package repository

import (
	"log"
	"token-hub/config"
	"token-hub/internal/model"

	sqlite "github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	var dialector gorm.Dialector

	switch config.AppConfig.Database.Driver {
	case "sqlite":
		dialector = sqlite.Open(config.AppConfig.Database.DSN)
	default:
		dialector = sqlite.Open(config.AppConfig.Database.DSN)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = migrateSchema()
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	initDefaultData()

	log.Println("Database initialized successfully")
}

func migrateSchema() error {
	err := DB.AutoMigrate(
		&model.User{},
		&model.Provider{},
		&model.Model{},
		&model.PointsRecord{},
		&model.PointsConfig{},
		&model.PaymentMethod{},
		&model.RechargeRecord{},
		&model.CallLog{},
		&model.APIKey{},
	)
	if err != nil {
		return err
	}
	return nil
}

func initDefaultData() {
	var adminCount int64
	DB.Model(&model.User{}).Where("role = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		admin := &model.User{
			Username: "admin",
			Password: "e10adc3949ba59abbe56e057f20f883e",
			Email:    "admin@token-hub.com",
			Nickname: "管理员",
			Role:     "admin",
			Status:   1,
		}
		DB.Create(admin)
		log.Println("Default admin user created: username=admin, password=123456")
	}

	var providerCount int64
	DB.Model(&model.Provider{}).Count(&providerCount)
	if providerCount == 0 {
		providers := []model.Provider{
			{
				Name:        "OpenAI",
				Code:        "openai",
				Description: "OpenAI GPT系列模型",
				Website:     "https://openai.com",
				APIEndpoint: "https://api.openai.com",
				Status:      1,
				Sort:        1,
			},
			{
				Name:        "Anthropic",
				Code:        "anthropic",
				Description: "Anthropic Claude系列模型",
				Website:     "https://anthropic.com",
				APIEndpoint: "https://api.anthropic.com",
				Status:      1,
				Sort:        2,
			},
			{
				Name:        "智谱AI",
				Code:        "zhipu",
				Description: "智谱AI GLM系列模型",
				Website:     "https://bigmodel.cn",
				APIEndpoint: "https://open.bigmodel.cn",
				Status:      1,
				Sort:        3,
			},
			{
				Name:        "阿里千问",
				Code:        "qwen",
				Description: "阿里云千问系列模型",
				Website:     "https://dashscope.console.aliyun.com",
				APIEndpoint: "https://dashscope.aliyuncs.com",
				Status:      1,
				Sort:        4,
			},
			{
				Name:        "腾讯混元",
				Code:        "hunyuan",
				Description: "腾讯混元大模型",
				Website:     "https://cloud.tencent.com",
				APIEndpoint: "https://hunyuan.tencentcloudapi.com",
				Status:      1,
				Sort:        5,
			},
			{
				Name:        "字节跳动豆包",
				Code:        "doubao",
				Description: "字节跳动豆包大模型",
				Website:     "https://www.volcengine.com",
				APIEndpoint: "https://ark.cn-beijing.volces.com",
				Status:      1,
				Sort:        6,
			},
		}
		DB.Create(&providers)
		log.Println("Default providers created")
	}

	var paymentCount int64
	DB.Model(&model.PaymentMethod{}).Count(&paymentCount)
	if paymentCount == 0 {
		payments := []model.PaymentMethod{
			{
				Name:        "支付宝",
				Code:        "alipay",
				Description: "支付宝在线支付",
				Status:      1,
				Sort:        1,
			},
			{
				Name:        "微信支付",
				Code:        "wechat",
				Description: "微信在线支付",
				Status:      1,
				Sort:        2,
			},
		}
		DB.Create(&payments)
		log.Println("Default payment methods created")
	}

	var pointsConfigCount int64
	DB.Model(&model.PointsConfig{}).Count(&pointsConfigCount)
	if pointsConfigCount == 0 {
		configs := []model.PointsConfig{
			{
				Key:         "points_rate",
				Value:       "100",
				Description: "积分汇率：1元兑换多少积分",
			},
		}
		DB.Create(&configs)
		log.Println("Default points config created")
	}
}
