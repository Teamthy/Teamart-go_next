"use client";

import { useState } from "react";
import Link from "next/link";
import SectionHeader from "@/components/ui/SectionHeader";
import ProductCard from "@/components/product/ProductCard";
import * as api from "@/lib/api";

export default function SearchPage() {
    const [query, setQuery] = useState("");
    const [results, setResults] = useState<any[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [hasSearched, setHasSearched] = useState(false);

    const handleSearch = async (searchQuery: string) => {
        if (!searchQuery.trim()) {
            setResults([]);
            setHasSearched(false);
            return;
        }

        setQuery(searchQuery);
        setIsLoading(true);
        setError(null);
        setHasSearched(true);

        try {
            const response = await api.searchProducts(searchQuery, 50, 0);
            setResults(response.products || []);
        } catch (err: any) {
            setError(err.message || "Search failed");
            console.error("Error searching:", err);
        } finally {
            setIsLoading(false);
        }
    };

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        handleSearch(query);
    };

    return (
        <div className="space-y-8">
            <SectionHeader
                title="Search products"
                description="Find products, creators, and livestream events in one place."
            />

            {/* Search Form */}
            <form onSubmit={handleSubmit} className="flex gap-2">
                <input
                    type="text"
                    placeholder="Search products, creators, events..."
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    className="flex-1 px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-900 text-black dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
                />
                <button
                    type="submit"
                    disabled={isLoading}
                    className="px-6 py-3 rounded-lg bg-indigo-600 hover:bg-indigo-500 text-white font-medium disabled:opacity-50 disabled:cursor-not-allowed transition"
                >
                    {isLoading ? "Searching..." : "Search"}
                </button>
            </form>

            {/* Error State */}
            {error && (
                <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
                    <p className="text-red-800 dark:text-red-200">Error: {error}</p>
                </div>
            )}

            {/* Loading State */}
            {isLoading && (
                <div className="flex justify-center items-center py-12">
                    <p className="text-gray-500 dark:text-gray-400">Searching...</p>
                </div>
            )}

            {/* No Search State */}
            {!hasSearched && !isLoading && (
                <div className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 p-8 shadow-sm">
                    <p className="text-sm text-slate-600 dark:text-slate-400">
                        Enter a search query above to find products.
                    </p>
                </div>
            )}

            {/* No Results State */}
            {hasSearched && !isLoading && results.length === 0 && !error && (
                <div className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 p-8 shadow-sm text-center">
                    <p className="text-sm text-slate-600 dark:text-slate-400">
                        No products found for "{query}". Try a different search term.
                    </p>
                </div>
            )}

            {/* Results Grid */}
            {results.length > 0 && (
                <div>
                    <p className="text-sm text-slate-600 dark:text-slate-400 mb-4">
                        Found {results.length} product{results.length !== 1 ? "s" : ""} for "{query}"
                    </p>
                    <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
                        {results.map((product) => (
                            <Link key={product.id} href={`/products/${product.id}`}>
                                <ProductCard product={product} showDetailLink={false} />
                            </Link>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}
