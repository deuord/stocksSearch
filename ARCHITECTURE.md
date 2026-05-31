# 全球股市情报面板 - 项目架构文档

> 最后更新: 2026-05-31

## 项目概述

全球股市实时行情监控面板，支持多市场（美股、港股、A股、韩股、澳股、日股、新加坡）分类查看，左侧数据面板 + 右侧3D地球可视化。

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.x + 标准库 net/http |
| 前端 | React 19 + TypeScript + Vite |
| 3D可视化 | Three.js + @react-three/fiber + @react-three/drei |
| 数据源 | iTick API (免费版) |
| 管理 | Git |

## 项目目录结构

```
stocksSearch/
├── backend/                    # Go 后端 (MVC 架构)
│   ├── main.go                 # 入口，启动服务
│   ├── config/                 # 配置层
│   │   ├── config.go           # 配置结构体与加载逻辑
│   │   ├── config.json         # 实际配置（含Token，gitignore）
│   │   └── config.example.json # 配置模板
│   ├── models/                 # Model 层 - 数据模型
│   │   ├── stock.go            # 股票模型 + 股票代码映射表
│   │   ├── category.go         # 分类模型（大分类/细分板块）
│   │   └── stock_name.go       # 股票中英文名称映射
│   ├── services/               # Service 层 - 业务逻辑
│   │   ├── itick.go            # iTick API 对接（kline接口）
│   │   └── stock.go            # 股票列表查询、排序、分页
│   ├── controllers/            # Controller 层 - 请求处理
│   │   └── stock.go            # API 接口处理器
│   └── router/                 # 路由层
│       └── router.go           # URL 路由映射
│
├── frontend/                   # React 前端 (SPA)
│   ├── src/
│   │   ├── types/index.ts      # TypeScript 类型定义
│   │   ├── hooks/              # 自定义 Hooks
│   │   │   ├── useCategories.ts # 分类数据获取
│   │   │   └── useStocks.ts    # 股票数据获取
│   │   ├── components/         # UI 组件
│   │   │   ├── Sidebar.tsx     # 侧边栏（分类选择 + 排序）
│   │   │   ├── StockList.tsx   # 股票列表表格
│   │   │   └── Globe.tsx       # 3D 线框地球
│   │   ├── App.tsx             # 主应用组件（状态管理）
│   │   ├── App.css             # 组件样式（黑白主题）
│   │   ├── index.css           # 全局基础样式
│   │   └── main.tsx            # 应用入口
│   ├── vite.config.ts          # Vite配置（含API代理）
│   └── tsconfig.json
│
├── .gitignore                  # Git忽略（含config.json）
├── README.md
├── ARCHITECTURE.md             # 本文件
└── CHANGELOG.md                # 需求流水线记录
```

## 架构分层

### 后端 MVC 架构

```
请求 → Router → Controller → Service → Model
                                      ↓
                                 iTick API
```

- **Router**: URL路由映射，CORS处理
- **Controller**: 请求参数解析、响应格式化
- **Service**: 业务逻辑（数据聚合、排序、分页） + 外部API调用
- **Model**: 数据结构定义、分类体系、名称映射

### 前端组件架构

```
App (状态管理中心)
├── Sidebar (分类选择、排序)
│   ├── 大分类列表
│   ├── 细分板块列表
│   └── 排序选项
├── StockList (数据展示)
│   ├── 股票表格
│   └── 分页
└── Globe (3D地球)
    ├── 线框球体
    ├── 高亮环 (响应分类选择)
    └── 粒子光环 (流光效果)
```

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/categories` | 获取全部分类树 |
| GET | `/api/categories/:major_id` | 获取子分类 |
| GET | `/api/stocks` | 股票列表查询 |
| GET | `/api/health` | 健康检查 |

### 股票列表参数

| 参数 | 类型 | 说明 |
|------|------|------|
| region | string | 市场区域(US/HK/CN/KR/AU/JP/SG) |
| sub_category | string | 细分板块ID |
| sort_by | string | 排序字段(market_cap/price/change_percent) |
| order | string | 排序方向(asc/desc) |
| page | int | 页码 |
| page_size | int | 每页数量 |

## 数据源

- **iTick API** (免费版): `https://api-free.itick.org`
- **使用端点**: `/stock/kline` (kType=1 分钟K线)
- **限制**: 免费版不支持 snapshot/reference，单次请求一只股票

## 分类体系

### 大分类（按市场区域）

US(美股) | HK(港股) | CN(A股) | KR(韩股) | AU(澳股) | JP(日股) | SG(新加坡)

### 细分板块（示例）

- 美股: 半导体、芯片、存储、消费、软件、AI、电动车、金融
- 港股: 互联网、房地产、金融、消费、半导体
- A股: 半导体、芯片、新能源、消费、金融
- 韩股: 半导体、电子、汽车

## 启动方式

```bash
# 后端
cd backend && go run main.go    # 默认 :8080

# 前端
cd frontend && npm run dev      # 默认 :3000，自动代理 /api → :8080
```
