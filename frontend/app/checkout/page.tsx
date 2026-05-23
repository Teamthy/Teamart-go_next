"use client";

import { useState } from "react";
import { z } from "zod";
import CheckoutSummary from "@/components/order/CheckoutSummary";
import SectionHeader from "@/components/ui/SectionHeader";
import * as api from "@/lib/api";

const checkoutSchema = z.object({
    fullName: z.string().trim().min(2, "Full name is required"),
    phone: z.string().trim().min(10, "Enter a valid phone number"),
    address: z.string().trim().min(5, "Enter your street address"),
    city: z.string().trim().min(2, "Enter your city"),
    postalCode: z.string().trim().min(4, "Enter your postal code"),
    cardNumber: z.string().trim().min(12, "Enter a valid card number"),
    expiryDate: z.string().trim().regex(/^(0[1-9]|1[0-2])\s*\/\s*\d{2}$/, "Enter expiry as MM / YY"),
    cvc: z.string().trim().regex(/^\d{3,4}$/, "Enter a valid CVC"),
    couponCode: z.string().trim().max(20).optional(),
});

type CheckoutInput = z.infer<typeof checkoutSchema>;

export default function CheckoutPage() {
    const [formData, setFormData] = useState<CheckoutInput>({
        fullName: "",
        phone: "",
        address: "",
        city: "",
        postalCode: "",
        cardNumber: "",
        expiryDate: "",
        cvc: "",
        couponCode: undefined,
    });

    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
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
        setFieldErrors({});

        const parsed = checkoutSchema.safeParse(formData);
        if (!parsed.success) {
            const errors = parsed.error.flatten().fieldErrors;
            const nextFieldErrors: Record<string, string> = {};
            Object.entries(errors).forEach(([key, messages]) => {
                if (messages?.length) {
                    nextFieldErrors[key] = messages[0];
                }
            });
            setFieldErrors(nextFieldErrors);
            setError("Please fix the highlighted fields before continuing.");
            setIsLoading(false);
            return;
        }

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
                    full_name: parsed.data.fullName,
                    phone: parsed.data.phone,
                    address: parsed.data.address,
                    city: parsed.data.city,
                    postal_code: parsed.data.postalCode,
                },
                payment_method: {
                    type: "card",
                    last_four: parsed.data.cardNumber.slice(-4),
                },
                coupon_code: parsed.data.couponCode || undefined,
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
                        <div>
                            <input
                                type="text"
                                name="fullName"
                                value={formData.fullName}
                                onChange={handleChange}
                                required
                                aria-invalid={!!fieldErrors.fullName}
                                className={`rounded-3xl border px-4 py-3 text-sm placeholder-slate-600 dark:placeholder-slate-400 w-full ${fieldErrors.fullName
                                    ? "border-rose-500 focus:border-rose-500"
                                    : "border-slate-200 dark:border-slate-700"
                                    } bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-white`}
                                placeholder="Full name"
                            />
                            {fieldErrors.fullName ? (
                                <p className="mt-2 text-sm text-rose-600">{fieldErrors.fullName}</p>
                            ) : null}
                        </div>
                        <div>
                            <input
                                type="tel"
                                name="phone"
                                value={formData.phone}
                                onChange={handleChange}
                                required
                                aria-invalid={!!fieldErrors.phone}
                                className={`rounded-3xl border px-4 py-3 text-sm placeholder-slate-600 dark:placeholder-slate-400 w-full ${fieldErrors.phone
                                    ? "border-rose-500 focus:border-rose-500"
                                    : "border-slate-200 dark:border-slate-700"
                                    } bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-white`}
                                placeholder="Phone number"
                            />
                            {fieldErrors.phone ? (
                                <p className="mt-2 text-sm text-rose-600">{fieldErrors.phone}</p>
                            ) : null}
                        </div>
                    </div>

                    <div>
                        <input
                            type="text"
                            name="address"
                            value={formData.address}
                            onChange={handleChange}
                            required
                            aria-invalid={!!fieldErrors.address}
                            className={`w-full rounded-3xl border px-4 py-3 text-sm placeholder-slate-600 dark:placeholder-slate-400 ${fieldErrors.address
                                ? "border-rose-500 focus:border-rose-500"
                                : "border-slate-200 dark:border-slate-700"
                                } bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-white`}
                            placeholder="Street address"
                        />
                        {fieldErrors.address ? (
                            <p className="mt-2 text-sm text-rose-600">{fieldErrors.address}</p>
                        ) : null}
                    </div>

                    <div className="grid gap-4 sm:grid-cols-2">
                        <div>
                            <input
                                type="text"
                                name="city"
                                value={formData.city}
                                onChange={handleChange}
                                required
                                aria-invalid={!!fieldErrors.city}
                                className={`rounded-3xl border px-4 py-3 text-sm placeholder-slate-600 dark:placeholder-slate-400 w-full ${fieldErrors.city
                                    ? "border-rose-500 focus:border-rose-500"
                                    : "border-slate-200 dark:border-slate-700"
                                    } bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-white`}
                                placeholder="City"
                            />
                            {fieldErrors.city ? (
                                <p className="mt-2 text-sm text-rose-600">{fieldErrors.city}</p>
                            ) : null}
                        </div>
                        <div>
                            <input
                                type="text"
                                name="postalCode"
                                value={formData.postalCode}
                                onChange={handleChange}
                                required
                                aria-invalid={!!fieldErrors.postalCode}
                                className={`rounded-3xl border px-4 py-3 text-sm placeholder-slate-600 dark:placeholder-slate-400 w-full ${fieldErrors.postalCode
                                    ? "border-rose-500 focus:border-rose-500"
                                    : "border-slate-200 dark:border-slate-700"
                                    } bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-white`}
                                placeholder="Postal code"
                            />
                            {fieldErrors.postalCode ? (
                                <p className="mt-2 text-sm text-rose-600">{fieldErrors.postalCode}</p>
                            ) : null}
                        </div>
                    </div>

                    {/* Payment Method */}
                    <div className="mt-10 space-y-6">
                        <div>
                            <h3 className="text-xl font-semibold text-slate-900 dark:text-white">Payment method</h3>
                            <p className="mt-2 text-sm text-slate-600 dark:text-slate-400">
                                Use the secure credit card or digital wallet integration below.
                            </p>
                        </div>

                        <div>
                            <input
                                type="text"
                                name="cardNumber"
                                value={formData.cardNumber}
                                onChange={handleChange}
                                required
                                aria-invalid={!!fieldErrors.cardNumber}
                                className={`w-full rounded-3xl border px-4 py-3 text-sm placeholder-slate-600 dark:placeholder-slate-400 ${fieldErrors.cardNumber
                                    ? "border-rose-500 focus:border-rose-500"
                                    : "border-slate-200 dark:border-slate-700"
                                    } bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-white`}
                                placeholder="Card number"
                                maxLength={19}
                            />
                            {fieldErrors.cardNumber ? (
                                <p className="mt-2 text-sm text-rose-600">{fieldErrors.cardNumber}</p>
                            ) : null}
                        </div>

                        <div className="grid gap-4 sm:grid-cols-2">
                            <div>
                                <input
                                    type="text"
                                    name="expiryDate"
                                    value={formData.expiryDate}
                                    onChange={handleChange}
                                    required
                                    aria-invalid={!!fieldErrors.expiryDate}
                                    className={`rounded-3xl border px-4 py-3 text-sm placeholder-slate-600 dark:placeholder-slate-400 w-full ${fieldErrors.expiryDate
                                        ? "border-rose-500 focus:border-rose-500"
                                        : "border-slate-200 dark:border-slate-700"
                                        } bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-white`}
                                    placeholder="MM / YY"
                                    maxLength={5}
                                />
                                {fieldErrors.expiryDate ? (
                                    <p className="mt-2 text-sm text-rose-600">{fieldErrors.expiryDate}</p>
                                ) : null}
                            </div>
                            <div>
                                <input
                                    type="text"
                                    name="cvc"
                                    value={formData.cvc}
                                    onChange={handleChange}
                                    required
                                    aria-invalid={!!fieldErrors.cvc}
                                    className={`rounded-3xl border px-4 py-3 text-sm placeholder-slate-600 dark:placeholder-slate-400 w-full ${fieldErrors.cvc
                                        ? "border-rose-500 focus:border-rose-500"
                                        : "border-slate-200 dark:border-slate-700"
                                        } bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-white`}
                                    placeholder="CVC"
                                    maxLength={4}
                                />
                                {fieldErrors.cvc ? (
                                    <p className="mt-2 text-sm text-rose-600">{fieldErrors.cvc}</p>
                                ) : null}
                            </div>
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
