"use client";

import SellerAppShell from "@/components/seller/SellerAppShell";
import PageHeader from "@/components/seller/PageHeader";
import MetricCard from "@/components/seller/MetricCard";
import { RevenueTrendChart } from "@/components/seller/Charts";

export default function AnalyticsRevenuePage() {
    return (
        <SellerAppShell>
            <PageHeader
                title="Revenue Analytics"
                description="Track revenue growth, conversion performance, and channel contributions."
                breadcrumbs={[
                    { label: "Dashboard", href: "/seller/dashboard" },
                    { label: "Analytics", href: "/seller/analytics" },
                    { label: "Revenue" },
                ]}
            />

            <div className="grid gap-6 md:grid-cols-3 mb-6">
                <MetricCard
                    label="Total Revenue"
                    value="$47,200"
                    subtext="Last 30 days"
                    trend={{ value: 14, direction: "up", period: "vs previous" }}
                />
                <MetricCard
                    label="Average Order"
                    value="$39.50"
                    subtext="Revenue per order"
                    trend={{ value: 6, direction: "up", period: "vs previous" }}
                />
                <MetricCard
                    label="Checkout Conversion"
                    value="3.7%"
                    subtext="Cart to payment"
                    trend={{ value: 1.2, direction: "up", period: "vs previous" }}
                />
            </div>

            <RevenueTrendChart />

            <div className="mt-6 rounded-3xl border border-[var(--border-light)] bg-[var(--bg-white)] p-6 shadow-sm shadow-slate-200/10">
                <h2 className="text-lg font-semibold text-[var(--text-primary)] mb-3">Revenue by Channel</h2>
                <p className="text-sm text-[var(--text-secondary)]">View the impact of livestream, search, social, and direct traffic on total seller revenue.</p>
            </div>
        </SellerAppShell>
    );
}
