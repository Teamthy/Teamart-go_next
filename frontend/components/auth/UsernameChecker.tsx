/**
 * UsernameChecker Component
 * Real-time username validation with suggestions
 */

"use client";

import React, { useState } from "react";
import { Loader2, AlertCircle, CheckCircle2 } from "lucide-react";
import { useUsername } from "@/hooks/useUsername";

interface UsernameCheckerProps {
    onUsernameChange?: (username: string) => void;
    onAvailabilityChange?: (available: boolean) => void;
    autoFocus?: boolean;
}

export default function UsernameChecker({
    onUsernameChange,
    onAvailabilityChange,
    autoFocus = false,
}: UsernameCheckerProps) {
    const { username, setUsername, isAvailable, isChecking, error, suggestions, isValid } = useUsername({
        onCheck: onAvailabilityChange,
    });
    const [copied, setCopied] = useState(false);

    const handleUsernameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value;
        setUsername(value);
        onUsernameChange?.(value);
    };

    const copySuggestion = (suggestion: string) => {
        setUsername(suggestion);
        onUsernameChange?.(suggestion);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    return (
        <div className="space-y-4">
            {/* Username input */}
            <div className="relative">
                <label htmlFor="username" className="mb-2 block text-sm font-semibold text-zinc-900">
                    Username
                    <span className="text-red-500">*</span>
                </label>
                <div className="relative">
                    <input
                        id="username"
                        type="text"
                        value={username}
                        onChange={handleUsernameChange}
                        placeholder="Enter your username"
                        autoFocus={autoFocus}
                        className={`w-full rounded-lg border-2 px-4 py-3 pr-10 font-mono text-sm transition-all focus:outline-none ${error || !isValid && username
                                ? "border-red-300 bg-red-50 focus:border-red-400 focus:ring-2 focus:ring-red-100"
                                : isAvailable === true
                                    ? "border-green-300 bg-green-50 focus:border-green-400 focus:ring-2 focus:ring-green-100"
                                    : "border-zinc-300 focus:border-pink-500 focus:ring-2 focus:ring-pink-100"
                            }`}
                    />

                    {/* Status icon */}
                    {isChecking ? (
                        <Loader2 className="absolute right-3 top-1/2 h-5 w-5 -translate-y-1/2 animate-spin text-blue-500" />
                    ) : isAvailable === true ? (
                        <CheckCircle2 className="absolute right-3 top-1/2 h-5 w-5 -translate-y-1/2 text-green-500" />
                    ) : isAvailable === false ? (
                        <AlertCircle className="absolute right-3 top-1/2 h-5 w-5 -translate-y-1/2 text-red-500" />
                    ) : null}
                </div>
            </div>

            {/* Status messages */}
            <div className="flex items-start gap-2">
                {isChecking ? (
                    <p className="text-xs text-blue-600 font-medium">Checking availability...</p>
                ) : isAvailable === true ? (
                    <p className="text-xs text-green-600 font-medium">✓ Username available!</p>
                ) : error ? (
                    <div className="text-xs text-red-600 space-y-1">
                        <p className="font-medium">{error}</p>
                    </div>
                ) : !username ? (
                    <p className="text-xs text-zinc-600">3-30 characters, letters, numbers, and underscores only</p>
                ) : null}
            </div>

            {/* Suggestions */}
            {suggestions.length > 0 && isAvailable === false && (
                <div className="rounded-lg border border-amber-200 bg-amber-50 p-4">
                    <p className="mb-3 text-sm font-semibold text-amber-900">Try these instead:</p>
                    <div className="space-y-2">
                        {suggestions.map((suggestion) => (
                            <button
                                key={suggestion}
                                onClick={() => copySuggestion(suggestion)}
                                className="w-full rounded-lg border border-amber-200 bg-white px-3 py-2 text-left text-sm font-mono text-amber-900 transition-all hover:bg-amber-100 hover:border-amber-300"
                            >
                                <div className="flex items-center justify-between">
                                    <span>@{suggestion}</span>
                                    <span className="text-xs text-amber-600">
                                        {copied && suggestion === username ? "✓ Copied" : "Click to use"}
                                    </span>
                                </div>
                            </button>
                        ))}
                    </div>
                </div>
            )}

            {/* Format helper */}
            {username && !isValid && !error && (
                <div className="rounded-lg bg-zinc-100 p-3">
                    <p className="text-xs text-zinc-700">
                        <strong>Username format:</strong>
                        <br />
                        • Start with a letter<br />
                        • Use only letters, numbers, and underscores<br />
                        • 3-30 characters
                    </p>
                </div>
            )}
        </div>
    );
}
