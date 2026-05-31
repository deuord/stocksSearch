export interface Stock {
  code: string
  name: string
  region: string
  price: number
  change: number
  change_percent: number
  volume: number
  market_cap: number
  currency: string
  sub_category: string
  updated_at: string
}

export interface Category {
  id: string
  name: string
  parent_id?: string
  region: string
  children?: Category[]
}

export interface StockListResponse {
  total: number
  stocks: Stock[]
  page: number
  pages: number
}

export interface GlobeInteraction {
  region: string
  subCategory: string
  stockCode: string
}
