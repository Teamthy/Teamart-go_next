"use client";

import { useQuery } from "@tanstack/react-query";
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
    const query = useQuery<FeedItem[], Error>({
        queryKey: ["feed", limit],
        queryFn: async () => {
            try {
                return await api.getFeed(limit);
            } catch (error) {
                const fallback = await api.listProducts(limit, 0);
                return fallback.products || [];
            }
        },
        staleTime: 1000 * 60 * 2,
        retry: 1,
        refetchOnWindowFocus: false,
    });

    return {
        items: (query.data ?? []) as FeedItem[],
        isLoading: query.isLoading,
        error: query.isError ? (query.error as Error)?.message ?? "Unable to load feed" : null,
        refetch: query.refetch,
    };
}
