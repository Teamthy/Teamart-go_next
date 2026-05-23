"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { z } from "zod";
import { useAuth } from "@/hooks/useAuth";

const authSchema = z.object({
    email: z.string().trim().email("Enter a valid email address"),
    password: z.string().min(8, "Password must be at least 8 characters"),
});

const mfaSchema = z.object({
    otp: z.string().trim().min(4, "Enter a valid verification code"),
});

export default function AuthTemplate({ variant = "login" }: { variant?: "login" | "register" | "mfa" }) {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [otp, setOtp] = useState("");
    const [remember, setRemember] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
    const router = useRouter();
    const { login, signup, verifyOTP, isLoading, error: authError } = useAuth();

    useEffect(() => {
        if (authError) {
            setError(authError);
        }
    }, [authError]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setFieldErrors({});

        try {
            if (variant === "mfa") {
                const parsed = mfaSchema.safeParse({ otp });
                if (!parsed.success) {
                    setFieldErrors({ otp: parsed.error.issues[0]?.message ?? "Enter your verification code" });
                    setError("Please fix the form errors before continuing.");
                    return;
                }

                const pending = sessionStorage.getItem("pendingSession");
                const sess = pending ? JSON.parse(pending) : null;
                const sessionId = sess?.session_id || sess?.sessionID || sess?.sessionID;
                if (!sessionId) throw new Error("Missing pending session for MFA");

                await verifyOTP(sessionId, otp);
                sessionStorage.removeItem("pendingSession");
                router.push("/");
                return;
            }

            const parsed = authSchema.safeParse({ email, password });
            if (!parsed.success) {
                const errors = parsed.error.flatten().fieldErrors;
                const nextFieldErrors: Record<string, string> = {};
                if (errors.email?.[0]) nextFieldErrors.email = errors.email[0];
                if (errors.password?.[0]) nextFieldErrors.password = errors.password[0];
                setFieldErrors(nextFieldErrors);
                setError("Please fix the form errors before continuing.");
                return;
            }

            if (variant === "login") {
                const res = await login(parsed.data.email, parsed.data.password);
                if (res?.requires_mfa || res?.requiresMFA) {
                    sessionStorage.setItem("pendingSession", JSON.stringify(res));
                    router.push("/auth/mfa");
                } else {
                    router.push("/");
                }
                return;
            }

            if (variant === "register") {
                await signup(parsed.data.email, parsed.data.password);
                router.push("/");
                return;
            }
        } catch (err: any) {
            setError(err?.message || "Request failed");
        }
    };

    return (
        <div className="flex h-[700px] w-full">
            <div className="w-full hidden md:inline-block">
                <img
                    className="h-full"
                    src="https://raw.githubusercontent.com/prebuiltui/prebuiltui/main/assets/login/leftSideImage.png"
                    alt="leftSideImage"
                />
            </div>

            <div className="w-full flex flex-col items-center justify-center">
                <form onSubmit={handleSubmit} className="md:w-96 w-80 flex flex-col items-center justify-center">
                    <h2 className="text-4xl text-gray-900 font-medium">{variant === "login" ? "Sign in" : variant === "register" ? "Create account" : "Verify"}</h2>
                    <p className="text-sm text-gray-500/90 mt-3">
                        {variant === "login"
                            ? "Welcome back! Please sign in to continue"
                            : variant === "register"
                                ? "Create an account to get started"
                                : "Enter the verification code sent to you"}
                    </p>

                    {error ? <div className="mt-4 rounded-md bg-red-50 px-4 py-2 text-sm text-red-700">{error}</div> : null}

                    {variant !== "mfa" && (
                        <>
                            <button type="button" className="w-full mt-8 bg-gray-500/10 flex items-center justify-center h-12 rounded-full">
                                <img src="https://raw.githubusercontent.com/prebuiltui/prebuiltui/main/assets/login/googleLogo.svg" alt="googleLogo" />
                            </button>

                            <div className="flex items-center gap-4 w-full my-5">
                                <div className="w-full h-px bg-gray-300/90" />
                                <p className="w-full text-nowrap text-sm text-gray-500/90">or sign in with email</p>
                                <div className="w-full h-px bg-gray-300/90" />
                            </div>

                            <div className="flex flex-col gap-2 w-full">
                                <div className="flex items-center w-full bg-transparent border border-gray-300/60 h-12 rounded-full overflow-hidden pl-6 gap-2">
                                    <svg width="16" height="11" viewBox="0 0 16 11" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <path fillRule="evenodd" clipRule="evenodd" d="M0 .55.571 0H15.43l.57.55v9.9l-.571.55H.57L0 10.45zm1.143 1.138V9.9h13.714V1.69l-6.503 4.8h-.697zM13.749 1.1H2.25L8 5.356z" fill="#6B7280" />
                                    </svg>
                                    <input
                                        type="email"
                                        placeholder="Email id"
                                        className="bg-transparent text-gray-500/80 placeholder-gray-500/80 outline-none text-sm w-full h-full"
                                        required
                                        aria-invalid={!!fieldErrors.email}
                                        value={email}
                                        onChange={(e) => setEmail(e.target.value)}
                                    />
                                </div>
                                {fieldErrors.email ? (
                                    <p className="text-sm text-rose-600">{fieldErrors.email}</p>
                                ) : null}
                            </div>

                            <div className="flex flex-col gap-2 w-full mt-6">
                                <div className="flex items-center w-full bg-transparent border border-gray-300/60 h-12 rounded-full overflow-hidden pl-6 gap-2">
                                    <svg width="13" height="17" viewBox="0 0 13 17" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <path d="M13 8.5c0-.938-.729-1.7-1.625-1.7h-.812V4.25C10.563 1.907 8.74 0 6.5 0S2.438 1.907 2.438 4.25V6.8h-.813C.729 6.8 0 7.562 0 8.5v6.8c0 .938.729 1.7 1.625 1.7h9.75c.896 0 1.625-.762 1.625-1.7zM4.063 4.25c0-1.406 1.093-2.55 2.437-2.55s2.438 1.144 2.438 2.55V6.8H4.061z" fill="#6B7280" />
                                    </svg>
                                    <input
                                        type="password"
                                        placeholder="Password"
                                        className="bg-transparent text-gray-500/80 placeholder-gray-500/80 outline-none text-sm w-full h-full"
                                        required
                                        aria-invalid={!!fieldErrors.password}
                                        value={password}
                                        onChange={(e) => setPassword(e.target.value)}
                                    />
                                </div>
                                {fieldErrors.password ? (
                                    <p className="text-sm text-rose-600">{fieldErrors.password}</p>
                                ) : null}
                            </div>

                            <div className="w-full flex items-center justify-between mt-8 text-gray-500/80">
                                <div className="flex items-center gap-2">
                                    <input className="h-5" type="checkbox" id="checkbox" checked={remember} onChange={(e) => setRemember(e.target.checked)} />
                                    <label className="text-sm" htmlFor="checkbox">
                                        Remember me
                                    </label>
                                </div>
                                <a className="text-sm underline" href="#">
                                    Forgot password?
                                </a>
                            </div>

                            <button disabled={isLoading} type="submit" className="mt-8 w-full h-11 rounded-full text-white bg-indigo-500 hover:opacity-90 transition-opacity">
                                {isLoading ? "Working…" : variant === "login" ? "Login" : "Create account"}
                            </button>

                            <p className="text-gray-500/90 text-sm mt-4">
                                {variant === "login" ? (
                                    <>Don’t have an account? <a className="text-indigo-400 hover:underline" href="/auth/register">Sign up</a></>
                                ) : (
                                    <>Already have an account? <a className="text-indigo-400 hover:underline" href="/auth/login">Sign in</a></>
                                )}
                            </p>
                        </>
                    )}

                    {variant === "mfa" && (
                        <>
                            <div className="flex flex-col gap-2 w-full">
                                <div className="flex items-center w-full bg-transparent border border-gray-300/60 h-12 rounded-full overflow-hidden pl-6 gap-2">
                                    <input
                                        type="text"
                                        placeholder="Enter verification code"
                                        className="bg-transparent text-gray-500/80 placeholder-gray-500/80 outline-none text-sm w-full h-full"
                                        required
                                        aria-invalid={!!fieldErrors.otp}
                                        value={otp}
                                        onChange={(e) => setOtp(e.target.value)}
                                    />
                                </div>
                                {fieldErrors.otp ? (
                                    <p className="text-sm text-rose-600">{fieldErrors.otp}</p>
                                ) : null}
                            </div>

                            <button disabled={isLoading} type="submit" className="mt-8 w-full h-11 rounded-full text-white bg-indigo-500 hover:opacity-90 transition-opacity">
                                {isLoading ? "Verifying…" : "Verify"}
                            </button>
                        </>
                    )}
                </form>
            </div>
        </div>
    );
}
