"use client";

import { useState } from "react";
import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { z } from "zod";
import SectionHeader from "@/components/ui/SectionHeader";
import ProductCard from "@/components/product/ProductCard";
import { ProductGridSkeleton } from "@/components/ui/Skeleton";
import * as api from "@/lib/api";

const searchSchema = z.object({
    query: z.string().trim().min(2, "Please enter at least 2 characters."),
});

type SearchSchema = z.infer<typeof searchSchema>;

export default function SearchPage() {
    const [query, setQuery] = useState("");
    const [submittedQuery, setSubmittedQuery] = useState("");
    const [hasSearched, setHasSearched] = useState(false);
    const [formError, setFormError] = useState<string | null>(null);

    const trimmedQuery = submittedQuery.trim();

    type ProductResponse = {
        products: {
            id: string;
            name: string;
            description: string;
            price: string;
            image: string;
        }[]
    };

    const searchQueryState = useQuery<ProductResponse, Error>({
        queryKey: ["search", trimmedQuery],
        queryFn: ({ queryKey }) => api.searchProducts(queryKey[1] as string, 50, 0),
        enabled: Boolean(trimmedQuery),
        staleTime: 1000 * 60 * 5,
        retry: 1,
    });

    const results = searchQueryState.data?.products || [];
    const isLoading = searchQueryState.isFetching;
    const error = searchQueryState.error as Error | null;

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setFormError(null);

        const result = searchSchema.safeParse({ query });
        if (!result.success) {
            setFormError(result.error.issues[0]?.message ?? "Please enter a valid search.");
            return;
        }

        setSubmittedQuery(result.data.query);
        setHasSearched(true);
    };

    return (
        <div className="space-y-8">
            <SectionHeader
                title="Search products"
                description="Find products, creators, and livestream events in one place."
            />

            <form onSubmit={handleSubmit} className="flex flex-col gap-2 sm:flex-row">
                <div className="flex-1">
                    <input
                        type="text"
                        placeholder="Search products, creators, events..."
                        value={query}
                        onChange={(e) => setQuery(e.target.value)}
                        className="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-900 text-black dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
                        aria-invalid={!!formError}
                    />
                    {formError ? (
                        <p className="mt-2 text-sm text-rose-600 dark:text-rose-400">{formError}</p>
                    ) : null}
                </div>
                <button
                    type="submit"
                    disabled={isLoading}
                    className="rounded-lg bg-indigo-600 px-6 py-3 text-white font-medium hover:bg-indigo-500 transition disabled:opacity-50 disabled:cursor-not-allowed"
                >
                    {isLoading ? "Searching..." : "Search"}
                </button>
            </form>

            {isLoading && <ProductGridSkeleton />}

            {error && (
                <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
                    <p className="text-red-800 dark:text-red-200">Error: {error.message}</p>
                </div>
            )}

            {!hasSearched && !isLoading && (
                <div className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 p-8 shadow-sm">
                    <p className="text-sm text-slate-600 dark:text-slate-400">
                        Enter a search query above to find products.
                    </p>
                </div>
            )}

            {hasSearched && !isLoading && results.length === 0 && !error && (
                <div className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 p-8 shadow-sm text-center">
                    <p className="text-sm text-slate-600 dark:text-slate-400">
                        No products found for "{query}". Try a different search term.
                    </p>
                </div>
            )}

            {results.length > 0 && (
                <div>
                    <p className="text-sm text-slate-600 dark:text-slate-400 mb-4">
                        Found {results.length} product{results.length !== 1 ? "s" : ""} for "{query}"
                    </p>
                    <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
                        {results.map((product) => (
                            <Link key={product.id} href={`/products/${product.id}`}>
                                <ProductCard product={product} />
                            </Link>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}
