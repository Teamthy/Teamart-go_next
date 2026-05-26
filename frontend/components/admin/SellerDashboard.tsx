"use client";

import { useEffect, useState } from "react";
import Badge from "@/components/ui/badge";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import DataTable from "@/components/admin/DataTable";
import SectionHeader from "@/components/ui/SectionHeader";
import StatusChip from "@/components/admin/StatusChip";
import * as api from "@/lib/api";

function statusTone(status?: string) {
    const value = status?.toUpperCase() || "";

    if (["PENDING", "PROCESSING", "AWAITING", "OPEN"].includes(value)) {
        return "warning";
    }

    if (["DELIVERED", "COMPLETED", "PAID"].includes(value)) {
        return "success";
    }

    if (["CANCELED", "CANCELLED", "REFUNDED"].includes(value)) {
        return "error";
    }

    return "info";
}

function formatCurrency(value?: number | string) {
    if (typeof value === "number") {
        return `$${value.toFixed(2)}`;
    }

    if (typeof value === "string") {
        return value.startsWith("$") ? value : `$${value}`;
    }

    return "$0.00";
}

function formatDate(value?: string) {
    if (!value) return "—";
    return new Date(value).toLocaleDateString();
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
    const [payouts, setPayouts] = useState<any[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchData = async () => {
            setIsLoading(true);
            setError(null);

            try {
                const productsRes = await api.listProducts(100, 0);
                const prods = productsRes.products || [];
                const ordersRes = await api.listOrders(100, 0);
                const ords = ordersRes.orders || [];

                setProducts(prods);
                setOrders(ords);
                setPayouts([
                    { id: "payout1", amount: "$2,100", period: "June 2025", status: "PENDING" },
                    { id: "payout2", amount: "$1,450", period: "May 2025", status: "COMPLETED" },
                ]);

                const inventoryValue = prods.reduce((sum: number, product: any) => {
                    const price = typeof product.price === "number" ? product.price : Number(product.price || 0);
                    const stock = Number(product.stock || 0);
                    return sum + price * stock;
                }, 0);

                const activeOrdersCount = ords.filter(
                    (order: any) =>
                        order.status &&
                        !["DELIVERED", "CANCELLED", "REFUNDED"].includes(order.status.toUpperCase())
                ).length;

                setStats({
                    totalProducts: prods.length,
                    inventoryValue,
                    activeOrders: activeOrdersCount,
                    pendingPayout: 0,
                });
            } catch (err: any) {
                setError(err.message || "Failed to load dashboard");
            } finally {
                setIsLoading(false);
            }
        };

        fetchData();
    }, []);

    const summaryCards = [
        { label: "Products", value: stats.totalProducts, detail: "Storefront items available" },
        { label: "Inventory value", value: formatCurrency(stats.inventoryValue), detail: "Projected stock value" },
        { label: "Open orders", value: stats.activeOrders, detail: "Orders in progress" },
        { label: "Pending payout", value: formatCurrency(stats.pendingPayout), detail: "Awaiting settlement" },
    ];

    const orderRows = orders.slice(0, 8).map((order) => ({
        id: order.id,
        total: formatCurrency(order.total_amount),
        status: order.status || "Processing",
        createdAt: formatDate(order.created_at),
    }));

    const productRows = products.slice(0, 8).map((product) => ({
        id: product.id,
        name: product.name,
        price: formatCurrency(product.price),
        stock: product.stock || 0,
        sku: product.sku || "N/A",
        status: product.status || "Live",
    }));

    if (isLoading) {
        return (
            <div className="mx-auto max-w-7xl px-4 py-10 sm:px-6 lg:px-8">
                <Card className="p-8 text-center text-sm text-slate-500">Loading merchant dashboard…</Card>
            </div>
        );
    }

    if (error) {
        return (
            <div className="mx-auto max-w-7xl px-4 py-10 sm:px-6 lg:px-8">
                <Card className="border-rose-200 bg-rose-50 p-4 text-sm text-rose-700">Error: {error}</Card>
            </div>
        );
    }

    return (
        <div className="mx-auto max-w-7xl space-y-8 px-4 py-10 sm:px-6 lg:px-8">
            <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
                <div className="space-y-3">
                    <Badge tone="info">Merchant dashboard</Badge>
                    <SectionHeader
                        title="Seller control center"
                        description="Keep your catalog, orders, and payout activity visible from one mobile-first workspace."
                    />
                </div>
                <div className="flex flex-wrap gap-3">
                    <Button asChild variant="secondary">
                        <a href="/feed">Preview shopper view</a>
                    </Button>
                    <Button asChild variant="primary">
                        <a href="/products">Manage catalog</a>
                    </Button>
                </div>
            </div>

            <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
                {summaryCards.map((card) => (
                    <Card key={card.label} className="p-5">
                        <p className="text-xs uppercase tracking-[0.2em] text-slate-500">{card.label}</p>
                        <p className="mt-4 text-3xl font-semibold text-slate-900">{card.value}</p>
                        <p className="mt-2 text-sm text-slate-600">{card.detail}</p>
                    </Card>
                ))}
            </div>

            <div className="grid gap-6 xl:grid-cols-[1.1fr_0.9fr]">
                <Card className="p-5">
                    <div className="mb-4">
                        <p className="text-sm font-semibold text-slate-900">Recent orders</p>
                        <p className="text-sm text-slate-500">The latest customer orders and their fulfillment lifecycle.</p>
                    </div>
                    <DataTable
                        columns={[
                            { header: "Order", accessor: "id" },
                            { header: "Total", accessor: "total" },
                            {
                                header: "Status",
                                accessor: (row) => <StatusChip label={row.status} tone={statusTone(row.status)} />,
                            },
                            { header: "Date", accessor: "createdAt" },
                        ]}
                        rows={orderRows}
                    />
                </Card>

                <div className="space-y-6">
                    <Card className="p-5">
                        <div className="mb-4">
                            <p className="text-sm font-semibold text-slate-900">Payouts</p>
                            <p className="text-sm text-slate-500">Recent payment activity and settlement status.</p>
                        </div>
                        <div className="space-y-3">
                            {payouts.map((payout) => (
                                <div key={payout.id} className="rounded-[20px] border border-slate-200 bg-slate-50 px-4 py-4">
                                    <div className="flex items-center justify-between gap-3">
                                        <div>
                                            <p className="text-sm font-semibold text-slate-900">{payout.amount}</p>
                                            <p className="text-xs text-slate-500">{payout.period}</p>
                                        </div>
                                        <StatusChip label={payout.status} tone={statusTone(payout.status)} />
                                    </div>
                                </div>
                            ))}
                        </div>
                    </Card>

                    <Card className="p-5">
                        <div className="mb-4">
                            <p className="text-sm font-semibold text-slate-900">Live campaign focus</p>
                            <p className="text-sm text-slate-500">Keep your top products featured while shoppers are active.</p>
                        </div>
                        <div className="space-y-3 text-sm text-slate-700">
                            <p>• Highlight best sellers in the hero carousel.</p>
                            <p>• Add bundle offers before peak livestream windows.</p>
                            <p>• Review order updates and refund flags daily.</p>
                        </div>
                    </Card>
                </div>
            </div>

            <Card className="p-5">
                <div className="mb-4">
                    <p className="text-sm font-semibold text-slate-900">Catalog snapshot</p>
                    <p className="text-sm text-slate-500">A concise view of the products currently on your storefront.</p>
                </div>
                <DataTable
                    columns={[
                        { header: "Product", accessor: "name" },
                        { header: "Price", accessor: "price" },
                        { header: "Stock", accessor: "stock" },
                        { header: "SKU", accessor: "sku" },
                        {
                            header: "Status",
                            accessor: (row) => <StatusChip label={row.status} tone={statusTone(row.status)} />,
                        },
                    ]}
                    rows={productRows}
                />
            </Card>
        </div>
    );
}
