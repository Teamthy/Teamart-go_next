"use client";

import { useState, useEffect } from "react";
import * as api from "@/lib/api";

export interface Product {
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

export function useProducts(limit = 20, offset = 0) {
    const [products, setProducts] = useState<Product[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchProducts = async () => {
            setIsLoading(true);
            setError(null);
            try {
                const response = await api.listProducts(limit, offset);
                setProducts(response.products || []);
            } catch (err: any) {
                setError(err.message || "Failed to fetch products");
            } finally {
                setIsLoading(false);
            }
        };

        fetchProducts();
    }, [limit, offset]);

    const search = async (query: string) => {
        setIsLoading(true);
        setError(null);
        try {
            const response = await api.searchProducts(query, limit, offset);
            setProducts(response.products || []);
            return response;
        } catch (err: any) {
            setError(err.message || "Search failed");
            throw err;
        } finally {
            setIsLoading(false);
        }
    };

    const getProduct = async (productId: number | string) => {
        setIsLoading(true);
        setError(null);
        try {
            const response = await api.getProduct(productId);
            return response;
        } catch (err: any) {
            setError(err.message || "Failed to fetch product");
            throw err;
        } finally {
            setIsLoading(false);
        }
    };

    return { products, isLoading, error, search, getProduct };
}

export function useProduct(productId: number | string | null) {
    const [product, setProduct] = useState<Product | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (!productId) return;

        const fetch = async () => {
            setIsLoading(true);
            setError(null);
            try {
                const response = await api.getProduct(productId);
                setProduct(response);
            } catch (err: any) {
                setError(err.message || "Failed to fetch product");
            } finally {
                setIsLoading(false);
            }
        };

        fetch();
    }, [productId]);

    return { product, isLoading, error };
}

export function useCreateProduct() {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const create = async (data: any) => {
        setIsLoading(true);
        setError(null);
        try {
            const response = await api.createProduct(data);
            return response;
        } catch (err: any) {
            setError(err.message || "Failed to create product");
            throw err;
        } finally {
            setIsLoading(false);
        }
    };

    return { create, isLoading, error };
}
