"use client";

import SellerAppShell from "@/components/seller/SellerAppShell";
import PageHeader from "@/components/seller/PageHeader";
import MetricCard from "@/components/seller/MetricCard";
import SellerButton from "@/components/seller/Button";
import DataTable from "@/components/seller/DataTable";
import Badge from "@/components/seller/Badge";
import {
  DollarSign,
  ShoppingCart,
  Box,
  TrendingUp,
  Clock,
  AlertCircle,
} from "lucide-react";
import { useState } from "react";

interface Order {
  id: string;
  orderNumber: string;
  customer: string;
  amount: number;
  status: "pending" | "paid" | "processing" | "shipped" | "delivered";
  date: string;
}

const mockOrders: Order[] = [
  {
    id: "1",
    orderNumber: "ORD-001",
    customer: "John Doe",
    amount: 299.99,
    status: "delivered",
    date: "2024-01-15",
  },
  {
    id: "2",
    orderNumber: "ORD-002",
    customer: "Jane Smith",
    amount: 149.99,
    status: "shipped",
    date: "2024-01-14",
  },
  {
    id: "3",
    orderNumber: "ORD-003",
    customer: "Bob Johnson",
    amount: 499.99,
    status: "processing",
    date: "2024-01-13",
  },
  {
    id: "4",
    orderNumber: "ORD-004",
    customer: "Alice Williams",
    amount: 89.99,
    status: "paid",
    date: "2024-01-12",
  },
  {
    id: "5",
    orderNumber: "ORD-005",
    customer: "Charlie Brown",
    amount: 599.99,
    status: "pending",
    date: "2024-01-11",
  },
];

const statusConfig = {
  pending: { variant: "warning", label: "Pending" },
  paid: { variant: "secondary", label: "Paid" },
  processing: { variant: "primary", label: "Processing" },
  shipped: { variant: "primary", label: "Shipped" },
  delivered: { variant: "success", label: "Delivered" },
} as const;

export default function DashboardPage() {
  const [sortColumn, setSortColumn] = useState<string>("");
  const [sortDirection, setSortDirection] = useState<"asc" | "desc">("asc");

  const handleSort = (column: string, direction: "asc" | "desc") => {
    setSortColumn(column);
    setSortDirection(direction);
  };

  return (
    <SellerAppShell>
      <PageHeader
        title="Dashboard"
        description="Monitor your store performance in real-time"
        actions={
          <SellerButton variant="primary">
            View Full Analytics
          </SellerButton>
        }
      />

      {/* KPI Row */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-6 gap-6 mb-8">
        <MetricCard
          label="Revenue"
          value="$12,345"
          subtext="Total"
          trend={{ value: 12, direction: "up", period: "last month" }}
          icon={<DollarSign className="w-5 h-5" />}
        />
        <MetricCard
          label="Orders"
          value="342"
          subtext="This month"
          trend={{ value: 8, direction: "up", period: "last month" }}
          icon={<ShoppingCart className="w-5 h-5" />}
        />
        <MetricCard
          label="Products Sold"
          value="1,204"
          subtext="This month"
          trend={{ value: 15, direction: "up", period: "last month" }}
          icon={<Box className="w-5 h-5" />}
        />
        <MetricCard
          label="Conversion Rate"
          value="3.2%"
          subtext="This month"
          trend={{ value: 2, direction: "down", period: "last month" }}
          icon={<TrendingUp className="w-5 h-5" />}
        />
        <MetricCard
          label="Avg Order Value"
          value="$89.50"
          subtext="Per order"
          trend={{ value: 5, direction: "up", period: "last month" }}
          icon={<DollarSign className="w-5 h-5" />}
        />
        <MetricCard
          label="Livestream Revenue"
          value="$3,450"
          subtext="This month"
          trend={{ value: 25, direction: "up", period: "last month" }}
          icon={<AlertCircle className="w-5 h-5" />}
        />
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-8">
        {/* Left Column - Charts */}
        <div className="lg:col-span-2 space-y-6">
          {/* Revenue Trend Chart */}
          <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
            <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-4">
              Revenue Trend
            </h3>
            <div className="h-64 bg-[var(--bg-lighter)] rounded-lg flex items-center justify-center text-[var(--text-tertiary)]">
              [Chart Placeholder - Revenue by Day]
            </div>
          </div>

          {/* Orders Trend Chart */}
          <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
            <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-4">
              Orders Trend
            </h3>
            <div className="h-64 bg-[var(--bg-lighter)] rounded-lg flex items-center justify-center text-[var(--text-tertiary)]">
              [Chart Placeholder - Orders by Day]
            </div>
          </div>
        </div>

        {/* Right Column - Widgets */}
        <div className="space-y-6">
          {/* Pending Orders Widget */}
          <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
            <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-4 flex items-center gap-2">
              <Clock className="w-5 h-5 text-[var(--warning)]" />
              Pending Orders
            </h3>
            <div className="space-y-3">
              <div className="flex items-center justify-between p-3 bg-[var(--bg-lighter)] rounded-lg">
                <span className="text-sm text-[var(--text-secondary)]">
                  ORD-005
                </span>
                <span className="font-semibold text-[var(--text-primary)]">
                  $599.99
                </span>
              </div>
              <div className="text-center">
                <SellerButton variant="secondary" size="sm">
                  View All Orders
                </SellerButton>
              </div>
            </div>
          </div>

          {/* Inventory Alerts Widget */}
          <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
            <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-4 flex items-center gap-2">
              <AlertCircle className="w-5 h-5 text-[var(--danger)]" />
              Inventory Alerts
            </h3>
            <div className="space-y-3">
              <div className="p-3 bg-[var(--danger-lighter)] rounded-lg">
                <p className="text-sm font-medium text-[var(--danger)] mb-1">
                  Product XYZ
                </p>
                <p className="text-xs text-[var(--danger)]">
                  Only 5 items left in stock
                </p>
              </div>
              <div className="text-center">
                <SellerButton variant="secondary" size="sm">
                  View Inventory
                </SellerButton>
              </div>
            </div>
          </div>

          {/* Livestream Schedule Widget */}
          <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
            <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-4">
              Upcoming Livestreams
            </h3>
            <div className="space-y-3">
              <div className="p-3 bg-[var(--bg-lighter)] rounded-lg">
                <p className="text-sm font-medium text-[var(--text-primary)]">
                  Product Launch
                </p>
                <p className="text-xs text-[var(--text-tertiary)]">
                  Tomorrow at 7:00 PM
                </p>
              </div>
              <div className="text-center">
                <SellerButton variant="secondary" size="sm">
                  Schedule Stream
                </SellerButton>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Recent Orders Table */}
      <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
        <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-4">
          Recent Orders
        </h3>
        <DataTable
          columns={[
            { key: "orderNumber", label: "Order ID", width: "120px" },
            { key: "customer", label: "Customer", sortable: true },
            {
              key: "amount",
              label: "Amount",
              render: (value) => `$${value.toFixed(2)}`,
              sortable: true,
            },
            {
              key: "status",
              label: "Status",
              render: (value) => (
                <Badge variant={statusConfig[value].variant as any}>
                  {statusConfig[value].label}
                </Badge>
              ),
            },
            { key: "date", label: "Date", sortable: true },
            {
              key: "id",
              label: "Action",
              render: () => (
                <SellerButton variant="ghost" size="sm">
                  View
                </SellerButton>
              ),
            },
          ]}
          data={mockOrders}
          sortColumn={sortColumn}
          sortDirection={sortDirection}
          onSort={handleSort}
          emptyMessage="No orders found"
        />
      </div>
    </SellerAppShell>
  );
}
