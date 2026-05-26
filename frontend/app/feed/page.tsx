"use client";

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
        </div>
    );
}
