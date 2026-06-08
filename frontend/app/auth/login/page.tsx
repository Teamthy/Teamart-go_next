/**
 * Login Page
 * Email + Password login with "Forgot Password" option
 */

"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { Controller } from "react-hook-form";
import { useFormValidation } from "@/hooks/useFormValidation";
import { loginSchema } from "@/schemas/auth.schema";
import { useAuthStore } from "@/store/useAuthStore";
import AuthShell from "@/components/auth/AuthShell";
import Input from "@/components/ui/input";
import { Eye, EyeOff, ArrowRight } from "lucide-react";
import { getErrorMessage } from "@/lib/form";

export default function LoginPage() {
    const router = useRouter();
    const { login } = useAuthStore();

    const [showPassword, setShowPassword] = useState(false);

    const { control, formState: { errors }, handleSubmit, onSubmitHandler, isSubmitting } = useFormValidation({
        schema: loginSchema,
        onSubmit: async (data) => {
            try {
                await login(data.email, data.password);
                router.push("/feed");
            } catch {
                // Error handling is done by the store
            }
        },
    });

    return (
        <AuthShell
            title="Welcome Back"
            description="Sign in to continue shopping and explore new products."
            variant="compact"
            showBackButton
            backHref="/auth"
        >
            <form onSubmit={handleSubmit(onSubmitHandler)} className="space-y-5">
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
                        />
                    )}
                />

                {/* Password field */}
                <div>
                    <div className="flex items-center justify-between mb-2">
                        <span className="text-sm font-semibold text-zinc-900">Password</span>
                        <button
                            type="button"
                            onClick={() => router.push("/auth/forgot-password")}
                            className="text-xs font-semibold text-pink-600 hover:text-pink-700"
                        >
                            Forgot?
                        </button>
                    </div>
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
                </div>

                {/* Remember me - optional feature */}
                <label className="flex items-center gap-2 cursor-pointer">
                    <input
                        type="checkbox"
                        defaultChecked={false}
                        className="h-4 w-4 rounded border-zinc-300 cursor-pointer"
                    />
                    <span className="text-sm text-zinc-700">Keep me signed in</span>
                </label>

                {/* Submit button */}
                <button
                    type="submit"
                    disabled={isSubmitting}
                    className="w-full rounded-lg bg-gradient-to-r from-pink-600 to-pink-500 px-4 py-3 font-semibold text-white transition-all hover:shadow-lg hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed disabled:scale-100 flex items-center justify-center gap-2"
                >
                    {isSubmitting ? (
                        <>
                            <div className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                            Signing in...
                        </>
                    ) : (
                        <>
                            Sign In
                            <ArrowRight className="h-4 w-4" />
                        </>
                    )}
                </button>

                {/* Sign up link */}
                <p className="text-center text-sm text-zinc-600">
                    Don&apos;t have an account?{" "}
                    <button
                        type="button"
                        onClick={() => router.push("/auth")}
                        className="font-semibold text-pink-600 hover:text-pink-700"
                    >
                        Create one
                    </button>
                </p>
            </form>

            {/* Security info */}
            <div className="rounded-lg bg-blue-50 border border-blue-200 p-4 text-xs text-blue-900 space-y-2 mt-6">
                <p className="font-semibold">🔒 We Keep Your Data Safe</p>
                <ul className="space-y-1 ml-2">
                    <li>• Passwords are encrypted end-to-end</li>
                    <li>• We never share your email</li>
                    <li>• 2FA available in account settings</li>
                </ul>
            </div>
        </AuthShell>
    );
}