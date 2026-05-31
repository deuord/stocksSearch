import { useState, useCallback, useEffect, useRef } from "react";
import type { Stock, StockListResponse } from "../types";

const API_BASE = "/api";

interface UseStocksParams {
  region: string;
  subCategory: string;
  sortBy: string;
  order: string;
  pageSize: number;
}

interface WSStockData {
  type: string;
  code: string;
  region: string;
  price: number;
  change: number;
  change_percent: number;
  volume: number;
  open: number;
  high: number;
  low: number;
  prev_close: number;
}

export function useStocks() {
  const [stocks, setStocks] = useState<Stock[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pages, setPages] = useState(0);
  const wsRef = useRef<WebSocket | null>(null);
  const codesRef = useRef<string[]>([]);
  const regionRef = useRef<string>("");

  useEffect(() => {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const ws = new WebSocket(`${protocol}//${window.location.host}/ws`);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log("✅ WebSocket 已连接");
      if (codesRef.current.length > 0 && regionRef.current) {
        ws.send(
          JSON.stringify({
            action: "subscribe",
            region: regionRef.current,
            codes: codesRef.current,
          }),
        );
      }
    };

    ws.onmessage = (event) => {
      try {
        const data: WSStockData = JSON.parse(event.data);
        if (data.type !== "quote") return;

        setStocks((prev) =>
          prev.map((s) =>
            s.code === data.code
              ? {
                  ...s,
                  price: data.price,
                  change: data.change,
                  change_percent: data.change_percent,
                  volume: Math.round(data.volume),
                }
              : s,
          ),
        );
      } catch {
        // 忽略解析错误
      }
    };

    ws.onclose = () => {
      console.log("🔌 WebSocket 已断开");
    };

    return () => {
      ws.close();
    };
  }, []);

  const fetchStocks = useCallback(
    async (params: UseStocksParams) => {
      if (!params.region) {
        setStocks([]);
        setLoading(false);
        codesRef.current = [];
        regionRef.current = "";
        return;
      }

      setLoading(true);

      const query = new URLSearchParams();
      query.set("region", params.region);
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

        // WebSocket 订阅实时数据
        const codes = (data.stocks || []).map((s: Stock) => s.code);
        codesRef.current = codes;
        regionRef.current = params.region;

        if (wsRef.current?.readyState === WebSocket.OPEN && codes.length > 0) {
          wsRef.current.send(
            JSON.stringify({
              action: "subscribe",
              region: params.region,
              codes,
            }),
          );
        }
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
