import Link from "next/link";
import Badge from "@/components/ui/badge";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import ProductCard from "@/components/product/ProductCard";
import { cartItems, checkoutConfidencePoints, checkoutPaymentSummary, checkoutRecoverySteps, checkoutReviewChecklist, checkoutShippingDetails, checkoutShippingNotes, checkoutSteps, launchReadinessChecklist, orderConfirmationItems, productEditorChecklist } from "@/lib/mock/orders";
import { products, recommendedProducts } from "@/lib/mock/products";
import { renderHero, titleCase } from "./common";

function renderProductsCollectionPage(title: string, description: string, badge: string, items: typeof products) {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title, description, badge })}
            <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
                {items.map((product) => (
                    <ProductCard key={product.id} product={product} />
                ))}
            </div>
        </div>
    );
}

function renderCartPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({
                title: "Cart",
                description: "Keep your order organized and ready for a smooth checkout flow.",
                badge: "Checkout-ready",
            })}
            <div className="grid gap-4 lg:grid-cols-[1.1fr_0.9fr]">
                <Card className="p-5">
                    <div className="space-y-3">
                        {cartItems.map((item) => (
                            <div key={item.id} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3">
                                <div className="flex items-center justify-between gap-4">
                                    <div>
                                        <p className="font-semibold text-zinc-900">{item.name}</p>
                                        <p className="text-sm text-zinc-600">Qty {item.qty}</p>
                                    </div>
                                    <p className="font-semibold text-zinc-900">{item.price}</p>
                                </div>
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Summary</p>
                    <h2 className="mt-3 text-xl font-semibold text-zinc-900">Order total</h2>
                    <div className="mt-4 space-y-3 text-sm text-zinc-600">
                        <div className="flex justify-between"><span>Items</span><span>$106</span></div>
                        <div className="flex justify-between"><span>Shipping</span><span>$8</span></div>
                        <div className="flex justify-between font-semibold text-zinc-900"><span>Total</span><span>$114</span></div>
                    </div>
                    <Button asChild variant="primary" className="mt-5 w-full">
                        <Link href="/checkout">Continue to checkout</Link>
                    </Button>
                </Card>
            </div>
        </div>
    );
}

function renderCheckoutPage(title: string, description: string) {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title, description, badge: "Checkout" })}
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Step</p>
                    <h2 className="mt-3 text-xl font-semibold text-zinc-900">Review your selected products</h2>
                    <div className="mt-4 space-y-3">
                        {cartItems.map((item) => (
                            <div key={item.id} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                {item.name} • {item.price}
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Progress</p>
                    <div className="mt-3 space-y-3">
                        {checkoutSteps.map((step) => (
                            <div key={step} className="rounded-[24px] bg-zinc-50 px-4 py-3 text-sm text-zinc-700">
                                {step}
                            </div>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderCheckoutShippingPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({
                title: "Shipping",
                description: "Capture address details and keep the handoff to payment smooth, clear, and fast.",
                badge: "Checkout shipping",
            })}
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Delivery details</p>
                    <div className="mt-4 space-y-3">
                        {checkoutShippingDetails.map(({ label, value }) => (
                            <div key={label} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                <span className="block text-[11px] uppercase tracking-[0.2em] text-zinc-500">{label}</span>
                                <span className="mt-1 block text-sm text-zinc-800">{value}</span>
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Shipping notes</p>
                    <div className="mt-4 space-y-3">
                        {checkoutShippingNotes.map((item) => (
                            <div key={item} className="rounded-[24px] bg-zinc-50 px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderCheckoutPaymentPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({
                title: "Payment",
                description: "Review methods, totals, and the final approval path before the order is locked in.",
                badge: "Checkout payment",
            })}
            <div className="grid gap-4 lg:grid-cols-[1.1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Payment summary</p>
                    <div className="mt-4 space-y-3">
                        {checkoutPaymentSummary.map(({ label, value }) => (
                            <div key={label} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                <span className="block text-[11px] uppercase tracking-[0.2em] text-zinc-500">{label}</span>
                                <span className="mt-1 block text-sm text-zinc-800">{value}</span>
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Checkout confidence</p>
                    <div className="mt-4 space-y-3">
                        {checkoutConfidencePoints.map((item) => (
                            <div key={item} className="rounded-[24px] bg-zinc-50 px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderCheckoutReviewPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({
                title: "Review",
                description: "Validate the order, shipping details, and promo selection before final confirmation.",
                badge: "Checkout review",
            })}
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Review checklist</p>
                    <div className="mt-4 space-y-3">
                        {checkoutReviewChecklist.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Order details</p>
                    <div className="mt-4 space-y-3">
                        {cartItems.map((item) => (
                            <div key={item.id} className="rounded-[24px] bg-zinc-50 px-4 py-3 text-sm text-zinc-700">
                                {item.name} • {item.price} • Qty {item.qty}
                            </div>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderCheckoutSuccessPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({
                title: "Success",
                description: "Your order is confirmed and the thank-you experience is ready to show.",
                badge: "Checkout success",
            })}
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Confirmation</p>
                    <div className="mt-4 space-y-3">
                        {orderConfirmationItems.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">What to do next</p>
                    <div className="mt-4 flex flex-wrap gap-3">
                        <Button asChild variant="primary"><Link href="/account/orders">View orders</Link></Button>
                        <Button asChild variant="secondary"><Link href="/explore">Keep shopping</Link></Button>
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderCheckoutFailedPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({
                title: "Failed",
                description: "A calm recovery space for shoppers who need to adjust or retry payment.",
                badge: "Checkout recovery",
            })}
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Recovery steps</p>
                    <div className="mt-4 space-y-3">
                        {checkoutRecoverySteps.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Next actions</p>
                    <div className="mt-4 flex flex-wrap gap-3">
                        <Button asChild variant="primary"><Link href="/checkout/payment">Retry payment</Link></Button>
                        <Button asChild variant="secondary"><Link href="/cart">Edit cart</Link></Button>
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderProductEditPage(productId: string) {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({
                title: `Edit product ${productId}`,
                description: "Update pricing, copy, and launch details for this catalog item while keeping the storefront styling consistent.",
                badge: "Product editor",
            })}
            <div className="grid gap-4 lg:grid-cols-[1.1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Draft checklist</p>
                    <div className="mt-4 space-y-3">
                        {productEditorChecklist.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Launch readiness</p>
                    <div className="mt-4 space-y-3">
                        {launchReadinessChecklist.map((item) => (
                            <div key={item} className="rounded-[24px] bg-zinc-50 px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

export function renderProducts(slug: string[]) {
    const route = slug[0] ?? "products";
    const second = slug[1];
    const third = slug[2];

    if (route === "cart") {
        return renderCartPage();
    }

    if (route === "checkout") {
        if (second === "shipping") return renderCheckoutShippingPage();
        if (second === "payment") return renderCheckoutPaymentPage();
        if (second === "review") return renderCheckoutReviewPage();
        if (second === "success") return renderCheckoutSuccessPage();
        if (second === "failed") return renderCheckoutFailedPage();
        return renderCheckoutPage("Checkout", "A step-based checkout experience built to feel polished and conversion-focused.");
    }

    if (route === "wishlist") {
        return renderProductsCollectionPage("Wishlist", "Keep favorite products saved and ready to come back to at any time.", "Saved for later", products);
    }

    if (route === "products") {
        if (second === "edit" && third) {
            return renderProductEditPage(third);
        }
        if (second === "new") {
            return renderProductsCollectionPage("Create product", "Design a polished new catalog entry with the same look and feel as the rest of the storefront.", "Product builder", products);
        }
        if (second === "drafts") {
            return renderProductsCollectionPage("Drafts", "A clean drafting workspace for products that are still being prepared for launch.", "Draft mode", products);
        }
        if (second === "featured") {
            return renderProductsCollectionPage("Featured", "Highlight your best-performing products in a premium storefront section.", "Featured", recommendedProducts);
        }
        if (second === "trending") {
            return renderProductsCollectionPage("Trending", "A live pulse of the products that are moving fastest across your marketplace.", "Trending now", products);
        }
        if (second === "recommended") {
            return renderProductsCollectionPage("Recommended", "A curated recommendation rail designed for follow-up browsing and conversion.", "Recommendation mode", recommendedProducts);
        }
        if (second === "recent") {
            return renderProductsCollectionPage("Recent", "A clear view of the items shoppers returned to most recently.", "Recently viewed", products);
        }
        if (second === "deals") {
            return renderProductsCollectionPage("Deals", "A promotional storefront section where savings are highlighted for quick decision-making.", "Deals", recommendedProducts);
        }
        if (second === "saved") {
            return renderProductsCollectionPage("Saved", "A simple saved-items page so customers can come back to favorite products.", "Saved items", products);
        }
        if (second === "compare") {
            return renderProductsCollectionPage("Compare", "Compare product highlights side-by-side in a shopper-friendly layout.", "Compare mode", products);
        }
        if (second === "reviews" || second === "questions" || second === "related" || second === "share") {
            return renderProductsCollectionPage(`${titleCase(second)} view`, `A dedicated ${titleCase(second)} section that keeps the product journey visually aligned with the rest of the marketplace.`, `Product detail`, products);
        }
        return renderProductsCollectionPage("Products", "Browse ongoing drops and curated product collections in a pink, conversion-ready layout.", "Catalog", products);
    }

    return renderProductsCollectionPage("Products", "Browse ongoing drops and curated product collections in a pink, conversion-ready layout.", "Catalog", products);
}
