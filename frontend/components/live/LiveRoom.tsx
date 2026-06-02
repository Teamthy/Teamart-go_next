"use client";

import { useMemo } from "react";
import type { LiveRoomDetails } from "@/types/commerce";
import LiveBadge from "@/components/live/LiveBadge";
import ViewerCount from "@/components/live/ViewerCount";
import LiveChat from "@/components/live/LiveChat";
import ProductDrawer from "@/components/live/ProductDrawer";
import FloatingCartButton from "@/components/live/FloatingCartButton";
import ReactionBubble from "@/components/live/ReactionBubble";
import { Play, Heart, Star } from "lucide-react";

interface LiveRoomProps {
    room: LiveRoomDetails;
    muted: boolean;
    onToggleMute: () => void;
    onReact: (reaction: string) => void;
    onOpenCart: () => void;
    cartCount: number;
    drawerOpen: boolean;
    onCloseDrawer: () => void;
    onAddToCart: (quantity: number) => void;
    onBuyNow: () => void;
    selectedProduct: LiveRoomDetails["pinned_products"][number] | null;
    onSelectProduct: (index: number) => void;
    chatProps: {
        messages: LiveRoomDetails["messages"];
        onSendMessage: (message: string) => void;
        isSending: boolean;
    };
}

export default function LiveRoom({
    room,
    muted,
    onToggleMute,
    onReact,
    onOpenCart,
    cartCount,
    drawerOpen,
    onCloseDrawer,
    onAddToCart,
    onBuyNow,
    selectedProduct,
    onSelectProduct,
    chatProps,
}: LiveRoomProps) {
    const product = selectedProduct ?? room.pinned_products[0] ?? null;
    const progressLabel = useMemo(() => {
        if (!room.elapsed_seconds) return "Starting soon";
        const minutes = Math.floor(room.elapsed_seconds / 60);
        return `${minutes}m ${room.elapsed_seconds % 60}s`;
    }, [room.elapsed_seconds]);

    const reactionEntries = useMemo(() => {
        const defaultReactions: Record<string, number> = {
            "\u2764\uFE0F": 0,
            "\u2728": 0,
            "\uD83D\uDD25": 0,
        };
        return Object.entries(room.reactions ?? defaultReactions);
    }, [room.reactions]);

    return (
        <div className="grid gap-6 xl:grid-cols-[1.45fr_0.95fr]">
            <section className="space-y-6">
                <div className="relative overflow-hidden rounded-[32px] bg-slate-950 shadow-2xl shadow-black/40">
                    <img src={room.thumbnail_url} alt={room.title} className="h-[540px] w-full object-cover" />
                    <div className="absolute inset-0 bg-gradient-to-t from-slate-950/95 via-transparent to-transparent" />
                    <div className="absolute inset-x-0 top-6 px-6">
                        <div className="flex flex-wrap items-center gap-3">
                            <LiveBadge status="LIVE" />
                            <ViewerCount viewers={room.viewers} />
                            <span className="rounded-full bg-white/10 px-3 py-1 text-xs uppercase tracking-[0.22em] text-white">
                                {progressLabel}
                            </span>
                        </div>
                    </div>
                    <div className="absolute inset-x-0 bottom-6 px-6">
                        <div className="flex flex-col gap-4 rounded-[32px] border border-white/10 bg-black/55 p-5 text-white shadow-2xl shadow-black/40">
                            <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                                <div>
                                    <p className="text-sm uppercase tracking-[0.22em] text-slate-300">{room.host.name} — Host</p>
                                    <h1 className="mt-2 text-3xl font-semibold text-white">{room.title}</h1>
                                </div>
                                <button
                                    type="button"
                                    onClick={onToggleMute}
                                    className="inline-flex items-center gap-2 rounded-full bg-white/10 px-4 py-3 text-sm font-semibold text-white transition hover:bg-white/20"
                                >
                                    <Play className="h-4 w-4" />
                                    {muted ? "Unmute" : "Mute"}
                                </button>
                            </div>
                            <div className="flex flex-wrap gap-3 text-sm text-slate-200">
                                <span>{room.description}</span>
                            </div>
                            <div className="grid gap-3 sm:grid-cols-3">
                                <div className="rounded-3xl bg-white/10 p-4">
                                    <p className="text-xs uppercase tracking-[0.24em] text-slate-400">Host message</p>
                                    <p className="mt-2 text-sm text-white">{room.host_message ?? "Live demonstrations and exclusive product drops running now."}</p>
                                </div>
                                <div className="rounded-3xl bg-white/10 p-4">
                                    <p className="text-xs uppercase tracking-[0.24em] text-slate-400">Top reaction</p>
                                    <p className="mt-2 text-sm text-white">{Object.entries(room.reactions ?? {}).sort((a, b) => b[1] - a[1])[0]?.[0] ?? "❤️"}</p>
                                </div>
                                <div className="rounded-3xl bg-white/10 p-4">
                                    <p className="text-xs uppercase tracking-[0.24em] text-slate-400">Pins</p>
                                    <p className="mt-2 text-sm text-white">{room.pinned_products.length} products ready to shop</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div className="grid gap-4 xl:grid-cols-[1fr_0.8fr]">
                    <div className="space-y-4 rounded-[32px] border border-slate-200/70 bg-white p-5 shadow-lg">
                        <div className="flex items-center justify-between gap-3">
                            <div>
                                <p className="text-xs uppercase tracking-[0.24em] text-slate-500">Pinned merchandise</p>
                                <h2 className="mt-2 text-xl font-semibold text-slate-900">{product?.name ?? "No pinned product"}</h2>
                            </div>
                            <div className="rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold uppercase tracking-[0.2em] text-slate-700">
                                {product?.stock ? `${product.stock} in stock` : "Out of stock"}
                            </div>
                        </div>
                        <div className="grid gap-4 sm:grid-cols-[0.9fr_0.5fr]">
                            <img src={product?.image_url ?? "/placeholder-product.png"} alt={product?.name ?? "Product"} className="h-72 w-full rounded-[28px] object-cover" />
                            <div className="space-y-4">
                                <div className="space-y-2 rounded-3xl bg-slate-50 p-4">
                                    <p className="text-sm text-slate-500">Price</p>
                                    <p className="text-2xl font-semibold text-slate-900">${product?.price.toFixed(2) ?? "0.00"}</p>
                                </div>
                                <div className="space-y-2 rounded-3xl bg-slate-50 p-4">
                                    <p className="text-sm text-slate-500">Merchant</p>
                                    <p className="font-semibold text-slate-900">{product?.merchant_name ?? "Brand"}</p>
                                </div>
                                <button
                                    type="button"
                                    onClick={() => onReact("❤️")}
                                    className="inline-flex items-center gap-2 rounded-full bg-pink-500 px-4 py-3 text-sm font-semibold text-white transition hover:bg-pink-600"
                                >
                                    <Heart className="h-4 w-4" />
                                    React
                                </button>
                                <button
                                    type="button"
                                    onClick={() => onReact("✨")}
                                    className="inline-flex items-center gap-2 rounded-full border border-slate-200 bg-white px-4 py-3 text-sm font-semibold text-slate-900 transition hover:bg-slate-50"
                                >
                                    <Star className="h-4 w-4" />
                                    Spotlight
                                </button>
                            </div>
                        </div>
                    </div>
                    <LiveChat messages={chatProps.messages} onSendMessage={chatProps.onSendMessage} isSending={chatProps.isSending} />
                </div>
            </section>
            <aside className="space-y-4">
                <div className="rounded-[32px] border border-slate-200/70 bg-white p-5 shadow-lg">
                    <p className="text-xs uppercase tracking-[0.24em] text-slate-500">Pinned products</p>
                    <div className="mt-4 grid gap-4">
                        {room.pinned_products.map((feedProduct, index) => (
                            <button
                                key={feedProduct.id}
                                type="button"
                                onClick={() => onSelectProduct(index)}
                                className={`flex items-center gap-3 rounded-3xl border px-4 py-3 text-left transition ${selectedProduct?.id === feedProduct.id ? "border-pink-500 bg-pink-50" : "border-slate-200 bg-white hover:bg-slate-50"
                                    }`}
                            >
                                <img src={feedProduct.image_url} alt={feedProduct.name} className="h-16 w-16 rounded-3xl object-cover" />
                                <div>
                                    <p className="text-sm font-semibold text-slate-900">{feedProduct.name}</p>
                                    <p className="text-xs text-slate-500">{feedProduct.merchant_name}</p>
                                </div>
                            </button>
                        ))}
                    </div>
                </div>
                <div className="rounded-[32px] border border-slate-200/70 bg-white p-5 shadow-lg">
                    <p className="text-xs uppercase tracking-[0.24em] text-slate-500">Live reactions</p>
                    <div className="mt-4 flex flex-wrap gap-3">
                        {reactionEntries.map(([emoji, count]) => (
                            <ReactionBubble key={emoji} emoji={emoji} label={`${count} reactions`} />
                        ))}
                    </div>
                </div>
            </aside>
            <FloatingCartButton itemCount={cartCount} onOpen={onOpenCart} />
            <ProductDrawer
                open={drawerOpen}
                product={product}
                onClose={onCloseDrawer}
                onAddToCart={onAddToCart}
                onBuyNow={onBuyNow}
            />
        </div>
    );
}
