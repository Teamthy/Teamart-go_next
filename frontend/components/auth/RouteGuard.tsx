"use client";

import { useEffect, useState } from "react";
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
    const [ready, setReady] = useState(false);

    useEffect(() => {
        const customer = getStoredCustomer();

        if (!customer) {
            router.replace(redirectTo);
            return;
        }

        if (requiredRole && !hasRole(requiredRole)) {
            router.replace(redirectTo);
            return;
        }

        setReady(true);
    }, [redirectTo, requiredRole, router]);

    if (!ready) {
        return (
            <div className="py-10 text-sm text-zinc-500">
                Checking your access...
            </div>
        );
    }

    return <>{children}</>;
}
