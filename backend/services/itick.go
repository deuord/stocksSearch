package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"stocksSearch/backend/config"
	"stocksSearch/backend/models"
)

// ITickService iTick API服务
type ITickService struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// iTickKlineItem kline K线数据项
type iTickKlineItem struct {
	Close  float64 `json:"c"`  // 收盘价
	High   float64 `json:"h"`  // 最高价
	Low    float64 `json:"l"`  // 最低价
	Open   float64 `json:"o"`  // 开盘价
	Volume float64 `json:"v"`  // 成交量
	Time   int64   `json:"t"`  // 时间戳(ms)
}

// iTickKlineResponse kline接口响应
type iTickKlineResponse struct {
	Code int              `json:"code"`
	Msg  *string          `json:"msg"`
	Data []iTickKlineItem `json:"data"`
}

// iTickSnapResponse snapshot接口响应（尝试使用，可能不可用）
type iTickSnapResponse struct {
	Code int             `json:"code"`
	Msg  *string         `json:"msg"`
	Data []iTickSnapItem `json:"data"`
}

// iTickSnapItem snapshot数据项
type iTickSnapItem struct {
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Open  float64 `json:"open"`
	PrevClose float64 `json:"prevClose"`
	Volume    float64 `json:"volume"`
	MarketCap float64 `json:"marketValue"`
}

// NewITickService 创建iTick服务实例
func NewITickService(cfg config.ITickConfig) *ITickService {
	return &ITickService{
		baseURL: cfg.BaseURL,
		token:   cfg.Token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchStocksByRegion 根据市场区域获取股票行情
// 使用kline接口并发获取最新价格数据
func (s *ITickService) FetchStocksByRegion(region string, codes []string) ([]models.Stock, error) {
	if len(codes) == 0 {
		return nil, fmt.Errorf("codes 不能为空")
	}

	var (
		stocks []models.Stock
		mu     sync.Mutex
		wg     sync.WaitGroup
	)

	// 限制并发数，避免触发API频率限制
	sem := make(chan struct{}, 5)

	for _, code := range codes {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			stock, err := s.fetchSingleStock(region, c)
			if err != nil {
				return
			}

			mu.Lock()
			stocks = append(stocks, stock)
			mu.Unlock()
		}(code)
	}

	wg.Wait()
	return stocks, nil
}

// fetchSingleStock 获取单只股票行情
func (s *ITickService) fetchSingleStock(region, code string) (models.Stock, error) {
	apiURL := fmt.Sprintf("%s/stock/kline", s.baseURL)
	params := url.Values{}
	params.Set("region", region)
	params.Set("code", code)
	params.Set("kType", "1") // 1分钟K线，获取最新价格

	reqURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return models.Stock{}, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("token", s.token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return models.Stock{}, fmt.Errorf("请求iTick API失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Stock{}, fmt.Errorf("读取响应失败: %w", err)
	}

	var klineResp iTickKlineResponse
	if err := json.Unmarshal(body, &klineResp); err != nil {
		return models.Stock{}, fmt.Errorf("解析响应失败: %w, body: %s", err, string(body))
	}

	if klineResp.Code != 0 {
		return models.Stock{}, fmt.Errorf("iTick API错误: code=%d", klineResp.Code)
	}

	if len(klineResp.Data) == 0 {
		return models.Stock{}, fmt.Errorf("无数据: %s.%s", region, code)
	}

	// 取最新一条K线数据作为当前行情
	latest := klineResp.Data[len(klineResp.Data)-1]
	prevClose := latest.Open
	if len(klineResp.Data) >= 2 {
		prevClose = klineResp.Data[len(klineResp.Data)-2].Close
	}

	change := latest.Close - prevClose
	changePercent := 0.0
	if prevClose != 0 {
		changePercent = (change / prevClose) * 100
	}

	return models.Stock{
		Code:          code,
		Name:          models.GetStockName(code),
		Region:        region,
		Price:         latest.Close,
		Change:        change,
		ChangePercent: changePercent,
		Volume:        int64(latest.Volume),
		MarketCap:     0,
		Currency:      s.getCurrency(region),
		UpdatedAt:     time.Now().Format(time.RFC3339),
	}, nil
}

// getCurrency 根据区域返回货币单位
func (s *ITickService) getCurrency(region string) string {
	currencies := map[string]string{
		"US": "USD",
		"HK": "HKD",
		"CN": "CNY",
		"KR": "KRW",
		"AU": "AUD",
		"JP": "JPY",
		"SG": "SGD",
	}
	if c, ok := currencies[region]; ok {
		return c
	}
	return "USD"
}
