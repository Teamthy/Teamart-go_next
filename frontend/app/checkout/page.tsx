"use client";

import { useState } from "react";
import CheckoutSummary from "@/components/order/CheckoutSummary";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import Input from "@/components/ui/input";
import Badge from "@/components/ui/badge";
import SectionHeader from "@/components/ui/SectionHeader";
import * as api from "@/lib/api";

const initialForm = {
    fullName: "",
    phone: "",
    address: "",
    city: "",
    postalCode: "",
    cardNumber: "",
    expiryDate: "",
    cvc: "",
    couponCode: "",
};

export default function CheckoutPage() {
    const [formData, setFormData] = useState(initialForm);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState(false);

    const readUserId = () => {
        if (typeof window === "undefined") {
            return null;
        }

        try {
            const sessionRaw = localStorage.getItem("session");
            const userRaw = localStorage.getItem("user");
            const session = sessionRaw ? JSON.parse(sessionRaw) : null;
            const user = userRaw ? JSON.parse(userRaw) : null;

            const userId = Number(session?.user_id ?? session?.user?.id ?? user?.id ?? 0);
            return Number.isFinite(userId) && userId > 0 ? userId : null;
        } catch {
            return null;
        }
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData((prev) => ({ ...prev, [name]: value }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        setError(null);

        try {
            const cartItems = JSON.parse(localStorage.getItem("cart") || "[]");
            const userId = readUserId();

            if (!cartItems.length) {
                throw new Error("Your cart is empty. Add products before checkout.");
            }

            if (!userId) {
                throw new Error("Please sign in to place an order.");
            }

            const totalAmount = cartItems.reduce((sum: number, item: any) => sum + Number(item.price || 0) * Number(item.quantity || 1), 0);

            const response = await api.createOrder({
                user_id: userId,
                total_amount: totalAmount,
                status: "pending",
            });

            localStorage.removeItem("cart");
            setSuccess(true);
            setFormData(initialForm);
            setTimeout(() => {
                window.location.href = `/dashboard/orders/${response.id}`;
            }, 1500);
        } catch (err: any) {
            setError(err.message || "Failed to place order");
        } finally {
            setIsLoading(false);
        }
    };

    if (success) {
        return (
            <div className="mx-auto max-w-2xl rounded-[28px] bg-emerald-50 p-8 text-emerald-900">
                <Badge tone="success">Order placed</Badge>
                <h2 className="mt-3 text-2xl font-semibold">Your order is confirmed</h2>
                <p className="mt-2 text-sm leading-7">
                    You’re being redirected to your order details page so you can track fulfillment in real time.
                </p>
            </div>
        );
    }

    return (
        <div className="space-y-8 pb-10">
            <section className="rounded-[28px] bg-[linear-gradient(135deg,#ecfdf5_0%,#ffffff_55%,#fef3c7_100%)] p-5 sm:p-6">
                <SectionHeader
                    title="Checkout"
                    description="Secure your order with a fast and polished checkout flow that feels native to the mobile app."
                />
                <div className="mt-4 flex flex-wrap gap-3">
                    <Badge tone="info">Secure payment</Badge>
                    <Badge tone="success">Express shipping</Badge>
                </div>
            </section>

            {error && (
                <div className="rounded-[24px] border border-rose-200 bg-rose-50 p-4 text-rose-700">
                    {error}
                </div>
            )}

            <div className="grid gap-8 xl:grid-cols-[1.1fr_0.9fr]">
                <Card className="p-6">
                    <form onSubmit={handleSubmit} className="space-y-5">
                        <div>
                            <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Shipping details</p>
                            <h2 className="mt-2 text-xl font-semibold text-slate-900">Where should it go?</h2>
                        </div>

                        <input
                            type="text"
                            name="cardNumber"
                            value={formData.cardNumber}
                            onChange={handleChange}
                            required
                            className="w-full rounded-3xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 px-4 py-3 text-sm text-slate-900 dark:text-white placeholder-slate-600 dark:placeholder-slate-400"
                            placeholder="Card number"
                            maxLength={19}
                        />

                        <div className="grid gap-4 sm:grid-cols-2">
                            <input
                                type="text"
                                name="expiryDate"
                                value={formData.expiryDate}
                                onChange={handleChange}
                                required
                                className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 px-4 py-3 text-sm text-slate-900 dark:text-white placeholder-slate-600 dark:placeholder-slate-400"
                                placeholder="MM / YY"
                                maxLength={5}
                            />
                            <input
                                type="text"
                                name="cvc"
                                value={formData.cvc}
                                onChange={handleChange}
                                required
                                className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 px-4 py-3 text-sm text-slate-900 dark:text-white placeholder-slate-600 dark:placeholder-slate-400"
                                placeholder="CVC"
                                maxLength={3}
                            />
                        </div>

                        <div className="pt-2">
                            <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Payment</p>
                            <h2 className="mt-2 text-xl font-semibold text-slate-900">Secure card entry</h2>
                        </div>
                        <Input label="Card number" name="cardNumber" value={formData.cardNumber} onChange={handleChange} required maxLength={19} placeholder="4242 4242 4242 4242" />
                        <div className="grid gap-4 sm:grid-cols-2">
                            <Input label="Expiry" name="expiryDate" value={formData.expiryDate} onChange={handleChange} required placeholder="MM / YY" maxLength={5} />
                            <Input label="CVC" name="cvc" value={formData.cvc} onChange={handleChange} required maxLength={3} />
                        </div>
                        <Input label="Coupon code" name="couponCode" value={formData.couponCode} onChange={handleChange} placeholder="Optional" />

                        <Button type="submit" variant="primary" className="w-full py-4">
                            {isLoading ? "Processing…" : "Place order securely"}
                        </Button>
                    </form>
                </Card>

                <CheckoutSummary />
            </div>
        </div>
    );
}
