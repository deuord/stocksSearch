package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"stocksSearch/backend/config"
	"stocksSearch/backend/controllers"
	"stocksSearch/backend/router"
	"stocksSearch/backend/services"
)

func main() {
	// 加载配置
	cfg := loadConfig()

	// 初始化服务层
	itickService := services.NewITickService(cfg.ITick)
	stockService := services.NewStockService(itickService)

	// 初始化控制器
	stockController := controllers.NewStockController(stockService)

	// 设置路由
	handler := router.SetupRouter(stockController)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("🚀 全球股市情报面板启动：http://%s", addr)
	log.Printf("📡 iTick API Base: %s", cfg.ITick.BaseURL)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// loadConfig 加载配置，优先从环境变量或配置文件读取
func loadConfig() *config.Config {
	// 尝试从配置文件加载
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config/config.json"
	}

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Printf("警告: 无法加载配置文件 %s，使用默认配置: %v", cfgPath, err)
		cfg = config.DefaultConfig()
	}

	// 环境变量覆盖
	if envPort := os.Getenv("PORT"); envPort != "" {
		cfg.Server.Port = envPort
	}
	if envHost := os.Getenv("HOST"); envHost != "" {
		cfg.Server.Host = envHost
	}
	if envToken := os.Getenv("ITICK_TOKEN"); envToken != "" {
		cfg.ITick.Token = envToken
	}
	if envBaseURL := os.Getenv("ITICK_BASE_URL"); envBaseURL != "" {
		cfg.ITick.BaseURL = envBaseURL
	}

	return cfg
}
