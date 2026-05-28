"use client";

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
        </div>
    );
}
