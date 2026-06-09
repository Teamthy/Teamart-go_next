"use client";

import SellerAppShell from "@/components/seller/SellerAppShell";
import PageHeader from "@/components/seller/PageHeader";
import MetricCard from "@/components/seller/MetricCard";
import { TrafficSourcePie } from "@/components/seller/Charts";

const creators = [
    { name: "Sarah Chen", revenue: 12450, orders: 234 },
    { name: "Emma Davis", revenue: 15678, orders: 189 },
    { name: "Mike Johnson", revenue: 8234, orders: 145 },
];

export default function AnalyticsCreatorsPage() {
    return (
        <SellerAppShell>
            <PageHeader
                title="Creator Analytics"
                description="Monitor creator performance and revenue contribution across livestreams."
                breadcrumbs={[
                    { label: "Dashboard", href: "/seller/dashboard" },
                    { label: "Analytics", href: "/seller/analytics" },
                    { label: "Creators" },
                ]}
            />

            <div className="grid gap-6 md:grid-cols-3 mb-6">
                <MetricCard
                    label="Top Creator"
                    value="Emma Davis"
                    subtext="Most revenue"
                    trend={{ value: 18, direction: "up", period: "over last period" }}
                />
                <MetricCard
                    label="Creator Revenue"
                    value="$38,762"
                    subtext="Total this month"
                    trend={{ value: 14, direction: "up", period: "vs previous" }}
                />
                <MetricCard
                    label="Active Creators"
                    value="12"
                    subtext="This month"
                    trend={{ value: 3, direction: "up", period: "growth" }}
                />
            </div>

            <TrafficSourcePie />

            <div className="mt-6 grid gap-4 lg:grid-cols-3">
                {creators.map((creator) => (
                    <div key={creator.name} className="rounded-3xl border border-[var(--border-light)] bg-[var(--bg-white)] p-5 shadow-sm shadow-slate-200/10">
                        <p className="text-sm text-[var(--text-secondary)] mb-2">{creator.name}</p>
                        <p className="text-2xl font-semibold text-[var(--text-primary)]">${creator.revenue.toLocaleString()}</p>
                        <p className="text-sm text-[var(--primary)] mt-2">{creator.orders} orders</p>
                    </div>
                ))}
            </div>
        </SellerAppShell>
    );
}
