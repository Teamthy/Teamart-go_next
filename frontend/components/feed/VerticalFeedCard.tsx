"use client";

import Link from "next/link";
import { ArrowRight, Play, Pause } from "lucide-react";
import type { FeedItem } from "@/types/commerce";
import CreatorMeta from "@/components/feed/CreatorMeta";
import FeedActions from "@/components/feed/FeedActions";
import PinnedProductCard from "@/components/feed/PinnedProductCard";
import VideoCard from "@/components/feed/VideoCard";

interface VerticalFeedCardProps {
    item: FeedItem;
    muted: boolean;
    isActive: boolean;
    onLike: () => void;
    onBookmark: () => void;
    onAddToCart: () => void;
    onShare: () => void;
    onComment: () => void;
    onToggleMute: () => void;
    registerVideo: (id: number) => (node: HTMLVideoElement | null) => void;
}

export default function VerticalFeedCard({
    item,
    muted,
    isActive,
    onLike,
    onBookmark,
    onAddToCart,
    onShare,
    onComment,
    onToggleMute,
    registerVideo,
}: VerticalFeedCardProps) {
    return (
        <article className="relative h-[88vh] min-h-[600px] overflow-hidden rounded-[32px] border border-white/10 bg-slate-950 shadow-2xl shadow-black/20">
            <div className="absolute inset-0 overflow-hidden">
                <VideoCard item={item} muted={muted} registerVideo={registerVideo} />
                <div className="absolute inset-0 bg-gradient-to-t from-black/65 via-transparent to-transparent" />
                <div className="absolute inset-x-0 bottom-0 px-5 pb-5">
                    <CreatorMeta
                        creator={item.creator}
                        caption={item.caption}
                        hashtags={item.hashtags}
                        productName={item.product.name}
                    />
                    <div className="mt-4 grid gap-4 lg:grid-cols-[1.4fr_0.95fr]">
                        <PinnedProductCard
                            product={item.product}
                            onAddToCart={onAddToCart}
                            onVisitStore={() => {
                                window.location.href = `/stores/${encodeURIComponent(item.product.merchant_name)}`;
                            }}
                        />
                        <div className="flex flex-col justify-between gap-4 rounded-3xl bg-black/55 p-4 text-white shadow-xl shadow-black/25">
                            <div className="space-y-3">
                                <div className="inline-flex items-center gap-2 rounded-full bg-pink-500/20 px-3 py-1 text-xs uppercase tracking-[0.2em] text-pink-100">
                                    {item.is_live ? "LIVE" : item.kind.toUpperCase()}
                                </div>
                                <div className="flex items-center gap-3 text-sm text-slate-200">
                                    <span className="font-semibold text-white">{item.live_title ?? item.title}</span>
                                    {item.viewers ? <span>{item.viewers.toLocaleString()} viewers</span> : null}
                                </div>
                                <p className="text-sm leading-6 text-slate-300">{item.product.merchant_name} · {item.product.name}</p>
                            </div>
                            <div className="flex items-center justify-between gap-3">
                                <button
                                    type="button"
                                    className="inline-flex items-center gap-2 rounded-full bg-white px-4 py-2 text-sm font-semibold text-slate-900 transition hover:bg-slate-100"
                                    onClick={onToggleMute}
                                >
                                    {muted ? <Play className="h-4 w-4" /> : <Pause className="h-4 w-4" />}
                                    {muted ? "Unmute" : "Mute"}
                                </button>
                                <Link
                                    href={`/products/${item.product.id}`}
                                    className="inline-flex items-center gap-2 rounded-full border border-white/20 bg-white/10 px-4 py-2 text-sm font-semibold text-white transition hover:bg-white/15"
                                >
                                    View product <ArrowRight className="h-4 w-4" />
                                </Link>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div className="absolute right-5 top-1/2 z-20 flex -translate-y-1/2 flex-col items-center gap-3">
                <FeedActions
                    likeCount={item.like_count}
                    commentCount={item.comment_count}
                    liked={item.liked}
                    bookmarked={item.bookmarked}
                    onLike={onLike}
                    onBookmark={onBookmark}
                    onShare={onShare}
                    onComment={onComment}
                    onAddToCart={onAddToCart}
                />
            </div>
            <div className="absolute left-5 top-5 rounded-full bg-black/50 px-3 py-2 text-xs uppercase tracking-[0.24em] text-white shadow-lg shadow-black/30">
                {isActive ? "Watching" : "Swipe up"}
            </div>
        </article>
    );
}
