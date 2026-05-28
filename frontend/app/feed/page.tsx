"use client";

<<<<<<< HEAD
import { useMemo, useState } from "react";
import Link from "next/link";
import PageHeader from "@/components/ui/PageHeader";
import FeedCard from "@/components/ui/FeedCard";
import StatCard from "@/components/ui/StatCard";
import Tabs from "@/components/ui/Tabs";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import Badge from "@/components/ui/badge";
import { getStoredCustomer } from "@/lib/auth-state";
import { feedItems } from "@/lib/mock/feed";
import { creators } from "@/lib/mock/creators";
import { products } from "@/lib/mock/products";
import { stores } from "@/lib/mock/stores";

const tabOptions = [
    { label: "For you", value: "all" },
    { label: "Live", value: "livestream now" },
    { label: "Creators", value: "creator post" },
    { label: "Merchants", value: "merchant spotlight" },
    { label: "Reviews", value: "customer review" },
    { label: "Promos", value: "product promo" },
];

const trendingTerms = ["#livebundle", "#creatordrop", "#beauty", "#giftable", "#limitedstock"];

export default function FeedPage() {
    const [activeTab, setActiveTab] = useState("all");
    const customer = getStoredCustomer();

    const filteredItems = useMemo(() => {
        return activeTab === "all" ? feedItems : feedItems.filter((item) => item.kind === activeTab);
    }, [activeTab]);

    const liveCount = feedItems.filter((item) => item.kind === "livestream now").length;

    return (
        <div className="space-y-8 pb-10">
            <PageHeader
                title="For you feed"
                description={customer ? `Welcome back, ${customer.firstName}. Your ${customer.favoriteCategory ?? "fashion"} feed is tuned for live rooms, creator drops, and checkout-ready moments.` : "A TikTok-native social shopping experience with live rooms, creator bundles, merchant spotlights, and instant checkout paths."}
                actions={
                    <>
                        <Button asChild variant="primary">
                            <Link href="/live">Jump into live</Link>
                        </Button>
                        <Button asChild variant="secondary">
                            <Link href="/search">Search now</Link>
                        </Button>
                    </>
                }
            />

            <div className="grid gap-4 md:grid-cols-3">
                <StatCard label="Live rooms" value={String(liveCount)} helper="Creator and merchant rooms are on the move right now." />
                <StatCard label="Storefront picks" value="20" helper="Fresh content is curated to keep discovery and checkout aligned." />
                <StatCard label="Saved moments" value="8" helper="Your social commerce feed stays tuned for fast, useful shopping decisions." />
            </div>

            <Tabs tabs={tabOptions} active={activeTab} onChange={setActiveTab} />

            <div className="grid gap-6 xl:grid-cols-[0.9fr_1.1fr_0.9fr]">
                <div className="space-y-4 xl:sticky xl:top-24 xl:self-start">
                    <Card className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Discover</p>
                        <h2 className="mt-3 text-xl font-semibold text-zinc-900">Swipe through the feed and shop in one motion</h2>
                        <p className="mt-3 text-sm leading-6 text-zinc-600">
                            The left rail keeps discovery warm with trending terms, live suggestions, and creator-led shortcuts for the next scroll.
                        </p>
                        <div className="mt-4 flex flex-wrap gap-2">
                            {trendingTerms.map((term) => (
                                <Link key={term} href="/search" className="rounded-full border border-zinc-200 bg-zinc-50 px-3 py-2 text-sm font-semibold text-zinc-700">
                                    {term}
                                </Link>
                            ))}
                        </div>
                    </Card>

                    <Card className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Live now</p>
                        <div className="mt-3 space-y-3">
                            {feedItems.filter((item) => item.kind === "livestream now").slice(0, 3).map((item) => (
                                <Link key={item.id} href={item.actionHref} className="block rounded-[24px] bg-[#FFF8FB] p-3">
                                    <div className="flex items-center justify-between gap-3">
                                        <div>
                                            <p className="text-sm font-semibold text-zinc-900">{item.title}</p>
                                            <p className="mt-1 text-xs text-zinc-600">{item.description}</p>
                                        </div>
                                        <Badge tone="success">LIVE</Badge>
                                    </div>
                                </Link>
                            ))}
                        </div>
                    </Card>
                </div>

                <div className="space-y-4">
                    {filteredItems.length === 0 ? (
                        <Card className="p-8 text-center text-zinc-600">
                            No moments match this view yet. Try another tab to keep the feed moving.
                        </Card>
                    ) : (
                        <div className="max-h-[calc(100vh-16rem)] space-y-4 overflow-y-auto pr-1 snap-y snap-mandatory">
                            {filteredItems.map((item) => (
                                <FeedCard key={item.id} item={item} />
                            ))}
                        </div>
                    )}
                </div>

                <div className="space-y-4 xl:sticky xl:top-24 xl:self-start">
                    <Card className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Sticky cart</p>
                        <div className="mt-3 rounded-[24px] bg-zinc-950 p-4 text-white">
                            <p className="text-sm text-white/75">Quick shop</p>
                            <p className="mt-3 text-2xl font-semibold">$86</p>
                            <p className="mt-2 text-sm text-white/80">Three high-conversion picks are ready for checkout in one tap.</p>
                            <Button asChild variant="primary" className="mt-4 w-full">
                                <Link href="/cart">Open cart</Link>
                            </Button>
                        </div>
                    </Card>

                    <Card className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Recommended stores</p>
                        <div className="mt-3 space-y-3">
                            {stores.slice(0, 3).map((store) => (
                                <Link key={store.slug} href={`/stores/${store.slug}`} className="block rounded-[24px] bg-[#FFF8FB] p-3">
                                    <p className="text-sm font-semibold text-zinc-900">{store.name}</p>
                                    <p className="mt-1 text-xs text-zinc-600">{store.category} • {store.live}</p>
                                </Link>
                            ))}
                        </div>
                    </Card>

                    <Card className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Top creators</p>
                        <div className="mt-3 space-y-3">
                            {creators.map((creator) => (
                                <Link key={creator.id} href={`/creator/${creator.id}`} className="flex items-center gap-3 rounded-[24px] bg-zinc-50 p-3">
                                    <img src={creator.avatar} alt={creator.name} className="h-10 w-10 rounded-full object-cover" />
                                    <div>
                                        <p className="text-sm font-semibold text-zinc-900">{creator.name}</p>
                                        <p className="text-xs text-zinc-500">{creator.followers}</p>
                                    </div>
                                </Link>
                            ))}
                        </div>
                    </Card>

                    <Card className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Best in cart</p>
                        <div className="mt-3 space-y-3">
                            {products.slice(0, 3).map((product) => (
                                <Link key={product.id} href={`/products/${product.id}`} className="block rounded-[24px] bg-zinc-50 p-3">
                                    <p className="text-sm font-semibold text-zinc-900">{product.name}</p>
                                    <p className="mt-1 text-xs text-zinc-600">{product.price} • {product.category}</p>
                                </Link>
                            ))}
                        </div>
                    </Card>
                </div>
            </div>
=======
import Link from "next/link";
import ProductCard from "@/components/product/ProductCard";
import { useFeed } from "@/hooks/useFeed";
import { ArrowRight } from "lucide-react";
import { useState } from "react";

const tabs = ["For you", "Live drops", "Trending", "Creators"];

export default function FeedPage() {
    const { items, isLoading, error } = useFeed(50);
    const [activeTab, setActiveTab] = useState(tabs[0]);

    const feedProducts = items.map((item) => ({
        id: String(item.id),
        name: item.name,
        description: item.description ?? "Creator spotlight product",
        price: `$${item.price.toFixed(2)}`,
        image: item.image_url ?? "/images/placeholder-product.png",
    }));

    return (
        <div className="min-h-screen bg-slate-950 text-white">
            <section className="bg-[radial-gradient(circle_at_top,_rgba(192,132,252,0.18),transparent_30%),linear-gradient(180deg,#020617_0%,#0f172a_100%)] px-4 py-16 sm:px-6 lg:px-8">
                <div className="mx-auto max-w-6xl">
                    <div className="grid gap-10 lg:grid-cols-[0.95fr_0.55fr] lg:items-end">
                        <div className="space-y-5">
                            <span className="inline-flex items-center rounded-full bg-fuchsia-500/10 px-4 py-2 text-sm font-semibold text-fuchsia-200 tracking-[0.28em]">
                                Live feed
                            </span>
                            <h1 className="text-4xl font-semibold tracking-tight text-white sm:text-5xl">
                                Discover products, creators, and live drops in one immersive feed.
                            </h1>
                            <p className="max-w-2xl text-base leading-8 text-slate-300">
                                Teamart surfaces the freshest creator commerce, livestream shopping, and trending storefront experiences in a mobile-first discovery flow.
                            </p>
                            <div className="flex flex-wrap gap-3">
                                <Link href="/auth/login" className="rounded-full border border-white/10 bg-white/5 px-6 py-3 text-sm font-semibold text-white transition hover:bg-white/10">
                                    Sign in
                                </Link>
                                <Link href="/auth/register" className="rounded-full bg-fuchsia-500 px-6 py-3 text-sm font-semibold text-white transition hover:bg-fuchsia-400">
                                    Create account
                                </Link>
                            </div>
                        </div>
                        <div className="rounded-[2.5rem] border border-white/10 bg-slate-900/90 p-6 shadow-2xl shadow-fuchsia-500/10">
                            <p className="text-sm uppercase tracking-[0.35em] text-slate-400">Now playing</p>
                            <div className="mt-6 space-y-4">
                                {feedProducts.slice(0, 3).map((product) => (
                                    <div key={product.id} className="rounded-3xl border border-white/10 bg-slate-950/80 p-4">
                                        <p className="text-sm text-slate-400">{product.name}</p>
                                        <p className="mt-2 text-xl font-semibold text-white">{product.price}</p>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </div>
                </div>
            </section>

            <section className="px-4 py-12 sm:px-6 lg:px-8">
                <div className="mx-auto max-w-6xl rounded-[2.5rem] border border-white/10 bg-slate-950/95 p-6 shadow-2xl shadow-slate-950/30">
                    <div className="flex flex-wrap items-center justify-between gap-4 border-b border-white/10 pb-4">
                        <div>
                            <h2 className="text-xl font-semibold text-white">Feed controls</h2>
                            <p className="text-sm text-slate-400">Filter what you see in the discovery stream.</p>
                        </div>
                        <div className="flex flex-wrap gap-2">
                            {tabs.map((tab) => (
                                <button
                                    key={tab}
                                    type="button"
                                    onClick={() => setActiveTab(tab)}
                                    className={`rounded-full px-4 py-2 text-sm font-semibold transition ${activeTab === tab ? "bg-fuchsia-500 text-white" : "bg-white/5 text-slate-300 hover:bg-white/10"}`}
                                >
                                    {tab}
                                </button>
                            ))}
                        </div>
                    </div>

                    {isLoading && (
                        <div className="py-16 text-center text-slate-400">Loading feed...</div>
                    )}

                    {error && !isLoading && (
                        <div className="rounded-3xl border border-yellow-400/20 bg-yellow-400/10 p-6 text-yellow-100">
                            Unable to load the feed right now. Showing available products instead.
                        </div>
                    )}

                    {!isLoading && feedProducts.length === 0 && (
                        <div className="py-16 text-center text-slate-400">No products available yet.</div>
                    )}

                    <div className="mt-8 grid gap-6 md:grid-cols-2 xl:grid-cols-3">
                        {feedProducts.map((product) => (
                            <ProductCard key={product.id} product={product} />
                        ))}
                    </div>

                    <div className="mt-10 flex items-center justify-center gap-2 text-sm text-slate-400">
                        <span>{feedProducts.length} items</span>
                        <ArrowRight className="h-4 w-4" />
                    </div>
                </div>
            </section>
>>>>>>> 8018627 (feat(feed): tiktok-style feed, search and landing experience)
        </div>
    );
}
