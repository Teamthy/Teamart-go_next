"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { getStoredCustomer, hasRole, type RoleKey } from "@/lib/auth-state";

interface RouteGuardProps {
    children: React.ReactNode;
    requiredRole?: RoleKey;
    redirectTo?: string;
}

export default function RouteGuard({
    children,
    requiredRole,
    redirectTo = "/auth/onboarding",
}: RouteGuardProps) {
    const router = useRouter();
    const [hydrated, setHydrated] = useState(false);

    useEffect(() => {
        setHydrated(true);
    }, []);

    const customer = useMemo(() => {
        if (typeof window === "undefined") return null;
        return getStoredCustomer();
    }, []);

    const hasAccess = useMemo(
        () => Boolean(customer && (!requiredRole || hasRole(requiredRole))),
        [customer, requiredRole]
    );

    useEffect(() => {
        if (hydrated && !hasAccess) {
            router.replace(redirectTo);
        }
    }, [hydrated, hasAccess, redirectTo, router]);

    if (!hydrated || !hasAccess) {
        return (
            <div className="py-10 text-sm text-zinc-500">Checking your access…</div>
        );
    }

    return <>{children}</>;
}
