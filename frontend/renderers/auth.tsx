import Link from "next/link";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import { renderHero, SummaryList, titleCase } from "./common";

export function renderAuth(slug: string[]) {
    const route = slug[0] ?? "auth";
    const second = slug[1];
    const third = slug[2];

    const titles: Record<string, string> = {
        "verify-email": "Verify your email",
        "email-sent": "Email sent",
        "social/google": "Continue with Google",
        "social/facebook": "Continue with Facebook",
        "social/callback": "Social sign-in callback",
        "invite": "Accept your invite",
        "complete-profile": "Complete your profile",
        "logout": "You’ve been signed out",
    };

    const descriptions: Record<string, string> = {
        "verify-email": "Confirm your inbox and unlock the full creator and merchant experience.",
        "email-sent": "Check your inbox for the link that will continue your setup.",
        "social/google": "Finish the flow with your Google account in a few taps.",
        "social/facebook": "Use your Facebook profile to keep the onboarding experience moving.",
        "social/callback": "Your social account is being connected securely.",
        "invite": "Review the invitation details and confirm your place in the team.",
        "complete-profile": "Add the final profile details so your storefront is ready for shoppers.",
        "logout": "Your session is closed and your account is safely signed out.",
    };

    const resolvedRoute = route === "auth" && second === "social" && third ? `social/${third}` : route === "auth" ? second ?? "complete-profile" : route;

    return (
        <div className="space-y-8 pb-10">
            {renderHero({
                title: titles[resolvedRoute] ?? titleCase(resolvedRoute),
                description: descriptions[resolvedRoute] ?? "Continue your journey with a polished Teamart flow.",
                badge: "Auth",
            })}
            <Card className="p-5 sm:p-6">
                <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                    <div>
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Review</p>
                        <h2 className="mt-3 text-xl font-semibold text-zinc-900">A clean, guided experience</h2>
                        <p className="mt-3 text-sm leading-7 text-zinc-600">
                            This page is styled to match the live shopping product system while keeping the user flow clear and action-focused.
                        </p>
                        <div className="mt-4 flex flex-wrap gap-3">
                            <Button asChild variant="primary">
                                <Link href="/auth/login">Back to login</Link>
                            </Button>
                            <Button asChild variant="secondary">
                                <Link href="/">Go home</Link>
                            </Button>
                        </div>
                    </div>
                    <SummaryList items={["Secure account steps", "Creator-ready onboarding", "Fast support handoff"]} />
                </div>
            </Card>
        </div>
    );
}
