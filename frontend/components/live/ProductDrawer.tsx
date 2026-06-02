"use client";

import { useMemo, useState } from "react";
import type { FeedProduct } from "@/types/commerce";
import { X } from "lucide-react";

interface ProductDrawerProps {
    open: boolean;
    product: FeedProduct | null;
    onClose: () => void;
    onAddToCart: (quantity: number) => void;
    onBuyNow: () => void;
}

export default function ProductDrawer({ open, product, onClose, onAddToCart, onBuyNow }: ProductDrawerProps) {
    const [quantity, setQuantity] = useState(1);

    const totalPrice = useMemo(() => {
        if (!product) return 0;
        return product.price * quantity;
    }, [product, quantity]);

    if (!open || !product) {
        return null;
    }

    return (
        <div className="fixed inset-0 z-50 flex items-end justify-center bg-black/40 p-4 backdrop-blur-sm">
            <div className="w-full max-w-2xl rounded-[32px] bg-white p-6 shadow-2xl">
                <div className="flex items-start justify-between gap-4">
                    <div>
                        <p className="text-xs uppercase tracking-[0.28em] text-slate-500">Product details</p>
                        <h2 className="mt-3 text-2xl font-semibold text-slate-900">{product.name}</h2>
                    </div>
                    <button
                        type="button"
                        onClick={onClose}
                        className="rounded-full p-2 text-slate-700 transition hover:bg-slate-100"
                        aria-label="Close drawer"
                    >
                        <X className="h-5 w-5" />
                    </button>
                </div>
                <div className="mt-6 grid gap-6 lg:grid-cols-[1fr_0.95fr]">
                    <img
                        src={product.image_url}
                        alt={product.name}
                        className="h-72 w-full rounded-[28px] object-cover"
                    />
                    <div className="space-y-5">
                        <div className="rounded-3xl bg-slate-50 p-4">
                            <p className="text-sm text-slate-500">Merchant</p>
                            <p className="mt-2 text-base font-semibold text-slate-900">{product.merchant_name}</p>
                        </div>
                        <div className="rounded-3xl bg-slate-50 p-4">
                            <p className="text-sm text-slate-500">Price</p>
                            <p className="mt-2 text-2xl font-semibold text-slate-900">${product.price.toFixed(2)}</p>
                        </div>
                        <div className="rounded-3xl bg-slate-50 p-4">
                            <p className="text-sm text-slate-500">Quantity</p>
                            <div className="mt-3 flex items-center gap-3">
                                <button
                                    type="button"
                                    onClick={() => setQuantity((value) => Math.max(1, value - 1))}
                                    className="inline-flex h-11 w-11 items-center justify-center rounded-full border border-slate-200 text-slate-900"
                                >
                                    -
                                </button>
                                <span className="text-base font-semibold text-slate-900">{quantity}</span>
                                <button
                                    type="button"
                                    onClick={() => setQuantity((value) => value + 1)}
                                    className="inline-flex h-11 w-11 items-center justify-center rounded-full border border-slate-200 text-slate-900"
                                >
                                    +
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
                <div className="mt-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                    <div>
                        <p className="text-sm text-slate-500">Total</p>
                        <p className="mt-1 text-2xl font-semibold text-slate-900">${totalPrice.toFixed(2)}</p>
                    </div>
                    <div className="flex flex-col gap-3 sm:flex-row">
                        <button
                            type="button"
                            className="rounded-full bg-white px-5 py-3 text-sm font-semibold text-slate-900 transition hover:bg-slate-100"
                            onClick={() => onAddToCart(quantity)}
                        >
                            Add to cart
                        </button>
                        <button
                            type="button"
                            className="rounded-full bg-pink-500 px-5 py-3 text-sm font-semibold text-white transition hover:bg-pink-600"
                            onClick={onBuyNow}
                        >
                            Buy now
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
