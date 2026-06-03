/**
 * Customer Onboarding - Follow Creators Page
 * Discover and follow recommended creators
 */

"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import PremiumAuthLayout from "@/components/auth/PremiumAuthLayout";
import ProgressStepper from "@/components/auth/ProgressStepper";
import { Users, Heart } from "lucide-react";

interface Creator {
    id: string;
    name: string;
    handle: string;
    category: string;
    followers: number;
    avatar: string;
    featured?: boolean;
}

const RECOMMENDED_CREATORS: Creator[] = [
    {
        id: "1",
        name: "Emma Chen",
        handle: "@emmachen",
        category: "Fashion",
        followers: 245600,
        avatar: "EC",
        featured: true,
    },
    {
        id: "2",
        name: "Alex Rivera",
        handle: "@alexrivera",
        category: "Fitness",
        followers: 189300,
        avatar: "AR",
    },
    {
        id: "3",
        name: "Sophia Lee",
        handle: "@sophialee",
        category: "Beauty",
        followers: 156800,
        avatar: "SL",
        featured: true,
    },
    {
        id: "4",
        name: "Marcus Johnson",
        handle: "@marcusjohnson",
        category: "Tech",
        followers: 198400,
        avatar: "MJ",
    },
    {
        id: "5",
        name: "Jessica Wang",
        handle: "@jessicawang",
        category: "Home Decor",
        followers: 112900,
        avatar: "JW",
    },
    {
        id: "6",
        name: "David Kim",
        handle: "@davidkim",
        category: "Gaming",
        followers: 267500,
        avatar: "DK",
        featured: true,
    },
];

export default function FollowCreatorsPage() {
    const router = useRouter();
    const [followedCreators, setFollowedCreators] = useState<string[]>([]);
    const [isLoading, setIsLoading] = useState(false);

    const toggleFollow = (creatorId: string) => {
        setFollowedCreators((prev) =>
            prev.includes(creatorId)
                ? prev.filter((id) => id !== creatorId)
                : [...prev, creatorId]
        );
    };

    const handleContinue = async () => {
        setIsLoading(true);
        try {
            // Would save follows to backend
            await new Promise((resolve) => setTimeout(resolve, 500));
            router.push("/auth/onboarding/customer/notifications");
        } catch (error) {
            console.error("Failed to save follows:", error);
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
            title="Follow Creators"
            description="Discover amazing creators and follow your favorites. You can find more anytime."
            variant="compact"
            showBackButton
            backHref="/auth/onboarding/customer/interests"
        >
            <div className="space-y-6">
                {/* Progress */}
                <ProgressStepper steps={steps} currentStep={2} variant="linear" showLabels={true} />

                {/* Creators Grid */}
                <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                    {RECOMMENDED_CREATORS.map((creator) => {
                        const isFollowed = followedCreators.includes(creator.id);

                        return (
                            <div
                                key={creator.id}
                                className={`group overflow-hidden rounded-xl border-2 p-4 transition-all ${isFollowed
                                        ? "border-pink-500 bg-gradient-to-br from-pink-50 to-pink-50 shadow-lg"
                                        : "border-zinc-200 bg-white hover:border-zinc-300 hover:shadow-md"
                                    } ${creator.featured ? "sm:col-span-2" : ""}`}
                            >
                                <div className="flex items-start justify-between gap-3">
                                    {/* Creator info */}
                                    <div className="flex gap-3 flex-1 min-w-0">
                                        {/* Avatar */}
                                        <div className="h-12 w-12 flex-shrink-0 rounded-full bg-gradient-to-br from-pink-400 to-pink-600 flex items-center justify-center text-white font-bold text-lg shadow-md">
                                            {creator.avatar}
                                        </div>

                                        {/* Details */}
                                        <div className="flex-1 min-w-0">
                                            <div className="flex items-start justify-between gap-2">
                                                <div className="min-w-0">
                                                    <p className="font-semibold text-zinc-900 truncate">{creator.name}</p>
                                                    <p className="text-xs text-zinc-600 truncate">{creator.handle}</p>
                                                </div>
                                                {creator.featured && (
                                                    <span className="text-xs font-bold bg-amber-100 text-amber-700 px-2 py-1 rounded-full flex-shrink-0 whitespace-nowrap">
                                                        ⭐ Trending
                                                    </span>
                                                )}
                                            </div>

                                            <div className="mt-2 space-y-1">
                                                <p className="text-xs text-zinc-600">{creator.category}</p>
                                                <div className="flex items-center gap-1 text-xs text-zinc-600">
                                                    <Users className="h-3 w-3" />
                                                    <span>{(creator.followers / 1000).toFixed(0)}K followers</span>
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                    {/* Follow button */}
                                    <button
                                        onClick={() => toggleFollow(creator.id)}
                                        className={`flex-shrink-0 px-3 py-2 rounded-lg font-semibold text-sm transition-all whitespace-nowrap ${isFollowed
                                                ? "bg-pink-600 text-white hover:bg-pink-700"
                                                : "border-2 border-pink-300 text-pink-600 hover:bg-pink-50"
                                            }`}
                                    >
                                        {isFollowed ? (
                                            <span className="flex items-center gap-1">
                                                <Heart className="h-4 w-4 fill-current" />
                                                Following
                                            </span>
                                        ) : (
                                            "Follow"
                                        )}
                                    </button>
                                </div>
                            </div>
                        );
                    })}
                </div>

                {/* Follow counter */}
                <div className="text-center">
                    <p className="text-sm font-semibold text-zinc-900">
                        {followedCreators.length} creator{followedCreators.length !== 1 ? "s" : ""} followed
                    </p>
                    <p className="text-xs text-zinc-600">
                        {followedCreators.length === 0
                            ? "Follow creators to get started"
                            : "You can follow more creators anytime"}
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
                            Saving follows...
                        </>
                    ) : (
                        "Continue to Notifications"
                    )}
                </button>

                {/* Skip option */}
                <button
                    onClick={() => router.push("/auth/onboarding/customer/notifications")}
                    className="w-full text-sm font-semibold text-zinc-600 hover:text-zinc-900 py-2"
                >
                    Skip for now
                </button>

                {/* Info box */}
                <div className="rounded-lg bg-blue-50 border border-blue-200 p-4 text-xs text-blue-900 space-y-2">
                    <p className="font-semibold">💡 Creator Tips</p>
                    <p>Following creators helps them grow and lets you get notified about new products and livestreams. Unfollow anytime from your profile.</p>
                </div>
            </div>
        </PremiumAuthLayout>
    );
}
