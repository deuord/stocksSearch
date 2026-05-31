package models

// DerivativeItem 衍生品通用数据模型
type DerivativeItem struct {
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	Category      string  `json:"category"`       // CRYPTO/CS2/FOREX/INDICES
	SubCategory   string  `json:"sub_category"`
	Price         float64 `json:"price"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
	Volume        float64 `json:"volume"`
	MarketCap     float64 `json:"market_cap"`
	Currency      string  `json:"currency"`
	UpdatedAt     string  `json:"updated_at"`
}

// DerivativeListRequest 衍生品列表请求
type DerivativeListRequest struct {
	Category    string `json:"category"`
	SubCategory string `json:"sub_category"`
	SortBy      string `json:"sort_by"`
	Order       string `json:"order"`
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
}

// DerivativeListResponse 衍生品列表响应
type DerivativeListResponse struct {
	Total int              `json:"total"`
	Items []DerivativeItem `json:"items"`
	Page  int              `json:"page"`
	Pages int              `json:"pages"`
}

// 衍生品代码映射
var CryptoCodes = map[string][]string{
	"CRYPTO_MAJOR": {"BTC", "ETH", "BNB", "SOL", "XRP", "ADA", "DOGE", "AVAX", "DOT", "LINK"},
	"CRYPTO_DEFI":  {"UNI", "AAVE", "MKR", "COMP", "SNX", "CRV"},
}

// iTick crypto 需要特殊符号格式
var CryptoSymbols = map[string]string{
	"BTC": "BTCUSDT", "ETH": "ETHUSDT", "BNB": "BNBUSDT", "SOL": "SOLUSDT",
	"XRP": "XRPUSDT", "ADA": "ADAUSDT", "DOGE": "DOGEUSDT", "AVAX": "AVAXUSDT",
	"DOT": "DOTUSDT", "LINK": "LINKUSDT", "UNI": "UNIUSDT", "AAVE": "AAVEUSDT",
	"MKR": "MKRUSDT", "COMP": "COMPUSDT", "SNX": "SNXUSDT", "CRV": "CRVUSDT",
}

// iTick forex 符号
var ForexCodes = map[string][]string{
	"FOREX_MAJOR": {"EUR", "GBP", "JPY", "AUD", "NZD", "CAD", "CHF"},
	"FOREX_METAL": {"XAU", "XAG"},
}

var ForexSymbols = map[string]string{
	"EUR": "EURUSD", "GBP": "GBPUSD", "JPY": "USDJPY", "AUD": "AUDUSD",
	"NZD": "NZDUSD", "CAD": "USDCAD", "CHF": "USDCHF",
	"XAU": "XAUUSD", "XAG": "XAGUSD",
}

// iTick indices 符号
var IndicesCodes = map[string][]string{
	"INDICES_US":   {"SPX", "NAS", "DJI"},
	"INDICES_ASIA": {"HIS", "JPN225", "A50"},
	"INDICES_EU":   {"UKX", "DAX", "CAC"},
}

var IndicesSymbols = map[string]string{
	"SPX": "SPX", "NAS": "NAS", "DJI": "DJI",
	"HIS": "HIS", "JPN225": "JPN225", "A50": "A50",
	"UKX": "UKX", "DAX": "DAX", "CAC": "CAC",
}

// CS2 皮肤代码映射
var CS2SkinCodes = map[string][]string{
	"CS2_RIFLE": {
		"AK-47 | Redline (Field-Tested)",
		"AK-47 | Nightwish (Minimal Wear)",
		"M4A1-S | Nightmare (Minimal Wear)",
		"M4A4 | Desolate Space (Field-Tested)",
		"AK-47 | Neon Rider (Field-Tested)",
	},
	"CS2_SNIPER": {
		"AWP | Duality (Field-Tested)",
		"AWP | Containment Breach (Field-Tested)",
		"SSG 08 | Blood in the Water (Minimal Wear)",
		"AWP | Crakow! (Minimal Wear)",
	},
	"CS2_PISTOL": {
		"Desert Eagle | Code Red (Minimal Wear)",
		"Desert Eagle | Hypnotic (Factory New)",
		"USP-S | Road Rash (Field-Tested)",
		"Glock-18 | Water Elemental (Factory New)",
	},
	"CS2_KNIFE": {
		"★ Shadow Daggers | Ultraviolet (Field-Tested)",
		"Driver Gloves | King Snake (Field-Tested)",
		"★ Flip Knife | Rust Coat (Battle-Scarred)",
		"Specialist Gloves | Crimson Kimono (Field-Tested)",
	},
}

// 衍生品名称映射
var DerivativeNameMap = map[string]string{
	// 加密货币
	"BTC": "比特币", "ETH": "以太坊", "BNB": "币安币", "SOL": "Solana",
	"XRP": "瑞波币", "ADA": "艾达币", "DOGE": "狗狗币", "AVAX": "Avalanche",
	"DOT": "波卡", "LINK": "Chainlink", "UNI": "Uniswap", "AAVE": "Aave",
	"MKR": "Maker", "COMP": "Compound", "SNX": "Synthetix", "CRV": "Curve",
	// 外汇
	"EUR": "欧元/美元", "GBP": "英镑/美元", "JPY": "美元/日元",
	"AUD": "澳元/美元", "NZD": "纽元/美元", "CAD": "美元/加元",
	"CHF": "美元/瑞郎", "XAU": "黄金/美元", "XAG": "白银/美元",
	// 指数
	"SPX": "标普500", "NAS": "纳斯达克", "DJI": "道琼斯",
	"HIS": "恒生指数", "JPN225": "日经225", "A50": "富时A50",
	"UKX": "英国富时100", "DAX": "德国DAX", "CAC": "法国CAC40",
}
