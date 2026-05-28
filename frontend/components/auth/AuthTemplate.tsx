"use client";

import Link from "next/link";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { ArrowRight, Lock, Mail, ShieldCheck, Sparkles } from "lucide-react";
import * as api from "@/lib/api";

const registerRoles = ["Shopper", "Creator", "Merchant"] as const;

type Role = (typeof registerRoles)[number];

function AuthIllustration() {
    return (
        <div className="mx-auto mb-8 w-full max-w-xs rounded-[2.5rem] bg-white/5 p-6 text-white shadow-2xl shadow-fuchsia-500/10 sm:mb-0">
            <div className="relative h-72 overflow-hidden rounded-[2rem] bg-gradient-to-br from-fuchsia-500/10 via-transparent to-transparent p-4">
                <svg viewBox="0 0 180 180" className="h-full w-full">
                    <circle cx="90" cy="90" r="86" fill="#FCE4EC" />
                    <rect x="50" y="36" width="80" height="104" rx="22" fill="#111827" />
                    <rect x="58" y="44" width="64" height="88" rx="16" fill="#f8fafc" />
                    <path d="M82 72c0-9 7.5-16 16-16s16 7 16 16c0 13-16 24-16 24s-16-11-16-24Z" fill="#E91E63" />
                    <path d="M70 44c0-3.5 2.8-6.2 6.2-6.2h37.6c3.4 0 6.2 2.8 6.2 6.2v4.4H70v-4.4Z" fill="#E91E63" />
                    <circle cx="44" cy="140" r="8" fill="#E91E63" />
                    <circle cx="136" cy="40" r="6" fill="#F8BBD0" />
                    <circle cx="30" cy="58" r="4" fill="#E91E63" />
                    <path d="M124 80c0-3.5 2.8-6.4 6.3-6.4s6.3 2.8 6.3 6.4-2.8 6.4-6.3 6.4-6.3-2.8-6.3-6.4Z" fill="#E91E63" />
                </svg>
            </div>
        </div>
    );
}

export default function AuthTemplate({ variant = "login" }: { variant?: "login" | "register" | "mfa" }) {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [otp, setOtp] = useState("");
    const [remember, setRemember] = useState(false);
    const [selectedRole, setSelectedRole] = useState<Role>("Shopper");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const router = useRouter();

    const isLogin = variant === "login";
    const isRegister = variant === "register";
    const isMfa = variant === "mfa";

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setLoading(true);

        try {
            if (isMfa) {
                const pending = sessionStorage.getItem("pendingSession");
                const sess = pending ? JSON.parse(pending) : null;
                const sessionId = sess?.session_id || sess?.sessionID || sess?.sessionID;
                if (!sessionId) throw new Error("Missing pending session for MFA");

                await api.verifyOTP(sessionId, otp);
                sessionStorage.removeItem("pendingSession");
                router.push("/");
                return;
            }

            if (isLogin) {
                const res = await api.login(email, password);
                sessionStorage.setItem("session", JSON.stringify(res));
                if (res.requires_mfa || res.requiresMFA) {
                    sessionStorage.setItem("pendingSession", JSON.stringify(res));
                    router.push("/auth/mfa");
                } else {
                    router.push("/");
                }
                return;
            }

            if (isRegister) {
                const res = await api.signup(email, password);
                sessionStorage.setItem("signupResult", JSON.stringify(res));
                router.push("/auth/login");
                return;
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
                <div className="space-y-8 rounded-[3rem] border border-white/10 bg-white/5 p-8 shadow-2xl shadow-fuchsia-500/10 backdrop-blur-xl sm:p-10">
                    <div className="flex flex-col gap-4">
                        <span className="inline-flex items-center gap-2 rounded-full border border-fuchsia-400/30 bg-fuchsia-500/10 px-4 py-2 text-sm text-fuchsia-200">
                            <Sparkles className="h-4 w-4" />
                            Built for social commerce, live shopping, and creator growth
                        </span>
                        <div className="space-y-4">
                            <h1 className="text-4xl font-semibold tracking-tight text-white sm:text-5xl">
                                A modern marketplace experience for creators, shoppers, and merchants.
                            </h1>
                            <p className="max-w-2xl text-base leading-8 text-slate-300">
                                Teamart lets you discover creator drops, launch livestream commerce, and manage storefront operations in one polished platform.
                            </p>
                        </div>
                        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                            {[
                                { label: "Creator-first feed", value: "Discover trending drops." },
                                { label: "Live commerce rooms", value: "Shop while watching streams." },
                                { label: "Role-based onboarding", value: "Shopper, creator, or merchant." },
                            ].map((item) => (
                                <div key={item.label} className="rounded-3xl border border-white/10 bg-slate-950/80 p-5">
                                    <p className="text-sm uppercase tracking-[0.28em] text-slate-400">{item.label}</p>
                                    <p className="mt-3 text-base font-semibold text-white">{item.value}</p>
                                </div>
                            ))}
                        </div>
                    </div>
                    <div className="grid gap-4 rounded-[2.5rem] border border-white/10 bg-slate-950/80 p-6 sm:p-8">
                        <AuthIllustration />
                        <div className="space-y-3">
                            <p className="text-sm uppercase tracking-[0.35em] text-fuchsia-300">Fast start</p>
                            <h2 className="text-2xl font-semibold text-white">Keep your account flow simple and onboarding confident.</h2>
                            <p className="text-sm leading-6 text-slate-400">
                                Whether you’re here to shop, sell, or stream, Teamart helps you move from browsing to commerce with confidence.
                            </p>
                        </div>
                    </div>
                </div>
                <div className="rounded-[2.5rem] border border-white/10 bg-slate-950/95 p-8 shadow-2xl shadow-slate-950/40 sm:p-10">
                    <div className="mb-8 space-y-4">
                        <div className="flex items-center justify-between gap-4">
                            <div>
                                <p className="text-xs uppercase tracking-[0.35em] text-fuchsia-400">Teamart</p>
                                <h2 className="mt-3 text-3xl font-semibold text-white">
                                    {isLogin ? "Sign in" : isRegister ? "Create account" : "Verify identity"}
                                </h2>
                            </div>
                            <div className="rounded-3xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-slate-300">
                                {isLogin ? "Returning user" : isRegister ? "Role-driven signup" : "Two-step verification"}
                            </div>
                        </div>
                        <p className="text-sm leading-6 text-slate-400">
                            {isLogin
                                ? "Use your email and password to access your marketplace dashboard."
                                : isRegister
                                    ? "Choose a role and start building your creator, merchant, or shopper experience."
                                    : "Enter the one-time code from your authentication device."}
                        </p>
                    </div>

                    {error ? (
                        <div className="rounded-3xl border border-red-500/20 bg-red-500/10 p-4 text-sm text-red-200">{error}</div>
                    ) : null}

                    <form onSubmit={handleSubmit} className="space-y-6">
                        {!isMfa && (
                            <button type="button" className="flex w-full items-center justify-center gap-3 rounded-full border border-white/10 bg-white/5 px-5 py-3 text-sm font-semibold text-white transition hover:bg-white/10">
                                <span className="grid h-5 w-5 place-items-center rounded-full bg-[#EA4335] text-[11px] font-semibold">G</span>
                                Continue with Google
                            </button>
                        )}

                        {!isMfa && (
                            <div className="relative flex items-center gap-3 rounded-full border border-slate-800 bg-slate-900/90 px-4 py-3 text-slate-400">
                                <Mail className="h-4 w-4" />
                                <span className="text-sm">or continue with email</span>
                            </div>
                        )}

                        {!isMfa && (
                            <div className="grid gap-4">
                                <label className="grid gap-2 text-sm text-slate-300">
                                    Email address
                                    <div className="flex items-center gap-3 rounded-3xl border border-slate-800 bg-slate-900/90 px-4 py-3">
                                        <Mail className="h-5 w-5 text-fuchsia-400" />
                                        <input
                                            type="email"
                                            placeholder="you@example.com"
                                            className="w-full bg-transparent text-sm text-white outline-none placeholder:text-slate-500"
                                            required
                                            value={email}
                                            onChange={(e) => setEmail(e.target.value)}
                                        />
                                    </div>
                                </label>
                                <label className="grid gap-2 text-sm text-slate-300">
                                    Password
                                    <div className="flex items-center gap-3 rounded-3xl border border-slate-800 bg-slate-900/90 px-4 py-3">
                                        <Lock className="h-5 w-5 text-slate-400" />
                                        <input
                                            type="password"
                                            placeholder="Enter your password"
                                            className="w-full bg-transparent text-sm text-white outline-none placeholder:text-slate-500"
                                            required
                                            value={password}
                                            onChange={(e) => setPassword(e.target.value)}
                                        />
                                    </div>
                                </label>
                            </div>
                        )}

                        {isRegister && (
                            <div className="space-y-4 rounded-3xl border border-white/10 bg-slate-900/80 p-4">
                                <p className="text-xs uppercase tracking-[0.35em] text-fuchsia-300">Role selection</p>
                                <div className="grid gap-3 sm:grid-cols-3">
                                    {registerRoles.map((role) => (
                                        <button
                                            key={role}
                                            type="button"
                                            onClick={() => setSelectedRole(role)}
                                            className={`rounded-3xl border px-4 py-3 text-sm font-semibold transition ${selectedRole === role ? "border-fuchsia-400 bg-fuchsia-500/10 text-white" : "border-white/10 bg-slate-950/80 text-slate-300 hover:border-fuchsia-300 hover:text-white"}`}
                                        >
                                            {role}
                                        </button>
                                    ))}
                                </div>
                                <p className="text-sm leading-6 text-slate-400">
                                    {selectedRole} access gives you the right tools for your marketplace journey.
                                </p>
                            </div>
                        )}

                        {isMfa && (
                            <label className="grid gap-2 text-sm text-slate-300">
                                Verification code
                                <div className="flex items-center gap-3 rounded-3xl border border-slate-800 bg-slate-900/90 px-4 py-3">
                                    <ShieldCheck className="h-5 w-5 text-fuchsia-400" />
                                    <input
                                        type="text"
                                        placeholder="Enter verification code"
                                        className="w-full bg-transparent text-sm text-white outline-none placeholder:text-slate-500"
                                        required
                                        value={otp}
                                        onChange={(e) => setOtp(e.target.value)}
                                    />
                                </div>
                            </label>
                        )}

                        {!isMfa && (
                            <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                                <label className="inline-flex items-center gap-2 text-sm text-slate-400">
                                    <input
                                        type="checkbox"
                                        className="h-4 w-4 rounded border-slate-700 bg-slate-800 text-fuchsia-500"
                                        checked={remember}
                                        onChange={(e) => setRemember(e.target.checked)}
                                    />
                                    Remember me
                                </label>
                                <Link href="/auth/forgot-password" className="text-sm text-fuchsia-300 hover:text-white">
                                    Forgot password?
                                </Link>
                            </div>
                        )}

                        <button
                            disabled={loading}
                            type="submit"
                            className="inline-flex w-full items-center justify-center gap-3 rounded-full bg-fuchsia-500 px-6 py-3 text-sm font-semibold text-white shadow-xl shadow-fuchsia-500/20 transition hover:bg-fuchsia-400 disabled:cursor-not-allowed disabled:opacity-70"
                        >
                            {loading ? "Working…" : isLogin ? "Sign in" : isRegister ? "Create account" : "Verify"}
                            <ArrowRight className="h-4 w-4" />
                        </button>
                    </form>

                    {!isMfa && (
                        <p className="mt-6 text-center text-sm text-slate-400">
                            {isLogin ? (
                                <>New to Teamart? <Link href="/auth/register" className="text-fuchsia-300 hover:text-white">Create account</Link></>
                            ) : (
                                <>Already have an account? <Link href="/auth/login" className="text-fuchsia-300 hover:text-white">Sign in</Link></>
                            )}
                        </p>
                    )}
                </div>
            </div>
        </div>
    );
}
