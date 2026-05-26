import Link from "next/link";
import Badge from "@/components/ui/badge";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import ProductCard from "@/components/product/ProductCard";
import { sellerOrders } from "@/lib/mock/orders";
import {
    accountAddressCards,
    accountBillingHistory,
    accountDownloadItems,
    accountMemberships,
    accountOrderSnapshot,
    accountPaymentMethods,
    accountPreferences,
    accountProfile,
    accountQuickActions,
    accountReturnRequests,
    accountSecurityChecks,
    accountSupportTopics,
    accountWalletSummary,
    savedProducts,
    sessionHistory,
} from "@/lib/mock/users";
import { accountNotifications } from "@/lib/mock/notifications";
import { renderHero } from "./common";

function renderAccountPage(title: string, description: string) {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title, description, badge: "Account" })}
            <div className="grid gap-4 lg:grid-cols-[0.9fr_1.1fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Profile</p>
                    <h2 className="mt-3 text-xl font-semibold text-zinc-900">{accountProfile.name}</h2>
                    <p className="mt-2 text-sm text-zinc-600">{accountProfile.email}</p>
                    <p className="mt-1 text-sm text-zinc-600">{accountProfile.favoriteStorefront}</p>
                    <div className="mt-4 space-y-3">
                        {accountQuickActions.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Recent sessions</p>
                    <div className="mt-4 space-y-3">
                        {sessionHistory.map((session) => (
                            <div key={session.id} className="rounded-[24px] bg-zinc-50 px-4 py-3 text-sm text-zinc-700">
                                <div className="flex items-center justify-between gap-3">
                                    <span>{session.device}</span>
                                    <Badge tone={session.active ? "success" : "default"}>{session.active ? "Active" : "Inactive"}</Badge>
                                </div>
                                <p className="mt-2 text-xs text-zinc-500">{session.location}</p>
                            </div>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderAccountSecurityPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Security settings", description: "Keep your account protected with a streamlined security overview and focused recovery actions.", badge: "Account security" })}
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Protection status</p>
                    <div className="mt-4 space-y-3">
                        {accountSecurityChecks.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Recommended actions</p>
                    <div className="mt-4 space-y-3">
                        {[
                            "Review your saved payment methods",
                            "Confirm your primary email address",
                            "Use password reset if you suspect unusual activity",
                        ].map((item) => (
                            <div key={item} className="rounded-[24px] bg-zinc-50 px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                    <div className="mt-5 flex flex-wrap gap-3">
                        <Button asChild variant="primary">
                            <Link href="/auth/forgot-password">Reset password</Link>
                        </Button>
                        <Button asChild variant="secondary">
                            <Link href="/account">Back to account</Link>
                        </Button>
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderAccountOrdersPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Orders", description: "Review recent orders, delivery status, and fulfillment milestones from a calm, customer-focused view.", badge: "Account orders" })}
            <div className="grid gap-4 lg:grid-cols-[1.1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Recent orders</p>
                    <div className="mt-4 space-y-3">
                        {sellerOrders.slice(0, 4).map((order) => (
                            <div key={order.id} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                <div className="flex items-center justify-between gap-3">
                                    <div>
                                        <p className="font-semibold text-zinc-900">{order.id}</p>
                                        <p className="mt-1 text-xs text-zinc-500">{order.customer} • {order.date}</p>
                                    </div>
                                    <Badge tone={order.status === "Delivered" ? "success" : "default"}>{order.status}</Badge>
                                </div>
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Order snapshot</p>
                    <div className="mt-4 space-y-3">
                        {accountOrderSnapshot.map((item) => (
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

function renderAccountNotificationsPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Notifications", description: "Keep up with live drops, order updates, and creator activity without leaving the pink storefront experience.", badge: "Account notifications" })}
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Inbox</p>
                    <div className="mt-4 space-y-3">
                        {accountNotifications.map((item) => (
                            <div key={item.id} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                <div className="flex items-center justify-between gap-3">
                                    <div>
                                        <p className="font-semibold text-zinc-900">{item.title}</p>
                                        <p className="mt-1 text-xs text-zinc-500">{item.detail}</p>
                                    </div>
                                    <Badge tone={item.unread ? "warning" : "default"}>{item.unread ? "Unread" : "Read"}</Badge>
                                </div>
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Quick actions</p>
                    <div className="mt-4 space-y-3">
                        {accountQuickActions.map((item) => (
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

function renderAccountPaymentsPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Payments", description: "Keep billing and payment preferences tidy with a focused, commerce-first summary.", badge: "Account payments" })}
            <div className="grid gap-4 lg:grid-cols-[1.1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Payment methods</p>
                    <div className="mt-4 space-y-3">
                        {accountPaymentMethods.map((item) => (
                            <div key={item.label} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                {item.label} • {item.value}
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">This month</p>
                    <div className="mt-4 space-y-3">
                        {accountBillingHistory.map((item) => (
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

function renderAccountSavedItemsPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Saved items", description: "Keep favorite products, creator picks, and wishlist moments ready to revisit in a curated, shopper-friendly area.", badge: "Account saved items" })}
            <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
                {savedProducts.map((savedItem) => (
                    <Card key={savedItem.id} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Saved</p>
                        <p className="mt-3 text-lg font-semibold text-zinc-900">{savedItem.name}</p>
                        <p className="mt-2 text-sm text-zinc-600">{savedItem.note}</p>
                        <p className="mt-3 text-sm font-semibold text-zinc-900">{savedItem.price}</p>
                    </Card>
                ))}
            </div>
            <div className="grid gap-4 md:grid-cols-3">
                {savedProducts.map((savedItem) => (
                    <ProductCard key={savedItem.id} product={{
                        id: savedItem.id,
                        name: savedItem.name,
                        price: savedItem.price,
                        description: savedItem.note,
                        image: "/product-tote.jpg",
                    }} />
                ))}
            </div>
        </div>
    );
}

function renderAccountWalletPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Wallet", description: "Review balance, recent payouts, and payment activity in a dedicated wallet view for shopper and creator finances.", badge: "Account wallet" })}
            <div className="grid gap-4 lg:grid-cols-[1.1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Balance snapshot</p>
                    <div className="mt-4 space-y-3">
                        {accountWalletSummary.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Recent activity</p>
                    <div className="mt-4 space-y-3">
                        {accountOrderSnapshot.map((item) => (
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

function renderAccountReturnsPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Returns", description: "Review return requests, exchange status, and refund readiness in a dedicated customer support view.", badge: "Account returns" })}
            <div className="grid gap-4 lg:grid-cols-[1.1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Open requests</p>
                    <div className="mt-4 space-y-3">
                        {accountReturnRequests.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">What to know</p>
                    <div className="mt-4 space-y-3">
                        {[
                            "Most returns are processed within 5 business days",
                            "Exchange requests can be updated from this page",
                            "Support can help when an item is damaged or delayed",
                        ].map((item) => (
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

function renderAccountSubscriptionsPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Subscriptions", description: "Track recurring creator drops, renewal dates, and premium perks that keep your shopping journey active.", badge: "Account subscriptions" })}
            <div className="grid gap-4 lg:grid-cols-2">
                {accountMemberships.map((item) => (
                    <Card key={item.name} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.name}</p>
                        <p className="mt-2 text-sm font-semibold text-zinc-900">{item.level}</p>
                        <div className="mt-3 space-y-2">
                            {item.perks.map((perk) => (
                                <div key={perk} className="rounded-[20px] bg-[#FFF8FB] px-3 py-2 text-sm text-zinc-700">
                                    {perk}
                                </div>
                            ))}
                        </div>
                    </Card>
                ))}
            </div>
        </div>
    );
}

function renderAccountProfilePage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Profile", description: "Update your brand identity, personal details, and storefront preferences in a focused account editor.", badge: "Account profile" })}
            <div className="grid gap-4 lg:grid-cols-[1.1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Profile details</p>
                    <div className="mt-4 space-y-3">
                        {[
                            `Name: ${accountProfile.name}`,
                            `Email: ${accountProfile.email}`,
                            `Location: ${accountProfile.location}`,
                            `Favorite storefront: ${accountProfile.favoriteStorefront}`,
                        ].map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">What to adjust</p>
                    <div className="mt-4 space-y-3">
                        {[
                            "Update avatar and display name",
                            "Switch primary shipping address",
                            "Refine storefront interests and discovery preferences",
                        ].map((item) => (
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

function renderAccountAddressesPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Addresses", description: "Keep shipping and billing addresses tidy and ready for quick order creation, returns, or store credit.", badge: "Account addresses" })}
            <div className="grid gap-4 lg:grid-cols-2">
                {accountAddressCards.map((item) => (
                    <Card key={item.label} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.label}</p>
                        <p className="mt-3 text-sm leading-7 text-zinc-700">{item.value}</p>
                    </Card>
                ))}
            </div>
        </div>
    );
}

function renderAccountBillingPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Billing", description: "Review invoices, payment history, and renewal timing from a polished billing overview.", badge: "Account billing" })}
            <div className="grid gap-4 lg:grid-cols-[1.1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Current plan</p>
                    <div className="mt-4 space-y-3">
                        {[
                            `Plan: ${accountProfile.plan}`,
                            "Next charge: May 31",
                            "Payment method: Visa ending in 2412",
                        ].map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Billing activity</p>
                    <div className="mt-4 space-y-3">
                        {accountBillingHistory.map((item) => (
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

function renderAccountDownloadsPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Downloads", description: "Keep order receipts, product manuals, and valuable purchase documents in one easy-to-revisit place.", badge: "Account downloads" })}
            <div className="grid gap-4 lg:grid-cols-2">
                {accountDownloadItems.map((item) => (
                    <Card key={item.label} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.label}</p>
                        <p className="mt-3 text-sm text-zinc-700">{item.value}</p>
                    </Card>
                ))}
            </div>
        </div>
    );
}

function renderAccountPreferencesPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Preferences", description: "Tune your discovery, marketing, and communication preferences to match how you shop and engage.", badge: "Account preferences" })}
            <div className="grid gap-4 lg:grid-cols-2">
                {accountPreferences.map((item) => (
                    <Card key={item.label} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.label}</p>
                        <p className="mt-3 text-sm text-zinc-700">{item.value}</p>
                    </Card>
                ))}
            </div>
        </div>
    );
}

function renderAccountSupportPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Support", description: "Access help topics, order support, and the most relevant ways to get a quick response from the Teamart team.", badge: "Account support" })}
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Common help topics</p>
                    <div className="mt-4 space-y-3">
                        {accountSupportTopics.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Reach out</p>
                    <div className="mt-4 space-y-3">
                        {[
                            "Live support available during business hours",
                            "Order-specific questions are routed quickly",
                            "The team can help with refunds, returns, and shipping updates",
                        ].map((item) => (
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

export function renderAccount(slug: string[]) {
    const route = slug[0] ?? "account";
    const second = slug[1];

    if (route !== "account") {
        return renderAccountPage("Account", "Manage your profile, security, and order history in one consistent workspace.");
    }

    switch (second) {
        case "profile":
            return renderAccountProfilePage();
        case "addresses":
            return renderAccountAddressesPage();
        case "billing":
            return renderAccountBillingPage();
        case "downloads":
            return renderAccountDownloadsPage();
        case "preferences":
            return renderAccountPreferencesPage();
        case "support":
            return renderAccountSupportPage();
        case "security":
            return renderAccountSecurityPage();
        case "orders":
            return renderAccountOrdersPage();
        case "notifications":
            return renderAccountNotificationsPage();
        case "payments":
            return renderAccountPaymentsPage();
        case "saved-items":
            return renderAccountSavedItemsPage();
        case "returns":
            return renderAccountReturnsPage();
        case "subscriptions":
            return renderAccountSubscriptionsPage();
        case "wallet":
            return renderAccountWalletPage();
        default:
            return renderAccountPage("Account", "Manage your profile, security, and order history in one consistent workspace.");
    }
}
