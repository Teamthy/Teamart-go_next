/**
 * Signup Page
 * Email + Password + Role selection with real-time validation
 */

"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Controller } from "react-hook-form";
import { useFormValidation } from "@/hooks/useFormValidation";
import { signupSchema } from "@/schemas/auth.schema";
import { useAuthStore } from "@/store/useAuthStore";
import PremiumAuthLayout from "@/components/auth/PremiumAuthLayout";
import { Eye, EyeOff, AlertCircle, CheckCircle2, Mail } from "lucide-react";
import { getErrorMessage } from "@/lib/form";

export default function SignupPage() {
    const router = useRouter();
    const { signup } = useAuthStore();

    const [role, setRole] = useState<"customer" | "creator" | "merchant">("customer");

    // Read role from URL on client-side to avoid SSR/prerender issues
    useEffect(() => {
        try {
            const params = new URLSearchParams(window.location.search);
            const r = (params.get("role") || "customer") as "customer" | "creator" | "merchant";
            setRole(r);
        } catch (e) {
            setRole("customer");
        }
    }, []);
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [passwordStrength, setPasswordStrength] = useState<"weak" | "medium" | "strong">("weak");

    const { control, formState: { errors }, handleSubmit, watch, onSubmitHandler, isSubmitting } = useFormValidation({
        schema: signupSchema,
        onSubmit: async (data) => {
            await signup(data.email, data.password, role);
            // Navigate to email verification
            router.push(`/auth/verify-email?email=${encodeURIComponent(data.email)}`);
        },
    });

    const password = watch("password");

    // Calculate password strength
    useEffect(() => {
        if (!password) {
            setPasswordStrength("weak");
            return;
        }

        let strength = 0;
        if (password.length >= 12) strength++;
        if (/[a-z]/.test(password)) strength++;
        if (/[A-Z]/.test(password)) strength++;
        if (/[0-9]/.test(password)) strength++;
        if (/[^A-Za-z0-9]/.test(password)) strength++;

        if (strength <= 2) setPasswordStrength("weak");
        else if (strength <= 3) setPasswordStrength("medium");
        else setPasswordStrength("strong");
    }, [password]);

    const getPasswordStrengthColor = () => {
        switch (passwordStrength) {
            case "weak":
                return "text-red-600";
            case "medium":
                return "text-yellow-600";
            case "strong":
                return "text-green-600";
        }
    };

    const getRoleLabel = () => {
        switch (role) {
            case "customer":
                return "Buyer";
            case "creator":
                return "Creator";
            case "merchant":
                return "Seller";
        }
    };

    return (
        <PremiumAuthLayout
            title="Create Your Account"
            description={`Join as a ${getRoleLabel()}. You'll receive an email to verify your account.`}
            variant="compact"
            showBackButton
            backHref="/auth"
        >
            <form onSubmit={handleSubmit(onSubmitHandler)} className="space-y-5">
                {/* Role badge */}
                <div className="inline-block rounded-full bg-pink-100 px-3 py-1 text-xs font-semibold text-pink-700">
                    {getRoleLabel()} Account
                </div>

                {/* Email field */}
                <div>
                    <label htmlFor="email" className="block text-sm font-semibold text-zinc-900 mb-2">
                        Email Address
                        <span className="text-red-500">*</span>
                    </label>
                    <div className="relative">
                        <Controller
                            name="email"
                            control={control}
                            render={({ field }) => (
                                <input
                                    {...field}
                                    type="email"
                                    id="email"
                                    placeholder="you@example.com"
                                    className={`w-full rounded-lg border-2 px-4 py-3 pr-10 transition-all focus:outline-none ${errors.email
                                        ? "border-red-300 bg-red-50 focus:border-red-400 focus:ring-2 focus:ring-red-100"
                                        : "border-zinc-300 focus:border-pink-500 focus:ring-2 focus:ring-pink-100"
                                        }`}
                                    aria-invalid={errors.email ? "true" : "false"}
                                />
                            )}
                        />
                        {!errors.email && password && (
                            <Mail className="absolute right-3 top-1/2 h-5 w-5 -translate-y-1/2 text-blue-500" />
                        )}
                    </div>
                    {errors.email && (
                        (() => {
                            const msg = getErrorMessage(errors.email, "Invalid email");
                            return msg ? (
                                <div className="mt-2 flex items-start gap-2 text-sm text-red-600">
                                    <AlertCircle className="h-4 w-4 mt-0.5 flex-shrink-0" />
                                    <span>{msg}</span>
                                </div>
                            ) : null;
                        })()
                    )}
                </div>

                {/* Password field */}
                <div>
                    <label htmlFor="password" className="block text-sm font-semibold text-zinc-900 mb-2">
                        Password
                        <span className="text-red-500">*</span>
                    </label>
                    <div className="space-y-2">
                        <div className="relative">
                            <Controller
                                name="password"
                                control={control}
                                render={({ field }) => (
                                    <input
                                        {...field}
                                        type={showPassword ? "text" : "password"}
                                        id="password"
                                        placeholder="••••••••••••"
                                        className={`w-full rounded-lg border-2 px-4 py-3 pr-10 transition-all focus:outline-none font-mono text-sm ${errors.password
                                            ? "border-red-300 bg-red-50 focus:border-red-400 focus:ring-2 focus:ring-red-100"
                                            : "border-zinc-300 focus:border-pink-500 focus:ring-2 focus:ring-pink-100"
                                            }`}
                                        aria-invalid={errors.password ? "true" : "false"}
                                    />
                                )}
                            />
                            <button
                                type="button"
                                onClick={() => setShowPassword(!showPassword)}
                                className="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-zinc-700"
                            >
                                {showPassword ? (
                                    <EyeOff className="h-5 w-5" />
                                ) : (
                                    <Eye className="h-5 w-5" />
                                )}
                            </button>
                        </div>

                        {/* Password strength meter */}
                        {password && (
                            <div className="space-y-1.5">
                                <div className="flex gap-1">
                                    {["weak", "medium", "strong"].map((level) => (
                                        <div
                                            key={level}
                                            className={`h-1 flex-1 rounded-full transition-colors ${(level === "weak" && passwordStrength) ||
                                                (level === "medium" && ["medium", "strong"].includes(passwordStrength)) ||
                                                (level === "strong" && passwordStrength === "strong")
                                                ? level === "weak"
                                                    ? "bg-red-500"
                                                    : level === "medium"
                                                        ? "bg-yellow-500"
                                                        : "bg-green-500"
                                                : "bg-zinc-200"
                                                }`}
                                        />
                                    ))}
                                </div>
                                <p className={`text-xs font-semibold ${getPasswordStrengthColor()}`}>
                                    {passwordStrength.charAt(0).toUpperCase() + passwordStrength.slice(1)} password
                                </p>
                            </div>
                        )}
                    </div>
                    {errors.password && (
                        <div className="mt-2 text-xs text-red-600 space-y-1">
                            <p className="font-semibold">Password must include:</p>
                            <ul className="space-y-0.5 ml-2">
                                <li className={password?.length >= 12 ? "line-through text-green-600" : ""}>
                                    ✓ At least 12 characters
                                </li>
                                <li className={/[A-Z]/.test(password) ? "line-through text-green-600" : ""}>
                                    ✓ One uppercase letter
                                </li>
                                <li className={/[a-z]/.test(password) ? "line-through text-green-600" : ""}>
                                    ✓ One lowercase letter
                                </li>
                                <li className={/[0-9]/.test(password) ? "line-through text-green-600" : ""}>
                                    ✓ One number
                                </li>
                                <li className={/[^A-Za-z0-9]/.test(password) ? "line-through text-green-600" : ""}>
                                    ✓ One special character
                                </li>
                            </ul>
                        </div>
                    )}
                </div>

                {/* Confirm password field */}
                <div>
                    <label htmlFor="confirmPassword" className="block text-sm font-semibold text-zinc-900 mb-2">
                        Confirm Password
                        <span className="text-red-500">*</span>
                    </label>
                    <div className="relative">
                        <Controller
                            name="confirmPassword"
                            control={control}
                            render={({ field }) => (
                                <input
                                    {...field}
                                    type={showConfirmPassword ? "text" : "password"}
                                    id="confirmPassword"
                                    placeholder="••••••••••••"
                                    className={`w-full rounded-lg border-2 px-4 py-3 pr-10 transition-all focus:outline-none font-mono text-sm ${errors.confirmPassword
                                        ? "border-red-300 bg-red-50 focus:border-red-400 focus:ring-2 focus:ring-red-100"
                                        : "border-zinc-300 focus:border-pink-500 focus:ring-2 focus:ring-pink-100"
                                        }`}
                                    aria-invalid={errors.confirmPassword ? "true" : "false"}
                                />
                            )}
                        />
                        <button
                            type="button"
                            onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                            className="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-zinc-700"
                        >
                            {showConfirmPassword ? (
                                <EyeOff className="h-5 w-5" />
                            ) : (
                                <Eye className="h-5 w-5" />
                            )}
                        </button>
                    </div>
                    {errors.confirmPassword && (
                        (() => {
                            const msg = getErrorMessage(errors.confirmPassword, "Passwords do not match");
                            return msg ? (
                                <div className="mt-2 flex items-start gap-2 text-sm text-red-600">
                                    <AlertCircle className="h-4 w-4 mt-0.5 flex-shrink-0" />
                                    <span>{msg}</span>
                                </div>
                            ) : null;
                        })()
                    )}
                </div>

                {/* Submit button */}
                <button
                    type="submit"
                    disabled={isSubmitting}
                    className="w-full rounded-lg bg-gradient-to-r from-pink-600 to-pink-500 px-4 py-3 font-semibold text-white transition-all hover:shadow-lg hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed disabled:scale-100"
                >
                    {isSubmitting ? (
                        <span className="flex items-center justify-center gap-2">
                            <div className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                            Creating account...
                        </span>
                    ) : (
                        "Continue to Email Verification"
                    )}
                </button>

                {/* Terms */}
                <p className="text-xs text-center text-zinc-600">
                    By signing up, you agree to our{" "}
                    <a href="/terms" className="font-semibold text-pink-600 hover:underline">
                        Terms of Service
                    </a>{" "}
                    and{" "}
                    <a href="/privacy" className="font-semibold text-pink-600 hover:underline">
                        Privacy Policy
                    </a>
                </p>
            </form>
        </PremiumAuthLayout>
    );
}
