"use client";

import { useEffect, useRef } from "react";
import { useQueryClient } from "@tanstack/react-query";
import useAppStore from "@/store/useAppStore";
import type { FeedItem } from "@/hooks/useFeed";
import { BASE } from "@/lib/api";

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

    return {
        id,
        name: String(entry.name ?? entry.title ?? "Unknown product"),
        description: typeof entry.description === "string" ? entry.description : "",
        price: typeof entry.price === "number" ? entry.price : Number(entry.price ?? 0),
        image_url: typeof entry.image_url === "string" ? entry.image_url : typeof entry.image === "string" ? entry.image : "",
        category: typeof entry.category === "string" ? entry.category : undefined,
        score: typeof entry.score === "number" ? entry.score : undefined,
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
                        queryClient.setQueryData<FeedItem[]>(["feed", 50], (current) => {
                            const existing = current ?? [];
                            const alreadyExists = existing.some((item) => item.id === feedItem.id);
                            return alreadyExists ? existing : [feedItem, ...existing];
                        });
                        addFeedUpdates(1);
                        return;
                    }

                    if (Array.isArray(payload)) {
                        const normalizedItems = payload
                            .map(normalizeFeedItem)
                            .filter((item): item is FeedItem => item !== null);
                        if (normalizedItems.length > 0) {
                            queryClient.setQueryData<FeedItem[]>(["feed", 50], (current) => {
                                const existing = current ?? [];
                                const existingIds = new Set(existing.map((item) => item.id));
                                return [
                                    ...normalizedItems.filter((item) => !existingIds.has(item.id)),
                                    ...existing,
                                ];
                            });
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
