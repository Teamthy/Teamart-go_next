"use client";

import Badge from "@/components/ui/badge";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import DataTable from "@/components/admin/DataTable";
import SectionHeader from "@/components/ui/SectionHeader";
import StatusChip from "@/components/admin/StatusChip";
import { useAdminDashboard, useAuditLogs, useDisputes, useFraudAlerts } from "@/hooks/useAdmin";

function formatDate(value?: string) {
    if (!value) return "—";

    return new Date(value).toLocaleString();
}

function severityTone(severity?: string) {
    const value = severity?.toLowerCase() || "default";

    if (value.includes("high") || value.includes("critical")) return "error";
    if (value.includes("medium")) return "warning";
    if (value.includes("low")) return "success";
    return "info";
}

function disputeTone(status?: string) {
    const value = status?.toLowerCase() || "open";

    if (value.includes("resolved") || value.includes("closed")) return "success";
    if (value.includes("pending") || value.includes("open")) return "warning";
    return "info";
}

export default function AdminPage() {
    const { dashboard, isLoading: dashboardLoading, error: dashboardError } = useAdminDashboard();
    const { disputes, isLoading: disputesLoading } = useDisputes();
    const { alerts, isLoading: alertsLoading } = useFraudAlerts();
    const { logs, isLoading: logsLoading } = useAuditLogs();

    const isLoading = dashboardLoading || disputesLoading || alertsLoading || logsLoading;

    const summaryCards = [
        {
            label: "Open disputes",
            value: dashboard?.open_disputes ?? 0,
            note: "Escalations waiting on action",
        },
        {
            label: "Pending payouts",
            value: dashboard?.pending_payouts ?? 0,
            note: "Scheduled for review",
        },
        {
            label: "Fraud alerts",
            value: dashboard?.fraud_alerts ?? 0,
            note: "Risk signals across live traffic",
        },
    ];

    const disputeRows = disputes.slice(0, 5).map((dispute) => ({
        id: dispute.id,
        status: dispute.status || "Open",
        createdAt: formatDate(dispute.created_at),
    }));

    const alertRows = alerts.slice(0, 5).map((alert) => ({
        id: alert.id,
        userId: alert.user_id,
        reason: alert.reason,
        severity: alert.severity || "Medium",
        createdAt: formatDate(alert.created_at),
    }));

    const auditRows = logs.slice(0, 8).map((log, index) => ({
        key: log.id || `${log.action}-${index}`,
        action: log.action || "Unknown action",
        createdAt: formatDate(log.created_at),
    }));

    return (
        <div className="space-y-8 pb-10">
            <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
                <div className="space-y-3">
                    <Badge tone="info">Platform operations</Badge>
                    <SectionHeader
                        title="Admin console"
                        description="Monitor disputes, fraud signals, and activity across the marketplace with a calmer, more actionable overview."
                    />
                </div>
                <div className="flex flex-wrap gap-3">
                    <Button asChild variant="secondary">
                        <a href="/dashboard">Open merchant view</a>
                    </Button>
                </div>
            </div>

            {dashboardError ? (
                <Card className="border-rose-200 bg-rose-50 p-4">
                    <p className="text-sm font-semibold text-rose-700">Error: {dashboardError}</p>
                </Card>
            ) : null}

            {isLoading ? (
                <Card className="p-8 text-center text-sm text-slate-500">
                    Loading admin overview…
                </Card>
            ) : (
                <>
                    <div className="grid gap-4 md:grid-cols-3">
                        {summaryCards.map((card) => (
                            <Card key={card.label} className="p-5">
                                <p className="text-xs uppercase tracking-[0.2em] text-slate-500">{card.label}</p>
                                <p className="mt-4 text-3xl font-semibold text-slate-900">{card.value}</p>
                                <p className="mt-2 text-sm text-slate-600">{card.note}</p>
                            </Card>
                        ))}
                    </div>

                    <div className="grid gap-6 xl:grid-cols-[1.1fr_0.9fr]">
                        <Card className="p-5">
                            <div className="mb-4 flex items-center justify-between gap-4">
                                <div>
                                    <p className="text-sm font-semibold text-slate-900">Recent disputes</p>
                                    <p className="text-sm text-slate-500">A quick view of the newest customer support cases.</p>
                                </div>
                            </div>
                            <DataTable
                                columns={[
                                    { header: "ID", accessor: "id" },
                                    {
                                        header: "Status",
                                        accessor: (row) => <StatusChip label={row.status} tone={disputeTone(row.status)} />,
                                    },
                                    { header: "Created", accessor: "createdAt" },
                                ]}
                                rows={disputeRows}
                            />
                        </Card>

                        <Card className="p-5">
                            <div className="mb-4 flex items-center justify-between gap-4">
                                <div>
                                    <p className="text-sm font-semibold text-slate-900">Fraud alerts</p>
                                    <p className="text-sm text-slate-500">High-signal events that need immediate review.</p>
                                </div>
                            </div>
                            <DataTable
                                columns={[
                                    { header: "Alert", accessor: "id" },
                                    { header: "User", accessor: "userId" },
                                    { header: "Reason", accessor: "reason" },
                                    {
                                        header: "Severity",
                                        accessor: (row) => <StatusChip label={row.severity} tone={severityTone(row.severity)} />,
                                    },
                                    { header: "Created", accessor: "createdAt" },
                                ]}
                                rows={alertRows}
                            />
                        </Card>
                    </div>

                    <Card className="p-5">
                        <div className="mb-4">
                            <p className="text-sm font-semibold text-slate-900">Audit trail</p>
                            <p className="text-sm text-slate-500">Recent platform activity and administrative actions.</p>
                        </div>
                        <DataTable
                            columns={[
                                { header: "Action", accessor: "action" },
                                { header: "Timestamp", accessor: "createdAt" },
                            ]}
                            rows={auditRows}
                        />
                    </Card>
                </>
            )}
        </div>
    );
}
