"use client";

import { useMemo, useState, useEffect } from "react";
import SellerAppShell from "@/components/seller/SellerAppShell";
import PageHeader from "@/components/seller/PageHeader";
import SellerButton from "@/components/seller/Button";
import DataTable from "@/components/seller/DataTable";
import Badge from "@/components/seller/Badge";
import { Download, Printer, Search } from "lucide-react";

interface Order {
    id: string;
    orderNumber: string;
    customer: string;
    amount: number;
    items: number;
    status: "pending" | "paid" | "processing" | "packed" | "shipped" | "delivered" | "refunded" | "cancelled";
    date: string;
    paymentMethod: string;
}

import { fetchPendingOrders, updateOrderStatus } from "@/lib/sellerApi";
import { useToast } from "@/components/seller/ToastProvider";

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

export default function PendingOrdersPage() {
    const [searchQuery, setSearchQuery] = useState("");
    const [sortColumn, setSortColumn] = useState<string>("");
    const [sortDirection, setSortDirection] = useState<"asc" | "desc">("asc");

    const handleSort = (column: string, direction: "asc" | "desc") => {
        setSortColumn(column);
        setSortDirection(direction);
    };

    const [orders, setOrders] = useState<Order[]>([]);
    const { toast } = useToast();

    useEffect(() => {
        let mounted = true;
        fetchPendingOrders().then((data) => {
            if (mounted) setOrders(data);
        });
        return () => {
            mounted = false;
        };
    }, []);

    const pendingOrders = useMemo(
        () =>
            orders.filter(
                (order) =>
                    order.status === "pending" &&
                    (order.orderNumber.toLowerCase().includes(searchQuery.toLowerCase()) ||
                        order.customer.toLowerCase().includes(searchQuery.toLowerCase()))
            ),
        [searchQuery, orders]
    );

    return (
        <SellerAppShell>
            <PageHeader
                title="Pending Orders"
                description="Review and take action on orders waiting for fulfillment."
                breadcrumbs={[
                    { label: "Dashboard", href: "/seller/dashboard" },
                    { label: "Orders", href: "/seller/orders" },
                    { label: "Pending" },
                ]}
                actions={
                    <div className="flex gap-3">
                        <SellerButton variant="secondary" size="md" icon={<Download className="w-4 h-4" />}>
                            Export CSV
                        </SellerButton>
                        <SellerButton variant="secondary" size="md" icon={<Printer className="w-4 h-4" />}>
                            Print Labels
                        </SellerButton>
                    </div>
                }
            />

            <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6 mb-6">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div className="md:col-span-2">
                        <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                            Search Pending Orders
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
                    <div className="flex items-end justify-end">
                        <div className="text-right">
                            <p className="text-sm text-[var(--text-secondary)]">Open pending orders</p>
                            <p className="text-2xl font-semibold text-[var(--text-primary)]">{pendingOrders.length}</p>
                        </div>
                    </div>
                </div>
            </div>

            <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] overflow-hidden">
                <div className="p-6 border-b border-[var(--border-light)]">
                    <h3 className="text-lg font-semibold text-[var(--text-primary)]">Pending Orders</h3>
                </div>
                <DataTable
                    columns={[
                        { key: "orderNumber", label: "Order ID", width: "120px", sortable: true },
                        { key: "customer", label: "Customer", sortable: true },
                        { key: "items", label: "Items", width: "80px", sortable: true },
                        {
                            key: "amount",
                            label: "Amount",
                            render: (value) => `$${value.toFixed(2)}`,
                            sortable: true,
                            width: "120px",
                        },
                        {
                            key: "status",
                            label: "Status",
                            width: "140px",
                            render: (value) => {
                                const config = statusConfig[value as keyof typeof statusConfig];
                                return <Badge variant={config.variant}>{config.label}</Badge>;
                            },
                        },
                        { key: "date", label: "Date", sortable: true, width: "140px" },
                        {
                            key: "id",
                            label: "Action",
                            width: "140px",
                            render: (value, row: Order) => (
                                <div className="flex gap-2">
                                    <SellerButton
                                        variant="secondary"
                                        size="sm"
                                        onClick={async () => {
                                            // Optimistic update: remove from UI immediately
                                            setOrders((prev) => prev.filter((o) => o.id !== row.id));
                                            try {
                                                await updateOrderStatus(row.id, "processing");
                                                toast({ title: "Order updated", description: `${row.orderNumber} is now processing.`, type: "success" });
                                            } catch (e) {
                                                // Rollback on error
                                                setOrders((prev) => [row, ...prev]);
                                                toast({ title: "Failed", description: `Could not update ${row.orderNumber}.`, type: "error" });
                                            }
                                        }}
                                    >
                                        Pack
                                    </SellerButton>
                                    <SellerButton variant="ghost" size="sm">
                                        View
                                    </SellerButton>
                                </div>
                            ),
                        },
                    ]}
                    data={pendingOrders}
                    sortColumn={sortColumn}
                    sortDirection={sortDirection}
                    onSort={handleSort}
                    emptyMessage="No pending orders found"
                />
            </div>
        </SellerAppShell>
    );
}
