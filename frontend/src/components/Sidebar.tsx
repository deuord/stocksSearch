import type { Category } from '../types'

interface SidebarProps {
  categories: Category[]
  selectedMajor: string
  selectedSub: string
  sortBy: string
  order: string
  loading: boolean
  onSelectMajor: (id: string) => void
  onSelectSub: (id: string) => void
  onSortChange: (sortBy: string, order: string) => void
}

export default function Sidebar({
  categories,
  selectedMajor,
  selectedSub,
  sortBy,
  order,
  loading,
  onSelectMajor,
  onSelectSub,
  onSortChange,
}: SidebarProps) {

  const activeCategory = categories.find(c => c.id === selectedMajor)

  return (
    <div className="sidebar">
      <div className="sidebar-header">
        <h1 className="logo">🌐 全球股市情报</h1>
      </div>

      {/* 大分类 */}
      <div className="category-section">
        <div
          className={`category-item major ${selectedMajor === '' ? 'active' : ''}`}
          onClick={() => onSelectMajor('')}
        >
          全部市场
        </div>
        {categories.map(cat => (
          <div
            key={cat.id}
            className={`category-item major ${selectedMajor === cat.id ? 'active' : ''}`}
            onClick={() => onSelectMajor(cat.id)}
          >
            {cat.name}
          </div>
        ))}
      </div>

      {/* 细分板块 */}
      {activeCategory && activeCategory.children && activeCategory.children.length > 0 && (
        <div className="category-section sub-section">
          <div className="section-title">细分板块</div>
          <div
            className={`category-item sub ${selectedSub === '' ? 'active' : ''}`}
            onClick={() => onSelectSub('')}
          >
            全部
          </div>
          {activeCategory.children.map(sub => (
            <div
              key={sub.id}
              className={`category-item sub ${selectedSub === sub.id ? 'active' : ''}`}
              onClick={() => onSelectSub(sub.id)}
            >
              {sub.name}
            </div>
          ))}
        </div>
      )}

      {/* 排序选项 */}
      <div className="sort-section">
        <div className="section-title">排序方式</div>
        <div className="sort-options">
          <button
            className={`sort-btn ${sortBy === 'market_cap' ? 'active' : ''}`}
            onClick={() => onSortChange('market_cap', sortBy === 'market_cap' ? (order === 'desc' ? 'asc' : 'desc') : 'desc')}
          >
            市值 {sortBy === 'market_cap' && (order === 'desc' ? '↓' : '↑')}
          </button>
          <button
            className={`sort-btn ${sortBy === 'price' ? 'active' : ''}`}
            onClick={() => onSortChange('price', sortBy === 'price' ? (order === 'desc' ? 'asc' : 'desc') : 'desc')}
          >
            价格 {sortBy === 'price' && (order === 'desc' ? '↓' : '↑')}
          </button>
          <button
            className={`sort-btn ${sortBy === 'change_percent' ? 'active' : ''}`}
            onClick={() => onSortChange('change_percent', sortBy === 'change_percent' ? (order === 'desc' ? 'asc' : 'desc') : 'desc')}
          >
            涨跌幅 {sortBy === 'change_percent' && (order === 'desc' ? '↓' : '↑')}
          </button>
        </div>
      </div>

      {loading && <div className="loading-indicator">加载中...</div>}
    </div>
  )
}
