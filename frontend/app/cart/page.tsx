"use client";

import Link from "next/link";
import { useMemo, useState } from "react";
import CartItemRow from "@/components/cart/CartItemRow";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import Input from "@/components/ui/input";
import Badge from "@/components/ui/badge";
import SectionHeader from "@/components/ui/SectionHeader";

const starterItems = [
    { id: "1", name: "Sustainable Canvas Tote", subtitle: "Recycled canvas · Sand", price: 32, quantity: 1 },
    { id: "2", name: "Artist Collaboration Hoodie", subtitle: "Limited drop · Olive", price: 74, quantity: 2 },
];

export default function CartPage() {
    const [items, setItems] = useState(starterItems);
    const [coupon, setCoupon] = useState("");

    const subtotal = useMemo(() => items.reduce((sum, item) => sum + item.price * item.quantity, 0), [items]);
    const shipping = 8;
    const discount = coupon.trim().toLowerCase() === "teamart10" ? Math.round(subtotal * 0.1) : 0;
    const total = Math.max(subtotal + shipping - discount, 0);

    const updateQty = (id: string, quantity: number) => {
        setItems((current) => current.map((item) => (item.id === id ? { ...item, quantity } : item)));
    };

    const removeItem = (id: string) => {
        setItems((current) => current.filter((item) => item.id !== id));
    };

    return (
        <div className="space-y-8 pb-10">
            <section className="rounded-[28px] bg-[linear-gradient(135deg,#ecfdf5_0%,#ffffff_45%,#fef3c7_100%)] p-5 sm:p-6">
                <SectionHeader
                    title="Your cart"
                    description="Build your order, apply a coupon, and move from discovery to checkout in a few taps."
                />
                <div className="mt-4 flex flex-wrap gap-3">
                    <Badge tone="success">{items.length} items ready</Badge>
                    <Badge tone="info">Express checkout available</Badge>
                </div>
            </section>

            <div className="grid gap-8 xl:grid-cols-[1.1fr_0.9fr]">
                <div className="space-y-4">
                    {items.map((item) => (
                        <CartItemRow
                            key={item.id}
                            name={item.name}
                            subtitle={item.subtitle}
                            price={item.price}
                            quantity={item.quantity}
                            onChange={(quantity) => updateQty(item.id, quantity)}
                            onRemove={() => removeItem(item.id)}
                        />
                    ))}
                    {items.length === 0 && (
                        <Card className="p-8 text-center text-slate-600">
                            Your cart is empty. Add products from the feed or search page to get started.
                        </Card>
                    )}
                </div>

                <div className="space-y-4">
                    <Card className="p-5">
                        <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Offer code</p>
                        <div className="mt-4">
                            <Input
                                label="Coupon"
                                placeholder="teamart10"
                                value={coupon}
                                onChange={(e) => setCoupon(e.target.value)}
                                helperText="Try teamart10 for 10% off your order"
                            />
                        </div>
                    </Card>
                    <Card className="p-5">
                        <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Summary</p>
                        <div className="mt-4 space-y-3 text-sm text-slate-700">
                            <div className="flex items-center justify-between">
                                <span>Subtotal</span>
                                <span>${subtotal.toFixed(2)}</span>
                            </div>
                            <div className="flex items-center justify-between">
                                <span>Shipping</span>
                                <span>${shipping.toFixed(2)}</span>
                            </div>
                            {discount > 0 && (
                                <div className="flex items-center justify-between text-emerald-600">
                                    <span>Coupon</span>
                                    <span>-${discount.toFixed(2)}</span>
                                </div>
                            )}
                        </div>
                        <div className="mt-4 border-t border-slate-200 pt-4">
                            <div className="flex items-center justify-between text-base font-semibold text-slate-900">
                                <span>Total</span>
                                <span>${total.toFixed(2)}</span>
                            </div>
                        </div>
                        <div className="mt-5 flex flex-col gap-3 sm:flex-row">
                            <Button asChild variant="primary" className="flex-1">
                                <Link href="/checkout">Continue to checkout</Link>
                            </Button>
                            <Button asChild variant="secondary" className="flex-1">
                                <Link href="/feed">Keep browsing</Link>
                            </Button>
                        </div>
                    </Card>
                </div>
            </div>
        </div>
    );
}
