import Link from "next/link";
import Badge from "@/components/ui/badge";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import { sellerProducts } from "@/lib/mock/products";
import { creatorAnalyticsStats, creatorDashboardStats, creatorQuickLinks, creatorStudioFocus, creatorWorkspaceActions } from "@/lib/mock/creators";
import { creatorPerformancePulse, creatorRecommendedActions } from "@/lib/mock/analytics";
import { renderHero } from "./common";

function renderCreatorPage(title: string, description: string) {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title, description, badge: "Creator studio" })}
            <div className="grid gap-4 md:grid-cols-3">
                {creatorDashboardStats.map((item) => (
                    <Card key={item.label} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.label}</p>
                        <p className="mt-3 text-2xl font-semibold text-zinc-900">{item.value}</p>
                    </Card>
                ))}
            </div>
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Featured actions</p>
                    <div className="mt-4 space-y-3">
                        {creatorWorkspaceActions.map((item) => (
                            <div key={item.label} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                <p className="font-semibold text-zinc-900">{item.label}</p>
                                <p className="mt-1 text-xs text-zinc-500">{item.value}</p>
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Quick links</p>
                    <div className="mt-4 flex flex-wrap gap-3">
                        {creatorQuickLinks.map((item) => (
                            <Button key={item.href} asChild variant="secondary">
                                <Link href={item.href}>{item.label}</Link>
                            </Button>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderCreatorAnalyticsPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Creator analytics", description: "Track conversion, audience momentum, and creator performance across newest launches and livestream sessions.", badge: "Analytics" })}
            <div className="grid gap-4 md:grid-cols-3">
                {creatorAnalyticsStats.map((item) => (
                    <Card key={item.label} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.label}</p>
                        <p className="mt-3 text-2xl font-semibold text-zinc-900">{item.value}</p>
                    </Card>
                ))}
            </div>
            <div className="grid gap-4 lg:grid-cols-2">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Performance pulse</p>
                    <div className="mt-4 space-y-3">
                        {creatorPerformancePulse.map((item) => (
                            <div key={item.title} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                <p className="font-semibold text-zinc-900">{item.title}</p>
                                <p className="mt-1 text-zinc-600">{item.description}</p>
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Recommended actions</p>
                    <div className="mt-4 space-y-3">
                        {creatorRecommendedActions.map((item) => (
                            <div key={item.title} className="rounded-[24px] bg-zinc-50 px-4 py-3 text-sm text-zinc-700">
                                <p className="font-semibold text-zinc-900">{item.title}</p>
                                <p className="mt-1 text-zinc-600">{item.description}</p>
                            </div>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderCreatorProductsPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Creator products", description: "Manage the catalog experience from one place with live-ready items, launch checklists, and a clean storefront preview.", badge: "Creator catalog" })}
            <div className="grid gap-4 md:grid-cols-3">
                {[{ label: "Catalog items", value: String(sellerProducts.length) }, { label: "Live sellers", value: "12" }, { label: "Upcoming drop", value: "Today 7:00 PM" }].map((item) => (
                    <Card key={item.label} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.label}</p>
                        <p className="mt-3 text-2xl font-semibold text-zinc-900">{item.value}</p>
                    </Card>
                ))}
            </div>
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Launch-ready items</p>
                    <div className="mt-4 space-y-3">
                        {sellerProducts.slice(0, 4).map((item) => (
                            <div key={item.id} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                <div className="flex items-start justify-between gap-3">
                                    <div>
                                        <p className="font-semibold text-zinc-900">{item.name}</p>
                                        <p className="mt-1 text-xs text-zinc-500">{item.sku} • {item.stock} in stock</p>
                                    </div>
                                    <Badge tone="success">{item.status}</Badge>
                                </div>
                            </div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Creator checklist</p>
                    <div className="mt-4 space-y-3">
                        {creatorStudioFocus.map((item) => (
                            <div key={item} className="rounded-[24px] bg-zinc-50 px-4 py-3 text-sm text-zinc-700">{item}</div>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderCreatorLivestreamPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Creator livestream", description: "Build a live shopping experience with a clear schedule, audience insights, and an always-ready product pin strategy.", badge: "Live studio" })}
            <div className="grid gap-4 md:grid-cols-3">
                {[{ label: "Live viewers", value: "1.8k" }, { label: "Avg engagement", value: "8.9%" }, { label: "Next stream", value: "Tonight • 8PM" }].map((item) => (
                    <Card key={item.label} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.label}</p>
                        <p className="mt-3 text-2xl font-semibold text-zinc-900">{item.value}</p>
                    </Card>
                ))}
            </div>
            <div className="grid gap-4 lg:grid-cols-[1.1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Today’s run of show</p>
                    <div className="mt-4 space-y-3">
                        {["06:30 PM — Open live room and welcome viewers", "07:00 PM — Feature the top bundle and pins", "07:20 PM — Highlight customer questions and reviews"].map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">{item}</div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Audience pulse</p>
                    <div className="mt-4 space-y-3">
                        {["Most requested item: Signature Graphic Tee", "Asked about bundle pricing and shipping options", "Strong interest in creator-exclusive drops"].map((item) => (
                            <div key={item} className="rounded-[24px] bg-zinc-50 px-4 py-3 text-sm text-zinc-700">{item}</div>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderCreatorStudioPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Creator studio", description: "Coordinate your launch plan, product placement, and live moments from one polished creator workspace.", badge: "Creator studio" })}
            <div className="grid gap-4 md:grid-cols-3">
                {[{ label: "Active campaigns", value: "4" }, { label: "Live launches", value: "2" }, { label: "Follower lift", value: "+12%" }].map((item) => (
                    <Card key={item.label} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.label}</p>
                        <p className="mt-3 text-2xl font-semibold text-zinc-900">{item.value}</p>
                    </Card>
                ))}
            </div>
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Studio focus</p>
                    <div className="mt-4 space-y-3">
                        {["Review today’s performance and content priorities", "Pin the hero product for the next livestream", "Share the latest campaign recap with the team"].map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">{item}</div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Quick links</p>
                    <div className="mt-4 flex flex-wrap gap-3">
                        <Button asChild variant="secondary"><Link href="/creator/products">Products</Link></Button>
                        <Button asChild variant="secondary"><Link href="/creator/analytics">Analytics</Link></Button>
                        <Button asChild variant="secondary"><Link href="/creator/livestream">Livestream</Link></Button>
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderCreatorStudioOverviewPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Studio overview", description: "A deeper view into the behind-the-scenes launch plan, product readiness, and live campaign rhythm.", badge: "Creator studio overview" })}
            <div className="grid gap-4 md:grid-cols-3">
                {[{ label: "Active playbooks", value: "3" }, { label: "Ready-to-launch items", value: String(sellerProducts.length) }, { label: "Campaign confidence", value: "92%" }].map((item) => (
                    <Card key={item.label} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.label}</p>
                        <p className="mt-3 text-2xl font-semibold text-zinc-900">{item.value}</p>
                    </Card>
                ))}
            </div>
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">What’s live now</p>
                    <div className="mt-4 space-y-3">
                        {["Bundle spotlight: Spring drop set", "Live CTA pinned across the feed", "Team notes synced for the evening stream"].map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">{item}</div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Next steps</p>
                    <div className="mt-4 flex flex-wrap gap-3">
                        <Button asChild variant="primary"><Link href="/creator/livestream">Open livestream planner</Link></Button>
                        <Button asChild variant="secondary"><Link href="/creator/products">Review catalog</Link></Button>
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderCreatorStudioLaunchesPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Studio launches", description: "Preview launch-ready products, campaign sequencing, and the content plan that keeps each drop organized.", badge: "Creator studio launches" })}
            <div className="grid gap-4 lg:grid-cols-2">
                {sellerProducts.slice(0, 4).map((item) => (
                    <Card key={item.id} className="p-5">
                        <div className="flex items-start justify-between gap-3">
                            <div>
                                <p className="text-lg font-semibold text-zinc-900">{item.name}</p>
                                <p className="mt-1 text-sm text-zinc-600">{item.sku} • {item.stock} in stock</p>
                            </div>
                            <Badge tone="success">{item.status}</Badge>
                        </div>
                    </Card>
                ))}
            </div>
        </div>
    );
}

export function renderCreator(slug: string[]) {
    const route = slug[0] ?? "creator";
    const second = slug[1];
    const third = slug[2];

    if (route !== "creator") {
        return renderCreatorPage("Creator", "A polished creator workspace for analytics, products, livestreams, and payouts.");
    }

    if (second === "analytics") return renderCreatorAnalyticsPage();
    if (second === "products") return renderCreatorProductsPage();
    if (second === "livestream") return renderCreatorLivestreamPage();
    if (second === "studio") {
        if (third === "overview") return renderCreatorStudioOverviewPage();
        if (third === "launches") return renderCreatorStudioLaunchesPage();
        return renderCreatorStudioPage();
    }

    return renderCreatorPage("Creator", "A polished creator workspace for analytics, products, livestreams, and payouts.");
}
