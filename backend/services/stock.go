package services

import (
	"sort"

	"stocksSearch/backend/models"
)

// StockService 股票业务服务
type StockService struct {
	itick *ITickService
}

// NewStockService 创建股票服务实例
func NewStockService(itick *ITickService) *StockService {
	return &StockService{itick: itick}
}

// GetStockList 获取股票列表，支持按分类过滤和排序
func (s *StockService) GetStockList(req models.StockListRequest) (*models.StockListResponse, error) {
	// 收集需要查询的股票代码
	codesMap := s.collectCodes(req.Region, req.SubCategory)

	var allStocks []models.Stock

	for region, codes := range codesMap {
		stocks, err := s.itick.FetchStocksByRegion(region, codes)
		if err != nil {
			// 记录错误但继续处理其他市场
			continue
		}
		allStocks = append(allStocks, stocks...)
	}

	// 按市值排序（默认降序）
	s.sortStocks(allStocks, req.SortBy, req.Order)

	// 分页
	total := len(allStocks)
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

	return &models.StockListResponse{
		Total:  total,
		Stocks: allStocks[start:end],
		Page:   page,
		Pages:  pages,
	}, nil
}

// collectCodes 收集需要查询的股票代码
func (s *StockService) collectCodes(region, subCategory string) map[string][]string {
	codesMap := make(map[string][]string)

	if region != "" {
		// 指定了大分类
		if subCategory != "" {
			// 指定了细分板块，只查该板块下的股票
			if codes, ok := models.RegionStockCodes[region][subCategory]; ok {
				codesMap[region] = codes
			}
		} else {
			// 不指定细分板块，查该市场所有板块
			if subMap, ok := models.RegionStockCodes[region]; ok {
				var allCodes []string
				for _, codes := range subMap {
					allCodes = append(allCodes, codes...)
				}
				codesMap[region] = allCodes
			}
		}
	} else {
		// 不指定大分类，查询所有市场
		if subCategory != "" {
			// 有细分板块，需要遍历查找属于该板块的市场
			for region, subMap := range models.RegionStockCodes {
				if codes, ok := subMap[subCategory]; ok {
					codesMap[region] = codes
				}
			}
		} else {
			// 全查所有市场所有板块
			for region, subMap := range models.RegionStockCodes {
				var allCodes []string
				for _, codes := range subMap {
					allCodes = append(allCodes, codes...)
				}
				codesMap[region] = allCodes
			}
		}
	}

	return codesMap
}

// sortStocks 对股票列表排序
func (s *StockService) sortStocks(stocks []models.Stock, sortBy, order string) {
	if sortBy == "" {
		sortBy = "market_cap"
	}
	if order == "" {
		order = "desc"
	}

	desc := order == "desc"

	sort.Slice(stocks, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "price":
			less = stocks[i].Price < stocks[j].Price
		case "change_percent":
			less = stocks[i].ChangePercent < stocks[j].ChangePercent
		case "volume":
			less = stocks[i].Volume < stocks[j].Volume
		default:
			less = stocks[i].MarketCap < stocks[j].MarketCap
		}
		if desc {
			return !less
		}
		return less
	})
}
