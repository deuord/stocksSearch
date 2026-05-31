import { useState, useEffect } from "react";
import type { Category } from "../types";

const API_BASE = "/api";

export function useCategories() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`${API_BASE}/categories`)
      .then((res) => res.json())
      .then((data) => {
        setCategories(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error("获取分类失败:", err);
        setLoading(false);
      });
  }, []);

  return { categories, loading };
}
