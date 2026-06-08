"use client";

import { useState } from "react";
import SellerAppShell from "@/components/seller/SellerAppShell";
import PageHeader from "@/components/seller/PageHeader";
import MetricCard from "@/components/seller/MetricCard";
import SellerButton from "@/components/seller/Button";
import { Download, Calendar } from "lucide-react";

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

      {/* Charts Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        {/* Revenue Trend */}
        <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-[var(--text-primary)]">
              Revenue Trend
            </h3>
            <span className="text-sm text-[var(--text-tertiary)]">
              {dateRangeOptions.find((d) => d.value === dateRange)?.label}
            </span>
          </div>
          <div className="h-80 bg-[var(--bg-lighter)] rounded-lg flex items-center justify-center text-[var(--text-tertiary)] border border-[var(--border-light)]">
            [Line Chart - Revenue by Day]
          </div>
        </div>

        {/* Orders Trend */}
        <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-[var(--text-primary)]">
              Orders Trend
            </h3>
            <span className="text-sm text-[var(--text-tertiary)]">
              {dateRangeOptions.find((d) => d.value === dateRange)?.label}
            </span>
          </div>
          <div className="h-80 bg-[var(--bg-lighter)] rounded-lg flex items-center justify-center text-[var(--text-tertiary)] border border-[var(--border-light)]">
            [Bar Chart - Orders by Day]
          </div>
        </div>
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

        {/* Traffic Sources */}
        <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
          <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-4">
            Traffic Sources
          </h3>
          <div className="space-y-4">
            {[
              { source: "Livestream", percentage: 45 },
              { source: "Search", percentage: 28 },
              { source: "Social", percentage: 18 },
              { source: "Direct", percentage: 9 },
            ].map((source, index) => (
              <div key={index}>
                <div className="flex items-center justify-between mb-1">
                  <p className="text-sm font-medium text-[var(--text-secondary)]">
                    {source.source}
                  </p>
                  <p className="text-sm font-semibold text-[var(--text-primary)]">
                    {source.percentage}%
                  </p>
                </div>
                <div className="h-2 bg-[var(--bg-lighter)] rounded-full overflow-hidden">
                  <div
                    className="h-full bg-[var(--primary)]"
                    style={{ width: `${source.percentage}%` }}
                  />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Livestream Analytics */}
      <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
        <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-4">
          Livestream Analytics
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div>
            <p className="text-sm text-[var(--text-tertiary)] mb-1">
              Total Live Hours
            </p>
            <p className="text-2xl font-bold text-[var(--text-primary)]">
              45.5h
            </p>
          </div>
          <div>
            <p className="text-sm text-[var(--text-tertiary)] mb-1">
              Average Viewers
            </p>
            <p className="text-2xl font-bold text-[var(--text-primary)]">
              1,234
            </p>
          </div>
          <div>
            <p className="text-sm text-[var(--text-tertiary)] mb-1">
              Livestream Revenue
            </p>
            <p className="text-2xl font-bold text-[var(--primary)]">
              $12,450
            </p>
          </div>
          <div>
            <p className="text-sm text-[var(--text-tertiary)] mb-1">
              Conversion Rate
            </p>
            <p className="text-2xl font-bold text-[var(--text-primary)]">
              4.2%
            </p>
          </div>
        </div>
        <div className="h-80 bg-[var(--bg-lighter)] rounded-lg flex items-center justify-center text-[var(--text-tertiary)] border border-[var(--border-light)]">
          [Area Chart - Livestream Revenue Trend]
        </div>
      </div>
    </SellerAppShell>
  );
}
