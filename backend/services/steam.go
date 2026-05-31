package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// SteamService Steam Market CS2服务
type SteamService struct {
	httpClient *http.Client
}

// SteamPriceOverview Steam价格响应
type SteamPriceOverview struct {
	Success      bool   `json:"success"`
	LowestPrice  string `json:"lowest_price"`
	Volume       string `json:"volume"`
	MedianPrice  string `json:"median_price"`
}

// NewSteamService 创建Steam服务
func NewSteamService() *SteamService {
	return &SteamService{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// GetCS2SkinPrice 获取CS2皮肤价格
func (s *SteamService) GetCS2SkinPrice(marketHashName string) (price float64, volume int, err error) {
	apiURL := "https://steamcommunity.com/market/priceoverview/"
	params := url.Values{}
	params.Set("appid", "730")
	params.Set("currency", "1")
	params.Set("market_hash_name", marketHashName)

	reqURL := apiURL + "?" + params.Encode()

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	var overview SteamPriceOverview
	if err := json.Unmarshal(body, &overview); err != nil {
		return 0, 0, fmt.Errorf("解析Steam响应失败: %w", err)
	}

	if !overview.Success {
		return 0, 0, fmt.Errorf("Steam返回失败")
	}

	price = parseSteamPrice(overview.LowestPrice)
	volume = parseSteamVolume(overview.Volume)

	return price, volume, nil
}

// parseSteamPrice 解析Steam价格字符串 "$1,234.56" → 1234.56
func parseSteamPrice(s string) float64 {
	if s == "" {
		return 0
	}
	cleaned := ""
	for _, ch := range s {
		if (ch >= '0' && ch <= '9') || ch == '.' {
			cleaned += string(ch)
		}
	}
	var result float64
	fmt.Sscanf(cleaned, "%f", &result)
	return result
}

// parseSteamVolume 解析Steam销量
func parseSteamVolume(s string) int {
	if s == "" {
		return 0
	}
	cleaned := ""
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			cleaned += string(ch)
		}
	}
	var result int
	fmt.Sscanf(cleaned, "%d", &result)
	return result
}

// simplifySkinName 简化皮肤名称显示
func simplifySkinName(fullName string) string {
	parts := splitSkinName(fullName)
	if len(parts) >= 2 {
		return parts[0] + " | " + parts[1]
	}
	return fullName
}

func splitSkinName(name string) []string {
	var result []string
	current := ""
	depth := 0
	for _, ch := range name {
		if ch == '(' {
			depth++
			if depth == 1 && current != "" {
				result = append(result, current)
				current = ""
			}
		} else if ch == ')' {
			depth--
		} else if depth == 0 {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
