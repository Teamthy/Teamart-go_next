/**
 * Customer Onboarding - Complete Page
 * Celebration screen and final redirect to feed
 */

"use client";

import React, { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/useAuthStore";
import PremiumAuthLayout from "@/components/auth/PremiumAuthLayout";
import { Sparkles, ArrowRight } from "lucide-react";

export default function CompletePage() {
    const router = useRouter();
    const { user } = useAuthStore();

    // Auto-redirect to feed after 3 seconds
    useEffect(() => {
        const timer = setTimeout(() => {
            router.push("/feed");
        }, 3000);

        return () => clearTimeout(timer);
    }, [router]);

    const handleGoToFeed = () => {
        router.push("/feed");
    };

    return (
        <PremiumAuthLayout
            title="Welcome to Teamart! 🎉"
            description={`You're all set, ${user?.firstName || "friend"}!`}
            variant="compact"
        >
            <div className="space-y-8">
                {/* Celebration animation */}
                <div className="relative h-32 flex items-center justify-center">
                    {/* Animated confetti-like elements */}
                    {[...Array(6)].map((_, i) => (
                        <div
                            key={i}
                            className="absolute h-2 w-2 rounded-full bg-pink-500 animate-bounce"
                            style={{
                                left: `${20 + i * 15}%`,
                                animationDelay: `${i * 0.1}s`,
                            }}
                        />
                    ))}

                    {/* Main celebration icon */}
                    <div className="relative z-10 h-24 w-24 rounded-full bg-gradient-to-br from-pink-400 to-pink-600 flex items-center justify-center shadow-2xl">
                        <Sparkles className="h-12 w-12 text-white animate-spin" />
                    </div>
                </div>

                {/* Completion message */}
                <div className="text-center space-y-3">
                    <h2 className="text-2xl font-bold text-zinc-900">
                        You're Ready to Shop!
                    </h2>
                    <p className="text-zinc-600 max-w-sm mx-auto">
                        Your profile is complete and personalized. Start discovering amazing products and creators on your feed.
                    </p>
                </div>

                {/* Benefits summary */}
                <div className="space-y-2">
                    {[
                        "✨ Personalized product recommendations",
                        "🎬 Exclusive livestream access",
                        "❤️ Following your favorite creators",
                        "🔔 Stay updated with notifications",
                    ].map((benefit, i) => (
                        <div key={i} className="flex items-center gap-3 text-sm text-zinc-700 bg-gradient-to-r from-pink-50 to-transparent p-3 rounded-lg">
                            <span className="text-lg">{benefit.split(" ")[0]}</span>
                            <span>{benefit.substring(2)}</span>
                        </div>
                    ))}
                </div>

                {/* Next steps */}
                <div className="space-y-3">
                    {/* Primary button */}
                    <button
                        onClick={handleGoToFeed}
                        className="w-full rounded-lg bg-gradient-to-r from-pink-600 to-pink-500 px-4 py-4 font-bold text-white text-lg transition-all hover:shadow-2xl hover:scale-105 flex items-center justify-center gap-2"
                    >
                        Start Shopping
                        <ArrowRight className="h-5 w-5" />
                    </button>

                    {/* Auto-redirect message */}
                    <p className="text-xs text-center text-zinc-600 animate-pulse">
                        Redirecting to your feed in 3 seconds...
                    </p>
                </div>

                {/* Quick tips */}
                <div className="rounded-xl bg-blue-50 border border-blue-200 p-5 space-y-3">
                    <p className="font-bold text-blue-900 text-sm">💡 Quick Tips to Get Started</p>
                    <ul className="space-y-2 text-xs text-blue-800">
                        <li className="flex gap-2">
                            <span className="flex-shrink-0">→</span>
                            <span>Browse your personalized feed of products</span>
                        </li>
                        <li className="flex gap-2">
                            <span className="flex-shrink-0">→</span>
                            <span>Follow more creators to customize your content</span>
                        </li>
                        <li className="flex gap-2">
                            <span className="flex-shrink-0">→</span>
                            <span>Join livestreams to see products in action</span>
                        </li>
                        <li className="flex gap-2">
                            <span className="flex-shrink-0">→</span>
                            <span>Check out your cart to keep track of items</span>
                        </li>
                    </ul>
                </div>

                {/* Footer message */}
                <p className="text-xs text-center text-zinc-600">
                    Questions? Visit our <a href="/help" className="font-semibold text-pink-600 hover:underline">Help Center</a> anytime.
                </p>
            </div>
        </PremiumAuthLayout>
    );
}
