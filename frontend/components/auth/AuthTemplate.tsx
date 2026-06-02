"use client";

import Link from "next/link";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { ArrowRight, Lock, Mail, ShieldCheck, Sparkles } from "lucide-react";

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

function persistAuthResponse(response: any) {
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
    const [remember, setRemember] = useState(false);
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
                const sess = pending ? JSON.parse(pending) : null;
                const sessionId = sess?.session_id || sess?.sessionID;

                if (!sessionId) throw new Error("Missing MFA session");

                const res = await api.verifyOTP(sessionId, otp);

                persistAuthResponse(res);

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
                await api.signup(email, password);

                router.push("/auth/login");
            }
        } catch (err: any) {
            setError(err?.message || "Request failed");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen bg-[radial-gradient(circle_at_top,_rgba(236,72,153,0.14),transparent_28%),linear-gradient(180deg,#050816_0%,#0b1124_100%)] px-4 py-12 text-white sm:px-6 lg:px-8">
            <div className="mx-auto grid max-w-6xl gap-10 lg:grid-cols-[1.05fr_0.95fr]">
                <div className="space-y-8 rounded-[3rem] border border-white/10 bg-white/5 p-8">
                    <span className="inline-flex items-center gap-2 rounded-full border border-fuchsia-400/30 bg-fuchsia-500/10 px-4 py-2 text-sm text-fuchsia-200">
                        <Sparkles className="h-4 w-4" />
                        Teamart social commerce
                    </span>

                    <AuthIllustration />
                </div>

                <div className="rounded-[2.5rem] border border-white/10 bg-slate-950/95 p-8">
                    <h2 className="mb-6 text-3xl font-semibold">
                        {isLogin
                            ? "Sign in"
                            : isRegister
                                ? "Create account"
                                : isMfa
                                    ? "Verify identity"
                                    : "Forgot password"}
                    </h2>

                    {error && (
                        <div className="mb-6 rounded-3xl border border-red-500/20 bg-red-500/10 p-4 text-sm text-red-200">
                            {error}
                        </div>
                    )}

                    <form onSubmit={handleSubmit} className="space-y-6">
                        {!isMfa && !isForgot && (
                            <>
                                <input
                                    type="email"
                                    placeholder="Email"
                                    className="w-full rounded-3xl border border-slate-800 bg-slate-900 px-4 py-3"
                                    required
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                />

                                <input
                                    type="password"
                                    placeholder="Password"
                                    className="w-full rounded-3xl border border-slate-800 bg-slate-900 px-4 py-3"
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
                                        className={`rounded-3xl border px-4 py-3 ${selectedRole === role
                                            ? "border-fuchsia-400 bg-fuchsia-500/10"
                                            : "border-white/10"
                                            }`}
                                    >
                                        {role}
                                    </button>
                                ))}
                            </div>
                        )}

                        {isMfa && (
                            <input
                                type="text"
                                placeholder="Verification code"
                                className="w-full rounded-3xl border border-slate-800 bg-slate-900 px-4 py-3"
                                required
                                value={otp}
                                onChange={(e) => setOtp(e.target.value)}
                            />
                        )}

                        {isForgot && (
                            <input
                                type="email"
                                placeholder="Email address"
                                className="w-full rounded-3xl border border-slate-800 bg-slate-900 px-4 py-3"
                            />
                        )}

                        <button
                            disabled={loading}
                            type="submit"
                            className="inline-flex w-full items-center justify-center gap-3 rounded-full bg-fuchsia-500 px-6 py-3"
                        >
                            {loading ? "Working…" : "Continue"}
                            <ArrowRight className="h-4 w-4" />
                        </button>
                    </form>

                    {!isMfa && (
                        <p className="mt-6 text-center text-sm text-slate-400">
                            {isLogin ? (
                                <>
                                    New to Teamart?{" "}
                                    <Link href="/auth/register" className="text-fuchsia-300">
                                        Create account
                                    </Link>
                                </>
                            ) : (
                                <>
                                    Already have an account?{" "}
                                    <Link href="/auth/login" className="text-fuchsia-300">
                                        Sign in
                                    </Link>
                                </>
                            )}
                        </p>
                    )}
                </div>
            </div>
        </div>
    );
}