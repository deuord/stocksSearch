import { useState, useEffect, useCallback } from "react";
import Sidebar from "./components/Sidebar";
import StockList from "./components/StockList";
import Globe from "./components/Globe";
import { useCategories } from "./hooks/useCategories";
import { useStocks } from "./hooks/useStocks";
import { useDerivatives } from "./hooks/useDerivatives";
import type { Stock } from "./types";
import "./App.css";

const DERIVATIVE_IDS = new Set(["CRYPTO", "CS2", "FOREX", "INDICES"]);

export default function App() {
  const { categories, loading: catLoading } = useCategories();
  const stockHook = useStocks();
  const derivHook = useDerivatives();

  const [selectedMajor, setSelectedMajor] = useState("");
  const [selectedSub, setSelectedSub] = useState("");
  const [sortBy, setSortBy] = useState("market_cap");
  const [order, setOrder] = useState("desc");
  const [activeStock, setActiveStock] = useState("");

  const isDerivative = DERIVATIVE_IDS.has(selectedMajor);

  const { stocks, loading, total, page, pages, setPage } = isDerivative
    ? { ...derivHook, setPage: derivHook.setPage }
    : { ...stockHook, setPage: stockHook.setPage };

  const handleSelectMajor = useCallback(
    (id: string) => {
      setSelectedMajor(id);
      setSelectedSub("");
      setPage(1);
      setActiveStock("");
      if (DERIVATIVE_IDS.has(id)) {
        setSortBy("price");
      } else {
        setSortBy("market_cap");
      }
    },
    [setPage],
  );

  const handleSelectSub = useCallback(
    (id: string) => {
      setSelectedSub(id);
      setPage(1);
      setActiveStock("");
    },
    [setPage],
  );

  const handleSortChange = useCallback(
    (newSortBy: string, newOrder: string) => {
      setSortBy(newSortBy);
      setOrder(newOrder);
      setPage(1);
    },
    [setPage],
  );

  const handleStockClick = useCallback((stock: Stock) => {
    setActiveStock((prev) => (prev === stock.code ? "" : stock.code));
  }, []);

  useEffect(() => {
    if (isDerivative) {
      derivHook.fetchDerivatives({
        category: selectedMajor,
        subCategory: selectedSub,
        sortBy,
        order,
        pageSize: 50,
      });
    } else {
      stockHook.fetchStocks({
        region: selectedMajor,
        subCategory: selectedSub,
        sortBy,
        order,
        pageSize: 50,
      });
    }
  }, [selectedMajor, selectedSub, sortBy, order, page]);

  const hasSelection = selectedMajor !== "";

  return (
    <div className="app">
      <div className="left-panel">
        <Sidebar
          categories={categories}
          selectedMajor={selectedMajor}
          selectedSub={selectedSub}
          sortBy={sortBy}
          order={order}
          loading={catLoading}
          onSelectMajor={handleSelectMajor}
          onSelectSub={handleSelectSub}
          onSortChange={handleSortChange}
          isDerivative={isDerivative}
        />
        {hasSelection ? (
          <StockList
            stocks={stocks}
            loading={loading}
            total={total}
            page={page}
            pages={pages}
            selectedMajor={selectedMajor}
            onPageChange={(p) => setPage(p)}
            onStockClick={handleStockClick}
            activeStock={activeStock}
          />
        ) : (
          <div className="stock-list">
            <div className="empty-state">请选择一个市场分类查看</div>
          </div>
        )}
      </div>
      <div className="right-panel">
        <Globe highlightRegion={selectedMajor} activeStock={activeStock} />
      </div>
    </div>
  );
}
