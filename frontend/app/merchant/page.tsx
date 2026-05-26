"use client";

import { useEffect } from "react";
import Link from "next/link";
import PageHeader from "@/components/ui/PageHeader";
import RoleCard from "@/components/ui/RoleCard";
import StatCard from "@/components/ui/StatCard";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import StoreCard from "@/components/ui/StoreCard";
import RouteGuard from "@/components/auth/RouteGuard";
import { getStoredCustomer, hasRole, setWorkspaceRole } from "@/lib/auth-state";
import { merchantProfiles } from "@/lib/mock/merchant";

export default function MerchantPage() {
    const customer = getStoredCustomer();
    const isMerchant = hasRole("merchant");

    useEffect(() => {
        if (isMerchant) {
            setWorkspaceRole("merchant");
        }
    }, [isMerchant]);

    return (
        <RouteGuard requiredRole="merchant" redirectTo="/auth/merchant">
            <div className="space-y-8 pb-10">
                <PageHeader
                    title="Merchant workspace"
                    description="Keep your storefront, orders, and live promotions organized in one merchant-first workspace for shoppers and campaigns."
                    actions={
                        <>
                            <Button asChild variant="primary">
                                <Link href="/merchant/orders">Review orders</Link>
                            </Button>
                            <Button asChild variant="secondary">
                                <Link href="/merchant/products">Manage products</Link>
                            </Button>
                        </>
                    }
                />

                <div className="grid gap-4 md:grid-cols-3">
                    <StatCard label="Active storefronts" value="320" helper="Your merchant network is expanding quickly." />
                    <StatCard label="Orders today" value="124" helper="Fast-moving bundles and live promotions are keeping momentum high." />
                    <StatCard label="Customer mood" value="92%" helper={customer ? `Welcome back, ${customer.firstName}.` : "Guest visitors are browsing merchant spotlighted stores."} />
                </div>

                <div className="grid gap-4 xl:grid-cols-[0.85fr_1.15fr]">
                    <Card className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Merchant signal</p>
                        <h2 className="mt-3 text-xl font-semibold text-zinc-900">Turn storefront activity into consistent conversions</h2>
                        <p className="mt-3 text-sm leading-7 text-zinc-600">
                            Use the merchant workspace to reinforce live promotions, optimize product readiness, and guide shoppers toward your strongest offers.
                        </p>
                        <div className="mt-5 space-y-3">
                            {[
                                "Review low-stock products before every livestream",
                                "Highlight bundles that convert well during live room moments",
                                "Use the payouts and inventory pages to keep the business running smoothly",
                            ].map((item) => (
                                <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                    {item}
                                </div>
                            ))}
                        </div>
                    </Card>

                    <div className="space-y-4">
                        <RoleCard
                            title={isMerchant ? "Merchant access is active" : "Open merchant access"}
                            description={isMerchant ? "Your storefront can now run live offers, merchant promos, and order operations from one workspace." : "Create a merchant account to unlock storefront management, inventory insights, and payouts."}
                            requirements={isMerchant ? ["Storefront synced", "Orders and payouts ready", "Live promotion controls enabled"] : ["Store name and owner details", "A business email address", "Inventory and payout preferences"]}
                            ctaLabel={isMerchant ? "Open merchant tools" : "Create merchant account"}
                            href={isMerchant ? "/merchant/products" : "/auth/register?role=merchant"}
                            tone={isMerchant ? "success" : "warning"}
                        />
                    </div>
                </div>

                <div className="space-y-3">
                    <div className="flex items-center justify-between gap-3">
                        <div>
                            <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Store spotlight</p>
                            <h2 className="mt-2 text-xl font-semibold text-zinc-900">Merchant storefronts worth watching</h2>
                        </div>
                        <Button asChild variant="secondary">
                            <Link href="/stores">Browse all stores</Link>
                        </Button>
                    </div>
                    <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
                        {merchantProfiles.slice(0, 3).map((profile) => (
                            <StoreCard
                                key={profile.id}
                                name={profile.name}
                                slug={profile.store}
                                category={profile.category}
                                rating={profile.rating}
                                banner={`https://images.unsplash.com/photo-1505693416388-ac5ce068fe85?w=1200&q=80`}
                                tagline={`A polished storefront with ${profile.products} products and a live-first selling strategy.`}
                                live={profile.live}
                                products={profile.products}
                            />
                        ))}
                    </div>
                </div>
            </div>
        </RouteGuard>
    );
}
