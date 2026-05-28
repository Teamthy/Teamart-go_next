"use client";

import { useMemo, useState } from "react";
import Link from "next/link";
import ProductCard from "@/components/product/ProductCard";
import SectionHeader from "@/components/ui/SectionHeader";
import Button from "@/components/ui/button";
import Input from "@/components/ui/input";
import Badge from "@/components/ui/badge";
import Card from "@/components/ui/card";
import { creators } from "@/lib/mock/creators";
import { feedItems } from "@/lib/mock/feed";
import { products } from "@/lib/mock/products";
import { stores } from "@/lib/mock/stores";

const tabs = ["Top", "Products", "Live", "Creators", "Videos", "Stores"];
const recentSearches = ["glow set", "ring light", "hoodie", "planner", "gift set"];
const trendingTerms = ["#livebundle", "#creatorpick", "#bestseller", "#giftable", "#fastdelivery"];

export default function SearchPage() {
    const [query, setQuery] = useState("");
    const [activeTab, setActiveTab] = useState("Top");

    const searchText = query.trim().toLowerCase();

    const filteredProducts = useMemo(() => {
        if (!searchText) {
            return products.slice(0, 4);
        }

        return products.filter((product) => `${product.name} ${product.description} ${product.category ?? ""}`.toLowerCase().includes(searchText));
    }, [searchText]);

    const liveResults = useMemo(() => {
        if (!searchText) {
            return feedItems.filter((item) => item.kind === "livestream now").slice(0, 3);
        }

        return feedItems.filter((item) => item.kind === "livestream now" && `${item.title} ${item.description}`.toLowerCase().includes(searchText));
    }, [searchText]);

    const creatorResults = useMemo(() => {
        if (!searchText) {
            return creators.slice(0, 3);
        }

        return creators.filter((creator) => `${creator.name} ${creator.bio}`.toLowerCase().includes(searchText));
    }, [searchText]);

    const videoResults = useMemo(() => {
        if (!searchText) {
            return feedItems.slice(0, 3);
        }

        return feedItems.filter((item) => `${item.title} ${item.description}`.toLowerCase().includes(searchText));
    }, [searchText]);

    const storeResults = useMemo(() => {
        if (!searchText) {
            return stores.slice(0, 4);
        }

        return stores.filter((store) => `${store.name} ${store.category} ${store.tagline}`.toLowerCase().includes(searchText));
    }, [searchText]);

    const topResults = useMemo(() => {
        return [
            ...filteredProducts.slice(0, 2),
            ...liveResults.slice(0, 1),
            ...creatorResults.slice(0, 1),
            ...storeResults.slice(0, 1),
        ];
    }, [creatorResults, filteredProducts, liveResults, storeResults]);

    const visibleResults = useMemo(() => {
        switch (activeTab) {
            case "Products":
                return filteredProducts;
            case "Live":
                return liveResults;
            case "Creators":
                return creatorResults;
            case "Videos":
                return videoResults;
            case "Stores":
                return storeResults;
            default:
                return topResults;
        }
    }, [activeTab, creatorResults, filteredProducts, liveResults, storeResults, topResults, videoResults]);

    const handleQuickSearch = (value: string) => {
        setQuery(value);
    };

    return (
        <div className="space-y-8 pb-10">
            <section className="rounded-[32px] bg-[linear-gradient(135deg,#fff8fb_0%,#f8fafc_55%,#ecfdf5_100%)] p-5 sm:p-6">
                <SectionHeader
                    title="TikTok-native search"
                    description="Search by product, creator, live room, or merchant, then jump into a polished social commerce path without leaving the feed."
                />
                <form className="mt-4 flex flex-col gap-3 lg:flex-row">
                    <div className="flex-1">
                        <Input
                            label="Discover your next favorite item"
                            placeholder="Search products, creators, live rooms, or stores"
                            value={query}
                            onChange={(event) => setQuery(event.target.value)}
                        />
                    </div>
                    <div className="flex items-end">
                        <Button type="button" variant="primary" className="w-full lg:w-auto" onClick={() => setActiveTab("Top")}>
                            Search
                        </Button>
                    </div>
                </form>

                <div className="mt-4 flex flex-wrap gap-2">
                    {recentSearches.map((search) => (
                        <button
                            key={search}
                            type="button"
                            onClick={() => handleQuickSearch(search)}
                            className="rounded-full border border-zinc-200 bg-white px-3 py-2 text-sm font-semibold text-zinc-700"
                        >
                            {search}
                        </button>
                    ))}
                </div>
            </section>

            <div className="grid gap-6 xl:grid-cols-[0.85fr_1.15fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Trending now</p>
                    <div className="mt-4 flex flex-wrap gap-2">
                        {trendingTerms.map((term) => (
                            <button
                                key={term}
                                type="button"
                                onClick={() => handleQuickSearch(term)}
                                className="rounded-full bg-[#FFF8FB] px-3 py-2 text-sm font-semibold text-[#E91E63]"
                            >
                                {term}
                            </button>
                        ))}
                    </div>

                    <div className="mt-5 space-y-4">
                        <div>
                            <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Live suggestions</p>
                            <div className="mt-3 space-y-2">
                                {liveResults.map((item) => (
                                    <Link key={item.id} href={item.actionHref} className="block rounded-[20px] bg-zinc-50 px-4 py-3 text-sm text-zinc-700">
                                        {item.title}
                                    </Link>
                                ))}
                            </div>
                        </div>
                        <div>
                            <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Merchant suggestions</p>
                            <div className="mt-3 space-y-2">
                                {storeResults.slice(0, 4).map((store) => (
                                    <Link key={store.slug} href={`/stores/${store.slug}`} className="block rounded-[20px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                        {store.name} • {store.category}
                                    </Link>
                                ))}
                            </div>
                        </div>
                    </div>
                </Card>

                <div className="space-y-4">
                    <div className="flex flex-wrap gap-2">
                        {tabs.map((tab) => (
                            <button
                                key={tab}
                                type="button"
                                onClick={() => setActiveTab(tab)}
                                className={`rounded-full px-4 py-2 text-sm font-semibold ${activeTab === tab ? "bg-zinc-950 text-white" : "border border-zinc-200 bg-white text-zinc-700"}`}
                            >
                                {tab}
                            </button>
                        ))}
                    </div>

                    <div className="grid gap-4 md:grid-cols-2">
                        {activeTab === "Products" && filteredProducts.map((product) => (
                            <Link key={product.id} href={`/products/${product.id}`}>
                                <ProductCard product={product} showDetailLink={false} />
                            </Link>
                        ))}

                        {activeTab === "Live" && liveResults.map((item) => (
                            <Link key={item.id} href={item.actionHref} className="block rounded-[28px] border border-zinc-200 bg-white p-4">
                                <Badge tone="success">LIVE</Badge>
                                <p className="mt-3 text-lg font-semibold text-zinc-900">{item.title}</p>
                                <p className="mt-2 text-sm leading-6 text-zinc-600">{item.description}</p>
                            </Link>
                        ))}

                        {activeTab === "Creators" && creatorResults.map((creator) => (
                            <Link key={creator.id} href={`/creator/${creator.id}`} className="block rounded-[28px] border border-zinc-200 bg-white p-4">
                                <div className="flex items-center gap-3">
                                    <img src={creator.avatar} alt={creator.name} className="h-12 w-12 rounded-full object-cover" />
                                    <div>
                                        <p className="font-semibold text-zinc-900">{creator.name}</p>
                                        <p className="text-sm text-zinc-500">{creator.followers}</p>
                                    </div>
                                </div>
                                <p className="mt-3 text-sm leading-6 text-zinc-600">{creator.bio}</p>
                            </Link>
                        ))}

                        {activeTab === "Videos" && videoResults.map((item) => (
                            <Link key={item.id} href={item.actionHref} className="block rounded-[28px] border border-zinc-200 bg-white p-4">
                                <div className="flex items-center justify-between gap-3">
                                    <p className="text-sm font-semibold text-zinc-900">{item.title}</p>
                                    <Badge tone="default">{item.tag}</Badge>
                                </div>
                                <p className="mt-3 text-sm leading-6 text-zinc-600">{item.description}</p>
                            </Link>
                        ))}

                        {activeTab === "Stores" && storeResults.map((store) => (
                            <Link key={store.slug} href={`/stores/${store.slug}`} className="block rounded-[28px] border border-zinc-200 bg-white p-4">
                                <p className="text-sm font-semibold text-zinc-900">{store.name}</p>
                                <p className="mt-2 text-sm text-zinc-500">{store.category}</p>
                                <p className="mt-3 text-sm leading-6 text-zinc-600">{store.tagline}</p>
                            </Link>
                        ))}

                        {activeTab === "Top" && visibleResults.map((result) => {
                            if ("price" in result) {
                                return (
                                    <Link key={result.id} href={`/products/${result.id}`}>
                                        <ProductCard product={result} />
                                    </Link>
                                );
                            }

                            if ("actionHref" in result) {
                                return (
                                    <Link key={result.id} href={result.actionHref} className="block rounded-[28px] border border-zinc-200 bg-white p-4">
                                        <Badge tone="success">LIVE</Badge>
                                        <p className="mt-3 text-lg font-semibold text-zinc-900">{result.title}</p>
                                        <p className="mt-2 text-sm leading-6 text-zinc-600">{result.description}</p>
                                    </Link>
                                );
                            }

                            if ("slug" in result) {
                                return (
                                    <Link key={result.slug} href={`/stores/${result.slug}`} className="block rounded-[28px] border border-zinc-200 bg-white p-4">
                                        <p className="text-sm font-semibold text-zinc-900">{result.name}</p>
                                        <p className="mt-2 text-sm text-zinc-500">{result.category}</p>
                                    </Link>
                                );
                            }

                            return null;
                        })}
                    </div>

                    <Card className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Product quick cards</p>
                        <div className="mt-3 grid gap-3 md:grid-cols-2">
                            {filteredProducts.slice(0, 4).map((product) => (
                                <Link key={product.id} href={`/products/${product.id}`} className="rounded-[24px] bg-[#FFF8FB] p-4">
                                    <p className="text-sm font-semibold text-zinc-900">{product.name}</p>
                                    <p className="mt-2 text-sm text-zinc-600">{product.price}</p>
                                </Link>
                            ))}
                        </div>
                    </Card>

                    <Card className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Creator quick cards</p>
                        <div className="mt-3 grid gap-3 md:grid-cols-2">
                            {creatorResults.slice(0, 4).map((creator) => (
                                <Link key={creator.id} href={`/creator/${creator.id}`} className="rounded-[24px] bg-zinc-50 p-4">
                                    <p className="text-sm font-semibold text-zinc-900">{creator.name}</p>
                                    <p className="mt-2 text-sm text-zinc-600">{creator.followers}</p>
                                </Link>
                            ))}
                        </div>
                    </Card>
                </div>
            </div>
        </div>
    );
}
