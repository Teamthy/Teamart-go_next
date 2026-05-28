"use client";

import { useState, type FormEvent } from "react";

interface AuthFormProps {
    mode: "login" | "register" | "forgot" | "otp" | "reset";
}

export default function AuthForm({ mode }: AuthFormProps) {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [otp, setOtp] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState(false);

    const actionLabel = {
        login: "Sign in",
        register: "Create account",
        forgot: "Send reset link",
        otp: "Verify code",
        reset: "Reset password",
    }[mode];

    const description = {
        login: "Use your email and password to access your storefront, feed, and creator tools.",
        register: "Create your Teamart account and unlock social commerce, livestream shopping, and creator features.",
        forgot: "Enter your email and we’ll send a password reset link so you can get back into your account.",
        otp: "Enter the one-time code sent to your device to continue securely.",
        reset: "Choose a new password for your account and keep your creator journey moving.",
    }[mode];

    const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        setLoading(true);
        setError(null);
        setSuccess(false);

        await new Promise((resolve) => setTimeout(resolve, 800));

        if (mode === "forgot" && !email.includes("@")) {
            setError("Please enter a valid email address.");
        } else if (mode === "reset" && (!password || password.length < 8 || password !== confirmPassword)) {
            setError("Your passwords must match and be at least 8 characters.");
        } else {
            setSuccess(true);
        }

        setLoading(false);
    };

    return (
        <div className="rounded-3xl border border-white/10 bg-white p-8 shadow-xl shadow-pink-500/5">
            <div className="space-y-3">
                <p className="text-sm uppercase tracking-[0.32em] text-pink-500">Teamart auth</p>
                <h2 className="text-3xl font-semibold text-slate-950">{actionLabel}</h2>
                <p className="text-sm text-slate-600">{description}</p>
            </div>

            <form onSubmit={handleSubmit} className="mt-8 space-y-5">
                {(mode === "login" || mode === "register" || mode === "forgot") && (
                    <label className="block text-sm text-slate-700">
                        Email address
                        <input
                            type="email"
                            value={email}
                            onChange={(event) => setEmail(event.target.value)}
                            placeholder="you@example.com"
                            className="mt-3 w-full rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-900 outline-none focus:border-pink-400 focus:ring-2 focus:ring-pink-100"
                        />
                    </label>
                )}

                {(mode === "login" || mode === "register" || mode === "reset") && (
                    <label className="block text-sm text-slate-700">
                        Password
                        <input
                            type="password"
                            value={password}
                            onChange={(event) => setPassword(event.target.value)}
                            placeholder="Enter your password"
                            className="mt-3 w-full rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-900 outline-none focus:border-pink-400 focus:ring-2 focus:ring-pink-100"
                        />
                    </label>
                )}

                {mode === "register" && (
                    <label className="block text-sm text-slate-700">
                        Confirm password
                        <input
                            type="password"
                            value={confirmPassword}
                            onChange={(event) => setConfirmPassword(event.target.value)}
                            placeholder="Confirm your password"
                            className="mt-3 w-full rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-900 outline-none focus:border-pink-400 focus:ring-2 focus:ring-pink-100"
                        />
                    </label>
                )}

                {mode === "otp" && (
                    <label className="block text-sm text-slate-700">
                        Verification code
                        <input
                            type="text"
                            value={otp}
                            onChange={(event) => setOtp(event.target.value)}
                            placeholder="Enter 6-digit code"
                            className="mt-3 w-full rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-900 outline-none focus:border-pink-400 focus:ring-2 focus:ring-pink-100"
                        />
                    </label>
                )}

                {error ? (
                    <div className="rounded-3xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{error}</div>
                ) : null}

                {success ? (
                    <div className="rounded-3xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
                        {mode === "forgot"
                            ? "Reset link sent. Check your inbox."
                            : mode === "reset"
                                ? "Password updated successfully."
                                : "Ready to continue."}
                    </div>
                ) : null}

                <button
                    type="submit"
                    disabled={loading}
                    className="w-full rounded-3xl bg-[#E91E63] px-4 py-3 text-sm font-semibold text-white transition hover:bg-[#d81b60] disabled:cursor-not-allowed disabled:opacity-60"
                >
                    {loading ? "Working..." : actionLabel}
                </button>
            </form>
        </div>
    );
}
