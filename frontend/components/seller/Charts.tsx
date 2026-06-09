"use client";

import {
    Area,
    AreaChart,
    Bar,
    BarChart,
    CartesianGrid,
    Cell,
    Legend,
    Line,
    LineChart,
    Pie,
    PieChart,
    ResponsiveContainer,
    Tooltip,
    XAxis,
    YAxis,
} from "recharts";
import type { ReactNode } from "react";

interface ChartCardProps {
    title: string;
    subtitle?: string;
    children: ReactNode;
}

export function ChartCard({ title, subtitle, children }: ChartCardProps) {
    return (
        <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
            <div className="flex items-start justify-between gap-4 mb-4">
                <div>
                    <h3 className="text-lg font-semibold text-[var(--text-primary)]">{title}</h3>
                    {subtitle && <p className="text-sm text-[var(--text-tertiary)] mt-1">{subtitle}</p>}
                </div>
            </div>
            <div className="h-64 w-full">{children}</div>
        </div>
    );
}

const revenueData = [
    { date: "Feb 1", revenue: 4200, orders: 38 },
    { date: "Feb 2", revenue: 5200, orders: 45 },
    { date: "Feb 3", revenue: 4800, orders: 41 },
    { date: "Feb 4", revenue: 6100, orders: 52 },
    { date: "Feb 5", revenue: 7300, orders: 61 },
    { date: "Feb 6", revenue: 6800, orders: 57 },
    { date: "Feb 7", revenue: 7600, orders: 64 },
];

const orderData = [
    { date: "Feb 1", processing: 18, shipped: 12, delivered: 8 },
    { date: "Feb 2", processing: 16, shipped: 14, delivered: 15 },
    { date: "Feb 3", processing: 20, shipped: 18, delivered: 12 },
    { date: "Feb 4", processing: 14, shipped: 20, delivered: 18 },
    { date: "Feb 5", processing: 10, shipped: 23, delivered: 28 },
    { date: "Feb 6", processing: 16, shipped: 19, delivered: 22 },
    { date: "Feb 7", processing: 18, shipped: 21, delivered: 25 },
];

const trafficSources = [
    { name: "Livestream", value: 45 },
    { name: "Search", value: 28 },
    { name: "Social", value: 18 },
    { name: "Direct", value: 9 },
];

const COLORS = ["#E91E63", "#F59E0B", "#22C55E", "#6B7280"];

export function RevenueTrendChart() {
    return (
        <ChartCard title="Revenue Trend" subtitle="Last 7 days">
            <ResponsiveContainer width="100%" height="100%">
                <LineChart data={revenueData} margin={{ top: 10, right: 16, left: 0, bottom: 0 }}>
                    <CartesianGrid strokeDasharray="3 3" stroke="var(--border-light)" />
                    <XAxis dataKey="date" tickLine={false} axisLine={false} tick={{ fill: "var(--text-tertiary)", fontSize: 12 }} />
                    <YAxis tickLine={false} axisLine={false} tick={{ fill: "var(--text-tertiary)", fontSize: 12 }} />
                    <Tooltip contentStyle={{ borderRadius: 12, borderColor: "var(--border-light)" }} />
                    <Line type="monotone" dataKey="revenue" stroke="#E91E63" strokeWidth={3} dot={{ r: 3 }} activeDot={{ r: 6 }} />
                </LineChart>
            </ResponsiveContainer>
        </ChartCard>
    );
}

export function OrdersTrendChart() {
    return (
        <ChartCard title="Orders Trend" subtitle="Last 7 days">
            <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={orderData} margin={{ top: 10, right: 16, left: 0, bottom: 0 }}>
                    <defs>
                        <linearGradient id="colorOrders" x1="0" y1="0" x2="0" y2="1">
                            <stop offset="5%" stopColor="#E91E63" stopOpacity={0.35} />
                            <stop offset="95%" stopColor="#E91E63" stopOpacity={0} />
                        </linearGradient>
                    </defs>
                    <CartesianGrid strokeDasharray="3 3" stroke="var(--border-light)" />
                    <XAxis dataKey="date" tickLine={false} axisLine={false} tick={{ fill: "var(--text-tertiary)", fontSize: 12 }} />
                    <YAxis tickLine={false} axisLine={false} tick={{ fill: "var(--text-tertiary)", fontSize: 12 }} />
                    <Tooltip contentStyle={{ borderRadius: 12, borderColor: "var(--border-light)" }} />
                    <Area type="monotone" dataKey="delivered" stroke="#E91E63" strokeWidth={2} fill="url(#colorOrders)" />
                </AreaChart>
            </ResponsiveContainer>
        </ChartCard>
    );
}

export function TrafficSourcePie() {
    return (
        <ChartCard title="Traffic Sources" subtitle="Performance by channel">
            <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                    <Pie data={trafficSources} dataKey="value" nameKey="name" innerRadius={52} outerRadius={84} paddingAngle={4}>
                        {trafficSources.map((entry, index) => (
                            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                        ))}
                    </Pie>
                    <Tooltip contentStyle={{ borderRadius: 12, borderColor: "var(--border-light)" }} />
                    <Legend verticalAlign="bottom" height={34} iconType="circle" wrapperStyle={{ fontSize: 12, color: "var(--text-secondary)" }} />
                </PieChart>
            </ResponsiveContainer>
        </ChartCard>
    );
}

export function LivestreamRevenueChart() {
    return (
        <ChartCard title="Livestream Revenue" subtitle="Daily performance">
            <ResponsiveContainer width="100%" height="100%">
                <BarChart data={revenueData} margin={{ top: 10, right: 16, left: 0, bottom: 0 }}>
                    <CartesianGrid strokeDasharray="3 3" stroke="var(--border-light)" />
                    <XAxis dataKey="date" tickLine={false} axisLine={false} tick={{ fill: "var(--text-tertiary)", fontSize: 12 }} />
                    <YAxis tickLine={false} axisLine={false} tick={{ fill: "var(--text-tertiary)", fontSize: 12 }} />
                    <Tooltip contentStyle={{ borderRadius: 12, borderColor: "var(--border-light)" }} />
                    <Bar dataKey="revenue" fill="#E91E63" radius={[10, 10, 0, 0]} />
                </BarChart>
            </ResponsiveContainer>
        </ChartCard>
    );
}
