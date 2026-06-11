"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import SectionHeader from "@/components/ui/SectionHeader";
import * as api from "@/lib/api";

interface OrderDetail {
    id: string;
    order_number?: string;
    status?: string;
    total_amount?: number;
    created_at?: string;
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

export default function CheckoutSuccessPage() {
    const searchParams = useSearchParams();
    const router = useRouter();
    const orderId = searchParams.get("orderId") ?? "";
    const [order, setOrder] = useState<OrderDetail | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (!orderId) {
            return;
        }

        setIsLoading(true);
        setError(null);

        api.getOrder(orderId)
            .then((response: any) => {
                const orderData = response?.order || response;
                setOrder({
                    id: String(orderData?.id ?? orderId),
                    order_number: orderData?.order_number || orderData?.orderNumber || String(orderData?.id ?? orderId),
                    status: orderData?.status || "pending",
                    total_amount: Number(orderData?.total_amount ?? orderData?.totalAmount ?? 0),
                    created_at: orderData?.created_at || orderData?.createdAt || "",
                });
            })
            .catch((err: any) => {
                setError(err?.message || "Unable to load order details.");
            })
            .finally(() => setIsLoading(false));
    }, [orderId]);

    return (
        <div className="space-y-8 px-4 py-10 sm:px-6 lg:px-8">
            <SectionHeader
                title="Order confirmed"
                description="Your order has been received. Review the confirmation details and track it from your account."
            />
            <Card className="space-y-6 p-8">
                {error ? (
                    <div className="rounded-3xl border border-rose-200 bg-rose-50 p-6 text-sm text-rose-700">
                        <p>{error}</p>
                    </div>
                ) : isLoading ? (
                    <div className="rounded-3xl border border-slate-200 bg-slate-50 p-6 text-sm text-slate-700">
                        Loading order details…
                    </div>
                ) : order ? (
                    <div className="space-y-6">
                        <div className="rounded-3xl border border-slate-200 bg-emerald-50 p-6 text-slate-900">
                            <p className="font-semibold text-lg">Thank you! Your order is on its way.</p>
                            <p className="mt-2 text-sm text-slate-700">Order #{order.order_number || order.id} has been placed successfully.</p>
                        </div>

                        <div className="grid gap-4 sm:grid-cols-3">
                            <div className="rounded-3xl border border-slate-200 bg-white p-5">
                                <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Order</p>
                                <p className="mt-3 text-lg font-semibold text-slate-900">{order.order_number || order.id}</p>
                            </div>
                            <div className="rounded-3xl border border-slate-200 bg-white p-5">
                                <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Total</p>
                                <p className="mt-3 text-lg font-semibold text-slate-900">{formatCurrency(order.total_amount)}</p>
                            </div>
                            <div className="rounded-3xl border border-slate-200 bg-white p-5">
                                <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Status</p>
                                <p className="mt-3 text-lg font-semibold text-slate-900">{order.status}</p>
                            </div>
                        </div>

                        <div className="rounded-3xl border border-slate-200 bg-white p-5">
                            <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Placed on</p>
                            <p className="mt-3 text-lg font-semibold text-slate-900">
                                {order.created_at ? new Date(order.created_at).toLocaleString() : "Pending"}
                            </p>
                        </div>
                    </div>
                ) : (
                    <div className="rounded-3xl border border-slate-200 bg-slate-50 p-6 text-slate-700">
                        <p>No order details were available. If you just completed checkout, try again in a few moments.</p>
                    </div>
                )}

                <div className="flex flex-wrap gap-4">
                    <Button onClick={() => router.push("/account/orders")}>View my orders</Button>
                    <Button variant="secondary" onClick={() => router.push("/")}>Continue shopping</Button>
                    <Button variant="secondary" asChild>
                        <Link href="/support">Contact support</Link>
                    </Button>
                </div>
            </Card>
        </div>
    );
}
