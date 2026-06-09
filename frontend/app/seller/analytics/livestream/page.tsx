"use client";

import SellerAppShell from "@/components/seller/SellerAppShell";
import PageHeader from "@/components/seller/PageHeader";
import MetricCard from "@/components/seller/MetricCard";
import { LivestreamRevenueChart } from "@/components/seller/Charts";

const streamHighlights = [
    { title: "Holiday Flash Sale", viewers: 1874, revenue: 825.42 },
    { title: "Product Launch Event", viewers: 2341, revenue: 1245.99 },
    { title: "Weekend Showcase", viewers: 1560, revenue: 732.15 },
];

export default function AnalyticsLivestreamPage() {
    return (
        <SellerAppShell>
            <PageHeader
                title="Livestream Analytics"
                description="Review live commerce performance metrics and revenue per stream."
                breadcrumbs={[
                    { label: "Dashboard", href: "/seller/dashboard" },
                    { label: "Analytics", href: "/seller/analytics" },
                    { label: "Livestream" },
                ]}
            />

            <div className="grid gap-6 md:grid-cols-3 mb-6">
                <MetricCard
                    label="Live Revenue"
                    value="$39,420"
                    subtext="This month"
                    trend={{ value: 21, direction: "up", period: "vs previous" }}
                />
                <MetricCard
                    label="Avg. Viewers"
                    value="1,672"
                    subtext="Active stream average"
                    trend={{ value: 5, direction: "up", period: "week over week" }}
                />
                <MetricCard
                    label="Conversion Rate"
                    value="4.7%"
                    subtext="Live product showcases"
                    trend={{ value: 0.8, direction: "up", period: "vs last month" }}
                />
            </div>

            <LivestreamRevenueChart />

            <div className="mt-6 grid gap-4 lg:grid-cols-3">
                {streamHighlights.map((stream) => (
                    <div key={stream.title} className="rounded-3xl border border-[var(--border-light)] bg-[var(--bg-white)] p-5 shadow-sm shadow-slate-200/10">
                        <p className="text-sm text-[var(--text-secondary)] mb-2">{stream.title}</p>
                        <p className="text-2xl font-semibold text-[var(--text-primary)]">{stream.viewers} viewers</p>
                        <p className="text-sm text-[var(--primary)] mt-2">${stream.revenue.toFixed(2)} revenue</p>
                    </div>
                ))}
            </div>
        </SellerAppShell>
    );
}
