"use client";

import ProductCard from "@/components/product/ProductCard";
import SectionHeader from "@/components/ui/SectionHeader";
import { useFeed } from "@/hooks/useFeed";

export default function FeedPage() {
    const { items, isLoading, error } = useFeed(50);

    return (
        <div className="space-y-8">
            <SectionHeader
                title="Product feed"
                description="Browse the latest products from creators and tailored recommendations."
            />

            {isLoading && (
                <div className="flex justify-center items-center py-12">
                    <p className="text-gray-500 dark:text-gray-400">Loading feed...</p>
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
                        <ProductCard key={product.id} product={product} />
                    ))}
                </div>
            )}
        </div>
    );
}
