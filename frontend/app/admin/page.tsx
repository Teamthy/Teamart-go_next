"use client";

import { useAdminDashboard, useDisputes, useFraudAlerts, useAuditLogs } from "@/hooks/useAdmin";
import SectionHeader from "@/components/ui/SectionHeader";

export default function AdminPage() {
    const { dashboard, isLoading: dashboardLoading, error: dashboardError } = useAdminDashboard();
    const { disputes, isLoading: disputesLoading } = useDisputes();
    const { alerts, isLoading: alertsLoading } = useFraudAlerts();
    const { logs, isLoading: logsLoading } = useAuditLogs();

    const isLoading = dashboardLoading || disputesLoading || alertsLoading || logsLoading;

    return (
        <div className="space-y-8">
            <SectionHeader
                title="Admin Console"
                description="Manage global platform settings, reports, and admin workflows."
            />

            {dashboardError && (
                <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
                    <p className="text-red-800 dark:text-red-200">Error: {dashboardError}</p>
                </div>
            )}

            {isLoading && (
                <div className="text-center py-12">
                    <p className="text-gray-500 dark:text-gray-400">Loading admin dashboard...</p>
                </div>
            )}

            {!isLoading && dashboard && (
                <>
                    {/* Dashboard Stats */}
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                        <div className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 px-6 py-8 shadow-sm">
                            <p className="text-xs uppercase tracking-[0.24em] text-slate-500 dark:text-slate-400">
                                Open Disputes
                            </p>
                            <p className="mt-3 text-3xl font-semibold text-slate-900 dark:text-white">
                                {dashboard.open_disputes || 0}
                            </p>
                        </div>

                        <div className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 px-6 py-8 shadow-sm">
                            <p className="text-xs uppercase tracking-[0.24em] text-slate-500 dark:text-slate-400">
                                Pending Payouts
                            </p>
                            <p className="mt-3 text-3xl font-semibold text-slate-900 dark:text-white">
                                {dashboard.pending_payouts || 0}
                            </p>
                        </div>

                        <div className="rounded-3xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 px-6 py-8 shadow-sm">
                            <p className="text-xs uppercase tracking-[0.24em] text-slate-500 dark:text-slate-400">
                                Fraud Alerts
                            </p>
                            <p className="mt-3 text-3xl font-semibold text-slate-900 dark:text-white">
                                {dashboard.fraud_alerts || 0}
                            </p>
                        </div>
                    </div>

                    {/* Disputes Table */}
                    <div className="space-y-4">
                        <h3 className="text-lg font-semibold text-slate-900 dark:text-white">Recent Disputes</h3>
                        {disputes.length === 0 ? (
                            <p className="text-slate-500 dark:text-slate-400">No disputes</p>
                        ) : (
                            <div className="overflow-x-auto rounded-lg border border-slate-200 dark:border-slate-700">
                                <table className="w-full text-sm">
                                    <thead className="border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900">
                                        <tr>
                                            <th className="px-6 py-3 text-left font-semibold text-slate-900 dark:text-white">
                                                ID
                                            </th>
                                            <th className="px-6 py-3 text-left font-semibold text-slate-900 dark:text-white">
                                                Status
                                            </th>
                                            <th className="px-6 py-3 text-left font-semibold text-slate-900 dark:text-white">
                                                Created
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {disputes.slice(0, 5).map((dispute: any) => (
                                            <tr
                                                key={dispute.id}
                                                className="border-b border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-900"
                                            >
                                                <td className="px-6 py-3 text-slate-900 dark:text-white">{dispute.id}</td>
                                                <td className="px-6 py-3">
                                                    <span className="inline-flex items-center rounded-full px-2.5 py-1 text-xs font-semibold bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-200">
                                                        {dispute.status || "Open"}
                                                    </span>
                                                </td>
                                                <td className="px-6 py-3 text-slate-600 dark:text-slate-400">
                                                    {new Date(dispute.created_at || "").toLocaleDateString()}
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        )}
                    </div>

                    {/* Fraud Alerts */}
                    <div className="space-y-4">
                        <h3 className="text-lg font-semibold text-slate-900 dark:text-white">Fraud Alerts</h3>
                        {alerts.length === 0 ? (
                            <p className="text-slate-500 dark:text-slate-400">No active fraud alerts</p>
                        ) : (
                            <div className="space-y-3">
                                {alerts.slice(0, 5).map((alert: any) => (
                                    <div
                                        key={alert.id}
                                        className="rounded-lg border border-red-200 dark:border-red-800 bg-red-50 dark:bg-red-900/20 p-4"
                                    >
                                        <div className="flex items-start justify-between">
                                            <div>
                                                <p className="font-semibold text-red-900 dark:text-red-200">
                                                    User {alert.user_id}: {alert.reason}
                                                </p>
                                                <p className="text-sm text-red-800 dark:text-red-300 mt-1">
                                                    Severity: {alert.severity}
                                                </p>
                                            </div>
                                            <span className="text-xs text-red-700 dark:text-red-300">
                                                {new Date(alert.created_at || "").toLocaleDateString()}
                                            </span>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>

                    {/* Audit Logs */}
                    <div className="space-y-4">
                        <h3 className="text-lg font-semibold text-slate-900 dark:text-white">Recent Audit Logs</h3>
                        {logs.length === 0 ? (
                            <p className="text-slate-500 dark:text-slate-400">No audit logs</p>
                        ) : (
                            <div className="space-y-2">
                                {logs.slice(0, 10).map((log: any, idx: number) => (
                                    <div
                                        key={idx}
                                        className="text-sm text-slate-600 dark:text-slate-400 py-2 border-b border-slate-200 dark:border-slate-700 last:border-0"
                                    >
                                        {log.action || "Unknown action"} - {new Date(log.created_at || "").toLocaleString()}
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                </>
            )}
        </div>
    );
}
