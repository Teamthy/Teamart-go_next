"use client";

import { useEffect, useRef } from "react";
import { useQueryClient, type InfiniteData } from "@tanstack/react-query";
import useAppStore from "@/store/useAppStore";
import type { FeedItem } from "@/types/commerce";
import { BASE } from "@/lib/api";

type FeedInfiniteData = InfiniteData<{ items: FeedItem[]; next_offset?: number; has_more?: boolean }>;

interface RealtimeEvent {
    type?: string;
    topic?: string;
    payload?: unknown;
    title?: string;
    body?: string;
    id?: number | string;
    name?: string;
    price?: number;
}

function normalizeFeedItem(value: unknown): FeedItem | null {
    if (!value || typeof value !== "object") return null;

    const entry = value as Record<string, unknown>;
    const id = Number(entry.id ?? entry.product_id ?? entry.item_id);
    if (Number.isNaN(id)) return null;

    const imageUrl = typeof entry.image_url === "string" ? entry.image_url : typeof entry.image === "string" ? entry.image : "";
    const name = String(entry.name ?? entry.title ?? "Unknown product");
    const price = typeof entry.price === "number" ? entry.price : Number(entry.price ?? 0);

    return {
        id,
        kind: "image",
        title: name,
        caption: String(entry.description ?? ""),
        hashtags: [],
        creator: {
            id: Number(entry.author_id ?? 0),
            name: String(entry.author ?? "Unknown"),
            handle: String(entry.author_handle ?? entry.author ?? "unknown"),
            avatar_url: String(entry.avatar_url ?? ""),
            verified: false,
        },
        product: {
            id,
            name,
            image_url: imageUrl,
            price,
            merchant_name: String(entry.merchant_name ?? entry.author ?? "Unknown"),
        },
        image_url: imageUrl,
        like_count: Number(entry.like_count ?? 0),
        comment_count: Number(entry.comment_count ?? 0),
        bookmarked: false,
        liked: false,
        viewers: typeof entry.viewers === "number" ? entry.viewers : Number(entry.viewer_count ?? 0),
        is_live: typeof entry.is_live === "boolean" ? entry.is_live : false,
        live_title: typeof entry.live_title === "string" ? entry.live_title : undefined,
        live_status: typeof entry.live_status === "string" ? entry.live_status : undefined,
    };
}

function isNotificationPayload(value: unknown): value is { title: string; body: string; type?: string;[key: string]: unknown } {
    return (
        !!value &&
        typeof value === "object" &&
        value !== null &&
        "title" in value &&
        "body" in value &&
        typeof (value as Record<string, unknown>).title === "string" &&
        typeof (value as Record<string, unknown>).body === "string"
    );
}

function mergeFeedItems(current: FeedInfiniteData | undefined, items: FeedItem[]): FeedInfiniteData {
    const normalizedItems = items.filter((item, index) => items.findIndex((other) => other.id === item.id) === index);
    if (normalizedItems.length === 0) {
        return current ?? { pages: [{ items: [], next_offset: undefined, has_more: true }], pageParams: [] };
    }

    if (!current || current.pages.length === 0) {
        return {
            pages: [{ items: normalizedItems, next_offset: undefined, has_more: true }],
            pageParams: [],
        };
    }

    const existingIds = new Set(current.pages.flatMap((page) => page.items.map((item) => item.id)));
    const uniqueNewItems = normalizedItems.filter((item) => !existingIds.has(item.id));
    if (uniqueNewItems.length === 0) {
        return current;
    }

    const firstPage = current.pages[0] ?? { items: [], next_offset: undefined, has_more: true };
    const updatedFirstPage = {
        ...firstPage,
        items: [...uniqueNewItems, ...firstPage.items],
    };

    return {
        ...current,
        pages: [updatedFirstPage, ...current.pages.slice(1)],
    };
}

export function useRealtime() {
    const queryClient = useQueryClient();
    const addNotification = useAppStore((state) => state.addNotification);
    const addFeedUpdates = useAppStore((state) => state.addFeedUpdates);
    const setSocketStatus = useAppStore((state) => state.setSocketStatus);

    const websocketRef = useRef<WebSocket | null>(null);
    const reconnectRef = useRef(0);

    useEffect(() => {
        if (typeof window === "undefined") return;

        const token = localStorage.getItem("access_token");
        if (!token) {
            setSocketStatus("idle");
            return;
        }

        let isCancelled = false;

        const connect = () => {
            if (isCancelled) return;

            setSocketStatus("connecting");
            const wsBase = BASE.replace(/^http/, "ws").replace(/^https/, "wss");
            const socket = new WebSocket(`${wsBase}/ws?token=${encodeURIComponent(token)}`);
            websocketRef.current = socket;

            socket.onopen = () => {
                setSocketStatus("open");
                reconnectRef.current = 0;
                socket.send(JSON.stringify({ type: "subscribe", topic: "notifications" }));
                socket.send(JSON.stringify({ type: "subscribe", topic: "feed" }));
            };

            socket.onmessage = (event) => {
                try {
                    const parsed = JSON.parse(event.data.toString()) as RealtimeEvent;
                    const payload = parsed.payload ?? parsed;

                    if (parsed.type === "error") {
                        return;
                    }

                    if (
                        parsed.type === "pong" ||
                        parsed.type === "subscribed" ||
                        parsed.type === "unsubscribed" ||
                        parsed.type === "published" ||
                        parsed.type === "message_sent"
                    ) {
                        return;
                    }

                    if (isNotificationPayload(payload)) {
                        addNotification({
                            id: String(payload.id ?? `notif-${Date.now()}`),
                            title: payload.title,
                            body: payload.body,
                            type: payload.type || parsed.type || "notification",
                            payload,
                            receivedAt: new Date().toISOString(),
                        });
                        return;
                    }

                    const feedItem = normalizeFeedItem(payload);
                    if (feedItem) {
                        queryClient.setQueriesData<FeedInfiniteData>({ queryKey: ["feed"], exact: false }, (current) =>
                            mergeFeedItems(current, [feedItem])
                        );
                        addFeedUpdates(1);
                        return;
                    }

                    if (Array.isArray(payload)) {
                        const normalizedItems = payload
                            .map(normalizeFeedItem)
                            .filter((item): item is FeedItem => item !== null);
                        if (normalizedItems.length > 0) {
                            queryClient.setQueriesData<FeedInfiniteData>({ queryKey: ["feed"], exact: false }, (current) =>
                                mergeFeedItems(current, normalizedItems)
                            );
                            addFeedUpdates(normalizedItems.length);
                        }
                    }
                } catch {
                    return;
                }
            };

            socket.onerror = () => {
                setSocketStatus("error");
            };

            socket.onclose = () => {
                setSocketStatus("closed");
                if (!isCancelled) {
                    const retryDelay = Math.min(10000, 1000 * 2 ** reconnectRef.current);
                    reconnectRef.current += 1;
                    window.setTimeout(connect, retryDelay);
                }
            };
        };

        connect();

        return () => {
            isCancelled = true;
            websocketRef.current?.close();
        };
    }, [addFeedUpdates, addNotification, queryClient, setSocketStatus]);
}
