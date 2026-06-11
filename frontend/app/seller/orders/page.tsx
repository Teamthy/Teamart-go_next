"use client";

import { useEffect, useState } from "react";
import SellerAppShell from "@/components/seller/SellerAppShell";
import PageHeader from "@/components/seller/PageHeader";
import SellerButton from "@/components/seller/Button";
import DataTable from "@/components/seller/DataTable";
import Badge from "@/components/seller/Badge";
import { Download, Printer, Search } from "lucide-react";
import * as api from "@/lib/api";

interface Order {
    id: string;
    orderNumber: string;
    customer: string;
    amount: number;
    items: number;
    status:
    | "pending"
    | "paid"
    | "processing"
    | "packed"
    | "shipped"
    | "delivered"
    | "refunded"
    | "cancelled";
    date: string;
    paymentMethod: string;
}

const statusConfig = {
    pending: { variant: "warning", label: "Pending" },
    paid: { variant: "secondary", label: "Paid" },
    processing: { variant: "primary", label: "Processing" },
    packed: { variant: "primary", label: "Packed" },
    shipped: { variant: "primary", label: "Shipped" },
    delivered: { variant: "success", label: "Delivered" },
    refunded: { variant: "secondary", label: "Refunded" },
    cancelled: { variant: "danger", label: "Cancelled" },
} as const;

export default function OrdersPage() {
    const [searchQuery, setSearchQuery] = useState("");
    const [statusFilter, setStatusFilter] = useState("");
    const [selectedIds, setSelectedIds] = useState<(string | number)[]>([]);
    const [sortColumn, setSortColumn] = useState<string>("");
    const [sortDirection, setSortDirection] = useState<"asc" | "desc">("asc");
    const [orders, setOrders] = useState<Order[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        setIsLoading(true);
        setError(null);

        api.listOrders(100, 0)
            .then((response: any) => {
                const loadedOrders = Array.isArray(response?.orders)
                    ? response.orders
                    : Array.isArray(response)
                        ? response
                        : [];

                setOrders(
                    loadedOrders.map((order: any) => ({
                        id: String(order.id || order.order_number || order.orderId || ""),
                        orderNumber: order.order_number || order.orderNumber || order.orderId || String(order.id || ""),
                        customer: order.customer_name || order.customer || "Guest",
                        amount: Number(order.total_amount ?? order.totalAmount ?? order.amount ?? 0),
                        items: Number(order.items_count ?? order.items ?? order.quantity ?? 0),
                        status: (order.status || "pending") as Order["status"],
                        date: order.created_at || order.createdAt || "",
                        paymentMethod: order.payment_method || order.paymentMethod || "Unknown",
                    }))
                );
            })
            .catch((error: any) => {
                setError(error?.message || "Failed to load orders.");
            })
            .finally(() => setIsLoading(false));
    }, []);

    const filteredOrders = orders.filter((order) => {
        const matchesSearch =
            order.orderNumber.toLowerCase().includes(searchQuery.toLowerCase()) ||
            order.customer.toLowerCase().includes(searchQuery.toLowerCase());
        const matchesStatus = !statusFilter || order.status === statusFilter;
        return matchesSearch && matchesStatus;
    });

    const totalAmount = filteredOrders.reduce((sum, order) => sum + order.amount, 0);

    return (
        <SellerAppShell>
            <PageHeader
                title="Orders"
                description="Manage and track customer orders"
                breadcrumbs={[
                    { label: "Dashboard", href: "/seller/dashboard" },
                    { label: "Orders" },
                ]}
                actions={
                    <div className="flex gap-3">
                        <SellerButton
                            variant="secondary"
                            size="md"
                            icon={<Download className="w-4 h-4" />}
                        >
                            Export CSV
                        </SellerButton>
                        <SellerButton
                            variant="secondary"
                            size="md"
                            icon={<Printer className="w-4 h-4" />}
                        >
                            Print Labels
                        </SellerButton>
                    </div>
                }
            />

            {/* Filters */}
            <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6 mb-6">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {/* Search */}
                    <div className="md:col-span-2">
                        <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                            Search Orders
                        </label>
                        <div className="relative">
                            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--text-tertiary)]" />
                            <input
                                type="text"
                                placeholder="Search by order ID or customer..."
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                className="w-full pl-10 pr-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                            />
                        </div>
                    </div>

                    {/* Status Filter */}
                    <div>
                        <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                            Status
                        </label>
                        <select
                            value={statusFilter}
                            onChange={(e) => setStatusFilter(e.target.value)}
                            className="w-full px-3 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                        >
                            <option value="">All Status</option>
                            <option value="pending">Pending</option>
                            <option value="paid">Paid</option>
                            <option value="processing">Processing</option>
                            <option value="packed">Packed</option>
                            <option value="shipped">Shipped</option>
                            <option value="delivered">Delivered</option>
                            <option value="refunded">Refunded</option>
                            <option value="cancelled">Cancelled</option>
                        </select>
                    </div>
                </div>

                {/* Summary & Actions */}
                <div className="flex items-center justify-between pt-4 border-t border-[var(--border-light)] mt-4">
                    <div>
                        <p className="text-sm text-[var(--text-secondary)]">
                            {filteredOrders.length} orders found • Total: ${totalAmount.toFixed(2)}
                        </p>
                    </div>
                    {selectedIds.length > 0 && (
                        <div className="flex gap-2">
                            <SellerButton variant="secondary" size="sm">
                                Mark as Shipped
                            </SellerButton>
                            <SellerButton variant="danger" size="sm">
                                Cancel Orders
                            </SellerButton>
                        </div>
                    )}
                </div>
            </div>

            {/* Orders Table */}
            <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] overflow-hidden">
                <div className="p-6 border-b border-[var(--border-light)]">
                    <h3 className="text-lg font-semibold text-[var(--text-primary)]">
                        All Orders
                    </h3>
                </div>
                <DataTable
                    columns={[
                        {
                            key: "orderNumber",
                            label: "Order ID",
                            width: "120px",
                            sortable: true,
                        },
                        {
                            key: "customer",
                            label: "Customer",
                            sortable: true,
                            render: (value, row: Order) => (
                                <div>
                                    <p className="font-medium text-[var(--text-primary)]">
                                        {value}
                                    </p>
                                    <p className="text-xs text-[var(--text-tertiary)]">
                                        {row.paymentMethod}
                                    </p>
                                </div>
                            ),
                        },
                        {
                            key: "items",
                            label: "Items",
                            width: "80px",
                            sortable: true,
                        },
                        {
                            key: "amount",
                            label: "Amount",
                            width: "100px",
                            render: (value) => (
                                <p className="font-semibold text-[var(--text-primary)]">
                                    ${value.toFixed(2)}
                                </p>
                            ),
                            sortable: true,
                        },
                        {
                            key: "status",
                            label: "Status",
                            width: "140px",
                            render: (value: Order["status"]) => (
                                <Badge variant={statusConfig[value].variant}>
                                    {statusConfig[value].label}
                                </Badge>
                            ),
                        },
                        {
                            key: "date",
                            label: "Date",
                            width: "120px",
                            sortable: true,
                        },
                        {
                            key: "id",
                            label: "Action",
                            width: "100px",
                            render: () => (
                                <SellerButton variant="ghost" size="sm">
                                    View
                                </SellerButton>
                            ),
                        },
                    ]}
                    data={filteredOrders}
                    selectable={true}
                    selectedIds={selectedIds}
                    onSelectChange={setSelectedIds}
                    sortColumn={sortColumn}
                    sortDirection={sortDirection}
                    onSort={handleSort}
                    emptyMessage="No orders found. Try adjusting your filters."
                />
            </div>
        </SellerAppShell>
    );
}
