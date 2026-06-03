/**
 * OTPInput Component
 * Premium 6-digit OTP input with auto-advance, paste support, and visual feedback
 */

"use client";

import React, { useEffect } from "react";
import { Loader2, AlertCircle, CheckCircle2 } from "lucide-react";
import { useOTP } from "@/hooks/useOTP";

interface OTPInputProps {
    onComplete: (code: string) => void;
    isLoading?: boolean;
    error?: string | null;
    length?: number;
    autoFocus?: boolean;
}

export default function OTPInput({
    onComplete,
    isLoading = false,
    error,
    length = 6,
    autoFocus = true,
}: OTPInputProps) {
    const { values, handleInputChange, handleKeyDown, handlePaste, handleFocus, isFull, clear, inputRefs } = useOTP({
        length,
        onComplete,
    });

    // Auto-focus first input
    useEffect(() => {
        if (autoFocus) {
            setTimeout(() => inputRefs.current[0]?.focus(), 100);
        }
    }, [autoFocus, inputRefs]);

    return (
        <div className="space-y-6">
            {/* OTP Input Fields */}
            <div className="flex justify-center gap-2 sm:gap-3">
                {values.map((value, index) => (
                    <input
                        key={index}
                        ref={(el) => {
                            inputRefs.current[index] = el;
                        }}
                        type="text"
                        inputMode="numeric"
                        maxLength={1}
                        value={value}
                        onChange={(e) => handleInputChange(index, e.target.value)}
                        onKeyDown={(e) => handleKeyDown(index, e)}
                        onPaste={handlePaste}
                        onFocus={() => handleFocus(index)}
                        disabled={isLoading}
                        className={`h-14 w-12 rounded-lg border-2 text-center text-2xl font-bold font-mono transition-all sm:h-16 sm:w-14 ${error
                                ? "border-red-300 bg-red-50 text-red-900 focus:border-red-400 focus:outline-none focus:ring-2 focus:ring-red-200"
                                : isFull && !isLoading
                                    ? "border-green-300 bg-green-50 text-green-900 focus:border-green-400 focus:outline-none focus:ring-2 focus:ring-green-200"
                                    : value
                                        ? "border-zinc-300 bg-white text-zinc-900 focus:border-pink-500 focus:outline-none focus:ring-2 focus:ring-pink-100"
                                        : "border-zinc-200 bg-zinc-50 text-zinc-900 placeholder-zinc-400 focus:border-pink-500 focus:outline-none focus:ring-2 focus:ring-pink-100"
                            }`}
                        placeholder="–"
                        aria-label={`Digit ${index + 1}`}
                    />
                ))}
            </div>

            {/* Status messages */}
            <div className="flex items-center justify-center gap-2 text-sm">
                {isLoading ? (
                    <>
                        <Loader2 className="h-4 w-4 animate-spin text-pink-500" />
                        <span className="text-zinc-600">Verifying code...</span>
                    </>
                ) : isFull && !error ? (
                    <>
                        <CheckCircle2 className="h-4 w-4 text-green-500" />
                        <span className="text-green-600">Code ready!</span>
                    </>
                ) : error ? (
                    <>
                        <AlertCircle className="h-4 w-4 text-red-500" />
                        <span className="text-red-600">{error}</span>
                    </>
                ) : (
                    <span className="text-zinc-500">Enter 6-digit code</span>
                )}
            </div>

            {/* Helper text */}
            <div className="space-y-2 text-center text-xs text-zinc-500">
                <p>💡 Tip: You can paste your code directly</p>
                {!isFull && !isLoading && (
                    <p>
                        Can't find it?{" "}
                        <button
                            type="button"
                            onClick={() => clear()}
                            className="font-medium text-pink-600 hover:text-pink-700"
                        >
                            Clear and try again
                        </button>
                    </p>
                )}
            </div>
        </div>
    );
}
