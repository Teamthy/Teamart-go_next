"use client";

import { useMemo, useState } from "react";
import PageHeader from "@/components/ui/PageHeader";
import StatCard from "@/components/ui/StatCard";
import Tabs from "@/components/ui/Tabs";
import StoreCard from "@/components/ui/StoreCard";
import Button from "@/components/ui/button";
import Link from "next/link";
import { stores } from "@/lib/mock/stores";

const categories = ["all", ...Array.from(new Set(stores.map((store) => store.category)))];

export default function StoresPage() {
    const [activeCategory, setActiveCategory] = useState("all");

    const filteredStores = useMemo(() => {
        if (activeCategory === "all") {
            return stores;
        }

        return stores.filter((store) => store.category === activeCategory);
    }, [activeCategory]);

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

            <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
                {filteredStores.map((store) => (
                    <StoreCard
                        key={store.id}
                        name={store.name}
                        slug={store.slug}
                        category={store.category}
                        rating={store.rating}
                        banner={store.banner}
                        tagline={store.tagline}
                        live={store.live}
                        products={store.products}
                    />
                ))}
            </div>
        </div>
    );
}
