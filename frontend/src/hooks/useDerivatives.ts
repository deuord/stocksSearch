import { useState, useCallback } from "react";
import type { Stock, DerivativeResponse, DerivativeItem } from "../types";

const API_BASE = "/api";

interface UseDerivativesParams {
  category: string;
  subCategory: string;
  sortBy: string;
  order: string;
  pageSize: number;
}

export function useDerivatives() {
  const [stocks, setStocks] = useState<Stock[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pages, setPages] = useState(0);

  const fetchDerivatives = useCallback(
    async (params: UseDerivativesParams) => {
      if (!params.category) {
        setStocks([]);
        setLoading(false);
        return;
      }

      setLoading(true);

      const query = new URLSearchParams();
      query.set("category", params.category);
      if (params.subCategory) query.set("sub_category", params.subCategory);
      query.set("sort_by", params.sortBy || "price");
      query.set("order", params.order || "desc");
      query.set("page", String(page));
      query.set("page_size", String(params.pageSize || 50));

      try {
        const res = await fetch(`${API_BASE}/derivatives?${query}`);
        const data: DerivativeResponse = await res.json();

        const converted: Stock[] = (data.items || []).map(
          (item: DerivativeItem) => ({
            code: item.code,
            name: item.name,
            region: item.category,
            price: item.price,
            change: item.change,
            change_percent: item.change_percent,
            volume: item.volume,
            market_cap: item.market_cap,
            currency: item.currency,
            sub_category: "",
            updated_at: item.updated_at,
          }),
        );

        setStocks(converted);
        setTotal(data.total);
        setPages(data.pages);
      } catch (err) {
        console.error("获取衍生品列表失败:", err);
        setStocks([]);
      } finally {
        setLoading(false);
      }
    },
    [page],
  );

  return { stocks, loading, total, page, pages, setPage, fetchDerivatives };
}
