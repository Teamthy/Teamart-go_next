"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import PageHeader from "@/components/ui/PageHeader";
import StatCard from "@/components/ui/StatCard";
import Card from "@/components/ui/card";
import Button from "@/components/ui/button";
import RouteGuard from "@/components/auth/RouteGuard";
import { getStoredCustomer, getStoredOnboarding, logout } from "@/lib/auth-state";

const accountSections = [
    {
        slug: "profile",
        title: "Profile",
        description: "Review your details, verification status, and saved preferences.",
    },
    {
        slug: "orders",
        title: "Orders",
        description: "Track recent orders, order status, and checkout notes.",
    },
    {
        slug: "wishlist",
        title: "Wishlist",
        description: "Keep your saved items close and resume shopping with one tap.",
    },
    {
        slug: "settings",
        title: "Settings",
        description: "Control your notifications, category preferences, and account visibility.",
    },
    {
        slug: "security",
        title: "Security",
        description: "Confirm your account status and manage your sign-in details.",
    },
];

export default function AccountPage() {
    const router = useRouter();
    const customer = getStoredCustomer();
    const onboarding = getStoredOnboarding();

    const completedSteps = onboarding?.progress.length ?? 0;

    return (
        <RouteGuard>
            <div className="space-y-8 pb-10">
                <PageHeader
                    title="My account"
                    description="Keep your customer profile, recent orders, and saved products aligned with the rest of your Teamart journey."
                    actions={
                        <>
                            <Button asChild variant="primary">
                                <Link href="/account/profile">Edit profile</Link>
                            </Button>
                            <Button
                                variant="secondary"
                                onClick={() => {
                                    logout();
                                    router.push("/");
                                }}
                            >
                                Log out
                            </Button>
                        </>
                    }
                />

                <div className="grid gap-4 md:grid-cols-3">
                    <StatCard label="Account status" value={customer?.verified ? "Verified" : "Pending"} helper="Your email verification state is synced from onboarding." />
                    <StatCard label="Saved steps" value={String(completedSteps)} helper="These steps were completed during onboarding and profile setup." />
                    <StatCard label="Favorite role" value={customer?.roles.creator ? "Creator" : customer?.roles.merchant ? "Merchant" : "Customer"} helper="Your active workspace role is surfaced from your saved profile." />
                </div>

                <div className="grid gap-4 xl:grid-cols-[0.95fr_1.05fr]">
                    <Card className="p-5 sm:p-6">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Profile snapshot</p>
                        <h2 className="mt-3 text-xl font-semibold text-zinc-900">{customer?.firstName} {customer?.lastName}</h2>
                        <p className="mt-3 text-sm leading-7 text-zinc-600">{customer?.email}</p>
                        <div className="mt-5 space-y-3">
                            <div className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                Verified account: {customer?.verified ? "Yes" : "Finish your onboarding verification"}
                            </div>
                            <div className="rounded-[24px] bg-zinc-50 px-4 py-3 text-sm text-zinc-700">
                                Customer access: {customer?.roles.customer ? "Enabled" : "Not active"}
                            </div>
                        </div>
                    </Card>

                    <Card className="p-5 sm:p-6">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Quick navigation</p>
                        <div className="mt-4 grid gap-3 md:grid-cols-2">
                            {accountSections.map((section) => (
                                <Link
                                    key={section.slug}
                                    href={`/account/${section.slug}`}
                                    className="rounded-[24px] border border-zinc-200 bg-white p-4 transition hover:border-[#E91E63] hover:bg-[#FFF8FB]"
                                >
                                    <p className="text-sm font-semibold text-zinc-900">{section.title}</p>
                                    <p className="mt-2 text-sm leading-6 text-zinc-600">{section.description}</p>
                                </Link>
                            ))}
                        </div>
                    </Card>
                </div>
            </div>
        </RouteGuard>
    );
}
