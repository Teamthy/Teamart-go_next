"use client";

import type { ReactNode } from "react";
import PremiumAuthLayout from "@/components/auth/PremiumAuthLayout";

interface AuthShellProps {
    title?: string;
    description?: string;
    variant?: "default" | "compact" | "wide";
    showBackButton?: boolean;
    backHref?: string;
    attentionText?: string;
    footer?: ReactNode;
    children: ReactNode;
}

export default function AuthShell({
    title,
    description,
    variant = "compact",
    showBackButton = false,
    backHref = "/auth",
    attentionText,
    footer,
    children,
}: AuthShellProps) {
    return (
        <PremiumAuthLayout
            title={title}
            description={description}
            variant={variant}
            showBackButton={showBackButton}
            backHref={backHref}
        >
            <div className="space-y-6">
                {attentionText ? (
                    <div className="rounded-3xl border border-[var(--surface-border)] bg-[var(--surface-muted)] px-4 py-3 text-sm text-[var(--text-muted)]">
                        {attentionText}
                    </div>
                ) : null}
                {children}
                {footer}
            </div>
        </PremiumAuthLayout>
    );
}
