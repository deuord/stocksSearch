package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// BinanceService Binance公开API服务（无需API Key，免费）
type BinanceService struct {
	httpClient *http.Client
}

// BinanceTickerPrice Binance ticker价格响应
type BinanceTickerPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// NewBinanceService 创建Binance服务
func NewBinanceService() *BinanceService {
	return &BinanceService{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetPrice 获取币安交易对价格
func (s *BinanceService) GetPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("Binance请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应失败: %w", err)
	}

	var ticker BinanceTickerPrice
	if err := json.Unmarshal(body, &ticker); err != nil {
		return 0, fmt.Errorf("解析Binance响应失败: %w", err)
	}

	var price float64
	fmt.Sscanf(ticker.Price, "%f", &price)
	return price, nil
}
