package router

import (
	"encoding/json"
	"net/http"
	"strings"

	"stocksSearch/backend/controllers"
	"stocksSearch/backend/models"
)

// SetupRouter 配置路由，返回 http.Handler
func SetupRouter(sc *controllers.StockController, wsc *controllers.WSController, dc *controllers.DerivativeController) http.Handler {
	mux := http.NewServeMux()

	// 股票列表接口
	mux.HandleFunc("/api/stocks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			setCORS(w)
			return
		}
		sc.GetStockList(w, r)
	})

	// 分类树接口
	mux.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			setCORS(w)
			return
		}

		// 解析子路由：/api/categories/US -> 返回子分类
		path := strings.TrimPrefix(r.URL.Path, "/api/categories")
		path = strings.TrimPrefix(path, "/")

		if path != "" && path != "categories" {
			subs := models.GetSubCategories(path)
			writeCategoriesJSON(w, subs)
			return
		}

		sc.GetCategories(w, r)
	})

	// 健康检查
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// WebSocket 实时行情端点
	mux.HandleFunc("/ws", wsc.HandleWS)

	// 衍生品列表接口
	mux.HandleFunc("/api/derivatives", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			setCORS(w)
			return
		}
		dc.GetDerivativeList(w, r)
	})

	return mux
}

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}

func writeCategoriesJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if data == nil {
		w.Write([]byte(`[]`))
		return
	}
	json.NewEncoder(w).Encode(data)
}
