import type { Stock } from '../types'

interface StockListProps {
  stocks: Stock[]
  loading: boolean
  total: number
  page: number
  pages: number
  selectedMajor: string
  onPageChange: (page: number) => void
  onStockClick: (stock: Stock) => void
  activeStock: string
}

function formatMarketCap(cap: number): string {
  if (cap >= 1e12) return `${(cap / 1e12).toFixed(2)}T`
  if (cap >= 1e8) return `${(cap / 1e8).toFixed(2)}亿`
  return cap.toLocaleString()
}

function formatVolume(vol: number): string {
  if (vol >= 1e8) return `${(vol / 1e8).toFixed(2)}亿`
  if (vol >= 1e4) return `${(vol / 1e4).toFixed(2)}万`
  return String(vol)
}

const REGION_NAMES: Record<string, string> = {
  US: '美股', HK: '港股', CN: 'A股', KR: '韩股', AU: '澳股', JP: '日股', SG: '新加坡',
}

export default function StockList({
  stocks,
  loading,
  total,
  page,
  pages,
  selectedMajor,
  onPageChange,
  onStockClick,
  activeStock,
}: StockListProps) {
  if (loading) {
    return (
      <div className="stock-list">
        <div className="loading-spinner">加载中...</div>
      </div>
    )
  }

  if (stocks.length === 0) {
    return (
      <div className="stock-list">
        <div className="empty-state">
          {selectedMajor ? '该分类暂无数据，请尝试其他分类' : '请选择分类查看股票'}
        </div>
      </div>
    )
  }

  return (
    <div className="stock-list">
      <div className="stock-count">共 {total} 只股票</div>
      <div className="stock-table-wrapper">
        <table className="stock-table">
          <thead>
            <tr>
              <th>代码</th>
              <th>名称</th>
              <th>市场</th>
              <th>价格</th>
              <th>涨跌</th>
              <th>涨跌幅</th>
              <th>成交量</th>
              <th>市值</th>
            </tr>
          </thead>
          <tbody>
            {stocks.map(stock => (
              <tr
                key={`${stock.region}-${stock.code}`}
                className={`stock-row ${activeStock === stock.code ? 'active' : ''} ${stock.change_percent > 0 ? 'up' : stock.change_percent < 0 ? 'down' : ''}`}
                onClick={() => onStockClick(stock)}
              >
                <td className="code">{stock.code}</td>
                <td className="name">{stock.name}</td>
                <td>
                  <span className={`region-badge region-${stock.region}`}>
                    {REGION_NAMES[stock.region] || stock.region}
                  </span>
                </td>
                <td className="price">{stock.price.toFixed(2)}</td>
                <td className={`change ${stock.change >= 0 ? 'positive' : 'negative'}`}>
                  {stock.change >= 0 ? '+' : ''}{stock.change.toFixed(2)}
                </td>
                <td className={`change-percent ${stock.change_percent >= 0 ? 'positive' : 'negative'}`}>
                  {stock.change_percent >= 0 ? '+' : ''}{stock.change_percent.toFixed(2)}%
                </td>
                <td className="volume">{formatVolume(stock.volume)}</td>
                <td className="market-cap">{formatMarketCap(stock.market_cap)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* 分页 */}
      {pages > 1 && (
        <div className="pagination">
          <button
            className="page-btn"
            disabled={page <= 1}
            onClick={() => onPageChange(page - 1)}
          >
            上一页
          </button>
          <span className="page-info">{page} / {pages}</span>
          <button
            className="page-btn"
            disabled={page >= pages}
            onClick={() => onPageChange(page + 1)}
          >
            下一页
          </button>
        </div>
      )}
    </div>
  )
}
