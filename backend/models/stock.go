package models

// Stock 股票数据模型
type Stock struct {
	Code         string  `json:"code"`          // 股票代码
	Name         string  `json:"name"`          // 股票名称
	Region       string  `json:"region"`        // 所属市场区域
	Price        float64 `json:"price"`         // 最新价格
	Change       float64 `json:"change"`        // 涨跌额
	ChangePercent float64 `json:"change_percent"` // 涨跌幅百分比
	Volume       int64   `json:"volume"`        // 成交量
	MarketCap    float64 `json:"market_cap"`    // 市值
	Currency     string  `json:"currency"`      // 货币单位
	SubCategory  string  `json:"sub_category"`  // 所属细分板块ID
	UpdatedAt    string  `json:"updated_at"`    // 更新时间
}

// StockListRequest 股票列表请求参数
type StockListRequest struct {
	Region      string `json:"region"`       // 市场区域，空表示全部
	SubCategory string `json:"sub_category"` // 细分板块，空表示全部
	SortBy      string `json:"sort_by"`      // 排序字段，默认 market_cap
	Order       string `json:"order"`        // 排序方向 asc/desc，默认 desc
	Page        int    `json:"page"`         // 页码
	PageSize    int    `json:"page_size"`    // 每页数量
}

// StockListResponse 股票列表响应
type StockListResponse struct {
	Total  int     `json:"total"`
	Stocks []Stock `json:"stocks"`
	Page   int     `json:"page"`
	Pages  int     `json:"pages"`
}

// RegionStockCodes 各市场热门前沿股票代码（用于iTick API查询）
// 按细分板块分组
var RegionStockCodes = map[string]map[string][]string{
	"US": {
		"US_SEMI":    {"NVDA", "AMD", "INTC", "QCOM", "TXN", "MU"},
		"US_CHIP":    {"AVGO", "MRVL", "ON", "MPWR", "NXPI"},
		"US_STORAGE": {"WDC", "STX", "PSTG", "NTAP"},
		"US_CONSUMER": {"AAPL", "AMZN", "TSLA", "NKE", "SBUX", "MCD"},
		"US_SOFTWARE": {"MSFT", "ADBE", "CRM", "ORCL", "NOW"},
		"US_AI":      {"GOOGL", "META", "NVDA", "MSFT", "PLTR"},
		"US_ELECTRIC": {"TSLA", "RIVN", "LCID", "F", "GM"},
		"US_FINANCE": {"JPM", "GS", "BAC", "MS", "V", "MA"},
	},
	"HK": {
		"HK_INTERNET":   {"0700", "9988", "9618", "3690", "9888"},
		"HK_REALESTATE": {"0016", "0012", "0083", "0017"},
		"HK_FINANCE":    {"0005", "0388", "0939", "1398", "1299"},
		"HK_CONSUMER":   {"2020", "2331", "1876", "0168"},
		"HK_SEMI":       {"0981", "1347", "2018"},
	},
	"CN": {
		"CN_SEMI":     {"688981", "002371", "603986", "600703"},
		"CN_CHIP":     {"688256", "688012", "688008"},
		"CN_NEWENERGY": {"300750", "601012", "002594", "600438"},
		"CN_CONSUMER": {"600519", "000858", "002304", "000568"},
		"CN_FINANCE":  {"600036", "601318", "600030", "000001"},
	},
	"KR": {
		"KR_SEMI":       {"005930", "000660", "035420"},
		"KR_ELECTRONICS": {"066570", "003550", "034220"},
		"KR_AUTO":       {"005380", "000270", "012330"},
	},
	"AU": {
		"AU_MINING":  {"BHP", "RIO", "FMG", "NCM"},
		"AU_FINANCE": {"CBA", "WBC", "NAB", "ANZ"},
		"AU_ENERGY":  {"WDS", "STO", "ORG"},
	},
	"JP": {
		"JP_AUTO":        {"7203", "7267", "7269", "7270"},
		"JP_ELECTRONICS": {"6758", "6752", "6753"},
		"JP_SEMI":        {"8035", "6723", "6501"},
	},
	"SG": {
		"SG_FINANCE": {"D05", "O39", "U11"},
		"SG_REIT":    {"C38U", "A17U", "ME8U"},
	},
}
