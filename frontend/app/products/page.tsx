"use client";

import { useState } from "react";
import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import ProductCard from "@/components/product/ProductCard";
import SectionHeader from "@/components/ui/SectionHeader";
import { ProductGridSkeleton } from "@/components/ui/Skeleton";
import * as api from "@/lib/api";

interface Product {
    id: number;
    name: string;
    description?: string;
    price: number;
    sku?: string;
    category?: string;
    image_url?: string;
    stock?: number;
    created_at?: string;
}

export default function ProductsPage() {
    const [searchQuery, setSearchQuery] = useState("");
    const [searchTerm, setSearchTerm] = useState("");

    type ProductResponse = { products: Product[] };

    const productsQuery = useQuery<ProductResponse, Error>({
        queryKey: ["products"],
        queryFn: () => api.listProducts(50, 0),
        staleTime: 1000 * 60 * 2,
    });

    const searchQueryState = useQuery<ProductResponse, Error>({
        queryKey: ["search", searchTerm],
        queryFn: ({ queryKey }) => api.searchProducts(queryKey[1] as string, 50, 0),
        enabled: Boolean(searchTerm),
        staleTime: 1000 * 60 * 5,
    });

    const products = productsQuery.data?.products || [];
    const searchResults = searchQueryState.data?.products || [];
    const isLoading = productsQuery.isLoading || searchQueryState.isFetching;
    const error = (productsQuery.error || searchQueryState.error) as Error | null;
    const displayedProducts = searchTerm ? searchResults : products;

    const handleSearch = (query: string) => {
        setSearchQuery(query);
        setSearchTerm(query.trim());
    };

    return (
        <div className="space-y-8">
            <SectionHeader
                title="Products"
                description="Browse the catalog and explore creator products."
            />

            <div className="flex gap-2">
                <input
                    type="text"
                    placeholder="Search products..."
                    value={searchQuery}
                    onChange={(e) => handleSearch(e.target.value)}
                    className="flex-1 px-4 py-2 rounded-lg border border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-900 text-black dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
                />
            </div>

            {isLoading && <ProductGridSkeleton />}

            {error && !isLoading && (
                <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
                    <p className="text-red-800 dark:text-red-200">Error: {error.message}</p>
                </div>
            )}

            {!isLoading && displayedProducts.length === 0 && !error && (
                <div className="text-center py-12">
                    <p className="text-gray-500 dark:text-gray-400">
                        {searchTerm ? "No products found matching your search." : "No products available."}
                    </p>
                </div>
            )}

            {!isLoading && displayedProducts.length > 0 && (
                <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
                    {displayedProducts.map((product) => (
                        <Link key={product.id} href={`/products/${product.id}`}>
                            <div className="h-full">
                                <ProductCard
                                    product={{
                                        id: String(product.id),
                                        name: product.name,
                                        description: product.description || "",
                                        price: String(product.price),
                                        image: product.image_url ?? "",
                                    }}
                                />
                            </div>
                        </Link>
                    ))}
                </div>
            )}
        </div>
    );
}
