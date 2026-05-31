package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stocksSearch/backend/models"
	"stocksSearch/backend/services"
)

// StockController 股票相关API控制器
type StockController struct {
	stockService *services.StockService
}

// NewStockController 创建股票控制器实例
func NewStockController(stockService *services.StockService) *StockController {
	return &StockController{stockService: stockService}
}

// GetStockList 获取股票列表 API
// GET /api/stocks?region=US&sub_category=US_SEMI&sort_by=market_cap&order=desc&page=1&page_size=20
func (c *StockController) GetStockList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(query.Get("page_size"))
	if pageSize <= 0 {
		pageSize = 20
	}

	req := models.StockListRequest{
		Region:      query.Get("region"),
		SubCategory: query.Get("sub_category"),
		SortBy:      query.Get("sort_by"),
		Order:       query.Get("order"),
		Page:        page,
		PageSize:    pageSize,
	}

	resp, err := c.stockService.GetStockList(req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// GetCategories 获取分类树 API
// GET /api/categories
func (c *StockController) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories := models.GetAllCategoriesWithChildren()
	writeJSON(w, http.StatusOK, categories)
}

// GetSubCategories 获取指定大分类下的子分类
// GET /api/categories/:major_id
func (c *StockController) GetSubCategories(w http.ResponseWriter, r *http.Request) {
	// 从路径中获取 major_id
	// 这里需要路由解析，暂时用 query 参数
	majorID := r.URL.Query().Get("major_id")
	if majorID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "major_id 不能为空"})
		return
	}

	subs := models.GetSubCategories(majorID)
	if subs == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "未找到该市场分类"})
		return
	}

	writeJSON(w, http.StatusOK, subs)
}

// writeJSON 统一JSON响应写入
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
