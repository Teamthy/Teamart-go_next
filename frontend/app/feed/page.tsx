"use client";

import { useEffect } from "react";
import ProductCard from "@/components/product/ProductCard";
import SectionHeader from "@/components/ui/SectionHeader";
import { ProductGridSkeleton } from "@/components/ui/Skeleton";
import { useFeed } from "@/hooks/useFeed";
import useAppStore from "@/store/useAppStore";

export default function FeedPage() {
    const { items, isLoading, error } = useFeed(50);
    const liveFeedCount = useAppStore((state) => state.liveFeedCount);
    const resetLiveFeedCount = useAppStore((state) => state.resetLiveFeedCount);

    useEffect(() => {
        if (liveFeedCount > 0 && items.length > 0) {
            const timer = window.setTimeout(() => resetLiveFeedCount(), 5000);
            return () => window.clearTimeout(timer);
        }
    }, [items.length, liveFeedCount, resetLiveFeedCount]);

    return (
        <div className="space-y-8">
            <SectionHeader
                title="Product feed"
                description="Browse the latest products from creators and tailored recommendations."
            />

            {isLoading && <ProductGridSkeleton />}

            {liveFeedCount > 0 && !isLoading && (
                <div className="rounded-3xl border border-slate-200 bg-sky-50 p-4 text-sm text-slate-700 shadow-sm">
                    Live feed updated with {liveFeedCount} new product{liveFeedCount !== 1 ? "s" : ""}. Changes are applied automatically.
                </div>
            )}

            {error && !isLoading && (
                <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
                    <p className="text-yellow-800 dark:text-yellow-200">
                        Unable to load personalized feed. Showing popular products instead.
                    </p>
                </div>
            )}

            {!isLoading && items.length === 0 && !error && (
                <div className="text-center py-12">
                    <p className="text-gray-500 dark:text-gray-400">No items in feed yet.</p>
                </div>
            )}

            {!isLoading && items.length > 0 && (
                <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
                    {items.map((product) => (
                        <ProductCard
                            key={product.id}
                            product={{
                                id: String(product.id),
                                name: product.name,
                                description: product.description || "",
                                price: `$${product.price?.toFixed(2) ?? "0.00"}`,
                                image: product.image_url || "",
                            }}
                        />
                    ))}
                </div>
            )}
        </div>
    );
}
