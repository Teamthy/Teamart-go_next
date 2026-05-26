"use client";

import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import PageHeader from "@/components/ui/PageHeader";
import ProgressIndicator from "@/components/ui/ProgressIndicator";
import OTPInput from "@/components/ui/OTPInput";
import RoleCard from "@/components/ui/RoleCard";
import StatCard from "@/components/ui/StatCard";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import { clearCustomerDraft, saveCustomer, saveCustomerDraft, setWorkspaceRole } from "@/lib/auth-state";

const progressSteps = ["Choose your path", "Verify your email", "Save your profile"];

const roleOptions = [
    {
        key: "customer" as const,
        title: "Customer",
        description: "Shop products, follow creators, and move quickly into checkout and live-room experiences.",
        requirements: ["A valid email address", "A first and last name", "A favorite category to personalize your feed"],
        ctaLabel: "Shop as a customer",
        href: "/feed",
        tone: "default" as const,
    },
    {
        key: "creator" as const,
        title: "Creator",
        description: "Build campaign-ready content, plan livestreams, and turn your audience into buyers.",
        requirements: ["Your creator name", "A social handle", "A niche and an audience-ready workflow"],
        ctaLabel: "Create as a creator",
        href: "/creator",
        tone: "warning" as const,
    },
    {
        key: "merchant" as const,
        title: "Merchant",
        description: "Manage storefronts, inventory, payouts, and promotions in a merchant-first workspace.",
        requirements: ["A store name", "Owner details", "Inventory and payout preferences"],
        ctaLabel: "Sell as a merchant",
        href: "/merchant",
        tone: "success" as const,
    },
];

export default function OnboardingPage() {
    const router = useRouter();
    const [selectedRole, setSelectedRole] = useState<(typeof roleOptions)[number]["key"]>("customer");
    const [firstName, setFirstName] = useState("");
    const [lastName, setLastName] = useState("");
    const [email, setEmail] = useState("");
    const [favoriteCategory, setFavoriteCategory] = useState("Fashion");
    const [otp, setOtp] = useState("");
    const [step, setStep] = useState(1);

    const currentStep = useMemo(() => Math.min(step, 3), [step]);

    useEffect(() => {
        saveCustomerDraft({
            firstName,
            lastName,
            email,
            favoriteCategory,
            role: selectedRole,
            verified: otp.length === 6,
        });
    }, [email, favoriteCategory, firstName, lastName, otp.length, selectedRole]);

    const handleContinue = () => {
        if (currentStep < 3) {
            setStep((value) => value + 1);
            return;
        }

        saveCustomer(
            {
                id: `${email.toLowerCase().replace(/[^a-z0-9]/g, "") || "customer"}`,
                email,
                firstName,
                lastName,
                verified: otp.length === 6,
                favoriteCategory,
                roles: {
                    customer: true,
                    creator: selectedRole === "creator",
                    merchant: selectedRole === "merchant",
                },
            },
            ["profile", "verified", selectedRole],
        );
        setWorkspaceRole(selectedRole);
        clearCustomerDraft();
        router.push(selectedRole === "merchant" ? "/merchant" : selectedRole === "creator" ? "/creator" : "/feed");
    };

    const isProfileComplete = Boolean(firstName && lastName && email && favoriteCategory);

    return (
        <div className="space-y-8 pb-10">
            <PageHeader
                title="Welcome to Teamart"
                description="Create your customer account, verify your email, and land in a personalized social commerce journey that fits your next move."
            />

            <div className="grid gap-4 md:grid-cols-3">
                <StatCard label="Onboarding progress" value={`${currentStep}/3`} helper="Stay in motion while you choose, verify, and save your profile." />
                <StatCard label="Customer flow" value="Live-ready" helper="Your feed and checkout path will feel tailored from day one." />
                <StatCard label="Role access" value={selectedRole} helper="Pick the route that best matches your experience and goals." />
            </div>

            <ProgressIndicator steps={progressSteps} currentStep={currentStep} />

            <div className="grid gap-4 xl:grid-cols-[0.9fr_1.1fr]">
                <Card className="p-5 sm:p-6">
                    {currentStep === 1 ? (
                        <div className="space-y-5">
                            <div>
                                <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Start here</p>
                                <h2 className="mt-2 text-xl font-semibold text-zinc-900">Choose the role you want to explore first</h2>
                            </div>
                            <div className="grid gap-4">
                                {roleOptions.map((role) => (
                                    <RoleCard
                                        key={role.key}
                                        title={role.title}
                                        description={role.description}
                                        requirements={role.requirements}
                                        ctaLabel={role.ctaLabel}
                                        href={role.key === selectedRole ? undefined : undefined}
                                        onClick={() => setSelectedRole(role.key)}
                                        tone={role.tone}
                                    />
                                ))}
                            </div>
                            <div className="rounded-[24px] bg-[#FFF8FB] p-4">
                                <p className="text-sm font-semibold text-zinc-900">Current role</p>
                                <p className="mt-2 text-sm text-zinc-600">You are setting up as <span className="font-semibold text-zinc-900">{selectedRole}</span>. Choose a role above to update your path.</p>
                            </div>
                        </div>
                    ) : currentStep === 2 ? (
                        <div className="space-y-5">
                            <div>
                                <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Verify your email</p>
                                <h2 className="mt-2 text-xl font-semibold text-zinc-900">Enter the 6-digit verification code</h2>
                                <p className="mt-2 text-sm leading-6 text-zinc-600">Use the OTP input below to confirm your email and unlock the personalized experience.</p>
                            </div>
                            <OTPInput value={otp} onChange={setOtp} />
                            <div className="rounded-[24px] bg-zinc-50 p-4 text-sm text-zinc-700">
                                We&apos;ll confirm your code in the next step and save your profile details.
                            </div>
                        </div>
                    ) : (
                        <div className="space-y-5">
                            <div>
                                <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Save your profile</p>
                                <h2 className="mt-2 text-xl font-semibold text-zinc-900">Complete your customer profile</h2>
                                <p className="mt-2 text-sm leading-6 text-zinc-600">These details personalize the feed, recommendations, and your first checkout flow.</p>
                            </div>
                            <div className="grid gap-4 sm:grid-cols-2">
                                <label className="text-sm font-semibold text-zinc-700">
                                    First name
                                    <input value={firstName} onChange={(event) => setFirstName(event.target.value)} className="mt-2 w-full rounded-[24px] border border-zinc-200 px-4 py-3 text-sm text-zinc-900" placeholder="Jordan" />
                                </label>
                                <label className="text-sm font-semibold text-zinc-700">
                                    Last name
                                    <input value={lastName} onChange={(event) => setLastName(event.target.value)} className="mt-2 w-full rounded-[24px] border border-zinc-200 px-4 py-3 text-sm text-zinc-900" placeholder="Lee" />
                                </label>
                            </div>
                            <label className="block text-sm font-semibold text-zinc-700">
                                Email address
                                <input type="email" value={email} onChange={(event) => setEmail(event.target.value)} className="mt-2 w-full rounded-[24px] border border-zinc-200 px-4 py-3 text-sm text-zinc-900" placeholder="you@example.com" />
                            </label>
                            <label className="block text-sm font-semibold text-zinc-700">
                                Favorite category
                                <input value={favoriteCategory} onChange={(event) => setFavoriteCategory(event.target.value)} className="mt-2 w-full rounded-[24px] border border-zinc-200 px-4 py-3 text-sm text-zinc-900" placeholder="Fashion" />
                            </label>
                            <div className="rounded-[24px] bg-[#FFF8FB] p-4 text-sm text-zinc-700">
                                Verified code: {otp || "pending"} • {isProfileComplete ? "Profile ready" : "Finish the form to continue"}
                            </div>
                        </div>
                    )}

                    <div className="mt-6 flex flex-wrap gap-3">
                        <Button variant="primary" onClick={handleContinue} disabled={currentStep === 3 && !isProfileComplete}>
                            {currentStep < 3 ? "Continue" : "Finish onboarding"}
                        </Button>
                        <Button variant="secondary" onClick={() => router.push("/")}>
                            Back home
                        </Button>
                    </div>
                </Card>

                <Card className="p-5 sm:p-6">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">What this unlocks</p>
                    <h2 className="mt-2 text-xl font-semibold text-zinc-900">A smooth path from onboarding to discovery</h2>
                    <div className="mt-4 space-y-3">
                        {[
                            "Personalized feed for product and creator discovery",
                            "Fast access to live rooms and checkout-ready bundles",
                            "Role-based navigation for customer, creator, and merchant paths",
                        ].map((item) => (
                            <div key={item} className="rounded-[24px] bg-zinc-50 px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                    <div className="mt-5 rounded-[24px] border border-dashed border-zinc-200 p-4">
                        <p className="text-sm font-semibold text-zinc-900">Role-specific entry points</p>
                        <div className="mt-3 flex flex-wrap gap-2">
                            <Link href="/auth/customer" className="rounded-full bg-[#FFF8FB] px-3 py-2 text-sm font-semibold text-zinc-900">Customer</Link>
                            <Link href="/auth/creator" className="rounded-full bg-[#FFF8FB] px-3 py-2 text-sm font-semibold text-zinc-900">Creator</Link>
                            <Link href="/auth/merchant" className="rounded-full bg-[#FFF8FB] px-3 py-2 text-sm font-semibold text-zinc-900">Merchant</Link>
                        </div>
                        <p className="mt-3 text-sm leading-6 text-zinc-600">Use the dedicated role routes when you want to follow a guided signup path tailored to your workspace.</p>
                    </div>
                </Card>
            </div>
        </div>
    );
}
