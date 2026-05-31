package controllers

import (
	"net/http"
	"strconv"

	"stocksSearch/backend/models"
	"stocksSearch/backend/services"
)

// DerivativeController 衍生品API控制器
type DerivativeController struct {
	derivativeService *services.DerivativeService
}

// NewDerivativeController 创建衍生品控制器
func NewDerivativeController(ds *services.DerivativeService) *DerivativeController {
	return &DerivativeController{derivativeService: ds}
}

// GetDerivativeList 获取衍生品列表
// GET /api/derivatives?category=CRYPTO&sub_category=CRYPTO_MAJOR&sort_by=price&order=desc
func (c *DerivativeController) GetDerivativeList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(query.Get("page_size"))
	if pageSize <= 0 {
		pageSize = 20
	}

	req := models.DerivativeListRequest{
		Category:    query.Get("category"),
		SubCategory: query.Get("sub_category"),
		SortBy:      query.Get("sort_by"),
		Order:       query.Get("order"),
		Page:        page,
		PageSize:    pageSize,
	}

	resp, err := c.derivativeService.GetDerivativeList(req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
