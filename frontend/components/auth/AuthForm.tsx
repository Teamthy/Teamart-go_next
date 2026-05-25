"use client";

import { useState } from "react";

interface AuthFormProps {
    mode: "login" | "register" | "mfa";
}

export default function AuthForm({ mode }: AuthFormProps) {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [otp, setOtp] = useState("");

    return (
        <div className="rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
            <div className="space-y-2">
                <h3 className="text-xl font-semibold text-slate-900">{mode === "login" ? "Sign in" : mode === "register" ? "Create account" : "Multi-factor authentication"}</h3>
                <p className="text-sm text-slate-600">
                    {mode === "login"
                        ? "Log in to access your customer dashboard and order history."
                        : mode === "register"
                            ? "Create your account and join the creator-first shopping experience."
                            : "Enter the code from your authenticator app or SMS."}
                </p>
            </div>

            <form className="mt-6 space-y-5">
                {(mode === "login" || mode === "register") && (
                    <>
                        <label className="block text-sm font-medium text-slate-700">
                            Email
                            <input
                                type="email"
                                value={email}
                                onChange={(event) => setEmail(event.target.value)}
                                className="mt-2 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-900 outline-none focus:border-slate-400 focus:ring-2 focus:ring-slate-200"
                                placeholder="you@example.com"
                            />
                        </label>

                        <label className="block text-sm font-medium text-slate-700">
                            Password
                            <input
                                type="password"
                                value={password}
                                onChange={(event) => setPassword(event.target.value)}
                                className="mt-2 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-900 outline-none focus:border-slate-400 focus:ring-2 focus:ring-slate-200"
                                placeholder="Enter your password"
                            />
                        </label>
                    </>
                )}

                {mode === "mfa" && (
                    <label className="block text-sm font-medium text-slate-700">
                        Verification code
                        <input
                            type="text"
                            value={otp}
                            onChange={(event) => setOtp(event.target.value)}
                            className="mt-2 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-900 outline-none focus:border-slate-400 focus:ring-2 focus:ring-slate-200"
                            placeholder="Enter 6-digit code"
                        />
                    </label>
                )}

                <button type="button" className="w-full rounded-2xl bg-slate-900 px-4 py-3 text-sm font-semibold text-white transition hover:bg-slate-700">
                    {mode === "login" ? "Sign in" : mode === "register" ? "Create account" : "Verify code"}
                </button>
            </form>
        </div>
    );
}
