"use client";

import { useState } from "react";
import CheckoutSummary from "@/components/order/CheckoutSummary";
import SectionHeader from "@/components/ui/SectionHeader";
import * as api from "@/lib/api";

export default function CheckoutPage() {
    const [formData, setFormData] = useState({
        fullName: "",
        phone: "",
        address: "",
        city: "",
        postalCode: "",
        cardNumber: "",
        expiryDate: "",
        cvc: "",
        couponCode: "",
    });

    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState(false);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData((prev) => ({
            ...prev,
            [name]: value,
        }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        setError(null);

        try {
            // Get cart from localStorage or calculate from context
            const cartItems = JSON.parse(localStorage.getItem("cart") || "[]");

            if (cartItems.length === 0) {
                setError("Your cart is empty");
                setIsLoading(false);
                return;
            }

            // Calculate total
            const totalAmount = cartItems.reduce(
                (sum: number, item: any) => sum + item.price * item.quantity,
                0
            );

            // Create order
            const orderData = {
                items: cartItems.map((item: any) => ({
                    product_id: item.id,
                    quantity: item.quantity,
                    price: item.price,
                })),
                total_amount: totalAmount,
                shipping_address: {
                    full_name: formData.fullName,
                    phone: formData.phone,
                    address: formData.address,
                    city: formData.city,
                    postal_code: formData.postalCode,
                },
                payment_method: {
                    type: "card",
                    last_four: formData.cardNumber.slice(-4),
                },
                coupon_code: formData.couponCode || undefined,
            };

            const response = await api.createOrder(orderData);

            // Clear cart
            localStorage.removeItem("cart");

            setSuccess(true);
            setFormData({
                fullName: "",
                phone: "",
                address: "",
                city: "",
                postalCode: "",
                cardNumber: "",
                expiryDate: "",
                cvc: "",
                couponCode: "",
            });

            // Redirect to orders page after success
            setTimeout(() => {
                window.location.href = `/dashboard/orders/${response.id}`;
            }, 2000);
        } catch (err: any) {
            setError(err.message || "Failed to place order");
            console.error("Error creating order:", err);
        } finally {
            setIsLoading(false);
        }
    };

    if (success) {
        return (
            <div className="space-y-8">
                <div className="bg-emerald-50 border border-emerald-200 rounded-lg p-6">
                    <h2 className="text-xl font-semibold text-emerald-900">Order Placed Successfully!</h2>
                    <p className="text-emerald-800 mt-2">
                        Your order has been confirmed. You'll be redirected to your orders page shortly.
                    </p>
                </div>
            </div>
        );
    }

    return (
        <div className="space-y-8">
            <SectionHeader
                title="Checkout"
                description="Complete your purchase with payment, shipping, and coupon options."
            />

            {error && (
                <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                    <p className="text-red-800">{error}</p>
                </div>
            )}

            <div className="grid gap-8 xl:grid-cols-[0.7fr_0.3fr]">
                <form onSubmit={handleSubmit} className="rounded-3xl border border-slate-200 bg-white dark:bg-slate-800 p-8 shadow-sm space-y-6">
                    {/* Shipping Details */}
                    <div>
                        <h3 className="text-xl font-semibold text-slate-900 dark:text-white">Shipping details</h3>
                        <p className="mt-2 text-sm text-slate-600 dark:text-slate-400">
                            Enter your delivery address for fast fulfillment.
                        </p>
                    </div>

                    <div className="grid gap-4 sm:grid-cols-2">
                        <input
                            type="text"
                            name="fullName"
                            value={formData.fullName}
                            onChange={handleChange}
                            required
                            className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 px-4 py-3 text-sm text-slate-900 dark:text-white placeholder-slate-600 dark:placeholder-slate-400"
                            placeholder="Full name"
                        />
                        <input
                            type="tel"
                            name="phone"
                            value={formData.phone}
                            onChange={handleChange}
                            required
                            className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 px-4 py-3 text-sm text-slate-900 dark:text-white placeholder-slate-600 dark:placeholder-slate-400"
                            placeholder="Phone number"
                        />
                    </div>

                    <input
                        type="text"
                        name="address"
                        value={formData.address}
                        onChange={handleChange}
                        required
                        className="w-full rounded-3xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 px-4 py-3 text-sm text-slate-900 dark:text-white placeholder-slate-600 dark:placeholder-slate-400"
                        placeholder="Street address"
                    />

                    <div className="grid gap-4 sm:grid-cols-2">
                        <input
                            type="text"
                            name="city"
                            value={formData.city}
                            onChange={handleChange}
                            required
                            className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 px-4 py-3 text-sm text-slate-900 dark:text-white placeholder-slate-600 dark:placeholder-slate-400"
                            placeholder="City"
                        />
                        <input
                            type="text"
                            name="postalCode"
                            value={formData.postalCode}
                            onChange={handleChange}
                            required
                            className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 px-4 py-3 text-sm text-slate-900 dark:text-white placeholder-slate-600 dark:placeholder-slate-400"
                            placeholder="Postal code"
                        />
                    </div>

                    {/* Payment Method */}
                    <div className="mt-10 space-y-6">
                        <div>
                            <h3 className="text-xl font-semibold text-slate-900 dark:text-white">Payment method</h3>
                            <p className="mt-2 text-sm text-slate-600 dark:text-slate-400">
                                Use the secure credit card or digital wallet integration below.
                            </p>
                        </div>

                        <input
                            type="text"
                            name="cardNumber"
                            value={formData.cardNumber}
                            onChange={handleChange}
                            required
                            className="w-full rounded-3xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 px-4 py-3 text-sm text-slate-900 dark:text-white placeholder-slate-600 dark:placeholder-slate-400"
                            placeholder="Card number"
                            maxLength="19"
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
                                maxLength="5"
                            />
                            <input
                                type="text"
                                name="cvc"
                                value={formData.cvc}
                                onChange={handleChange}
                                required
                                className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 px-4 py-3 text-sm text-slate-900 dark:text-white placeholder-slate-600 dark:placeholder-slate-400"
                                placeholder="CVC"
                                maxLength="3"
                            />
                        </div>

                        <input
                            type="text"
                            name="couponCode"
                            value={formData.couponCode}
                            onChange={handleChange}
                            className="w-full rounded-3xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 px-4 py-3 text-sm text-slate-900 dark:text-white placeholder-slate-600 dark:placeholder-slate-400"
                            placeholder="Coupon code (optional)"
                        />
                    </div>

                    <button
                        type="submit"
                        disabled={isLoading}
                        className="mt-8 w-full rounded-3xl bg-slate-900 dark:bg-indigo-600 px-6 py-4 text-sm font-semibold text-white transition hover:bg-slate-700 dark:hover:bg-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {isLoading ? "Processing..." : "Place order securely"}
                    </button>
                </form>

                <CheckoutSummary />
            </div>
        </div>
    );
}
