"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import ProductCard from "@/components/product/ProductCard";
import SectionHeader from "@/components/ui/SectionHeader";
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
    const [products, setProducts] = useState<Product[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [searchQuery, setSearchQuery] = useState("");
    const [isSearching, setIsSearching] = useState(false);

    useEffect(() => {
        const fetchProducts = async () => {
            setIsLoading(true);
            setError(null);
            try {
                const response = await api.listProducts(50, 0);
                setProducts(response.products || []);
            } catch (err: any) {
                setError(err.message || "Failed to fetch products");
                console.error("Error fetching products:", err);
            } finally {
                setIsLoading(false);
            }
        };

        fetchProducts();
    }, []);

    const handleSearch = async (query: string) => {
        if (!query.trim()) {
            setSearchQuery("");
            try {
                const response = await api.listProducts(50, 0);
                setProducts(response.products || []);
            } catch (err) {
                console.error("Error fetching products:", err);
            }
            return;
        }

        setSearchQuery(query);
        setIsSearching(true);
        try {
            const response = await api.searchProducts(query, 50, 0);
            setProducts(response.products || []);
        } catch (err: any) {
            setError(err.message || "Search failed");
            console.error("Error searching products:", err);
        } finally {
            setIsSearching(false);
        }
    };

    return (
        <div className="space-y-8">
            <SectionHeader
                title="Products"
                description="Browse the catalog and explore creator products."
            />

            {/* Search Bar */}
            <div className="flex gap-2">
                <input
                    type="text"
                    placeholder="Search products..."
                    value={searchQuery}
                    onChange={(e) => handleSearch(e.target.value)}
                    className="flex-1 px-4 py-2 rounded-lg border border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-900 text-black dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
                />
            </div>

            {/* Loading State */}
            {isLoading && (
                <div className="flex justify-center items-center py-12">
                    <p className="text-gray-500 dark:text-gray-400">Loading products...</p>
                </div>
            )}

            {/* Error State */}
            {error && !isLoading && (
                <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
                    <p className="text-red-800 dark:text-red-200">Error: {error}</p>
                </div>
            )}

            {/* Empty State */}
            {!isLoading && products.length === 0 && !error && (
                <div className="text-center py-12">
                    <p className="text-gray-500 dark:text-gray-400">
                        {searchQuery ? "No products found matching your search." : "No products available."}
                    </p>
                </div>
            )}

            {/* Products Grid */}
            {!isLoading && products.length > 0 && (
                <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
                    {products.map((product) => (
                        <Link key={product.id} href={`/products/${product.id}`}>
                            <div className="h-full">
                                <ProductCard
                                    product={{
                                        ...product,
                                        image: product.image_url,
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
