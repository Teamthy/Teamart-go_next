"use client";

import { useMemo, useState } from "react";
import Link from "next/link";
import PageHeader from "@/components/ui/PageHeader";
import ProductGrid from "@/components/ui/ProductGrid";
import StatCard from "@/components/ui/StatCard";
import Tabs from "@/components/ui/Tabs";
import EmptyState from "@/components/ui/EmptyState";
import Button from "@/components/ui/button";
import { products, productCategories } from "@/lib/mock/products";

const categories = ["all", ...productCategories];

export default function ProductsPage() {
    const [activeCategory, setActiveCategory] = useState("all");
    const [searchQuery, setSearchQuery] = useState("");

    const filteredProducts = useMemo(() => {
        const query = searchQuery.trim().toLowerCase();

        return products.filter((product) => {
            const matchesCategory = activeCategory === "all" || product.category === activeCategory;
            const matchesQuery = !query || `${product.name} ${product.description} ${product.creator ?? ""}`.toLowerCase().includes(query);

            return matchesCategory && matchesQuery;
        });
    }, [activeCategory, searchQuery]);

    return (
        <div className="space-y-8 pb-10">
            <PageHeader
                title="Products"
                description="Browse a curated catalog of creator picks, merchant bundles, and social-ready products built for fast discovery."
                actions={
                    <>
                        <Button asChild variant="primary">
                            <Link href="/feed">Explore feed</Link>
                        </Button>
                        <Button asChild variant="secondary">
                            <Link href="/live">Watch live</Link>
                        </Button>
                    </>
                }
            />

            <div className="grid gap-4 md:grid-cols-3">
                <StatCard label="Catalog size" value={String(products.length)} helper="A healthy mix of creator and merchant products is surfaced here." />
                <StatCard label="Top category" value="Accessories" helper="Accessory bundles continue to perform strongly in discovery and checkout." />
                <StatCard label="Live-ready" value="12" helper="Products are curated for sharing, pinning, and fast purchase moments." />
            </div>

            <div className="rounded-[28px] border border-zinc-200 bg-white p-4 sm:p-5">
                <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
                    <div>
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Search and filter</p>
                        <p className="mt-2 text-sm text-zinc-600">Search by name, description, or creator name to tune the catalog for your audience.</p>
                    </div>
                    <div className="flex-1 lg:max-w-xl">
                        <input
                            type="text"
                            placeholder="Search products, creators, or product cues"
                            value={searchQuery}
                            onChange={(event) => setSearchQuery(event.target.value)}
                            className="w-full rounded-[24px] border border-zinc-200 px-4 py-3 text-sm text-zinc-900 outline-none transition focus:border-[#E91E63]"
                        />
                    </div>
                </div>
            </div>

            <Tabs tabs={categories.map((category) => ({ label: category === "all" ? "All products" : category, value: category }))} active={activeCategory} onChange={setActiveCategory} />

            {filteredProducts.length === 0 ? (
                <EmptyState title="No matching products" description="Try a different category or search term to surface more items." />
            ) : (
                <ProductGrid products={filteredProducts} />
            )}
        </div>
    );
}
