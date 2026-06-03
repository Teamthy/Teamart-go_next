"use client";

import { useEffect, useMemo, useState } from "react";
import PageHeader from "@/components/ui/PageHeader";
import StatCard from "@/components/ui/StatCard";
import Tabs from "@/components/ui/Tabs";
import StoreCard from "@/components/ui/StoreCard";
import Button from "@/components/ui/button";
import Link from "next/link";
import * as api from "@/lib/api";
import type { StoreSummary } from "@/types/commerce";

export default function StoresPage() {
    const [stores, setStores] = useState<StoreSummary[]>([]);
    const [activeCategory, setActiveCategory] = useState("all");
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchStores = async () => {
            setIsLoading(true);
            setError(null);
            try {
                const response = await api.listPublicStores();
                setStores(response.stores || []);
            } catch (err: unknown) {
                const message = err instanceof Error ? err.message : "Unable to load stores";
                setError(message);
            } finally {
                setIsLoading(false);
            }
        };

        fetchStores();
    }, []);

    const categories = useMemo(() => ["all", ...Array.from(new Set(stores.map((store) => store.category)))], [stores]);

    const filteredStores = useMemo(() => {
        if (activeCategory === "all") {
            return stores;
        }

        return stores.filter((store) => store.category === activeCategory);
    }, [activeCategory, stores]);

    return (
        <div className="space-y-8 pb-10">
            <PageHeader
                title="Shop by merchant store"
                description="Explore standout storefronts, live moments, and categories that fit how shoppers browse and buy today."
                actions={
                    <>
                        <Button asChild variant="primary">
                            <Link href="/live">Join live rooms</Link>
                        </Button>
                        <Button asChild variant="secondary">
                            <Link href="/feed">Return to feed</Link>
                        </Button>
                    </>
                }
            />

            <div className="grid gap-4 md:grid-cols-3">
                <StatCard label="Featured stores" value="10" helper="A curated set of merchant storefronts with live and bundle-ready moments." />
                <StatCard label="Top category" value="Fashion" helper="Shoppers are gravitating toward style-first storefronts right now." />
                <StatCard label="Live previews" value="4" helper="Several stores are currently active with strong live room engagement." />
            </div>

            <Tabs
                tabs={categories.map((category) => ({ label: category === "all" ? "All stores" : category, value: category }))}
                active={activeCategory}
                onChange={setActiveCategory}
            />

            {isLoading && (
                <div className="py-12 text-center text-slate-500">Loading stores…</div>
            )}

            {error && !isLoading && (
                <div className="rounded-3xl border border-rose-100 bg-rose-50 p-6 text-sm text-rose-700">
                    {error}
                </div>
            )}

            {!isLoading && !error && filteredStores.length === 0 && (
                <div className="rounded-3xl border border-slate-200 bg-white p-8 text-center text-slate-500">
                    No stores are available right now. Please check back later.
                </div>
            )}

            {!isLoading && !error && filteredStores.length > 0 && (
                <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
                    {filteredStores.map((store) => (
                        <StoreCard
                            key={store.id}
                            name={store.name}
                            slug={store.slug}
                            category={store.category}
                            rating={store.rating}
                            banner={store.banner_url}
                            tagline={store.tagline}
                            live={store.live_status}
                            products={String(store.products)}
                        />
                    ))}
                </div>
            )}
        </div>
    );
}
