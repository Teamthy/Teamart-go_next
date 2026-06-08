/**
 * PremiumAuthLayout Component
 * Glassmorphic auth layout matching TikTok Shop / Stripe design
 */

"use client";

import React from "react";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";

interface PremiumAuthLayoutProps {
    children: React.ReactNode;
    title?: string;
    description?: string;
    illustration?: React.ReactNode;
    showBackButton?: boolean;
    backHref?: string;
    onBack?: () => void;
    variant?: "default" | "compact" | "wide";
}

export default function PremiumAuthLayout({
    children,
    title,
    description,
    illustration,
    showBackButton = false,
    backHref = "/auth",
    onBack,
    variant = "default",
}: PremiumAuthLayoutProps) {
    const handleBack = () => {
        if (onBack) {
            onBack();
        } else if (backHref && typeof window !== "undefined") {
            window.history.back();
        }
    };

    return (
        <div
            className="min-h-screen bg-gradient-to-br from-[var(--surface-muted)] to-[var(--background)] px-4 py-8 sm:px-6 lg:px-8"
            role="main"
        >
            {/* Background gradient accent */}
            <div className="pointer-events-none fixed inset-0 -z-10 overflow-hidden">
                <div className="absolute -top-40 right-0 h-80 w-80 rounded-full bg-pink-500/5 blur-3xl" />
                <div className="absolute -bottom-40 left-20 h-80 w-80 rounded-full bg-purple-500/5 blur-3xl" />
            </div>

            <div className="mx-auto max-w-7xl">
                {/* Header with back button */}
                {showBackButton && (
                    <button
                        onClick={handleBack}
                        className="mb-6 inline-flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium text-[var(--text-muted)] transition hover:bg-[var(--surface)] hover:text-[var(--foreground)]"
                        aria-label="Go back"
                    >
                        <ArrowLeft className="h-4 w-4" />
                        Back
                    </button>
                )}

                {/* Main content */}
                {variant === "wide" ? (
                    <div className="grid gap-12 lg:grid-cols-2 lg:items-center">
                        {/* Left column - Illustration or branding */}
                        {illustration ? (
                            <div className="hidden lg:block">{illustration}</div>
                        ) : (
                            <div className="hidden space-y-8 lg:block">
                                <div>
                                    <h1 className="text-4xl font-bold tracking-tight text-zinc-900 sm:text-5xl">
                                        Teamart
                                    </h1>
                                    <p className="mt-4 text-lg text-zinc-600">
                                        The social commerce platform for creators and merchants.
                                    </p>
                                </div>

                                {/* Premium features list */}
                                <div className="space-y-4">
                                    {[
                                        "Launch your store in minutes",
                                        "Host live shopping events",
                                        "Earn commissions as a creator",
                                        "Connect with millions of shoppers",
                                    ].map((feature) => (
                                        <div key={feature} className="flex items-center gap-3">
                                            <div className="h-2 w-2 rounded-full bg-pink-500" />
                                            <span className="text-sm text-zinc-600">{feature}</span>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}

                        {/* Right column - Auth form */}
                        <div className="rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)]/95 p-8 shadow-lg backdrop-blur-xl sm:p-10">
                            {title && (
                                <div className="mb-8 text-center">
                                    <h2 className="text-2xl font-bold text-[var(--foreground)]">{title}</h2>
                                    {description && (
                                        <p className="mt-2 text-sm text-[var(--text-muted)]">{description}</p>
                                    )}
                                </div>
                            )}
                            {children}
                        </div>
                    </div>
                ) : variant === "compact" ? (
                    <div className="mx-auto max-w-sm">
                        <div className="text-center mb-8">
                            <h2 className="text-2xl font-bold text-zinc-900">{title}</h2>
                            {description && (
                                <p className="mt-2 text-sm text-zinc-600">{description}</p>
                            )}
                        </div>
                        <div className="rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)]/95 p-6 shadow-lg backdrop-blur-xl">
                            {children}
                        </div>
                    </div>
                ) : (
                    /* Default layout */
                    <div className="mx-auto max-w-2xl">
                        <div className="text-center mb-8">
                            {title && <h2 className="text-3xl font-bold text-zinc-900">{title}</h2>}
                            {description && (
                                <p className="mt-3 text-base text-zinc-600">{description}</p>
                            )}
                        </div>

                        <div className="rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)]/95 p-8 shadow-lg backdrop-blur-xl sm:p-10">
                            {children}
                        </div>
                    </div>
                )}
            </div>

            {/* Footer */}
            <div className="mx-auto mt-8 max-w-7xl text-center text-xs text-[var(--text-muted)]">
                <p>
                    © 2026 Teamart.{" "}
                    <Link href="/privacy" className="underline hover:text-zinc-700">
                        Privacy
                    </Link>
                    {" · "}
                    <Link href="/terms" className="underline hover:text-zinc-700">
                        Terms
                    </Link>
                </p>
            </div>
        </div>
    );
}
