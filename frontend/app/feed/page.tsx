"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import type { FeedItem } from "@/types/commerce";
import Link from "next/link";
import * as api from "@/lib/api";
import { useFeed } from "@/hooks/useFeed";
import { useLiveFeed } from "@/hooks/useLiveFeed";
import { useVideoPlayback } from "@/hooks/useVideoPlayback";
import useAppStore from "@/store/useAppStore";
import VerticalFeedCard from "@/components/feed/VerticalFeedCard";
import LiveCard from "@/components/live/LiveCard";
import Button from "@/components/ui/button";

export default function FeedPage() {
    const { items, isLoading, hasNextPage, fetchNextPage, isFetchingNextPage, toggleLike, toggleBookmark } = useFeed(6);
    const { rooms } = useLiveFeed();
    const { registerVideo, muted, toggleMuted } = useVideoPlayback();
    const liveFeedCount = useAppStore((state) => state.liveFeedCount);

    const [activeIndex, setActiveIndex] = useState(0);
    const [cartCount, setCartCount] = useState(0);
    const [errorMessage, setErrorMessage] = useState<string | null>(null);
    const loadMoreRef = useRef<HTMLDivElement | null>(null);
    const cardRefs = useRef(new Map<number, HTMLDivElement>());

    const setCardRef = useCallback((cardId: number) => (node: HTMLDivElement | null) => {
        if (node) {
            cardRefs.current.set(cardId, node);
            return;
        }
        cardRefs.current.delete(cardId);
    }, []);

    useEffect(() => {
        if (!loadMoreRef.current || !hasNextPage) return;

        const observer = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting) {
                    fetchNextPage();
                }
            },
            {
                rootMargin: "240px",
            }
        );

        observer.observe(loadMoreRef.current);
        return () => observer.disconnect();
    }, [fetchNextPage, hasNextPage]);

    const handleAddToCart = useCallback(async (productId: number) => {
        try {
            await api.addCartItem({ product_id: productId, quantity: 1 });
            setCartCount((count) => count + 1);
        } catch (error) {
            setErrorMessage(error instanceof Error ? error.message : "Failed to add product to cart.");
        }
    }, []);

    const handleShare = useCallback(() => {
        if (typeof window === "undefined") return;
        navigator.clipboard
            .writeText(window.location.href)
            .catch(() => setErrorMessage("Unable to copy share link."));
    }, []);

    const handleComment = useCallback(() => {
        setErrorMessage("Comment features are available inside the live room experience.");
    }, []);

    const handleWheel = useCallback(
        (event: React.WheelEvent<HTMLDivElement>) => {
            if (Math.abs(event.deltaY) < 24 || items.length === 0) {
                return;
            }

            event.preventDefault();
            const nextIndex = event.deltaY > 0 ? Math.min(items.length - 1, activeIndex + 1) : Math.max(0, activeIndex - 1);
            setActiveIndex(nextIndex);
            const card = cardRefs.current.get(items[nextIndex]?.id);
            card?.scrollIntoView({ behavior: "smooth", block: "center" });
        },
        [activeIndex, items]
    );

    const visibleItems = useMemo<FeedItem[]>(() => items, [items]);

    return (
        <div className="space-y-8 pb-10">
            <section className="rounded-[32px] border border-slate-200 bg-white p-8 shadow-xl shadow-slate-200/50">
                <div className="flex flex-col gap-6 xl:flex-row xl:items-center xl:justify-between">
                    <div className="max-w-3xl">
                        <p className="text-sm uppercase tracking-[0.3em] text-pink-600">For you</p>
                        <h1 className="mt-4 text-4xl font-semibold tracking-tight text-slate-950">Vertical social commerce feed</h1>
                        <p className="mt-4 max-w-2xl text-base leading-8 text-slate-600">
                            Discover creator-driven product drops, immersive video moments, and live shopping rooms inside a premium TikTok-style experience.
                        </p>
                    </div>
                    <div className="flex flex-wrap gap-3">
                        <Button variant="secondary" onClick={() => window.location.assign("/live")}>Open live</Button>
                        <Button variant="primary" asChild>
                            <Link href="/search">Search shop</Link>
                        </Button>
                    </div>
                </div>
                <div className="mt-6 grid gap-4 sm:grid-cols-3">
                    <div className="rounded-[28px] bg-slate-950 p-5 text-white shadow-lg">
                        <p className="text-sm uppercase tracking-[0.2em] text-slate-400">Live pulse</p>
                        <p className="mt-3 text-3xl font-semibold">{liveFeedCount}</p>
                        <p className="mt-2 text-sm text-slate-300">Realtime feed updates and fresh commerce moments.</p>
                    </div>
                    <div className="rounded-[28px] bg-slate-950 p-5 text-white shadow-lg">
                        <p className="text-sm uppercase tracking-[0.2em] text-slate-400">Active posts</p>
                        <p className="mt-3 text-3xl font-semibold">{visibleItems.length}</p>
                        <p className="mt-2 text-sm text-slate-300">Swipe vertically, settle on the product, and checkout faster.</p>
                    </div>
                    <div className="rounded-[28px] bg-slate-950 p-5 text-white shadow-lg">
                        <p className="text-sm uppercase tracking-[0.2em] text-slate-400">Cart</p>
                        <p className="mt-3 text-3xl font-semibold">{cartCount}</p>
                        <p className="mt-2 text-sm text-slate-300">Optimistic add-to-cart updates from the feed.</p>
                    </div>
                </div>
            </section>

            <div className="grid gap-6 xl:grid-cols-[1.45fr_0.9fr]">
                <div className="rounded-[32px] border border-slate-200 bg-slate-950 p-4 shadow-2xl shadow-black/10">
                    <div className="mb-4 flex items-center justify-between gap-3 px-4 text-sm text-slate-200">
                        <span className="font-semibold text-white">Swipe feed</span>
                        <span>{isLoading ? "Loading feed..." : `${visibleItems.length} posts`}</span>
                    </div>
                    <div className="h-[83vh] overflow-y-auto snap-y snap-mandatory space-y-6 pr-4" onWheel={handleWheel}>
                        {visibleItems.map((item, index) => (
                            <div key={item.id} ref={setCardRef(item.id)} className="snap-center px-4">
                                <VerticalFeedCard
                                    item={item}
                                    isActive={activeIndex === index}
                                    muted={muted}
                                    registerVideo={registerVideo}
                                    onLike={() => toggleLike(item.id)}
                                    onBookmark={() => toggleBookmark(item.id)}
                                    onAddToCart={() => handleAddToCart(item.product.id)}
                                    onShare={handleShare}
                                    onComment={handleComment}
                                    onToggleMute={toggleMuted}
                                />
                            </div>
                        ))}
                        <div ref={loadMoreRef} className="flex h-24 items-center justify-center text-sm text-slate-300">
                            {hasNextPage ? (isFetchingNextPage ? "Loading more content..." : "Scroll to load more") : "You have reached the freshest feed content."}
                        </div>
                    </div>
                </div>
                <aside className="space-y-6">
                    <div className="rounded-[28px] border border-slate-200 bg-white p-6 shadow-lg">
                        <p className="text-sm uppercase tracking-[0.24em] text-pink-600">Live commerce</p>
                        <h2 className="mt-3 text-2xl font-semibold text-slate-950">Creator rooms and product drops</h2>
                        <p className="mt-4 text-sm leading-7 text-slate-600">
                            Curated live streams, product spotlight cards, and a clear path from discovery to checkout.
                        </p>
                        <div className="mt-5 flex flex-wrap gap-3">
                            <Button variant="primary" onClick={() => window.location.assign("/live")}>Watch now</Button>
                            <Button variant="secondary" asChild>
                                <Link href="/search">Browse products</Link>
                            </Button>
                        </div>
                    </div>
                    <div className="space-y-4">
                        {rooms.map((room) => (
                            <LiveCard key={room.id} room={room} />
                        ))}
                    </div>
                </aside>
            </div>
            {errorMessage ? (
                <div className="rounded-3xl bg-rose-500/10 px-4 py-3 text-sm text-rose-700">{errorMessage}</div>
            ) : null}
        </div>
    );
}
