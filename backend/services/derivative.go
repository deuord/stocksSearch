package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"stocksSearch/backend/config"
	"stocksSearch/backend/models"
)

// DerivativeService 衍生品统一服务
type DerivativeService struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// iTickCryptoKlineItem 衍生品K线数据项
type iTickCryptoKlineItem struct {
	Close  float64 `json:"c"`
	High   float64 `json:"h"`
	Low    float64 `json:"l"`
	Open   float64 `json:"o"`
	Volume float64 `json:"v"`
	Time   int64   `json:"t"`
}

// iTickGenericKlineResponse 通用K线响应
type iTickGenericKlineResponse struct {
	Code int                    `json:"code"`
	Msg  *string                `json:"msg"`
	Data []iTickCryptoKlineItem `json:"data"`
}

// NewDerivativeService 创建衍生品服务
func NewDerivativeService(cfg config.ITickConfig) *DerivativeService {
	return &DerivativeService{
		baseURL: cfg.BaseURL,
		token:   cfg.Token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetDerivativeList 获取衍生品列表
func (s *DerivativeService) GetDerivativeList(req models.DerivativeListRequest) (*models.DerivativeListResponse, error) {
	var items []models.DerivativeItem

	switch req.Category {
	case "CRYPTO":
		items = s.fetchCryptoFromBinance(req.SubCategory)
	case "FOREX":
		items = s.fetchForexFromITick(req.SubCategory)
	case "INDICES":
		items = s.fetchProduct("indices", models.IndicesCodes, models.IndicesSymbols, models.DerivativeNameMap, "INDICES", req.SubCategory)
	case "CS2":
		items = s.fetchCS2Skins(req.SubCategory)
	default:
		items = append(items, s.fetchCryptoFromBinance("")...)
	}

	// 排序
	sortDerivativeItems(items, req.SortBy, req.Order)

	// 分页
	total := len(items)
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pages := (total + pageSize - 1) / pageSize
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	return &models.DerivativeListResponse{
		Total: total,
		Items: items[start:end],
		Page:  page,
		Pages: pages,
	}, nil
}

// fetchProduct 通用衍生品数据获取（使用kline获取最新价格）并发请求
func (s *DerivativeService) fetchProduct(
	product string,
	codeMap map[string][]string,
	symbolMap map[string]string,
	nameMap map[string]string,
	category string,
	subCategory string,
) []models.DerivativeItem {
	codes := collectCodes(codeMap, subCategory)

	var (
		items []models.DerivativeItem
		mu    sync.Mutex
		wg    sync.WaitGroup
	)

	// 限制并发为2，避免触发免费版限流
	sem := make(chan struct{}, 2)

	for _, code := range codes {
		sym := symbolMap[code]
		if sym == "" {
			continue
		}

		wg.Add(1)
		go func(c, sym string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			price, volume, err := s.fetchKlinePrice(product, sym)
			if err != nil {
				log.Printf("⚠️ 获取 %s/%s 失败: %v", product, sym, err)
				return
			}

			name := nameMap[c]
			if name == "" {
				name = c
			}

			mu.Lock()
			items = append(items, models.DerivativeItem{
				Code:      c,
				Name:      name,
				Category:  category,
				Price:     price,
				Volume:    volume,
				Currency:  "USD",
				UpdatedAt: time.Now().Format(time.RFC3339),
			})
			mu.Unlock()
		}(code, sym)
	}

	wg.Wait()
	return items
}

// fetchForexFromITick 使用iTick获取外汇行情（需要region=FX参数）
func (s *DerivativeService) fetchForexFromITick(subCategory string) []models.DerivativeItem {
	codes := collectCodes(models.ForexCodes, subCategory)

	var (
		items []models.DerivativeItem
		mu    sync.Mutex
		wg    sync.WaitGroup
	)

	sem := make(chan struct{}, 2)

	for _, code := range codes {
		sym := models.ForexSymbols[code]
		if sym == "" {
			continue
		}

		wg.Add(1)
		go func(c, symbol string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			price, _, err := s.fetchKlinePrice("forex", symbol, "region", "FX")
			if err != nil {
				log.Printf("⚠️ 获取 forex/%s 失败: %v", symbol, err)
				return
			}

			name := models.DerivativeNameMap[c]
			if name == "" {
				name = c
			}

			mu.Lock()
			items = append(items, models.DerivativeItem{
				Code:      c,
				Name:      name,
				Category:  "FOREX",
				Price:     price,
				Currency:  "USD",
				UpdatedAt: time.Now().Format(time.RFC3339),
			})
			mu.Unlock()
		}(code, sym)
	}

	wg.Wait()
	return items
}

// fetchCS2Skins 获取CS2皮肤价格（从Steam Market）
func (s *DerivativeService) fetchCS2Skins(subCategory string) []models.DerivativeItem {
	codes := collectCodes(models.CS2SkinCodes, subCategory)

	var (
		items []models.DerivativeItem
		mu    sync.Mutex
		wg    sync.WaitGroup
	)

	steamSvc := NewSteamService()
	sem := make(chan struct{}, 2)

	for _, skinName := range codes {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			price, volume, err := steamSvc.GetCS2SkinPrice(name)
			if err != nil {
				return
			}

			mu.Lock()
			items = append(items, models.DerivativeItem{
				Code:      name,
				Name:      simplifySkinName(name),
				Category:  "CS2",
				Price:     price,
				Volume:    float64(volume),
				Currency:  "USD",
				UpdatedAt: time.Now().Format(time.RFC3339),
			})
			mu.Unlock()
		}(skinName)
	}

	wg.Wait()
	return items
}

// fetchKlinePrice 通用kline价格获取
func (s *DerivativeService) fetchKlinePrice(product, code string, extraParams ...string) (price, volume float64, err error) {
	apiURL := fmt.Sprintf("%s/%s/kline", s.baseURL, product)
	params := url.Values{}
	params.Set("code", code)
	params.Set("kType", "1")

	// 支持额外参数（如 forex 的 region 等）
	for i := 0; i+1 < len(extraParams); i += 2 {
		params.Set(extraParams[i], extraParams[i+1])
	}

	req, err := http.NewRequest("GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("token", s.token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var klineResp iTickGenericKlineResponse
	if err := json.Unmarshal(body, &klineResp); err != nil {
		return 0, 0, fmt.Errorf("解析失败: %w", err)
	}

	if klineResp.Code != 0 || len(klineResp.Data) == 0 {
		return 0, 0, fmt.Errorf("无数据: code=%d", klineResp.Code)
	}

	latest := klineResp.Data[len(klineResp.Data)-1]
	return latest.Close, latest.Volume, nil
}

// collectCodes 从map中收集代码
func collectCodes(codeMap map[string][]string, subCategory string) []string {
	if subCategory != "" {
		if codes, ok := codeMap[subCategory]; ok {
			return codes
		}
		return nil
	}
	var all []string
	for _, codes := range codeMap {
		all = append(all, codes...)
	}
	return all
}

// fetchCryptoFromBinance 使用Binance公开API获取加密货币价格
func (s *DerivativeService) fetchCryptoFromBinance(subCategory string) []models.DerivativeItem {
	codes := collectCodes(models.CryptoCodes, subCategory)

	var (
		items []models.DerivativeItem
		mu    sync.Mutex
		wg    sync.WaitGroup
	)

	binanceSvc := NewBinanceService()
	sem := make(chan struct{}, 3)

	for _, code := range codes {
		sym := models.CryptoSymbols[code]
		if sym == "" {
			continue
		}

		wg.Add(1)
		go func(c, symbol string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			price, err := binanceSvc.GetPrice(symbol)
			if err != nil {
				log.Printf("⚠️ Binance获取 %s 失败: %v", symbol, err)
				return
			}

			name := models.DerivativeNameMap[c]
			if name == "" {
				name = c
			}

			mu.Lock()
			items = append(items, models.DerivativeItem{
				Code:      c,
				Name:      name,
				Category:  "CRYPTO",
				Price:     price,
				Currency:  "USD",
				UpdatedAt: time.Now().Format(time.RFC3339),
			})
			mu.Unlock()
		}(code, sym)
	}

	wg.Wait()
	return items
}

// sortDerivativeItems 排序衍生品列表
func sortDerivativeItems(items []models.DerivativeItem, sortBy, order string) {
	if sortBy == "" {
		sortBy = "price"
	}
	desc := order == "desc"

	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			var shouldSwap bool
			switch sortBy {
			case "change_percent":
				if desc {
					shouldSwap = items[i].ChangePercent < items[j].ChangePercent
				} else {
					shouldSwap = items[i].ChangePercent > items[j].ChangePercent
				}
			default:
				if desc {
					shouldSwap = items[i].Price < items[j].Price
				} else {
					shouldSwap = items[i].Price > items[j].Price
				}
			}
			if shouldSwap {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}
