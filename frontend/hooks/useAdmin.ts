"use client";

import { useState, useEffect } from "react";
import * as api from "@/lib/api";

export interface DashboardSummary {
    open_disputes: number;
    pending_payouts: number;
    fraud_alerts: number;
}

export interface Dispute {
    id: string;
    status: string;
    reason?: string;
    created_at?: string;
}

export interface FraudAlert {
    id: string;
    user_id: number;
    reason: string;
    severity: string;
    created_at?: string;
}

export function useAdminDashboard() {
    const [dashboard, setDashboard] = useState<DashboardSummary | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetch = async () => {
            setIsLoading(true);
            setError(null);
            try {
                const response = await api.getAdminDashboard();
                setDashboard(response);
            } catch (err: any) {
                setError(err.message || "Failed to fetch dashboard");
            } finally {
                setIsLoading(false);
            }
        };

        fetch();
    }, []);

    return { dashboard, isLoading, error };
}

export function useDisputes() {
    const [disputes, setDisputes] = useState<Dispute[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetch = async () => {
            setIsLoading(true);
            setError(null);
            try {
                const response = await api.listDisputes();
                setDisputes(response || []);
            } catch (err: any) {
                setError(err.message || "Failed to fetch disputes");
            } finally {
                setIsLoading(false);
            }
        };

        fetch();
    }, []);

    return { disputes, isLoading, error };
}

export function useFraudAlerts() {
    const [alerts, setAlerts] = useState<FraudAlert[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetch = async () => {
            setIsLoading(true);
            setError(null);
            try {
                const response = await api.listFraudAlerts();
                setAlerts(response || []);
            } catch (err: any) {
                setError(err.message || "Failed to fetch fraud alerts");
            } finally {
                setIsLoading(false);
            }
        };

        fetch();
    }, []);

    return { alerts, isLoading, error };
}

export function useAuditLogs() {
    const [logs, setLogs] = useState<any[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetch = async () => {
            setIsLoading(true);
            setError(null);
            try {
                const response = await api.listAuditLogs();
                setLogs(response || []);
            } catch (err: any) {
                setError(err.message || "Failed to fetch audit logs");
            } finally {
                setIsLoading(false);
            }
        };

        fetch();
    }, []);

    return { logs, isLoading, error };
}
