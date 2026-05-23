"use client";

import { useState, useEffect } from "react";
import * as api from "@/lib/api";

export interface FeedItem {
    id: number;
    name: string;
    description?: string;
    price: number;
    image_url?: string;
    category?: string;
    score?: number;
}

export function useFeed(limit = 50) {
    const [items, setItems] = useState<FeedItem[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchFeed = async () => {
            setIsLoading(true);
            setError(null);
            try {
                const response = await api.getFeed(limit);
                setItems(response || []);
            } catch (err: any) {
                setError(err.message || "Failed to fetch feed");
                // Fallback to products list if feed fails
                try {
                    const productsRes = await api.listProducts(limit, 0);
                    setItems(productsRes.products || []);
                } catch (fallbackErr) {
                    console.error("Error fetching feed and products:", fallbackErr);
                }
            } finally {
                setIsLoading(false);
            }
        };

        fetchFeed();
    }, [limit]);

    return { items, isLoading, error };
}
