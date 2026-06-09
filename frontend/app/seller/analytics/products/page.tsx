"use client";

import SellerAppShell from "@/components/seller/SellerAppShell";
import PageHeader from "@/components/seller/PageHeader";
import MetricCard from "@/components/seller/MetricCard";
import { OrdersTrendChart } from "@/components/seller/Charts";

const topProducts = [
    { name: "Wireless Headphones Pro", sold: 342, revenue: 68458 },
    { name: "Smart Watch Elite", sold: 189, revenue: 54584 },
    { name: "USB-C Hub Universal", sold: 156, revenue: 32900 },
];

export default function AnalyticsProductsPage() {
    return (
        <SellerAppShell>
            <PageHeader
                title="Product Analytics"
                description="Understand which products are driving sales and where inventory should be focused."
                breadcrumbs={[
                    { label: "Dashboard", href: "/seller/dashboard" },
                    { label: "Analytics", href: "/seller/analytics" },
                    { label: "Products" },
                ]}
            />

            <div className="grid gap-6 md:grid-cols-3 mb-6">
                <MetricCard
                    label="Top Product"
                    value="Wireless Headphones Pro"
                    subtext="Best seller"
                    trend={{ value: 22, direction: "up", period: "sales growth" }}
                />
                <MetricCard
                    label="Out of Stock"
                    value="4"
                    subtext="SKUs needing restock"
                    trend={{ value: 8, direction: "down", period: "week over week" }}
                />
                <MetricCard
                    label="Return Rate"
                    value="1.9%"
                    subtext="Order returns"
                    trend={{ value: 0.3, direction: "down", period: "vs previous" }}
                />
            </div>

            <OrdersTrendChart />

            <div className="mt-6 grid gap-4 lg:grid-cols-3">
                {topProducts.map((product) => (
                    <div key={product.name} className="rounded-3xl border border-[var(--border-light)] bg-[var(--bg-white)] p-5 shadow-sm shadow-slate-200/10">
                        <p className="text-sm text-[var(--text-secondary)] mb-2">{product.name}</p>
                        <p className="text-2xl font-semibold text-[var(--text-primary)]">{product.sold} sold</p>
                        <p className="text-sm text-[var(--primary)] mt-2">${(product.revenue / 100).toFixed(0)} revenue</p>
                    </div>
                ))}
            </div>
        </SellerAppShell>
    );
}
