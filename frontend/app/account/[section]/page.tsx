"use client";

import Link from "next/link";
import PageHeader from "@/components/ui/PageHeader";
import StatCard from "@/components/ui/StatCard";
import Card from "@/components/ui/card";
import Button from "@/components/ui/button";
import RouteGuard from "@/components/auth/RouteGuard";
import { getStoredCustomer } from "@/lib/auth-state";
import { products } from "@/lib/mock/products";

const accountSections = [
    {
        slug: "profile",
        title: "Profile",
        description: "Keep your customer profile, verification status, and saved preferences up to date.",
    },
    {
        slug: "orders",
        title: "Orders",
        description: "Review recent orders, shipping, and fast checkout insights.",
    },
    {
        slug: "wishlist",
        title: "Wishlist",
        description: "Return to your saved products and creator picks at any time.",
    },
    {
        slug: "settings",
        title: "Settings",
        description: "Manage your category preferences, notifications, and account display.",
    },
    {
        slug: "security",
        title: "Security",
        description: "Check your account protection and sign-in readiness.",
    },
];

export default function AccountSectionPage({ params }: { params: { section?: string } }) {
    const customer = getStoredCustomer();
    const section = accountSections.find((item) => item.slug === params.section)?.slug ?? "profile";

    const sharedStats = [
        { label: "Saved items", value: "8", helper: "A quick summary of wishlist and saved checkout moments." },
        { label: "Orders", value: "3", helper: "Recent orders are ready for review and reorder." },
        { label: "Verified", value: customer?.verified ? "Yes" : "No", helper: "Verification is reflected from your onboarding step." },
    ];

    return (
        <RouteGuard>
            <div className="space-y-8 pb-10">
                <PageHeader
                    title={`Account • ${accountSections.find((item) => item.slug === section)?.title ?? "Profile"}`}
                    description={accountSections.find((item) => item.slug === section)?.description ?? "Manage your Teamart journey from one place."}
                />

                <div className="grid gap-4 md:grid-cols-3">
                    {sharedStats.map((stat) => (
                        <StatCard key={stat.label} label={stat.label} value={stat.value} helper={stat.helper} />
                    ))}
                </div>

                <div className="grid gap-4 xl:grid-cols-[1fr_0.9fr]">
                    <Card className="p-5 sm:p-6">
                        {section === "profile" ? (
                            <div className="space-y-5">
                                <div>
                                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Customer profile</p>
                                    <h2 className="mt-3 text-xl font-semibold text-zinc-900">{customer?.firstName} {customer?.lastName}</h2>
                                    <p className="mt-3 text-sm leading-7 text-zinc-600">{customer?.email}</p>
                                </div>
                                <div className="rounded-[24px] bg-[#FFF8FB] p-4 text-sm text-zinc-700">
                                    Your saved role mix includes customer access plus {customer?.roles.creator ? "creator" : "no creator"} and {customer?.roles.merchant ? "merchant" : "no merchant"} permissions.
                                </div>
                                <Button asChild variant="secondary">
                                    <Link href="/auth/onboarding">Update onboarding</Link>
                                </Button>
                            </div>
                        ) : section === "orders" ? (
                            <div className="space-y-4">
                                <div>
                                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Recent orders</p>
                                    <h2 className="mt-3 text-xl font-semibold text-zinc-900">Order activity</h2>
                                </div>
                                {[
                                    { id: "#TA-1021", status: "Delivered", total: "$48" },
                                    { id: "#TA-1024", status: "Processing", total: "$72" },
                                ].map((order) => (
                                    <div key={order.id} className="rounded-[24px] border border-zinc-200 p-4">
                                        <div className="flex items-center justify-between gap-4">
                                            <div>
                                                <p className="font-semibold text-zinc-900">{order.id}</p>
                                                <p className="text-sm text-zinc-500">{order.status}</p>
                                            </div>
                                            <p className="text-sm font-semibold text-zinc-900">{order.total}</p>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        ) : section === "wishlist" ? (
                            <div className="space-y-4">
                                <div>
                                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Saved products</p>
                                    <h2 className="mt-3 text-xl font-semibold text-zinc-900">Wishlist</h2>
                                </div>
                                <div className="grid gap-3 md:grid-cols-2">
                                    {products.slice(0, 2).map((product) => (
                                        <div key={product.id} className="rounded-[24px] border border-zinc-200 p-4">
                                            <p className="font-semibold text-zinc-900">{product.name}</p>
                                            <p className="mt-2 text-sm text-zinc-600">{product.description}</p>
                                            <p className="mt-3 text-sm font-semibold text-zinc-900">{product.price}</p>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        ) : section === "settings" ? (
                            <div className="space-y-4">
                                <div>
                                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Preferences</p>
                                    <h2 className="mt-3 text-xl font-semibold text-zinc-900">Notification and discovery settings</h2>
                                </div>
                                <div className="space-y-3 text-sm text-zinc-700">
                                    <div className="rounded-[24px] bg-zinc-50 px-4 py-3">Live room alerts enabled</div>
                                    <div className="rounded-[24px] bg-zinc-50 px-4 py-3">Creator recommendations enabled</div>
                                    <div className="rounded-[24px] bg-zinc-50 px-4 py-3">Merchant promotions synced</div>
                                </div>
                            </div>
                        ) : (
                            <div className="space-y-4">
                                <div>
                                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Security</p>
                                    <h2 className="mt-3 text-xl font-semibold text-zinc-900">Account protection</h2>
                                </div>
                                <div className="space-y-3 text-sm text-zinc-700">
                                    <div className="rounded-[24px] bg-[#FFF8FB] px-4 py-3">Email verification: {customer?.verified ? "Completed" : "Pending"}</div>
                                    <div className="rounded-[24px] bg-zinc-50 px-4 py-3">Sign-in details are secured with local session state.</div>
                                </div>
                            </div>
                        )}
                    </Card>

                    <Card className="p-5 sm:p-6">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Browse sections</p>
                        <div className="mt-4 space-y-3">
                            {accountSections.map((item) => (
                                <Link
                                    key={item.slug}
                                    href={`/account/${item.slug}`}
                                    className={`block rounded-[24px] border p-4 transition ${section === item.slug ? "border-[#E91E63] bg-[#FFF8FB]" : "border-zinc-200 bg-white hover:border-[#E91E63]"}`}
                                >
                                    <p className="font-semibold text-zinc-900">{item.title}</p>
                                    <p className="mt-2 text-sm text-zinc-600">{item.description}</p>
                                </Link>
                            ))}
                        </div>
                    </Card>
                </div>
            </div>
        </RouteGuard>
    );
}
