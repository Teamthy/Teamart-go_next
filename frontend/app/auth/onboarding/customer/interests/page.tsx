/**
 * Customer Onboarding - Interests Page
 * Select interests to personalize recommendations
 */

"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/useAuthStore";
import PremiumAuthLayout from "@/components/auth/PremiumAuthLayout";
import ProgressStepper from "@/components/auth/ProgressStepper";
import { Check } from "lucide-react";

interface Interest {
    id: string;
    name: string;
    emoji: string;
    description: string;
}

const INTERESTS: Interest[] = [
    { id: "fashion", name: "Fashion", emoji: "👗", description: "Clothing & accessories" },
    { id: "beauty", name: "Beauty", emoji: "💄", description: "Makeup & skincare" },
    { id: "home", name: "Home", emoji: "🏡", description: "Furniture & decor" },
    { id: "electronics", name: "Electronics", emoji: "📱", description: "Tech & gadgets" },
    { id: "sports", name: "Sports", emoji: "⚽", description: "Athletic gear" },
    { id: "books", name: "Books", emoji: "📚", description: "Reading & education" },
    { id: "gaming", name: "Gaming", emoji: "🎮", description: "Games & esports" },
    { id: "fitness", name: "Fitness", emoji: "💪", description: "Wellness & gym" },
    { id: "food", name: "Food", emoji: "🍕", description: "Cooking & snacks" },
    { id: "travel", name: "Travel", emoji: "✈️", description: "Trips & adventures" },
    { id: "music", name: "Music", emoji: "🎵", description: "Instruments & audio" },
    { id: "art", name: "Art", emoji: "🎨", description: "Crafts & supplies" },
];

export default function InterestsPage() {
    const router = useRouter();
    const { user } = useAuthStore();
    const [selectedInterests, setSelectedInterests] = useState<string[]>([]);
    const [isLoading, setIsLoading] = useState(false);

    const toggleInterest = (interestId: string) => {
        setSelectedInterests((prev) =>
            prev.includes(interestId)
                ? prev.filter((id) => id !== interestId)
                : [...prev, interestId]
        );
    };

    const handleContinue = async () => {
        if (selectedInterests.length === 0) return;

        setIsLoading(true);
        try {
            // Would save interests to backend
            await new Promise((resolve) => setTimeout(resolve, 500));
            router.push("/auth/onboarding/customer/follow-creators");
        } catch (error) {
            console.error("Failed to save interests:", error);
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

    return (
        <PremiumAuthLayout
            title="What Are You Interested In?"
            description="Select your interests to get personalized recommendations and discover creators you love."
            variant="compact"
            showBackButton
            backHref="/auth/onboarding/profile"
        >
            <div className="space-y-6">
                {/* Progress */}
                <ProgressStepper steps={steps} currentStep={1} variant="linear" showLabels={true} />

                {/* Interests Grid */}
                <div className="grid grid-cols-2 gap-3 sm:grid-cols-3">
                    {INTERESTS.map((interest) => {
                        const isSelected = selectedInterests.includes(interest.id);

                        return (
                            <button
                                key={interest.id}
                                onClick={() => toggleInterest(interest.id)}
                                className={`group relative overflow-hidden rounded-xl border-2 p-4 transition-all ${isSelected
                                        ? "border-pink-500 bg-gradient-to-br from-pink-50 to-pink-50 shadow-lg"
                                        : "border-zinc-200 bg-white hover:border-zinc-300 hover:shadow-md"
                                    }`}
                            >
                                {/* Selection indicator */}
                                {isSelected && (
                                    <div className="absolute top-2 right-2 h-6 w-6 rounded-full bg-pink-600 flex items-center justify-center">
                                        <Check className="h-4 w-4 text-white" />
                                    </div>
                                )}

                                {/* Content */}
                                <div className="space-y-2 text-center">
                                    <div className="text-3xl transition-transform group-hover:scale-110">
                                        {interest.emoji}
                                    </div>
                                    <div>
                                        <p className="font-semibold text-sm text-zinc-900">{interest.name}</p>
                                        <p className="text-xs text-zinc-600">{interest.description}</p>
                                    </div>
                                </div>
                            </button>
                        );
                    })}
                </div>

                {/* Selection counter */}
                <div className="text-center">
                    <p className="text-sm font-semibold text-zinc-900">
                        {selectedInterests.length} selected
                    </p>
                    <p className="text-xs text-zinc-600">
                        {selectedInterests.length === 0
                            ? "Select at least one interest"
                            : selectedInterests.length >= 5
                                ? "Great selection!"
                                : "Add more for better recommendations"}
                    </p>
                </div>

                {/* Continue button */}
                <button
                    onClick={handleContinue}
                    disabled={isLoading || selectedInterests.length === 0}
                    className="w-full rounded-lg bg-gradient-to-r from-pink-600 to-pink-500 px-4 py-3 font-semibold text-white transition-all hover:shadow-lg hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed disabled:scale-100 flex items-center justify-center gap-2"
                >
                    {isLoading ? (
                        <>
                            <div className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                            Saving interests...
                        </>
                    ) : (
                        "Continue to Follow Creators"
                    )}
                </button>

                {/* Skip option */}
                <button
                    onClick={() => router.push("/auth/onboarding/customer/follow-creators")}
                    className="w-full text-sm font-semibold text-zinc-600 hover:text-zinc-900 py-2"
                >
                    Skip for now
                </button>

                {/* Info box */}
                <div className="rounded-lg bg-blue-50 border border-blue-200 p-4 text-xs text-blue-900 space-y-2">
                    <p className="font-semibold">💡 Pro Tip</p>
                    <p>Your interests help us show you relevant products and creators. You can update them anytime in your settings.</p>
                </div>
            </div>
        </PremiumAuthLayout>
    );
}
