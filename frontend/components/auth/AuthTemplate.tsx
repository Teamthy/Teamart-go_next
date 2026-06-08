"use client";

import Link from "next/link";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { ArrowRight, Sparkles } from "lucide-react";

import AuthShell from "@/components/auth/AuthShell";
import Input from "@/components/ui/input";
import * as api from "@/lib/api";

type AuthVariant = "login" | "register" | "mfa" | "forgot";

const registerRoles = ["Shopper", "Creator", "Merchant"] as const;
type Role = (typeof registerRoles)[number];

function normalizeRole(role?: string): Role {
    if (!role) return "Shopper";

    switch (role.toLowerCase()) {
        case "creator":
            return "Creator";
        case "merchant":
            return "Merchant";
        case "shopper":
        default:
            return "Shopper";
    }
}

function persistAuthResponse(response: Record<string, unknown> | null) {
    if (!response || typeof window === "undefined") return;

    if (response.access_token) {
        localStorage.setItem("access_token", response.access_token);
    }

    if (response.refresh_token) {
        localStorage.setItem("refresh_token", response.refresh_token);
    }

    if (response.user) {
        localStorage.setItem("user", JSON.stringify(response.user));
    }

    sessionStorage.setItem("session", JSON.stringify(response));
}

function AuthIllustration() {
    return (
        <div className="mx-auto mb-8 w-full max-w-xs rounded-[2.5rem] bg-white/5 p-6 text-white shadow-2xl shadow-fuchsia-500/10 sm:mb-0">
            <div className="relative h-72 overflow-hidden rounded-[2rem] bg-gradient-to-br from-fuchsia-500/10 via-transparent to-transparent p-4">
                <svg viewBox="0 0 180 180" className="h-full w-full">
                    <circle cx="90" cy="90" r="86" fill="#FCE4EC" />
                    <rect x="50" y="36" width="80" height="104" rx="22" fill="#111827" />
                    <rect x="58" y="44" width="64" height="88" rx="16" fill="#f8fafc" />
                    <path
                        d="M82 72c0-9 7.5-16 16-16s16 7 16 16c0 13-16 24-16 24s-16-11-16-24Z"
                        fill="#E91E63"
                    />
                </svg>
            </div>
        </div>
    );
}

export default function AuthTemplate({
    variant = "login",
    initialRole,
}: {
    variant?: AuthVariant;
    initialRole?: Role | string;
}) {
    const router = useRouter();

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [otp, setOtp] = useState("");
    const [selectedRole, setSelectedRole] = useState<Role>(normalizeRole(initialRole));
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const isLogin = variant === "login";
    const isRegister = variant === "register";
    const isMfa = variant === "mfa";
    const isForgot = variant === "forgot";

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        setError(null);
        setLoading(true);

        try {
            if (isMfa) {
                const pending = sessionStorage.getItem("pendingSession");
                const sess = pending ? (JSON.parse(pending) as Record<string, unknown>) : null;
                const sessionId = sess?.session_id || sess?.sessionID;

                if (!sessionId) throw new Error("Missing MFA session");

                await api.verifyOTP(String(sessionId), otp);

                sessionStorage.removeItem("pendingSession");
                router.push("/");
                return;
            }

            if (isForgot) {
                router.push("/auth/login");
                return;
            }

            if (isLogin) {
                const res = await api.login(email, password);

                if (res.requires_mfa || res.requiresMFA) {
                    sessionStorage.setItem("pendingSession", JSON.stringify(res));
                    router.push("/auth/mfa");
                    return;
                }

                persistAuthResponse(res);

                router.push("/");

                return;
            }

            if (isRegister) {
                const apiRole = selectedRole === "Shopper" ? "customer" : selectedRole.toLowerCase();
                await api.signup(email, password, apiRole);

                router.push("/auth/login");
            }
        } catch (error: unknown) {
            setError(error instanceof Error ? error.message : "Request failed");
        } finally {
            setLoading(false);
        }
    };

    return (
        <AuthShell
            title={isLogin ? "Sign in" : isRegister ? "Create account" : isMfa ? "Verify identity" : "Forgot password"}
            description={
                isLogin
                    ? "Access your Teamart account and continue shopping."
                    : isRegister
                        ? "Create a secure account to buy, create, or sell with Teamart."
                        : isMfa
                            ? "Enter the verification code sent to your email."
                            : "Enter your email to reset your password."
            }
            variant="default"
            attentionText="One account can support shoppers, creators, and merchants across the platform."
            showBackButton={!isLogin}
            backHref="/auth"
        >
            <div className="space-y-8">
                <div className="rounded-[2rem] border border-zinc-200 bg-zinc-50 p-6">
                    <div className="inline-flex items-center gap-2 rounded-full bg-pink-50 px-4 py-2 text-sm font-semibold text-pink-700">
                        <Sparkles className="h-4 w-4" />
                        Teamart social commerce
                    </div>
                    <div className="mt-6 space-y-3 text-zinc-700">
                        <p className="text-lg font-semibold text-zinc-900">Fast, secure account access</p>
                        <p className="text-sm leading-6">
                            Use one growth-ready auth flow for buyers, creators, and sellers. Everything stays in sync across the app.
                        </p>
                    </div>
                    <div className="mt-6 hidden lg:block">
                        <AuthIllustration />
                    </div>
                </div>

                <div className="rounded-[2rem] border border-zinc-200 bg-white p-6 shadow-sm">
                    {error && (
                        <div className="mb-6 rounded-3xl border border-red-200 bg-red-50 p-4 text-sm text-red-700">
                            {error}
                        </div>
                    )}

                    <form onSubmit={handleSubmit} className="space-y-6">
                        {!isMfa && !isForgot && (
                            <>
                                <Input
                                    type="email"
                                    label="Email Address"
                                    placeholder="you@example.com"
                                    required
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                />

                                <Input
                                    type="password"
                                    label="Password"
                                    placeholder="••••••••••••"
                                    required
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                />
                            </>
                        )}

                        {isRegister && (
                            <div className="grid gap-3 sm:grid-cols-3">
                                {registerRoles.map((role) => (
                                    <button
                                        key={role}
                                        type="button"
                                        onClick={() => setSelectedRole(role)}
                                        className={`rounded-3xl border px-4 py-3 text-sm font-semibold ${selectedRole === role
                                            ? "border-pink-400 bg-pink-50 text-zinc-900"
                                            : "border-zinc-200 bg-white text-zinc-700"
                                            }`}
                                    >
                                        {role}
                                    </button>
                                ))}
                            </div>
                        )}

                        {isMfa && (
                            <Input
                                type="text"
                                label="Verification Code"
                                placeholder="Enter 6-digit code"
                                required
                                value={otp}
                                onChange={(e) => setOtp(e.target.value)}
                            />
                        )}

                        {isForgot && (
                            <Input
                                type="email"
                                label="Email Address"
                                placeholder="you@example.com"
                                required
                            />
                        )}

                        <button
                            disabled={loading}
                            type="submit"
                            className="inline-flex w-full items-center justify-center gap-3 rounded-full bg-pink-600 px-6 py-3 text-sm font-semibold text-white transition hover:bg-pink-700 disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            {loading ? "Working…" : "Continue"}
                            <ArrowRight className="h-4 w-4" />
                        </button>
                    </form>

                    {!isMfa && (
                        <p className="mt-6 text-center text-sm text-zinc-500">
                            {isLogin ? (
                                <>
                                    New to Teamart?{" "}
                                    <Link href="/auth/register" className="font-semibold text-pink-600 hover:text-pink-700">
                                        Create account
                                    </Link>
                                </>
                            ) : (
                                <>
                                    Already have an account?{" "}
                                    <Link href="/auth/login" className="font-semibold text-pink-600 hover:text-pink-700">
                                        Sign in
                                    </Link>
                                </>
                            )}
                        </p>
                    )}
                </div>
            </div>
        </AuthShell>
    );
}