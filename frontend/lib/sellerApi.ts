type ProductInput = {
    name: string;
    sku: string;
    category: string;
    price: string;
    inventory: string;
    description: string;
};

type Product = ProductInput & { id: string; createdAt: string };

type Order = {
    id: string;
    orderNumber: string;
    customer: string;
    amount: number;
    items: number;
    status: string;
    date: string;
    paymentMethod: string;
};

type Livestream = {
    id: string;
    title: string;
    creatorName: string;
    status: "live" | "scheduled" | "ended";
    viewers: number;
    revenue: number;
    startTime: string;
};

import * as api from "./api";

const wait = (ms = 500) => new Promise((res) => setTimeout(res, ms));

export async function createProduct(input: ProductInput): Promise<Product> {
    // Try backend createProduct first
    try {
        const res = await api.createProduct({
            name: input.name,
            sku: input.sku,
            category: input.category,
            price: input.price,
            inventory: input.inventory,
            description: input.description,
        });

        // Normalize response into Product shape when possible
        const obj = res as unknown as Record<string, unknown>;
        return {
            id: String(obj.id ?? `${Date.now()}-${Math.random().toString(36).slice(2, 6)}`),
            createdAt: String(obj.created_at ?? new Date().toISOString()),
            ...input,
        } as Product;
    } catch {
        // Fallback to localStorage mock if backend unavailable
        await wait(400);
        const stored = JSON.parse(localStorage.getItem("seller_products") || "[]") as Product[];
        const product: Product = {
            ...input,
            id: `${Date.now()}-${Math.random().toString(36).slice(2, 6)}`,
            createdAt: new Date().toISOString(),
        };
        stored.unshift(product);
        try {
            localStorage.setItem("seller_products", JSON.stringify(stored));
        } catch { }
        return product;
    }
}

export async function fetchPendingOrders(): Promise<Order[]> {
    try {
        const res = await api.listOrdersByStatus("pending", 50, 0);
        // Normalize to Order[] if possible
        if (Array.isArray(res)) {
            return (res as unknown[]).map((o: unknown) => {
                const obj = o as Record<string, unknown>;
                return {
                    id: String(obj.id ?? ""),
                    orderNumber: String(obj.order_number ?? obj.orderNumber ?? ""),
                    customer: String(obj.customer_name ?? obj.customer ?? obj.user_id ?? "Customer"),
                    amount: Number(obj.total_amount ?? obj.amount ?? 0),
                    items: Number(obj.items ?? obj.quantity ?? 1),
                    status: String(obj.status ?? "pending"),
                    date: String(obj.created_at ?? obj.date ?? ""),
                    paymentMethod: String(obj.payment_method ?? obj.paymentMethod ?? ""),
                } as Order;
            });
        }
    } catch {
        // ignore and fallback
    }

    await wait(300);
    return [
        { id: "5", orderNumber: "ORD-005", customer: "Charlie Brown", amount: 599.99, items: 4, status: "pending", date: "2024-01-11", paymentMethod: "Credit Card" },
        { id: "6", orderNumber: "ORD-006", customer: "Diana Prince", amount: 349.99, items: 2, status: "pending", date: "2024-01-10", paymentMethod: "Apple Pay" },
    ];
}

export async function updateOrderStatus(id: string, status: string): Promise<boolean> {
    try {
        // try to call backend (expects numeric id when possible)
        const numericId = Number(id);
        await api.updateOrderStatus(Number.isFinite(numericId) ? numericId : (id as unknown as number), status);
        return true;
    } catch {
        await wait(200);
        return true;
    }
}

export async function fetchActiveLivestreams(): Promise<Livestream[]> {
    try {
        const res = await api.getLiveRooms();
        if (Array.isArray(res)) {
            return (res as unknown[])
                .map((r: unknown) => {
                    const obj = r as Record<string, unknown>;
                    const creator = obj.creator as Record<string, unknown> | undefined;
                    return {
                        id: String(obj.id ?? ""),
                        title: String(obj.title ?? obj.name ?? "Livestream"),
                        creatorName: String(obj.creator_name ?? creator?.name ?? ""),
                        status: (String(obj.status ?? "live") as unknown) as "live" | "scheduled" | "ended",
                        viewers: Number(obj.viewers ?? obj.viewer_count ?? 0),
                        revenue: Number(obj.revenue ?? 0),
                        startTime: String(obj.start_time ?? obj.started_at ?? ""),
                    } as Livestream;
                })
                .filter((s: Livestream) => s.status === "live");
        }
    } catch {
        // fallback
    }

    await wait(300);
    return [
        { id: "1", title: "Product Launch Event", creatorName: "Sarah Chen", status: "live", viewers: 2341, revenue: 1245.99, startTime: "2024-01-15 19:00" },
        { id: "2", title: "Holiday Flash Sale", creatorName: "Mike Johnson", status: "live", viewers: 1874, revenue: 825.42, startTime: "2024-01-15 18:30" },
    ];
}

export type { Product, Order, Livestream };
