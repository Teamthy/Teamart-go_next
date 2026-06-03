/**
 * Customer Onboarding - Notifications Page
 * Configure notification preferences
 */

"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import PremiumAuthLayout from "@/components/auth/PremiumAuthLayout";
import ProgressStepper from "@/components/auth/ProgressStepper";
import { Bell, MessageSquare, Heart, ShoppingBag, Zap } from "lucide-react";

interface NotificationPreference {
    id: string;
    label: string;
    description: string;
    icon: React.ReactNode;
    enabled: boolean;
}

const DEFAULT_PREFERENCES: NotificationPreference[] = [
    {
        id: "livestreams",
        label: "Livestream Notifications",
        description: "Get notified when creators you follow go live",
        icon: <Zap className="h-5 w-5" />,
        enabled: true,
    },
    {
        id: "products",
        label: "New Products",
        description: "Be first to know when creators launch new drops",
        icon: <ShoppingBag className="h-5 w-5" />,
        enabled: true,
    },
    {
        id: "messages",
        label: "Direct Messages",
        description: "Get notified about new messages from creators",
        icon: <MessageSquare className="h-5 w-5" />,
        enabled: true,
    },
    {
        id: "likes",
        label: "Likes & Comments",
        description: "Get notified when others engage with your posts",
        icon: <Heart className="h-5 w-5" />,
        enabled: true,
    },
    {
        id: "recommendations",
        label: "Personalized Recommendations",
        description: "Discover products based on your interests",
        icon: <Bell className="h-5 w-5" />,
        enabled: false,
    },
];

export default function NotificationsPage() {
    const router = useRouter();
    const [preferences, setPreferences] = useState<NotificationPreference[]>(DEFAULT_PREFERENCES);
    const [isLoading, setIsLoading] = useState(false);

    const togglePreference = (id: string) => {
        setPreferences((prev) =>
            prev.map((pref) =>
                pref.id === id ? { ...pref, enabled: !pref.enabled } : pref
            )
        );
    };

    const handleContinue = async () => {
        setIsLoading(true);
        try {
            // Would save preferences to backend
            await new Promise((resolve) => setTimeout(resolve, 500));
            router.push("/auth/onboarding/customer/complete");
        } catch (error) {
            console.error("Failed to save preferences:", error);
        } finally {
            setIsLoading(false);
        }
    };

    const steps = [
        { id: "interests", label: "Your Interests", description: "What do you like?" },
        { id: "creators", label: "Follow Creators", description: "Discover creators" },
        { id: "notifications", label: "Notifications", description: "Stay updated" },
        { id: "complete", label: "Complete", description: "All set!" },
    ];

    const enabledCount = preferences.filter((p) => p.enabled).length;

    return (
        <PremiumAuthLayout
            title="Stay Updated"
            description="Choose how you'd like to be notified. You can change this anytime in your settings."
            variant="compact"
            showBackButton
            backHref="/auth/onboarding/customer/follow-creators"
        >
            <div className="space-y-6">
                {/* Progress */}
                <ProgressStepper steps={steps} currentStep={3} variant="linear" showLabels={true} />

                {/* Preference toggles */}
                <div className="space-y-3">
                    {preferences.map((preference) => (
                        <button
                            key={preference.id}
                            onClick={() => togglePreference(preference.id)}
                            className={`w-full flex items-start gap-3 p-4 rounded-xl border-2 transition-all text-left ${preference.enabled
                                    ? "border-pink-500 bg-gradient-to-r from-pink-50 to-pink-50"
                                    : "border-zinc-200 bg-white hover:border-zinc-300"
                                }`}
                        >
                            {/* Checkbox */}
                            <div
                                className={`mt-0.5 h-5 w-5 flex-shrink-0 rounded-full border-2 flex items-center justify-center transition-all ${preference.enabled
                                        ? "border-pink-600 bg-pink-600"
                                        : "border-zinc-300 bg-white"
                                    }`}
                            >
                                {preference.enabled && (
                                    <svg
                                        className="h-3 w-3 text-white"
                                        fill="none"
                                        stroke="currentColor"
                                        viewBox="0 0 24 24"
                                    >
                                        <path
                                            strokeLinecap="round"
                                            strokeLinejoin="round"
                                            strokeWidth={3}
                                            d="M5 13l4 4L19 7"
                                        />
                                    </svg>
                                )}
                            </div>

                            {/* Content */}
                            <div className="flex-1 min-w-0">
                                <div className="flex items-start justify-between gap-2 mb-1">
                                    <div className="flex items-center gap-2 min-w-0">
                                        <span
                                            className={`text-lg ${preference.enabled ? "text-pink-600" : "text-zinc-400"
                                                }`}
                                        >
                                            {preference.icon}
                                        </span>
                                        <p className="font-semibold text-zinc-900">{preference.label}</p>
                                    </div>
                                </div>
                                <p className="text-xs text-zinc-600">{preference.description}</p>
                            </div>

                            {/* Status badge */}
                            <div className="flex-shrink-0 mt-0.5">
                                <span
                                    className={`text-xs font-bold px-2 py-1 rounded-full ${preference.enabled
                                            ? "bg-pink-100 text-pink-700"
                                            : "bg-zinc-100 text-zinc-600"
                                        }`}
                                >
                                    {preference.enabled ? "On" : "Off"}
                                </span>
                            </div>
                        </button>
                    ))}
                </div>

                {/* Summary */}
                <div className="text-center">
                    <p className="text-sm font-semibold text-zinc-900">
                        {enabledCount} notification{enabledCount !== 1 ? "s" : ""} enabled
                    </p>
                    <p className="text-xs text-zinc-600 mt-1">
                        Customize these settings anytime from your account
                    </p>
                </div>

                {/* Continue button */}
                <button
                    onClick={handleContinue}
                    disabled={isLoading}
                    className="w-full rounded-lg bg-gradient-to-r from-pink-600 to-pink-500 px-4 py-3 font-semibold text-white transition-all hover:shadow-lg hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed disabled:scale-100 flex items-center justify-center gap-2"
                >
                    {isLoading ? (
                        <>
                            <div className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                            Saving preferences...
                        </>
                    ) : (
                        "Finish Setup"
                    )}
                </button>

                {/* Info box */}
                <div className="rounded-lg bg-green-50 border border-green-200 p-4 text-xs text-green-900 space-y-2">
                    <p className="font-semibold">✨ Almost there!</p>
                    <p>You're just one step away from exploring Teamart. Click "Finish Setup" to get started!</p>
                </div>
            </div>
        </PremiumAuthLayout>
    );
}
