import { useState, useCallback } from "react";
import type { Stock, StockListResponse } from "../types";

const API_BASE = "/api";

interface UseStocksParams {
  region: string;
  subCategory: string;
  sortBy: string;
  order: string;
  pageSize: number;
}

export function useStocks() {
  const [stocks, setStocks] = useState<Stock[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pages, setPages] = useState(0);

  const fetchStocks = useCallback(
    async (params: UseStocksParams) => {
      setLoading(true);

      const query = new URLSearchParams();
      if (params.region) query.set("region", params.region);
      if (params.subCategory) query.set("sub_category", params.subCategory);
      query.set("sort_by", params.sortBy || "market_cap");
      query.set("order", params.order || "desc");
      query.set("page", String(page));
      query.set("page_size", String(params.pageSize || 50));

      try {
        const res = await fetch(`${API_BASE}/stocks?${query}`);
        const data: StockListResponse = await res.json();
        setStocks(data.stocks || []);
        setTotal(data.total);
        setPages(data.pages);
      } catch (err) {
        console.error("获取股票列表失败:", err);
        setStocks([]);
      } finally {
        setLoading(false);
      }
    },
    [page],
  );

  return { stocks, loading, total, page, pages, setPage, fetchStocks };
}
