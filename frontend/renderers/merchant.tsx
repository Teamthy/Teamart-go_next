import Link from "next/link";
import Badge from "@/components/ui/badge";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import { sellerOrders } from "@/lib/mock/orders";
import { sellerProducts } from "@/lib/mock/products";
import { merchantActions, merchantInventoryHighlights, merchantMetrics, merchantOperationalFocus, merchantPayouts, merchantQuickActions, merchantSettings } from "@/lib/mock/merchant";
import { renderHero } from "./common";

function renderMerchantPage(title: string, description: string) {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title, description, badge: "Merchant dashboard" })}
            <div className="grid gap-4 md:grid-cols-3">
                {merchantMetrics.map((item) => (
                    <Card key={item.label} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.label}</p>
                        <p className="mt-3 text-2xl font-semibold text-zinc-900">{item.value}</p>
                    </Card>
                ))}
            </div>
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Operational focus</p>
                    <div className="mt-4 space-y-3">
                        {merchantOperationalFocus.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">{item}</div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Merchant actions</p>
                    <div className="mt-4 flex flex-wrap gap-3">
                        {merchantActions.map((item) => (
                            <Button key={item.href} asChild variant={item.href === "/merchant/orders" ? "primary" : "secondary"}>
                                <Link href={item.href}>{item.label}</Link>
                            </Button>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderMerchantOrdersPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Merchant orders", description: "Review fulfillment status, order routing, and merchandising readiness in one command view.", badge: "Orders" })}
            <div className="grid gap-4 lg:grid-cols-2">
                {sellerOrders.map((order) => (
                    <Card key={order.id} className="p-5">
                        <div className="flex items-center justify-between gap-3">
                            <div>
                                <p className="text-lg font-semibold text-zinc-900">{order.id}</p>
                                <p className="mt-1 text-sm text-zinc-600">{order.customer} • {order.date}</p>
                            </div>
                            <Badge tone={order.status === "Delivered" ? "success" : "default"}>{order.status}</Badge>
                        </div>
                    </Card>
                ))}
            </div>
        </div>
    );
}

function renderMerchantProductsPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Merchant products", description: "Keep your catalog healthy with inventory checks, pricing alignment, and ready-to-promote items.", badge: "Products" })}
            <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
                {sellerProducts.map((product) => (
                    <Card key={product.id} className="p-5">
                        <p className="text-lg font-semibold text-zinc-900">{product.name}</p>
                        <p className="mt-2 text-sm text-zinc-600">{product.sku} • {product.stock} in stock</p>
                        <div className="mt-3">
                            <Badge tone="success">{product.status}</Badge>
                        </div>
                    </Card>
                ))}
            </div>
        </div>
    );
}

function renderMerchantPayoutsPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Merchant payouts", description: "Track payouts, settlement timing, and financial status with a focused merchant overview.", badge: "Payouts" })}
            <div className="grid gap-4 lg:grid-cols-2">
                {merchantPayouts.map((payout) => (
                    <Card key={payout.id} className="p-5">
                        <div className="flex items-center justify-between gap-3">
                            <div>
                                <p className="text-lg font-semibold text-zinc-900">{payout.id}</p>
                                <p className="mt-1 text-sm text-zinc-600">{payout.status}</p>
                            </div>
                            <Badge tone={payout.status === "Completed" ? "success" : "default"}>{payout.amount}</Badge>
                        </div>
                    </Card>
                ))}
            </div>
        </div>
    );
}

function renderMerchantInventoryPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Merchant inventory", description: "Spot low-stock items and keep your inventory aligned with live drops, bundles, and upcoming campaigns.", badge: "Inventory" })}
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Inventory signal</p>
                    <div className="mt-4 space-y-3">
                        {merchantInventoryHighlights.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">{item}</div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Low stock products</p>
                    <div className="mt-4 space-y-3">
                        {sellerProducts.filter((product) => product.stock <= 12).slice(0, 4).map((product) => (
                            <div key={product.id} className="rounded-[24px] bg-zinc-50 px-4 py-3 text-sm text-zinc-700">
                                <div className="flex items-center justify-between gap-3">
                                    <div>
                                        <p className="font-semibold text-zinc-900">{product.name}</p>
                                        <p className="mt-1 text-xs text-zinc-500">{product.sku} • {product.stock} in stock</p>
                                    </div>
                                    <Badge tone="warning">Restock</Badge>
                                </div>
                            </div>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderMerchantSettingsPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Merchant settings", description: "Manage brand preferences, bank details, and operational defaults for the merchant workspace.", badge: "Settings" })}
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Preferences</p>
                    <div className="mt-4 space-y-3">
                        {merchantSettings.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">{item}</div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Quick actions</p>
                    <div className="mt-4 flex flex-wrap gap-3">
                        {merchantQuickActions.map((item) => (
                            <span key={item} className="rounded-full bg-zinc-100 px-3 py-2 text-sm font-semibold text-zinc-700">{item}</span>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

export function renderMerchant(slug: string[]) {
    const route = slug[0] ?? "merchant";
    const second = slug[1];

    if (route !== "merchant") {
        return renderMerchantPage("Merchant", "A merchant workspace for orders, products, payouts, and settings.");
    }

    if (second === "orders") return renderMerchantOrdersPage();
    if (second === "products") return renderMerchantProductsPage();
    if (second === "payouts") return renderMerchantPayoutsPage();
    if (second === "inventory") return renderMerchantInventoryPage();
    if (second === "settings") return renderMerchantSettingsPage();

    return renderMerchantPage("Merchant", "A merchant workspace for orders, products, payouts, and settings.");
}
