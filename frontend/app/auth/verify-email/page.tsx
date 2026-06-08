/**
 * Email Verification Page
 * OTP verification with resend logic
 */

"use client";

import React, { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/useAuthStore";
import * as api from "@/lib/api";
import AuthShell from "@/components/auth/AuthShell";
import Button from "@/components/ui/button";
import OTPInput from "@/components/auth/OTPInput";
import { Mail, AlertCircle } from "lucide-react";

export default function VerifyEmailPage() {
    const router = useRouter();
    const [email] = useState<string>(() => {
        if (typeof window === "undefined") {
            return "";
        }

        try {
            const params = new URLSearchParams(window.location.search);
            const queryEmail = params.get("email") || "";
            const stored = getPendingSessionData();
            return queryEmail || stored.email || "";
        } catch {
            return "";
        }
    });
    const { verifyOTP } = useAuthStore();

    function getPendingSessionData() {
        if (typeof window === "undefined") {
            return { session_id: null, email: null };
        }

        const raw = sessionStorage.getItem("pendingSession") || localStorage.getItem("session");
        if (!raw) {
            return { session_id: null, email: null };
        }

        try {
            const parsed = JSON.parse(raw) as Record<string, unknown>;
            return {
                session_id: String(parsed.session_id ?? parsed.sessionID ?? "") || null,
                email: String(parsed.email ?? parsed.user?.email ?? "") || null,
            };
        } catch {
            return { session_id: null, email: null };
        }
    }
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [statusMessage, setStatusMessage] = useState("Check your inbox and spam folder.");
    const [resendCountdown, setResendCountdown] = useState(0);
    const [resendAttempts, setResendAttempts] = useState(0);

    // Resend countdown timer
    useEffect(() => {
        if (resendCountdown > 0) {
            const timer = setTimeout(() => setResendCountdown(resendCountdown - 1), 1000);
            return () => clearTimeout(timer);
        }
    }, [resendCountdown]);

    const handleOTPComplete = async (code: string) => {
        setIsLoading(true);
        setError(null);

        try {
            const success = await verifyOTP(code);
            if (success) {
                // Navigate to profile setup
                router.push("/auth/onboarding/profile");
            } else {
                setError("Invalid verification code. Please try again.");
            }
        } catch (err) {
            setError(
                err instanceof Error ? err.message : "Failed to verify code. Please try again."
            );
        } finally {
            setIsLoading(false);
        }
    };

    const handleResendOTP = async () => {
        if (resendCountdown > 0 || resendAttempts >= 3) return;

        setIsLoading(true);
        setError(null);

        try {
            const pending = getPendingSessionData();
            const sessionId = pending.session_id;
            const targetEmail = email || pending.email;

            if (!sessionId || !targetEmail) {
                throw new Error("Unable to locate verification session or email.");
            }

            await api.resendOTP(sessionId, targetEmail);
            setResendCountdown(60);
            setResendAttempts(resendAttempts + 1);
            setStatusMessage(`Verification code resent to ${targetEmail}.`);
        } catch (err) {
            setError(
                err instanceof Error ? err.message : "Failed to resend code. Please try again."
            );
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <AuthShell
            title="Verify Your Email"
            description={`We sent a 6-digit code to ${email || "your email"}. Enter it below.`}
            variant="compact"
            showBackButton
            backHref="/auth/signup"
            attentionText="Code delivery can take a few moments. Check your inbox and spam folder."
        >
            <div className="space-y-6">
                {/* Email display */}
                <div className="flex items-center justify-center gap-2 rounded-lg bg-blue-50 px-4 py-3 text-sm text-blue-900 border border-blue-200" aria-live="polite">
                    <Mail className="h-4 w-4 flex-shrink-0" />
                    <div>
                        <p className="font-semibold">{email}</p>
                        <p className="text-xs text-blue-800">{statusMessage}</p>
                    </div>
                </div>

                {/* OTP Input */}
                <OTPInput
                    onComplete={handleOTPComplete}
                    isLoading={isLoading}
                    error={error}
                    length={6}
                    autoFocus
                />

                {/* Error message */}
                {error && (
                    <div className="flex items-start gap-2 rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-900">
                        <AlertCircle className="h-4 w-4 mt-0.5 flex-shrink-0" />
                        <div>
                            <p className="font-semibold">{error}</p>
                            {error.includes("Invalid") && (
                                <p className="text-xs text-red-800 mt-1">
                                    Make sure you entered all 6 digits correctly.
                                </p>
                            )}
                        </div>
                    </div>
                )}

                {/* Resend section */}
                <div className="space-y-3 border-t border-zinc-200 pt-4">
                    <p className="text-sm text-zinc-600">Didn&apos;t receive the code?</p>

                    {resendCountdown > 0 ? (
                        <div className="text-center">
                            <p className="text-sm font-semibold text-zinc-900">
                                Request new code in <span className="text-pink-600">{resendCountdown}s</span>
                            </p>
                        </div>
                    ) : resendAttempts >= 3 ? (
                        <div className="rounded-lg bg-amber-50 p-4 border border-amber-200">
                            <p className="text-sm text-amber-900 font-semibold">Too many attempts</p>
                            <p className="text-xs text-amber-800 mt-1">
                                Please contact support if you continue to have issues.
                            </p>
                        </div>
                    ) : (
                        <Button
                            type="button"
                            onClick={handleResendOTP}
                            disabled={isLoading || resendCountdown > 0}
                            variant="secondary"
                            className="w-full"
                        >
                            Send Code Again
                        </Button>
                    )}

                    {resendAttempts > 0 && resendAttempts < 3 && (
                        <p className="text-xs text-zinc-500 text-center">
                            Resend attempts: {resendAttempts}/3
                        </p>
                    )}
                </div>

                {/* Help text */}
                <div className="rounded-lg bg-zinc-100 p-4 text-xs text-zinc-700 space-y-2">
                    <p className="font-semibold">💡 Tips:</p>
                    <ul className="space-y-1 ml-2">
                        <li>• Check your spam or promotions folder</li>
                        <li>• Code expires in 10 minutes</li>
                        <li>• You can paste the code directly</li>
                    </ul>
                </div>
            </div>
        </AuthShell>
    );
}
