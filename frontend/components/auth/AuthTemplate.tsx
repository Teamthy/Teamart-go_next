"use client";

import { useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import * as api from "@/lib/api";

type AuthVariant = "login" | "register" | "mfa";

const socialProviders = [
    {
        name: "Google",
        slug: "google",
        icon: (
            <svg viewBox="0 0 24 24" className="h-5 w-5" aria-hidden="true">
                <path
                    fill="#4285F4"
                    d="M21.805 10.023h-9.68v3.908h5.504c-.239 1.44-1.701 4.226-5.504 4.226-3.312 0-6.006-2.744-6.006-6.128s2.694-6.128 6.006-6.128c1.887 0 3.155.8 3.88 1.5l2.65-2.57C16.74 3.9 14.8 3 12.125 3 7.5 3 3.7 6.8 3.7 11.5S7.5 20 12.125 20c5.93 0 7.8-4.156 7.8-7.94 0-.534-.057-1.03-.12-1.037z"
                />
            </svg>
        ),
    },
];

const stateConfig: Record<AuthVariant, { title: string; subtitle: string; panelTitle: string; panelCopy: string; actionLabel: string }> = {
    login: {
        title: "Sign in",
        subtitle: "Welcome back! Please sign in to continue",
        panelTitle: "Pick up where you left off",
        panelCopy: "A fast path back to your storefront, orders, and creator tools.",
        actionLabel: "Login",
    },
    register: {
        title: "Create account",
        subtitle: "Create an account to get started",
        panelTitle: "Build your storefront faster",
        panelCopy: "Start with a single account and move into your next campaign in minutes.",
        actionLabel: "Create account",
    },
    mfa: {
        title: "Verify",
        subtitle: "Enter the verification code sent to you",
        panelTitle: "Secure your next session",
        panelCopy: "Confirm your device to finish sign in and keep your account protected.",
        actionLabel: "Verify",
    },
};

function buildArt(variant: AuthVariant) {
    const palettes = {
        login: { bg: "#0f172a", accent: "#8b5cf6", glow: "#38bdf8" },
        register: { bg: "#052e2b", accent: "#34d399", glow: "#22d3ee" },
        mfa: { bg: "#1f123f", accent: "#f59e0b", glow: "#f472b6" },
    };

    const palette = palettes[variant];
    const svg = `
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 800">
            <defs>
                <linearGradient id="bg" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" stop-color="${palette.bg}" />
                    <stop offset="100%" stop-color="#111827" />
                </linearGradient>
            </defs>
            <rect width="640" height="800" rx="40" fill="url(#bg)" />
            <circle cx="480" cy="140" r="120" fill="${palette.glow}" fill-opacity="0.24" />
            <circle cx="160" cy="220" r="100" fill="${palette.accent}" fill-opacity="0.2" />
            <path d="M112 540C186 430 278 388 418 410C500 423 566 480 612 586" stroke="${palette.accent}" stroke-width="12" stroke-linecap="round" fill="none" stroke-opacity="0.9" />
            <rect x="116" y="140" width="340" height="180" rx="24" fill="white" fill-opacity="0.08" />
            <rect x="140" y="200" width="120" height="16" rx="8" fill="white" fill-opacity="0.9" />
            <rect x="140" y="232" width="200" height="12" rx="6" fill="white" fill-opacity="0.7" />
            <rect x="140" y="258" width="168" height="12" rx="6" fill="white" fill-opacity="0.4" />
            <rect x="118" y="418" width="404" height="220" rx="28" fill="white" fill-opacity="0.06" stroke="white" stroke-opacity="0.12" />
            <circle cx="180" cy="518" r="42" fill="${palette.glow}" fill-opacity="0.55" />
            <rect x="244" y="488" width="180" height="14" rx="7" fill="white" fill-opacity="0.85" />
            <rect x="244" y="520" width="160" height="10" rx="5" fill="white" fill-opacity="0.5" />
        </svg>`;

    return `data:image/svg+xml;charset=UTF-8,${encodeURIComponent(svg)}`;
}

function persistAuthResponse(response: any) {
    if (!response || typeof window === "undefined") {
        return;
    }

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

export default function AuthTemplate({ variant = "login" }: { variant?: AuthVariant }) {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [otp, setOtp] = useState("");
    const [remember, setRemember] = useState(false);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const router = useRouter();

    const currentState = stateConfig[variant];
    const artSrc = useMemo(() => buildArt(variant), [variant]);

    const handleSocialLogin = (provider: string) => {
        router.push(`/auth/social?provider=${provider}`);
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setLoading(true);

        try {
            if (variant === "mfa") {
                const pending = sessionStorage.getItem("pendingSession");
                const sess = pending ? JSON.parse(pending) : null;
                const sessionId = sess?.session_id || sess?.sessionID;
                if (!sessionId) throw new Error("Missing pending session for MFA");

                const res = await api.verifyOTP(sessionId, otp);
                persistAuthResponse(res);
                sessionStorage.removeItem("pendingSession");
                router.push("/");
                return;
            }

            if (variant === "login") {
                const res = await api.login(email, password);
                persistAuthResponse(res);
                sessionStorage.setItem("session", JSON.stringify(res));

                if (res.requires_mfa || res.requiresMFA) {
                    sessionStorage.setItem("pendingSession", JSON.stringify(res));
                    router.push("/auth/mfa");
                    return;
                }

                router.push("/");
                return;
            }

            const res = await api.signup(email, password);
            persistAuthResponse(res);
            sessionStorage.setItem("signupResult", JSON.stringify(res));
            router.push("/auth/login");
        } catch (err: any) {
            setError(err?.message || "Request failed");
        } finally {
            setLoading(false);
        }
    };

    const footerLink = variant === "login" ? "/auth/register" : "/auth/login";
    const footerPrompt = variant === "login" ? "Don’t have an account?" : "Already have an account?";
    const footerCta = variant === "login" ? "Sign up" : "Sign in";

    return (
        <div className="flex min-h-[700px] w-full overflow-hidden rounded-[32px] border border-gray-200 bg-white shadow-xl">
            <div className="hidden md:flex md:w-[46%] relative items-center justify-center bg-slate-950">
                <img src={artSrc} alt="Auth illustration" className="h-full w-full object-cover" />
                <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent" />
                <div className="absolute bottom-8 left-8 right-8 text-white">
                    <p className="text-xs uppercase tracking-[0.35em] text-white/70">{variant === "login" ? "Secure sign in" : variant === "register" ? "New account" : "Step-up verification"}</p>
                    <h2 className="mt-3 text-3xl font-semibold">{currentState.panelTitle}</h2>
                    <p className="mt-3 max-w-md text-sm leading-6 text-white/85">{currentState.panelCopy}</p>
                </div>
            </div>

            <div className="w-full md:w-[54%] flex flex-col items-center justify-center px-6 py-10 sm:px-10 lg:px-14">
                <form onSubmit={handleSubmit} className="w-full max-w-md flex flex-col items-center justify-center">
                    <h2 className="text-4xl text-gray-900 font-medium">{currentState.title}</h2>
                    <p className="text-sm text-gray-500/90 mt-3 text-center">{currentState.subtitle}</p>

                    {error ? <div className="mt-4 w-full rounded-md bg-red-50 px-4 py-2 text-sm text-red-700">{error}</div> : null}

                    {variant !== "mfa" && (
                        <>
                            <div className="mt-8 w-full space-y-3">
                                {socialProviders.map((provider) => (
                                    <button
                                        type="button"
                                        key={provider.slug}
                                        onClick={() => handleSocialLogin(provider.slug)}
                                        className="w-full h-12 rounded-full border border-gray-300 bg-white text-gray-700 hover:border-indigo-300 hover:text-indigo-600 transition-colors flex items-center justify-center gap-3"
                                    >
                                        {provider.icon}
                                        <span>Continue with {provider.name}</span>
                                    </button>
                                ))}
                            </div>

                            <div className="flex items-center gap-4 w-full my-5">
                                <div className="w-full h-px bg-gray-300/90" />
                                <p className="w-full text-nowrap text-sm text-gray-500/90">or sign in with email</p>
                                <div className="w-full h-px bg-gray-300/90" />
                            </div>

                            <div className="flex items-center w-full bg-transparent border border-gray-300/60 h-12 rounded-full overflow-hidden pl-6 gap-2">
                                <svg width="16" height="11" viewBox="0 0 16 11" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <path fillRule="evenodd" clipRule="evenodd" d="M0 .55.571 0H15.43l.57.55v9.9l-.571.55H.57L0 10.45zm1.143 1.138V9.9h13.714V1.69l-6.503 4.8h-.697zM13.749 1.1H2.25L8 5.356z" fill="#6B7280" />
                                </svg>
                                <input
                                    type="email"
                                    placeholder="Email id"
                                    className="bg-transparent text-gray-500/80 placeholder-gray-500/80 outline-none text-sm w-full h-full"
                                    required
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                />
                            </div>

                            <div className="flex items-center mt-6 w-full bg-transparent border border-gray-300/60 h-12 rounded-full overflow-hidden pl-6 gap-2">
                                <svg width="13" height="17" viewBox="0 0 13 17" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <path d="M13 8.5c0-.938-.729-1.7-1.625-1.7h-.812V4.25C10.563 1.907 8.74 0 6.5 0S2.438 1.907 2.438 4.25V6.8h-.813C.729 6.8 0 7.562 0 8.5v6.8c0 .938.729 1.7 1.625 1.7h9.75c.896 0 1.625-.762 1.625-1.7zM4.063 4.25c0-1.406 1.093-2.55 2.437-2.55s2.438 1.144 2.438 2.55V6.8H4.061z" fill="#6B7280" />
                                </svg>
                                <input
                                    type="password"
                                    placeholder="Password"
                                    className="bg-transparent text-gray-500/80 placeholder-gray-500/80 outline-none text-sm w-full h-full"
                                    required
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                />
                            </div>

                            <div className="w-full flex items-center justify-between mt-8 text-gray-500/80">
                                <div className="flex items-center gap-2">
                                    <input className="h-5" type="checkbox" id="checkbox" checked={remember} onChange={(e) => setRemember(e.target.checked)} />
                                    <label className="text-sm" htmlFor="checkbox">
                                        Remember me
                                    </label>
                                </div>
                                <button type="button" className="text-sm underline" onClick={() => router.push("/auth/forgot-password")}>
                                    Forgot password?
                                </button>
                            </div>

                            <button disabled={loading} type="submit" className="mt-8 w-full h-11 rounded-full text-white bg-indigo-500 hover:opacity-90 transition-opacity">
                                {loading ? "Working…" : currentState.actionLabel}
                            </button>

                            <p className="text-gray-500/90 text-sm mt-4">
                                {footerPrompt} <button type="button" className="text-indigo-400 hover:underline" onClick={() => router.push(footerLink)}>{footerCta}</button>
                            </p>
                        </>
                    )}

                    {variant === "mfa" && (
                        <>
                            <div className="flex items-center w-full bg-transparent border border-gray-300/60 h-12 rounded-full overflow-hidden pl-6 gap-2">
                                <input
                                    type="text"
                                    placeholder="Enter verification code"
                                    className="bg-transparent text-gray-500/80 placeholder-gray-500/80 outline-none text-sm w-full h-full"
                                    required
                                    value={otp}
                                    onChange={(e) => setOtp(e.target.value)}
                                />
                            </div>

                            <button disabled={loading} type="submit" className="mt-8 w-full h-11 rounded-full text-white bg-indigo-500 hover:opacity-90 transition-opacity">
                                {loading ? "Verifying…" : currentState.actionLabel}
                            </button>
                        </>
                    )}
                </form>
            </div>
        </div>
    );
}
