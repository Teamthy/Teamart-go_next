
/**
 * Auth Role Selector Page
 * Choose between Customer, Creator, or Merchant signup paths
 */

"use client";

import React from "react";
import { useRouter } from "next/navigation";
import PremiumAuthLayout from "@/components/auth/PremiumAuthLayout";
import RoleSelectorCard from "@/components/auth/RoleSelectorCard";
import { ShoppingBag, Sparkles, Store } from "lucide-react";

export default function AuthPage() {
    const router = useRouter();

    const roles = [
        {
            id: "customer",
            title: "Buyer",
            description: "Discover products and shop from creators",
            benefits: [
                "Browse curated product collections",
                "Follow favorite creators",
                "Get personalized recommendations",
                "Shop securely with buyer protection",
            ],
            icon: <ShoppingBag className="h-8 w-8" />,
        },
        {
            id: "creator",
            title: "Creator",
            description: "Build your audience and monetize your content",
            benefits: [
                "Go live and sell products in real-time",
                "Earn from live streams and sales",
                "Access creator analytics",
                "Get dedicated support",
            ],
            icon: <Sparkles className="h-8 w-8" />,
        },
        {
            id: "merchant",
            title: "Seller",
            description: "Reach millions of buyers and grow your business",
            benefits: [
                "Unlimited product listings",
                "Advanced seller analytics",
                "Fulfillment support",
                "Marketing tools included",
            ],
            icon: <Store className="h-8 w-8" />,
        },
    ];

    return (
        <PremiumAuthLayout
            title="Choose Your Path"
            description="Join Teamart as a buyer, creator, or seller. You can always add more roles later."
            variant="wide"
        >
            <div className="grid gap-6 sm:grid-cols-1 lg:grid-cols-3">
                {roles.map((role) => (
                    <RoleSelectorCard
                        key={role.id}
                        role={role.id as "customer" | "creator" | "merchant"}
                        title={role.title}
                        description={role.description}
                        benefits={role.benefits}
                        icon={role.icon}
                        href={`/auth/signup?role=${role.id}`}
                    />
                ))}
            </div>

            {/* Login link */}
            <div className="mt-8 text-center">
                <p className="text-sm text-zinc-600">
                    Already have an account?{" "}
                    <button
                        onClick={() => router.push("/auth/login")}
                        className="font-semibold text-pink-600 hover:text-pink-700"
                    >
                        Sign in
                    </button>
                </p>
            </div>
        </PremiumAuthLayout>
    );
}
