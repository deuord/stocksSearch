package models

// Category 分类模型，支持层级结构
type Category struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	ParentID string     `json:"parent_id,omitempty"` // 空表示大分类
	Region   string     `json:"region"`              // iTick 对应的 region 代码
	Children []Category `json:"children,omitempty"`  // 子分类列表
}

// MajorCategory 大分类（按市场区域）
type MajorCategory string

const (
	MajorUS   MajorCategory = "US"   // 美股
	MajorHK   MajorCategory = "HK"   // 港股
	MajorCN   MajorCategory = "CN"   // A股
	MajorKR   MajorCategory = "KR"   // 韩股
	MajorAU   MajorCategory = "AU"   // 澳股
	MajorJP   MajorCategory = "JP"   // 日股
	MajorSG   MajorCategory = "SG"   // 新加坡
	MajorEU   MajorCategory = "EU"   // 欧股
)

// GetMajorCategories 获取所有大分类
func GetMajorCategories() []Category {
	return []Category{
		{ID: "US", Name: "美股", Region: "US"},
		{ID: "HK", Name: "港股", Region: "HK"},
		{ID: "CN", Name: "A股", Region: "CN"},
		{ID: "KR", Name: "韩股", Region: "KR"},
		{ID: "AU", Name: "澳股", Region: "AU"},
		{ID: "JP", Name: "日股", Region: "JP"},
		{ID: "SG", Name: "新加坡", Region: "SG"},
		// 衍生品
		{ID: "CRYPTO", Name: "加密货币", Region: "CRYPTO"},
		{ID: "CS2", Name: "CS2 皮肤", Region: "CS2"},
		{ID: "FOREX", Name: "外汇/贵金属", Region: "FOREX"},
		{ID: "INDICES", Name: "全球指数", Region: "INDICES"},
	}
}

// GetSubCategories 根据大分类获取子分类（行业板块）
func GetSubCategories(majorID string) []Category {
	allSubs := map[string][]Category{
		"US": {
			{ID: "US_SEMI", Name: "半导体", ParentID: "US", Region: "US"},
			{ID: "US_CHIP", Name: "芯片", ParentID: "US", Region: "US"},
			{ID: "US_STORAGE", Name: "存储", ParentID: "US", Region: "US"},
			{ID: "US_CONSUMER", Name: "消费", ParentID: "US", Region: "US"},
			{ID: "US_SOFTWARE", Name: "软件", ParentID: "US", Region: "US"},
			{ID: "US_AI", Name: "人工智能", ParentID: "US", Region: "US"},
			{ID: "US_ELECTRIC", Name: "电动车", ParentID: "US", Region: "US"},
			{ID: "US_FINANCE", Name: "金融", ParentID: "US", Region: "US"},
		},
		"HK": {
			{ID: "HK_INTERNET", Name: "互联网", ParentID: "HK", Region: "HK"},
			{ID: "HK_REALESTATE", Name: "房地产", ParentID: "HK", Region: "HK"},
			{ID: "HK_FINANCE", Name: "金融", ParentID: "HK", Region: "HK"},
			{ID: "HK_CONSUMER", Name: "消费", ParentID: "HK", Region: "HK"},
			{ID: "HK_SEMI", Name: "半导体", ParentID: "HK", Region: "HK"},
		},
		"CN": {
			{ID: "CN_SEMI", Name: "半导体", ParentID: "CN", Region: "CN"},
			{ID: "CN_CHIP", Name: "芯片", ParentID: "CN", Region: "CN"},
			{ID: "CN_NEWENERGY", Name: "新能源", ParentID: "CN", Region: "CN"},
			{ID: "CN_CONSUMER", Name: "消费", ParentID: "CN", Region: "CN"},
			{ID: "CN_FINANCE", Name: "金融", ParentID: "CN", Region: "CN"},
		},
		"KR": {
			{ID: "KR_SEMI", Name: "半导体", ParentID: "KR", Region: "KR"},
			{ID: "KR_ELECTRONICS", Name: "电子", ParentID: "KR", Region: "KR"},
			{ID: "KR_AUTO", Name: "汽车", ParentID: "KR", Region: "KR"},
		},
		"AU": {
			{ID: "AU_MINING", Name: "矿业", ParentID: "AU", Region: "AU"},
			{ID: "AU_FINANCE", Name: "金融", ParentID: "AU", Region: "AU"},
			{ID: "AU_ENERGY", Name: "能源", ParentID: "AU", Region: "AU"},
		},
		"JP": {
			{ID: "JP_AUTO", Name: "汽车", ParentID: "JP", Region: "JP"},
			{ID: "JP_ELECTRONICS", Name: "电子", ParentID: "JP", Region: "JP"},
			{ID: "JP_SEMI", Name: "半导体", ParentID: "JP", Region: "JP"},
		},
		"SG": {
			{ID: "SG_FINANCE", Name: "金融", ParentID: "SG", Region: "SG"},
			{ID: "SG_REIT", Name: "REITs", ParentID: "SG", Region: "SG"},
		},
		// 衍生品子分类
		"CRYPTO": {
			{ID: "CRYPTO_MAJOR", Name: "主流币", ParentID: "CRYPTO", Region: "CRYPTO"},
			{ID: "CRYPTO_DEFI", Name: "DeFi", ParentID: "CRYPTO", Region: "CRYPTO"},
		},
		"CS2": {
			{ID: "CS2_RIFLE", Name: "步枪 AK/M4", ParentID: "CS2", Region: "CS2"},
			{ID: "CS2_SNIPER", Name: "狙击 AWP", ParentID: "CS2", Region: "CS2"},
			{ID: "CS2_PISTOL", Name: "手枪/其他", ParentID: "CS2", Region: "CS2"},
			{ID: "CS2_KNIFE", Name: "刀具/手套", ParentID: "CS2", Region: "CS2"},
		},
		"FOREX": {
			{ID: "FOREX_MAJOR", Name: "主要货币对", ParentID: "FOREX", Region: "FOREX"},
			{ID: "FOREX_METAL", Name: "贵金属", ParentID: "FOREX", Region: "FOREX"},
		},
		"INDICES": {
			{ID: "INDICES_US", Name: "美股指数", ParentID: "INDICES", Region: "INDICES"},
			{ID: "INDICES_ASIA", Name: "亚太指数", ParentID: "INDICES", Region: "INDICES"},
			{ID: "INDICES_EU", Name: "欧洲指数", ParentID: "INDICES", Region: "INDICES"},
		},
	}

	if subs, ok := allSubs[majorID]; ok {
		return subs
	}
	return nil
}

// GetAllCategoriesWithChildren 获取完整分类树
func GetAllCategoriesWithChildren() []Category {
	majors := GetMajorCategories()
	for i := range majors {
		majors[i].Children = GetSubCategories(majors[i].ID)
	}
	return majors
}
