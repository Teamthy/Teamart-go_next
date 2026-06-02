"use client";

import { useState, useEffect } from "react";
import * as api from "@/lib/api";

export interface OrderItem {
    product_id: number;
    quantity: number;
    price: number;
}

export interface Order {
    id: number;
    user_id: number;
    total_amount: number;
    status: string;
    items?: OrderItem[];
    created_at?: string;
    updated_at?: string;
}

function getErrorMessage(error: unknown, fallback: string) {
    if (error instanceof Error) return error.message;
    if (typeof error === "string") return error;
    return fallback;
}

export function useOrders(limit = 20, offset = 0) {
    const [orders, setOrders] = useState<Order[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchOrders = async () => {
            setIsLoading(true);
            setError(null);
            try {
                const response = await api.listOrders(limit, offset);
                setOrders(response.orders || []);
            } catch (err: unknown) {
                setError(getErrorMessage(err, "Failed to fetch orders"));
            } finally {
                setIsLoading(false);
            }
        };

        fetchOrders();
    }, [limit, offset]);

    return { orders, isLoading, error };
}

export function useUserOrders(userId: number | null, limit = 20, offset = 0) {
    const [orders, setOrders] = useState<Order[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (!userId) return;

        const fetchOrders = async () => {
            setIsLoading(true);
            setError(null);
            try {
                const response = await api.listUserOrders(userId, limit, offset);
                setOrders(response.orders || []);
            } catch (err: unknown) {
                setError(getErrorMessage(err, "Failed to fetch orders"));
            } finally {
                setIsLoading(false);
            }
        };

        fetchOrders();
    }, [userId, limit, offset]);

    return { orders, isLoading, error };
}

export function useOrder(orderId: number | string | null) {
    const [order, setOrder] = useState<Order | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (!orderId) return;

        const fetch = async () => {
            setIsLoading(true);
            setError(null);
            try {
                const response = await api.getOrder(orderId);
                setOrder(response);
            } catch (err: unknown) {
                setError(getErrorMessage(err, "Failed to fetch order"));
            } finally {
                setIsLoading(false);
            }
        };

        fetch();
    }, [orderId]);

    return { order, isLoading, error };
}

export function useCreateOrder() {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const create = async (data: Record<string, unknown>) => {
        setIsLoading(true);
        setError(null);
        try {
            const response = await api.createOrder(data);
            return response;
        } catch (err: unknown) {
            setError(getErrorMessage(err, "Failed to create order"));
            throw err;
        } finally {
            setIsLoading(false);
        }
    };

    return { create, isLoading, error };
}
