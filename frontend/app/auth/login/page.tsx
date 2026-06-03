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
import PremiumAuthLayout from "@/components/auth/PremiumAuthLayout";
import { Eye, EyeOff, AlertCircle, ArrowRight } from "lucide-react";

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
            } catch (error) {
                // Error handling is done by the store
            }
        },
    });

    return (
        <PremiumAuthLayout
            title="Welcome Back"
            description="Sign in to continue shopping and explore new products."
            variant="compact"
            showBackButton
            backHref="/auth"
        >
            <form onSubmit={handleSubmit(onSubmitHandler)} className="space-y-5">
                {/* Email field */}
                <div>
                    <label htmlFor="email" className="block text-sm font-semibold text-zinc-900 mb-2">
                        Email Address
                        <span className="text-red-500">*</span>
                    </label>
                    <Controller
                        name="email"
                        control={control}
                        render={({ field }) => (
                            <input
                                {...field}
                                type="email"
                                id="email"
                                placeholder="you@example.com"
                                className={`w-full rounded-lg border-2 px-4 py-3 transition-all focus:outline-none ${errors.email
                                        ? "border-red-300 bg-red-50 focus:border-red-400 focus:ring-2 focus:ring-red-100"
                                        : "border-zinc-300 focus:border-pink-500 focus:ring-2 focus:ring-pink-100"
                                    }`}
                                aria-invalid={errors.email ? "true" : "false"}
                            />
                        )}
                    />
                    {errors.email && (
                        <div className="mt-2 flex items-start gap-2 text-sm text-red-600">
                            <AlertCircle className="h-4 w-4 mt-0.5 flex-shrink-0" />
                            <span>{errors.email?.message || "Invalid email"}</span>
                        </div>
                    )}
                </div>

                {/* Password field */}
                <div>
                    <div className="flex items-center justify-between mb-2">
                        <label htmlFor="password" className="text-sm font-semibold text-zinc-900">
                            Password
                            <span className="text-red-500">*</span>
                        </label>
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
                    {errors.password && (
                        <div className="mt-2 flex items-start gap-2 text-sm text-red-600">
                            <AlertCircle className="h-4 w-4 mt-0.5 flex-shrink-0" />
                            <span>{errors.password?.message || "Invalid password"}</span>
                        </div>
                    )}
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
                    Don't have an account?{" "}
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
        </PremiumAuthLayout>
    );
}