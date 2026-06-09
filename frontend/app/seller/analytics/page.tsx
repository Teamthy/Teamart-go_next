"use client";

import { useState } from "react";
import SellerAppShell from "@/components/seller/SellerAppShell";
import PageHeader from "@/components/seller/PageHeader";
import MetricCard from "@/components/seller/MetricCard";
import SellerButton from "@/components/seller/Button";
import { Download } from "lucide-react";
import { RevenueTrendChart, OrdersTrendChart, TrafficSourcePie, LivestreamRevenueChart } from "@/components/seller/Charts";

export default function AnalyticsPage() {
    const [dateRange, setDateRange] = useState("last_30_days");

    const dateRangeOptions = [
        { value: "last_7_days", label: "Last 7 Days" },
        { value: "last_30_days", label: "Last 30 Days" },
        { value: "last_90_days", label: "Last 90 Days" },
        { value: "this_year", label: "This Year" },
        { value: "custom", label: "Custom Range" },
    ];

    return (
        <SellerAppShell>
            <PageHeader
                title="Analytics"
                description="Deep dive into your store performance metrics"
                breadcrumbs={[
                    { label: "Dashboard", href: "/seller/dashboard" },
                    { label: "Analytics" },
                ]}
                actions={
                    <div className="flex gap-3">
                        <select
                            value={dateRange}
                            onChange={(e) => setDateRange(e.target.value)}
                            className="px-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                        >
                            {dateRangeOptions.map((option) => (
                                <option key={option.value} value={option.value}>
                                    {option.label}
                                </option>
                            ))}
                        </select>
                        <SellerButton
                            variant="secondary"
                            size="md"
                            icon={<Download className="w-4 h-4" />}
                        >
                            Export Report
                        </SellerButton>
                    </div>
                }
            />

            {/* Key Metrics */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                <MetricCard
                    label="Total Revenue"
                    value="$45,234"
                    subtext="Last 30 days"
                    trend={{ value: 18, direction: "up", period: "previous period" }}
                    icon={<span className="text-xl">💰</span>}
                />
                <MetricCard
                    label="Orders"
                    value="1,240"
                    subtext="Last 30 days"
                    trend={{ value: 12, direction: "up", period: "previous period" }}
                    icon={<span className="text-xl">📦</span>}
                />
                <MetricCard
                    label="Avg. Order Value"
                    value="$36.48"
                    subtext="Last 30 days"
                    trend={{ value: 5, direction: "up", period: "previous period" }}
                    icon={<span className="text-xl">💵</span>}
                />
                <MetricCard
                    label="Conversion Rate"
                    value="3.45%"
                    subtext="Last 30 days"
                    trend={{ value: 2, direction: "down", period: "previous period" }}
                    icon={<span className="text-xl">📈</span>}
                />
            </div>

            <div className="mb-8">
                <h2 className="text-xl font-semibold text-[var(--text-primary)] mb-4">Explore analytics</h2>
                <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
                    {[
                        { label: "Revenue", href: "/seller/analytics/revenue", description: "View revenue performance and channel trends." },
                        { label: "Products", href: "/seller/analytics/products", description: "See your best sellers and inventory impact." },
                        { label: "Creators", href: "/seller/analytics/creators", description: "Track creator-driven live commerce results." },
                        { label: "Livestream", href: "/seller/analytics/livestream", description: "Analyze live commerce revenue and engagement." },
                    ].map((item) => (
                        <a
                            key={item.label}
                            href={item.href}
                            className="block rounded-3xl border border-[var(--border-light)] bg-[var(--bg-white)] p-6 text-[var(--text-primary)] hover:border-[var(--primary)] hover:bg-[var(--bg-muted)] transition"
                        >
                            <p className="text-sm text-[var(--text-tertiary)]">{item.label}</p>
                            <p className="mt-3 text-base font-semibold">{item.description}</p>
                        </a>
                    ))}
                </div>
            </div>

            {/* Charts Grid */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
                {/* Revenue Trend */}
                <RevenueTrendChart />
                <OrdersTrendChart />
            </div>

            {/* More Analytics Sections */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
                {/* Top Products */}
                <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
                    <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-4">
                        Top Products
                    </h3>
                    <div className="space-y-4">
                        {[
                            { name: "Wireless Headphones", sales: 342, revenue: 68458 },
                            { name: "Smart Watch Elite", sales: 156, revenue: 54584 },
                            { name: "USB-C Hub", sales: 412, revenue: 20588 },
                        ].map((product, index) => (
                            <div key={index} className="p-3 bg-[var(--bg-lighter)] rounded-lg">
                                <div className="flex items-center justify-between mb-2">
                                    <p className="font-medium text-[var(--text-primary)] text-sm">
                                        {product.name}
                                    </p>
                                    <p className="font-semibold text-[var(--primary)] text-sm">
                                        ${(product.revenue / 100).toFixed(0)}
                                    </p>
                                </div>
                                <p className="text-xs text-[var(--text-tertiary)]">
                                    {product.sales} sales
                                </p>
                            </div>
                        ))}
                    </div>
                </div>

                {/* Top Creators */}
                <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
                    <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-4">
                        Top Creators
                    </h3>
                    <div className="space-y-4">
                        {[
                            { name: "Sarah Chen", sales: 234, revenue: 12450 },
                            { name: "Emma Davis", sales: 189, revenue: 15678 },
                            { name: "Mike Johnson", sales: 145, revenue: 8234 },
                        ].map((creator, index) => (
                            <div key={index} className="p-3 bg-[var(--bg-lighter)] rounded-lg">
                                <div className="flex items-center justify-between mb-2">
                                    <p className="font-medium text-[var(--text-primary)] text-sm">
                                        {creator.name}
                                    </p>
                                    <p className="font-semibold text-[var(--primary)] text-sm">
                                        ${(creator.revenue / 100).toFixed(0)}
                                    </p>
                                </div>
                                <p className="text-xs text-[var(--text-tertiary)]">
                                    {creator.sales} orders
                                </p>
                            </div>
                        ))}
                    </div>
                </div>

                <TrafficSourcePie />
            </div>

            <LivestreamRevenueChart />
        </SellerAppShell>
    );
}
