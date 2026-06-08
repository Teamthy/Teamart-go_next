/**
 * Signup Page
 * Email + Password + Role selection with real-time validation
 */

"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { Controller } from "react-hook-form";
import { useFormValidation } from "@/hooks/useFormValidation";
import { signupSchema } from "@/schemas/auth.schema";
import { useAuthStore } from "@/store/useAuthStore";
import AuthShell from "@/components/auth/AuthShell";
import Button from "@/components/ui/button";
import Input from "@/components/ui/input";
import Link from "next/link";
import { Eye, EyeOff } from "lucide-react";
import { getErrorMessage } from "@/lib/form";

export default function SignupPage() {
    const router = useRouter();
    const { signup } = useAuthStore();

    const roleOptions = [
        { value: "customer", label: "Buyer", description: "Browse, save, and shop your favorite creators." },
        { value: "creator", label: "Creator", description: "Share work, grow an audience, and sell original products." },
        { value: "merchant", label: "Seller", description: "Open a store and connect directly with shoppers." },
    ] as const;

    const [selectedRole, setSelectedRole] = useState<"customer" | "creator" | "merchant">(() => {
        if (typeof window === "undefined") {
            return "customer";
        }

        try {
            const params = new URLSearchParams(window.location.search);
            return (params.get("role") || "customer") as "customer" | "creator" | "merchant";
        } catch {
            return "customer";
        }
    });
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);

    const { control, formState: { errors }, handleSubmit, watch, onSubmitHandler, isSubmitting } = useFormValidation({
        schema: signupSchema,
        onSubmit: async (data) => {
            await signup(data.email, data.password, selectedRole);
            router.push(`/auth/verify-email?email=${encodeURIComponent(data.email)}`);
        },
    });

    const password = watch("password");

    const getPasswordStrength = (passwordValue: string) => {
        let strength = 0;
        if (passwordValue.length >= 12) strength++;
        if (/[a-z]/.test(passwordValue)) strength++;
        if (/[A-Z]/.test(passwordValue)) strength++;
        if (/[0-9]/.test(passwordValue)) strength++;
        if (/[^A-Za-z0-9]/.test(passwordValue)) strength++;

        if (strength <= 2) return "weak";
        if (strength <= 3) return "medium";
        return "strong";
    };

    const passwordStrength = getPasswordStrength(password || "");

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
        switch (selectedRole) {
            case "customer":
                return "Buyer";
            case "creator":
                return "Creator";
            case "merchant":
                return "Seller";
        }
    };

    return (
        <AuthShell
            title="Create Your Account"
            description={`Join as a ${getRoleLabel()}. You'll receive an email to verify your account.`}
            variant="compact"
            showBackButton
            backHref="/auth"
        >
            <form onSubmit={handleSubmit(onSubmitHandler)} className="space-y-5">
                <div className="rounded-3xl border border-[var(--surface-border)] bg-[var(--surface)] p-4">
                    <p className="text-sm font-semibold text-[var(--foreground)]">Choose your account type</p>
                    <p className="mt-1 text-xs text-[var(--text-muted)]">
                        Pick the role that best matches how you plan to use Teamart today.
                    </p>

                    <div className="mt-4 grid gap-3 sm:grid-cols-3">
                        {roleOptions.map((option) => {
                            const active = selectedRole === option.value;
                            return (
                                <button
                                    key={option.value}
                                    type="button"
                                    onClick={() => setSelectedRole(option.value)}
                                    aria-pressed={active}
                                    className={`rounded-3xl border p-4 text-left transition ${active ? "border-[var(--primary)] bg-[var(--primary-muted)] text-[var(--foreground)]" : "border-[var(--surface-border)] bg-[var(--surface)] text-[var(--foreground)] hover:border-[var(--primary)] hover:bg-[var(--surface-muted)]"}`}
                                >
                                    <p className="text-sm font-semibold">{option.label}</p>
                                    <p className="mt-1 text-xs text-[var(--text-muted)]">{option.description}</p>
                                </button>
                            );
                        })}
                    </div>
                </div>

                <Controller
                    name="email"
                    control={control}
                    render={({ field }) => (
                        <Input
                            {...field}
                            type="email"
                            id="email"
                            label="Email Address"
                            placeholder="you@example.com"
                            error={errors.email ? getErrorMessage(errors.email, "Invalid email") : undefined}
                            aria-invalid={errors.email ? "true" : "false"}
                            className="pr-10"
                        />
                    )}
                />

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
                                    <Input
                                        {...field}
                                        type={showPassword ? "text" : "password"}
                                        id="password"
                                        placeholder="••••••••••••"
                                        helperText={password ? `${passwordStrength.charAt(0).toUpperCase() + passwordStrength.slice(1)} password` : "At least 12 characters, one uppercase letter, one number, and one special character."}
                                        helperTextId="password-strength"
                                        error={errors.password ? getErrorMessage(errors.password, "Invalid password") : undefined}
                                        className="pr-10 font-mono text-sm"
                                    />
                                )}
                            />
                            <button
                                type="button"
                                onClick={() => setShowPassword(!showPassword)}
                                aria-label={showPassword ? "Hide password" : "Show password"}
                                className="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-zinc-700"
                            >
                                {showPassword ? (
                                    <EyeOff className="h-5 w-5" />
                                ) : (
                                    <Eye className="h-5 w-5" />
                                )}
                            </button>
                        </div>

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
                                <Input
                                    {...field}
                                    type={showConfirmPassword ? "text" : "password"}
                                    id="confirmPassword"
                                    placeholder="••••••••••••"
                                    error={errors.confirmPassword ? getErrorMessage(errors.confirmPassword, "Passwords do not match") : undefined}
                                    className="pr-10 font-mono text-sm"
                                />
                            )}
                        />
                        <button
                            type="button"
                            onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                            aria-label={showConfirmPassword ? "Hide confirmation password" : "Show confirmation password"}
                            className="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-zinc-700"
                        >
                            {showConfirmPassword ? (
                                <EyeOff className="h-5 w-5" />
                            ) : (
                                <Eye className="h-5 w-5" />
                            )}
                        </button>
                    </div>
                </div>

                {/* Submit button */}
                <Button
                    type="submit"
                    disabled={isSubmitting}
                    variant="primary"
                    className="w-full"
                >
                    {isSubmitting ? (
                        <span className="flex items-center justify-center gap-2">
                            <div className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                            Creating account...
                        </span>
                    ) : (
                        "Continue to Email Verification"
                    )}
                </Button>

                {/* Terms */}
                <p className="text-xs text-center text-zinc-600">
                    By signing up, you agree to our{" "}
                    <Link href="/terms" className="font-semibold text-pink-600 hover:underline">
                        Terms of Service
                    </Link>{" "}
                    and{" "}
                    <Link href="/privacy" className="font-semibold text-pink-600 hover:underline">
                        Privacy Policy
                    </Link>
                </p>
            </form>
        </AuthShell>
    );
}
