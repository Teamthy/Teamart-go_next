"use client";

import { useCallback, useMemo } from "react";
import { useInfiniteQuery, useMutation, useQueryClient, type InfiniteData } from "@tanstack/react-query";
import * as api from "@/lib/api";
import type { FeedItem } from "@/types/commerce";

type FeedInfiniteData = InfiniteData<{ items: FeedItem[]; next_offset?: number; has_more?: boolean }>;

function updateFeedItem(
    data: FeedInfiniteData | undefined,
    itemId: number,
    updater: (item: FeedItem) => FeedItem
): FeedInfiniteData | undefined {
    if (!data) {
        return data;
    }

    return {
        ...data,
        pages: data.pages.map((page) => ({
            ...page,
            items: page.items.map((item) => (item.id === itemId ? updater(item) : item)),
        })),
    };
}

export function useFeed(pageSize = 5) {
    const queryClient = useQueryClient();

    const feedQuery = useInfiniteQuery<{
        items: FeedItem[];
        next_offset?: number;
        has_more?: boolean;
    }, Error, FeedInfiniteData, (string | number)[], number>({
        queryKey: ["feed", pageSize],
        queryFn: async ({ pageParam = 0 }) => {
            return await api.getFeedPage(pageSize, pageParam);
        },
        getNextPageParam: (lastPage) => {
            if (typeof lastPage.has_more === "boolean") {
                return lastPage.has_more ? lastPage.next_offset ?? undefined : undefined;
            }
            return lastPage.items.length === pageSize ? lastPage.next_offset ?? undefined : undefined;
        },
        initialPageParam: 0,
        staleTime: 1000 * 60 * 2,
        refetchOnWindowFocus: false,
    });

    const items = useMemo<FeedItem[]>(() => feedQuery.data?.pages.flatMap((page) => page.items) ?? [], [feedQuery.data]);

    const likeMutation = useMutation({
        mutationFn: async (itemId: number) => api.likeFeedItem(itemId),
        onMutate: async (itemId) => {
            await queryClient.cancelQueries({ queryKey: ["feed", pageSize] });
            const previousData = queryClient.getQueryData<FeedInfiniteData>(["feed", pageSize]);
            queryClient.setQueryData<FeedInfiniteData>(["feed", pageSize], (current) =>
                updateFeedItem(current, itemId, (item) => ({
                    ...item,
                    liked: !item.liked,
                    like_count: item.liked ? Math.max(0, item.like_count - 1) : item.like_count + 1,
                }))
            );
            return { previousData };
        },
        onError: (_error, _itemId, context) => {
            if (context?.previousData) {
                queryClient.setQueryData(["feed", pageSize], context.previousData);
            }
        },
        onSettled: () => {
            queryClient.invalidateQueries({ queryKey: ["feed", pageSize] });
        },
    });

    const bookmarkMutation = useMutation({
        mutationFn: async (itemId: number) => api.bookmarkFeedItem(itemId),
        onMutate: async (itemId) => {
            await queryClient.cancelQueries({ queryKey: ["feed", pageSize] });
            const previousData = queryClient.getQueryData<FeedInfiniteData>(["feed", pageSize]);
            queryClient.setQueryData<FeedInfiniteData>(["feed", pageSize], (current) =>
                updateFeedItem(current, itemId, (item) => ({
                    ...item,
                    bookmarked: !item.bookmarked,
                }))
            );
            return { previousData };
        },
        onError: (_error, _itemId, context) => {
            if (context?.previousData) {
                queryClient.setQueryData(["feed", pageSize], context.previousData);
            }
        },
        onSettled: () => {
            queryClient.invalidateQueries({ queryKey: ["feed", pageSize] });
        },
    });

    const toggleLike = useCallback((itemId: number) => {
        likeMutation.mutate(itemId);
    }, [likeMutation]);

    const toggleBookmark = useCallback((itemId: number) => {
        bookmarkMutation.mutate(itemId);
    }, [bookmarkMutation]);

    return {
        items,
        isLoading: feedQuery.isLoading,
        isFetchingNextPage: feedQuery.isFetchingNextPage,
        hasNextPage: Boolean(feedQuery.hasNextPage),
        fetchNextPage: feedQuery.fetchNextPage,
        isLiking: likeMutation.isPending,
        isBookmarking: bookmarkMutation.isPending,
        toggleLike,
        toggleBookmark,
        error: feedQuery.isError ? feedQuery.error?.message ?? "Unable to load feed" : null,
        refetch: feedQuery.refetch,
    };
}
