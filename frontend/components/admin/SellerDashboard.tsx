"use client";

import { useEffect, useState } from "react";
import SectionHeader from "@/components/ui/SectionHeader";
import * as api from "@/lib/api";

function statusBadge(status: string) {
    const base = "inline-flex items-center rounded-full px-2.5 py-1 text-xs font-semibold";
    switch (status?.toUpperCase()) {
        case "PENDING":
            return `${base} bg-yellow-100 text-yellow-800`;
        case "PROCESSING":
            return `${base} bg-sky-100 text-sky-800`;
        case "SHIPPED":
            return `${base} bg-indigo-100 text-indigo-800`;
        case "DELIVERED":
            return `${base} bg-emerald-100 text-emerald-800`;
        case "CANCELED":
        case "CANCELLED":
            return `${base} bg-rose-100 text-rose-800`;
        default:
            return `${base} bg-slate-100 text-slate-800`;
    }
}

export default function SellerDashboard() {
    const [stats, setStats] = useState({
        totalProducts: 0,
        inventoryValue: 0,
        activeOrders: 0,
        pendingPayout: 0,
    });
    const [products, setProducts] = useState<any[]>([]);
    const [orders, setOrders] = useState<any[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchData = async () => {
            setIsLoading(true);
            setError(null);
            try {
                // Fetch products
                const productsRes = await api.listProducts(100, 0);
                const prods = productsRes.products || [];
                setProducts(prods);

                // Fetch orders
                const ordersRes = await api.listOrders(100, 0);
                const ords = ordersRes.orders || [];
                setOrders(ords);

                // Calculate stats
                const inventoryValue = prods.reduce(
                    (sum: number, p: any) => sum + (p.price || 0) * (p.stock || 0),
                    0
                );
                const activeOrdersCount = ords.filter(
                    (o: any) =>
                        o.status &&
                        !["DELIVERED", "CANCELLED", "REFUNDED"].includes(o.status.toUpperCase())
                ).length;

                setStats({
                    totalProducts: prods.length,
                    inventoryValue,
                    activeOrders: activeOrdersCount,
                    pendingPayout: 0, // TODO: fetch from payouts API
                });
            } catch (err: any) {
                setError(err.message || "Failed to load dashboard");
                console.error("Error loading dashboard:", err);
            } finally {
                setIsLoading(false);
            }
        };

        fetchData();
    }, []);

    if (isLoading) {
        return (
            <div className="space-y-10">
                <div className="text-center py-12">
                    <p className="text-slate-500">Loading dashboard...</p>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="space-y-10">
                <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                    <p className="text-red-800">Error: {error}</p>
                </div>
            </div>
        );
    }

    return (
        <div className="space-y-10">
            <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
                <div>
                    <h1 className="text-3xl font-semibold text-slate-900 dark:text-white">Seller Dashboard</h1>
                    <p className="mt-2 max-w-2xl text-sm text-slate-600 dark:text-slate-400">
                        Manage your inventory, view recent orders, and track payouts in one place.
                    </p>
                </div>
                <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
                    <div className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 px-4 py-5 shadow-sm">
                        <p className="text-xs uppercase tracking-[0.24em] text-slate-500 dark:text-slate-400">
                            Products
                        </p>
                        <p className="mt-3 text-3xl font-semibold text-slate-900 dark:text-white">
                            {stats.totalProducts}
                        </p>
                    </div>
                    <div className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 px-4 py-5 shadow-sm">
                        <p className="text-xs uppercase tracking-[0.24em] text-slate-500 dark:text-slate-400">
                            Inventory value
                        </p>
                        <p className="mt-3 text-3xl font-semibold text-slate-900 dark:text-white">
                            ${stats.inventoryValue.toFixed(0)}
                        </p>
                    </div>
                    <div className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 px-4 py-5 shadow-sm">
                        <p className="text-xs uppercase tracking-[0.24em] text-slate-500 dark:text-slate-400">
                            Open orders
                        </p>
                        <p className="mt-3 text-3xl font-semibold text-slate-900 dark:text-white">
                            {stats.activeOrders}
                        </p>
                    </div>
                    <div className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 px-4 py-5 shadow-sm">
                        <p className="text-xs uppercase tracking-[0.24em] text-slate-500 dark:text-slate-400">
                            Pending payout
                        </p>
                        <p className="mt-3 text-3xl font-semibold text-slate-900 dark:text-white">
                            ${stats.pendingPayout.toFixed(2)}
                        </p>
                    </div>
                </div>
            </div>

            {/* Recent Orders */}
            <div className="space-y-4">
                <SectionHeader title="Recent Orders" description="Your latest orders and their status" />
                {orders.length === 0 ? (
                    <p className="text-slate-500 dark:text-slate-400">No orders yet</p>
                ) : (
                    <div className="overflow-x-auto rounded-lg border border-slate-200 dark:border-slate-700">
                        <table className="w-full text-sm">
                            <thead className="border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900">
                                <tr>
                                    <th className="px-6 py-3 text-left font-semibold text-slate-900 dark:text-white">
                                        Order ID
                                    </th>
                                    <th className="px-6 py-3 text-left font-semibold text-slate-900 dark:text-white">
                                        Amount
                                    </th>
                                    <th className="px-6 py-3 text-left font-semibold text-slate-900 dark:text-white">
                                        Status
                                    </th>
                                    <th className="px-6 py-3 text-left font-semibold text-slate-900 dark:text-white">
                                        Date
                                    </th>
                                </tr>
                            </thead>
                            <tbody>
                                {orders.slice(0, 10).map((order) => (
                                    <tr
                                        key={order.id}
                                        className="border-b border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-900"
                                    >
                                        <td className="px-6 py-3 text-slate-900 dark:text-white">#{order.id}</td>
                                        <td className="px-6 py-3 text-slate-900 dark:text-white">
                                            ${order.total_amount?.toFixed(2) || "0.00"}
                                        </td>
                                        <td className="px-6 py-3">{statusBadge(order.status)}</td>
                                        <td className="px-6 py-3 text-slate-600 dark:text-slate-400">
                                            {new Date(order.created_at || "").toLocaleDateString()}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>

            {/* Recent Products */}
            <div className="space-y-4">
                <SectionHeader title="Your Products" description="Products you've added to the platform" />
                {products.length === 0 ? (
                    <p className="text-slate-500 dark:text-slate-400">No products yet</p>
                ) : (
                    <div className="overflow-x-auto rounded-lg border border-slate-200 dark:border-slate-700">
                        <table className="w-full text-sm">
                            <thead className="border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900">
                                <tr>
                                    <th className="px-6 py-3 text-left font-semibold text-slate-900 dark:text-white">
                                        Product Name
                                    </th>
                                    <th className="px-6 py-3 text-left font-semibold text-slate-900 dark:text-white">
                                        Price
                                    </th>
                                    <th className="px-6 py-3 text-left font-semibold text-slate-900 dark:text-white">
                                        Stock
                                    </th>
                                    <th className="px-6 py-3 text-left font-semibold text-slate-900 dark:text-white">
                                        SKU
                                    </th>
                                </tr>
                            </thead>
                            <tbody>
                                {products.slice(0, 10).map((product) => (
                                    <tr
                                        key={product.id}
                                        className="border-b border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-900"
                                    >
                                        <td className="px-6 py-3 text-slate-900 dark:text-white font-medium">
                                            {product.name}
                                        </td>
                                        <td className="px-6 py-3 text-slate-900 dark:text-white">
                                            ${product.price?.toFixed(2) || "0.00"}
                                        </td>
                                        <td className="px-6 py-3 text-slate-900 dark:text-white">
                                            {product.stock || 0}
                                        </td>
                                        <td className="px-6 py-3 text-slate-600 dark:text-slate-400">
                                            {product.sku || "N/A"}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
}
<p className="mt-3 text-3xl font-semibold text-slate-900">${pendingPayout.toFixed(0)}</p>
                    </div >
                </div >
            </div >

    <div className="grid gap-8 lg:grid-cols-[1.5fr_1fr]">
        <section className="space-y-6 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <SectionHeader
                title="Inventory overview"
                description="Track stock levels, price, status, and forecast your next sale window."
            />

            <div className="overflow-hidden rounded-3xl border border-slate-200">
                <table className="min-w-full divide-y divide-slate-200 text-left text-sm">
                    <thead className="bg-slate-50">
                        <tr>
                            <th className="px-6 py-4 font-semibold text-slate-600">Product</th>
                            <th className="px-6 py-4 font-semibold text-slate-600">Price</th>
                            <th className="px-6 py-4 font-semibold text-slate-600">Inventory</th>
                            <th className="px-6 py-4 font-semibold text-slate-600">Status</th>
                            <th className="px-6 py-4 font-semibold text-slate-600">Sales</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-200 bg-white">
                        {sellerProducts.map((product) => (
                            <tr key={product.id} className="hover:bg-slate-50">
                                <td className="px-6 py-4">
                                    <div className="text-sm font-semibold text-slate-900">{product.name}</div>
                                    <div className="text-xs text-slate-500">{product.sku}</div>
                                </td>
                                <td className="px-6 py-4 text-slate-700">{product.price}</td>
                                <td className="px-6 py-4 text-slate-700">{product.stock}</td>
                                <td className="px-6 py-4">{statusBadge(product.status)}</td>
                                <td className="px-6 py-4 text-slate-700">{product.sales}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </section>

        <div className="space-y-6">
            <section className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
                <SectionHeader
                    title="Recent orders"
                    description="Review the latest customer orders and fulfillment status."
                />
                <div className="space-y-4">
                    {sellerOrders.map((order) => (
                        <div key={order.id} className="rounded-3xl border border-slate-200 bg-slate-50 p-4">
                            <div className="flex items-center justify-between gap-4">
                                <div>
                                    <div className="text-sm font-semibold text-slate-900">Order #{order.id}</div>
                                    <div className="text-xs text-slate-500">{order.customer}</div>
                                </div>
                                <span className={statusBadge(order.status)}>{order.status}</span>
                            </div>
                            <div className="mt-3 grid gap-2 text-sm text-slate-600">
                                <div>{order.items} items · {order.date}</div>
                                <div className="text-sm font-semibold text-slate-900">Total {order.total}</div>
                            </div>
                        </div>
                    ))}
                </div>
            </section>

            <section className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
                <SectionHeader
                    title="Payouts"
                    description="Keep track of completed and upcoming payments from recent sales."
                />
                <div className="space-y-4">
                    {sellerPayouts.map((payout) => (
                        <div key={payout.id} className="flex items-center justify-between rounded-3xl border border-slate-200 bg-slate-50 px-4 py-4">
                            <div>
                                <p className="text-sm font-semibold text-slate-900">{payout.amount}</p>
                                <p className="text-xs text-slate-500">{payout.period}</p>
                            </div>
                            <span className={statusBadge(payout.status)}>{payout.status}</span>
                        </div>
                    ))}
                </div>
            </section>
        </div>
    </div>
        </div >
    );
}
