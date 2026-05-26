"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import Illustration from "@/components/social/Illustration";
import Badge from "@/components/ui/badge";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import Input from "@/components/ui/input";
import * as api from "@/lib/api";
import { getStoredCustomer, saveCustomer } from "@/lib/auth-state";

const onboardingSlides = [
    {
        ill: "onboard1" as const,
        title: "The best social media app of the century",
        subtitle: "Pick creator drops, shop in-stream, and follow favorite trends without leaving the feed.",
    },
    {
        ill: "onboard2" as const,
        title: "Connect with the world in one scroll",
        subtitle: "See live creators, product pins, and community comments shaped around what you actually want.",
    },
    {
        ill: "onboard3" as const,
        title: "Everything you can do in one app",
        subtitle: "Live shopping, creator storefronts, secure checkout, and personalized recommendations from one pink-powered shell.",
    },
];

type AuthVariant = "onboarding" | "login" | "register" | "forgot" | "reset" | "mfa" | "success";
type AuthRole = "customer" | "merchant" | "creator";

type RoleOption = {
    value: AuthRole;
    label: string;
    helper: string;
    description: string;
    badge: string;
};

const roleOptions: RoleOption[] = [
    {
        value: "customer",
        label: "Customer",
        helper: "Shop the feed, save favorites, and checkout faster.",
        description: "Use the feed to discover live products, creator picks, and quick checkout moments.",
        badge: "Shop",
    },
    {
        value: "merchant",
        label: "Merchant",
        helper: "Launch storefronts, manage orders, and keep inventory moving.",
        description: "Grow your storefront, publish new products, and keep your buyers engaged.",
        badge: "Sell",
    },
    {
        value: "creator",
        label: "Creator",
        helper: "Open your studio, post drops, and go live with your audience.",
        description: "Create live moments, share curator picks, and turn your audience into buyers.",
        badge: "Create",
    },
];

const defaultRole: AuthRole = "customer";

const signupRequirements: Record<AuthRole, { title: string; details: string[] }> = {
    customer: {
        title: "Customer requirements",
        details: ["Your full name", "A valid email address", "A secure password", "A favorite category"],
    },
    creator: {
        title: "Creator requirements",
        details: ["Your creator name", "A valid email address", "A secure password", "Your niche and social handle"],
    },
    merchant: {
        title: "Merchant requirements",
        details: ["Your store name", "Owner details", "A valid email address", "A secure password and store category"],
    },
};

function normalizeRole(value?: string | null): AuthRole {
    if (value === "merchant" || value === "creator") {
        return value;
    }

    return defaultRole;
}

function getStoredRole(): AuthRole {
    if (typeof window === "undefined") {
        return defaultRole;
    }

    try {
        const stored = localStorage.getItem("auth_role");
        return normalizeRole(stored);
    } catch {
        return defaultRole;
    }
}

function getRoleDestination(role: AuthRole) {
    if (role === "merchant") {
        return "/merchant";
    }

    if (role === "creator") {
        return "/creator";
    }

    return "/feed";
}

type SocialProvider = "google" | "apple" | "tiktok";

function SocialButton({ icon, label, onClick }: { icon: string; label: string; onClick?: () => void }) {
    return (
        <button
            type="button"
            onClick={onClick}
            className="flex flex-1 items-center justify-center gap-2 rounded-[24px] border border-zinc-200 bg-white px-4 py-3 text-sm font-semibold text-zinc-700"
        >
            <span>{icon}</span>
            <span>{label}</span>
        </button>
    );
}

export default function AuthTemplate({
    variant = "login",
    initialRole,
}: {
    variant?: AuthVariant;
    initialRole?: AuthRole;
}) {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [otp, setOtp] = useState("");
    const [remember, setRemember] = useState(false);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [step, setStep] = useState(0);
    const [resetSent, setResetSent] = useState(false);
    const [passwordUpdated, setPasswordUpdated] = useState(false);
    const [currentVariant, setCurrentVariant] = useState<AuthVariant>(variant);
    const [selectedRole, setSelectedRole] = useState<AuthRole>(defaultRole);
    const [signupStep, setSignupStep] = useState(1);
    const [customerFirstName, setCustomerFirstName] = useState("");
    const [customerLastName, setCustomerLastName] = useState("");
    const [customerFavoriteCategory, setCustomerFavoriteCategory] = useState("");
    const [creatorName, setCreatorName] = useState("");
    const [creatorNiche, setCreatorNiche] = useState("");
    const [creatorHandle, setCreatorHandle] = useState("");
    const [merchantOwnerName, setMerchantOwnerName] = useState("");
    const [merchantStoreName, setMerchantStoreName] = useState("");
    const [merchantCategory, setMerchantCategory] = useState("");
    const [merchantWebsite, setMerchantWebsite] = useState("");
    const router = useRouter();

    useEffect(() => {
        setCurrentVariant(variant);
    }, [variant]);

    useEffect(() => {
        if (typeof window === "undefined") {
            return;
        }

        const params = new URLSearchParams(window.location.search);
        const queryRole = normalizeRole(params.get("role"));
        const nextRole = initialRole ? normalizeRole(initialRole) : queryRole || getStoredRole();
        setSelectedRole(nextRole);
        localStorage.setItem("auth_role", nextRole);
    }, [initialRole]);

    useEffect(() => {
        if (currentVariant !== "register") {
            setSignupStep(1);
        }
    }, [currentVariant]);

    const slide = onboardingSlides[step];
    const roleMeta = roleOptions.find((option) => option.value === selectedRole) ?? roleOptions[0];

    const roleActionLabel = useMemo(() => {
        if (selectedRole === "customer") {
            return "Continue to customer signup";
        }

        if (selectedRole === "creator") {
            return "Continue to creator signup";
        }

        return "Continue to merchant signup";
    }, [selectedRole]);

    const roleHeading = useMemo(() => {
        if (selectedRole === "customer") {
            return "Create your customer account";
        }

        if (selectedRole === "creator") {
            return "Create your creator account";
        }

        return "Create your merchant account";
    }, [selectedRole]);

    const persistAuthSession = (payload: Record<string, any>, workspaceRole?: AuthRole) => {
        if (typeof window === "undefined") {
            return;
        }

        const resolvedRole = normalizeRole(workspaceRole ?? payload?.user?.role ?? selectedRole);
        const fallbackName = payload?.user?.name || payload?.email?.split("@")[0] || roleMeta.label.split(" ")[0] || "Customer";
        const [firstName, ...lastNameParts] = fallbackName.split(" ");
        const resolvedEmail = payload?.email || payload?.user?.email || email;
        const profile = {
            id: String(
                payload?.user_id ??
                payload?.user?.id ??
                (resolvedEmail.replace(/[^a-z0-9]/gi, "") || "customer")
            ),
            email: resolvedEmail,
            firstName: firstName || "Customer",
            lastName: lastNameParts.join(" ") || "User",
            verified: true,
            roles: {
                customer: true,
                creator: resolvedRole === "creator",
                merchant: resolvedRole === "merchant",
            },
            favoriteCategory: customerFavoriteCategory || "Fashion",
        };

        const nextUser = {
            ...(payload?.user ?? {
                id: payload?.user_id ?? 0,
                email: resolvedEmail,
                created_at: payload?.created_at,
            }),
            role: resolvedRole,
            name: payload?.user?.name || `${profile.firstName} ${profile.lastName}`.trim(),
        };

        saveCustomer(profile, ["auth", resolvedRole]);
        localStorage.setItem("user", JSON.stringify(nextUser));
        localStorage.setItem("auth_role", resolvedRole);
        localStorage.setItem("workspace_role", resolvedRole);
        localStorage.setItem("session", JSON.stringify({ ...payload, user: nextUser }));

        if (payload?.session_id) {
            localStorage.setItem("session_id", payload.session_id);
        }

        if (payload?.access_token) {
            localStorage.setItem("access_token", payload.access_token);
        }

        if (payload?.refresh_token) {
            localStorage.setItem("refresh_token", payload.refresh_token);
        }

        sessionStorage.setItem("session", JSON.stringify({ ...payload, user: nextUser }));
        sessionStorage.setItem("auth_role", resolvedRole);
        sessionStorage.setItem("workspace_role", resolvedRole);
    };

    const persistSocialAuth = (provider: SocialProvider) => {
        const providerLabels: Record<SocialProvider, { label: string; email: string; firstName: string; lastName: string }> = {
            google: { label: "Google", email: "google@teamart.social", firstName: "Google", lastName: "User" },
            apple: { label: "Apple", email: "apple@teamart.social", firstName: "Apple", lastName: "User" },
            tiktok: { label: "TikTok", email: "tiktok@teamart.social", firstName: "TikTok", lastName: "User" },
        };

        const selected = providerLabels[provider];
        saveCustomer(
            {
                id: `${provider}-user`,
                email: selected.email,
                firstName: selected.firstName,
                lastName: selected.lastName,
                verified: true,
                roles: {
                    customer: true,
                    creator: false,
                    merchant: false,
                },
                favoriteCategory: "Fashion",
                provider,
            },
            ["social", provider],
        );

        localStorage.setItem("user", JSON.stringify({
            id: `${provider}-user`,
            email: selected.email,
            role: "customer",
            name: `${selected.firstName} ${selected.lastName}`,
        }));
        localStorage.setItem("auth_role", "customer");
        localStorage.setItem("workspace_role", "customer");
        localStorage.setItem("social_provider", provider);
        sessionStorage.setItem("auth_role", "customer");
        sessionStorage.setItem("workspace_role", "customer");
        router.push("/feed");
    };

    const goToVariant = (nextVariant: AuthVariant) => {
        setCurrentVariant(nextVariant);
        setError(null);
        setResetSent(false);
        setPasswordUpdated(false);
    };

    const resetRegisterFlow = () => {
        setSignupStep(1);
        setPassword("");
        setConfirmPassword("");
        setError(null);
    };

    const heading =
        currentVariant === "onboarding"
            ? "Welcome to Teamart"
            : currentVariant === "forgot"
                ? "Reset your password"
                : currentVariant === "reset"
                    ? "Create a new password"
                    : currentVariant === "mfa"
                        ? "Verify your access"
                        : currentVariant === "success"
                            ? "All set"
                            : currentVariant === "register"
                                ? roleHeading
                                : "Welcome back";

    const subheading =
        currentVariant === "onboarding"
            ? slide.subtitle
            : currentVariant === "forgot"
                ? "Enter your email and we’ll send you a secure reset link to get back into your account."
                : currentVariant === "reset"
                    ? "Choose a strong password and confirm it so your account stays protected."
                    : currentVariant === "mfa"
                        ? "Enter the verification code from your authenticator or email to continue."
                        : currentVariant === "success"
                            ? "Your password was updated and your account is ready for the next drop."
                            : currentVariant === "register"
                                ? roleMeta.helper
                                : "Sign in to keep your storefront, orders, and live drops moving.";

    const signupName = useMemo(() => {
        if (selectedRole === "customer") {
            return `${customerFirstName} ${customerLastName}`.trim() || email.split("@")[0] || "Customer";
        }

        if (selectedRole === "creator") {
            return creatorName || creatorHandle || email.split("@")[0] || "Creator";
        }

        return merchantStoreName || merchantOwnerName || email.split("@")[0] || "Merchant";
    }, [selectedRole, customerFirstName, customerLastName, creatorName, creatorHandle, merchantStoreName, merchantOwnerName, email]);

    const validateStepOne = () => {
        if (selectedRole === "customer") {
            if (!customerFirstName.trim() || !customerLastName.trim()) {
                throw new Error("Add your first and last name to continue.");
            }

            if (!email.trim()) {
                throw new Error("Add your email address to continue.");
            }

            if (!customerFavoriteCategory.trim()) {
                throw new Error("Add your favorite category to continue.");
            }

            return;
        }

        if (selectedRole === "creator") {
            if (!creatorName.trim() || !creatorNiche.trim() || !creatorHandle.trim()) {
                throw new Error("Add your creator name, niche, and social handle to continue.");
            }

            if (!email.trim()) {
                throw new Error("Add your email address to continue.");
            }

            return;
        }

        if (!merchantOwnerName.trim() || !merchantStoreName.trim() || !merchantCategory.trim()) {
            throw new Error("Add your owner name, store name, and store category to continue.");
        }

        if (!email.trim()) {
            throw new Error("Add your email address to continue.");
        }

        if (!merchantWebsite.trim()) {
            throw new Error("Add your store website to continue.");
        }
    };

    const validateStepTwo = () => {
        if (!password || password !== confirmPassword) {
            throw new Error("Enter a matching password for your account.");
        }
    };

    const renderRegisterFields = () => {
        if (selectedRole === "customer") {
            return (
                <>
                    <Input
                        label="First name"
                        type="text"
                        placeholder="Jordan"
                        required
                        value={customerFirstName}
                        onChange={(event) => setCustomerFirstName(event.target.value)}
                    />
                    <Input
                        label="Last name"
                        type="text"
                        placeholder="Lee"
                        required
                        value={customerLastName}
                        onChange={(event) => setCustomerLastName(event.target.value)}
                    />
                    <Input
                        label="Email address"
                        type="email"
                        placeholder="you@example.com"
                        required
                        value={email}
                        onChange={(event) => setEmail(event.target.value)}
                    />
                    <Input
                        label="Favorite category"
                        type="text"
                        placeholder="Fashion, wellness, or home"
                        required
                        value={customerFavoriteCategory}
                        onChange={(event) => setCustomerFavoriteCategory(event.target.value)}
                    />
                </>
            );
        }

        if (selectedRole === "creator") {
            return (
                <>
                    <Input
                        label="Creator name"
                        type="text"
                        placeholder="Ava Studio"
                        required
                        value={creatorName}
                        onChange={(event) => setCreatorName(event.target.value)}
                    />
                    <Input
                        label="Email address"
                        type="email"
                        placeholder="you@example.com"
                        required
                        value={email}
                        onChange={(event) => setEmail(event.target.value)}
                    />
                    <Input
                        label="Primary niche"
                        type="text"
                        placeholder="Beauty, wellness, or home"
                        required
                        value={creatorNiche}
                        onChange={(event) => setCreatorNiche(event.target.value)}
                    />
                    <Input
                        label="TikTok / Instagram handle"
                        type="text"
                        placeholder="@yourbrand"
                        required
                        value={creatorHandle}
                        onChange={(event) => setCreatorHandle(event.target.value)}
                    />
                </>
            );
        }

        return (
            <>
                <Input
                    label="Owner name"
                    type="text"
                    placeholder="Alex Morgan"
                    required
                    value={merchantOwnerName}
                    onChange={(event) => setMerchantOwnerName(event.target.value)}
                />
                <Input
                    label="Store name"
                    type="text"
                    placeholder="Luna Market"
                    required
                    value={merchantStoreName}
                    onChange={(event) => setMerchantStoreName(event.target.value)}
                />
                <Input
                    label="Email address"
                    type="email"
                    placeholder="you@example.com"
                    required
                    value={email}
                    onChange={(event) => setEmail(event.target.value)}
                />
                <Input
                    label="Store category"
                    type="text"
                    placeholder="Fashion, home, wellness"
                    required
                    value={merchantCategory}
                    onChange={(event) => setMerchantCategory(event.target.value)}
                />
                <Input
                    label="Store website"
                    type="url"
                    placeholder="https://yourstore.com"
                    required
                    value={merchantWebsite}
                    onChange={(event) => setMerchantWebsite(event.target.value)}
                />
            </>
        );
    };

    const renderRegisterSummary = () => {
        if (selectedRole === "customer") {
            return (
                <div className="rounded-[24px] border border-zinc-200 bg-zinc-50 p-4 text-sm text-zinc-700">
                    <p className="font-semibold text-zinc-900">Customer summary</p>
                    <p className="mt-2">{customerFirstName} {customerLastName}</p>
                    <p>{email}</p>
                    <p className="mt-2">Favorite category: {customerFavoriteCategory || "Not set"}</p>
                </div>
            );
        }

        if (selectedRole === "creator") {
            return (
                <div className="rounded-[24px] border border-zinc-200 bg-zinc-50 p-4 text-sm text-zinc-700">
                    <p className="font-semibold text-zinc-900">Creator summary</p>
                    <p className="mt-2">{creatorName}</p>
                    <p>{email}</p>
                    <p className="mt-2">Niche: {creatorNiche}</p>
                    <p>Handle: {creatorHandle}</p>
                </div>
            );
        }

        return (
            <div className="rounded-[24px] border border-zinc-200 bg-zinc-50 p-4 text-sm text-zinc-700">
                <p className="font-semibold text-zinc-900">Merchant summary</p>
                <p className="mt-2">{merchantStoreName}</p>
                <p>{email}</p>
                <p className="mt-2">Owner: {merchantOwnerName}</p>
                <p>Category: {merchantCategory}</p>
                <p>Website: {merchantWebsite}</p>
            </div>
        );
    };

    const handleSubmit = async (event: React.FormEvent) => {
        event.preventDefault();
        setError(null);

        if (currentVariant === "forgot") {
            setResetSent(true);
            return;
        }

        if (currentVariant === "reset") {
            if (!password || password !== confirmPassword) {
                setError("Enter a matching password.");
                return;
            }

            setPasswordUpdated(true);
            return;
        }

        if (currentVariant === "register") {
            try {
                if (signupStep === 1) {
                    validateStepOne();
                    setSignupStep(2);
                    return;
                }

                if (signupStep === 2) {
                    validateStepTwo();
                    setSignupStep(3);
                    return;
                }

                const signupResponse = await api.signup(email, password);
                const loginResponse = await api.login(email, password);

                persistAuthSession(loginResponse, selectedRole);
                localStorage.setItem(
                    "onboarding_profile",
                    JSON.stringify({
                        role: selectedRole,
                        email,
                        name: signupName,
                        createdAt: new Date().toISOString(),
                    })
                );
                sessionStorage.setItem(
                    "signupResult",
                    JSON.stringify({
                        ...signupResponse,
                        user: {
                            ...signupResponse.user,
                            role: selectedRole,
                            name: signupName,
                        },
                    })
                );

                router.push(getRoleDestination(selectedRole));
                return;
            } catch (err: unknown) {
                setError(err instanceof Error ? err.message : "Request failed");
                return;
            }
        }

        if (currentVariant === "login") {
            if (selectedRole !== "customer" && !getStoredCustomer()) {
                setError("Create a customer account first to unlock creator or merchant access.");
                setSelectedRole("customer");
                resetRegisterFlow();
                goToVariant("register");
                return;
            }

            setLoading(true);

            try {
                const res = await api.login(email, password);
                persistAuthSession(res, selectedRole);
                sessionStorage.setItem("session", JSON.stringify(res));

                if (res.requires_mfa || res.requiresMFA) {
                    sessionStorage.setItem("pendingSession", JSON.stringify(res));
                    router.push("/auth/mfa");
                    return;
                }

                sessionStorage.removeItem("pendingSession");
                router.push(getRoleDestination(selectedRole));
            } catch (err: unknown) {
                setError(err instanceof Error ? err.message : "Request failed");
            } finally {
                setLoading(false);
            }

            return;
        }

        setLoading(true);

        try {
            if (currentVariant === "mfa") {
                const pending = sessionStorage.getItem("pendingSession");
                const sess = pending ? JSON.parse(pending) : null;
                const sessionId = sess?.session_id || sess?.sessionID;
                if (!sessionId) {
                    throw new Error("Missing pending session for MFA");
                }

                const response = await api.verifyOTP(sessionId, otp);
                persistAuthSession(response, selectedRole);
                sessionStorage.removeItem("pendingSession");
                router.push(getRoleDestination(selectedRole));
                return;
            }

            setError("Unsupported auth path.");
        } catch (err: unknown) {
            setError(err instanceof Error ? err.message : "Request failed");
        } finally {
            setLoading(false);
        }
    };

    if (currentVariant === "onboarding") {
        return (
            <div className="mx-auto max-w-[420px] px-4 py-6 sm:px-6">
                <Card className="p-5 sm:p-6">
                    <div className="flex items-center justify-between">
                        <Badge tone="default">Onboarding</Badge>
                        <span className="text-[11px] text-zinc-500">
                            {step + 1} / {onboardingSlides.length}
                        </span>
                    </div>
                    <div className="mt-4 flex justify-center">
                        <Illustration variant={slide.ill} />
                    </div>
                    <div className="mt-5 space-y-3 text-center">
                        <h1 className="text-[24px] font-semibold tracking-tight text-zinc-900">{slide.title}</h1>
                        <p className="text-sm leading-6 text-zinc-600">{slide.subtitle}</p>
                    </div>
                    <div className="mt-5 flex justify-center gap-2">
                        {onboardingSlides.map((item, index) => (
                            <button
                                key={item.title}
                                type="button"
                                onClick={() => setStep(index)}
                                className={`h-2.5 rounded-full ${index === step ? "w-8 bg-[#E91E63]" : "w-2.5 bg-zinc-200"}`}
                            />
                        ))}
                    </div>
                    <div className="mt-6 rounded-[24px] border border-zinc-200 bg-zinc-50 p-4">
                        <div className="flex items-center justify-between gap-3">
                            <div>
                                <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Choose your path</p>
                                <p className="mt-2 text-sm font-semibold text-zinc-900">{roleMeta.label}</p>
                            </div>
                            <Badge tone="default">{roleMeta.badge}</Badge>
                        </div>
                        <p className="mt-2 text-sm leading-6 text-zinc-600">{roleMeta.description}</p>
                        <div className="mt-4 grid gap-2 sm:grid-cols-3">
                            {roleOptions.map((option) => {
                                const active = option.value === selectedRole;

                                return (
                                    <button
                                        key={option.value}
                                        type="button"
                                        onClick={() => setSelectedRole(option.value)}
                                        className={`rounded-[20px] border px-3 py-3 text-left ${active
                                            ? "border-[#E91E63] bg-white"
                                            : "border-transparent bg-white/70"
                                            }`}
                                    >
                                        <p className="text-sm font-semibold text-zinc-900">{option.label}</p>
                                        <p className="mt-1 text-xs leading-5 text-zinc-600">{option.helper}</p>
                                    </button>
                                );
                            })}
                        </div>
                    </div>
                    <div className="mt-6 flex flex-col gap-3">
                        <Button type="button" variant="primary" className="w-full" onClick={() => goToVariant("register")}>
                            {roleActionLabel}
                        </Button>
                        <Button type="button" variant="secondary" className="w-full" onClick={() => goToVariant("login")}>
                            {selectedRole === "customer" ? "I already have an account" : `Sign in to ${roleMeta.label} workspace`}
                        </Button>
                    </div>
                </Card>
            </div>
        );
    }

    return (
        <div className="mx-auto max-w-[420px] px-4 py-6 sm:px-6">
            <Card className="p-5 sm:p-6">
                <div className="space-y-3">
                    <Badge tone="default">
                        {currentVariant === "success"
                            ? "Success"
                            : currentVariant === "forgot"
                                ? "Recover access"
                                : currentVariant === "reset"
                                    ? "Secure update"
                                    : currentVariant === "mfa"
                                        ? "Secure verification"
                                        : currentVariant === "register"
                                            ? "New account"
                                            : "Returning customer"}
                    </Badge>
                    <div className="flex justify-center py-2">
                        {currentVariant === "forgot" ? (
                            <Illustration variant="forgot" />
                        ) : currentVariant === "success" ? (
                            <Illustration variant="success" />
                        ) : currentVariant === "mfa" ? (
                            <Illustration variant="phone" />
                        ) : (
                            <Illustration variant="phone" />
                        )}
                    </div>
                    <div className="space-y-2 text-center">
                        <h1 className="text-[24px] font-semibold tracking-tight text-zinc-900">{heading}</h1>
                        <p className="text-sm leading-6 text-zinc-600">{subheading}</p>
                    </div>
                </div>

                {currentVariant === "register" ? (
                    <div className="mt-4 rounded-[24px] border border-zinc-200 bg-zinc-50 p-4">
                        <div className="flex items-center justify-between gap-3">
                            <div>
                                <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Current role</p>
                                <p className="mt-2 text-sm font-semibold text-zinc-900">{roleMeta.label}</p>
                            </div>
                            <Badge tone="default">Step {signupStep} of 3</Badge>
                        </div>
                        <div className="mt-3">
                            <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{signupRequirements[selectedRole].title}</p>
                            <ul className="mt-2 space-y-1 text-sm text-zinc-700">
                                {signupRequirements[selectedRole].details.map((item) => (
                                    <li key={item}>• {item}</li>
                                ))}
                            </ul>
                        </div>
                    </div>
                ) : null}

                {error ? (
                    <div className="mt-4 rounded-[24px] border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-700">
                        {error}
                    </div>
                ) : null}

                {resetSent ? (
                    <div className="mt-4 rounded-[24px] border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
                        A reset link has been sent to your inbox.
                    </div>
                ) : null}

                {passwordUpdated ? (
                    <div className="mt-4 rounded-[24px] border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
                        Your password has been updated. You can sign in now.
                    </div>
                ) : null}

                <form onSubmit={handleSubmit} className="mt-5 space-y-4">
                    {currentVariant === "register" ? (
                        <>
                            {signupStep === 1 ? renderRegisterFields() : null}
                            {signupStep === 2 ? (
                                <>
                                    <Input
                                        label="Password"
                                        type="password"
                                        placeholder="Create a secure password"
                                        required
                                        value={password}
                                        onChange={(event) => setPassword(event.target.value)}
                                    />
                                    <Input
                                        label="Confirm password"
                                        type="password"
                                        placeholder="Confirm your password"
                                        required
                                        value={confirmPassword}
                                        onChange={(event) => setConfirmPassword(event.target.value)}
                                    />
                                </>
                            ) : null}
                            {signupStep === 3 ? renderRegisterSummary() : null}
                            <div className="flex flex-wrap gap-3">
                                {signupStep > 1 ? (
                                    <Button
                                        type="button"
                                        variant="secondary"
                                        className="flex-1"
                                        onClick={() => setSignupStep((value) => Math.max(1, value - 1))}
                                    >
                                        Back
                                    </Button>
                                ) : null}
                                <Button type="submit" variant="primary" className="flex-1">
                                    {signupStep === 1
                                        ? "Continue"
                                        : signupStep === 2
                                            ? "Review details"
                                            : "Create account"}
                                </Button>
                            </div>
                        </>
                    ) : currentVariant === "login" ? (
                        <>
                            <div className="grid gap-3">
                                {roleOptions.map((option) => {
                                    const active = option.value === selectedRole;

                                    return (
                                        <button
                                            key={option.value}
                                            type="button"
                                            onClick={() => setSelectedRole(option.value)}
                                            className={`rounded-[24px] border px-4 py-3 text-left ${active
                                                ? "border-[#E91E63] bg-[#FFF8FB]"
                                                : "border-zinc-200 bg-white"
                                                }`}
                                        >
                                            <div className="flex items-center justify-between gap-3">
                                                <div>
                                                    <p className="text-sm font-semibold text-zinc-900">{option.label}</p>
                                                    <p className="mt-1 text-xs leading-5 text-zinc-600">{option.helper}</p>
                                                </div>
                                                <Badge tone="default">{option.badge}</Badge>
                                            </div>
                                            {option.value !== "customer" ? (
                                                <p className="mt-2 text-[11px] font-semibold text-rose-600">
                                                    Requires a customer account first
                                                </p>
                                            ) : null}
                                        </button>
                                    );
                                })}
                            </div>
                            <Input
                                label="Email address"
                                type="email"
                                placeholder="you@example.com"
                                required
                                value={email}
                                onChange={(event) => setEmail(event.target.value)}
                            />
                            <Input
                                label="Password"
                                type="password"
                                placeholder="Enter your password"
                                required
                                value={password}
                                onChange={(event) => setPassword(event.target.value)}
                            />
                            <div className="flex items-center justify-between gap-3 text-[12px] text-zinc-600">
                                <label className="flex items-center gap-2">
                                    <input
                                        type="checkbox"
                                        checked={remember}
                                        onChange={(event) => setRemember(event.target.checked)}
                                        className="h-4 w-4 rounded border-zinc-300 text-[#E91E63]"
                                    />
                                    Remember me
                                </label>
                                <Link href="/auth/forgot-password" className="font-semibold text-[#E91E63]">
                                    Forgot password?
                                </Link>
                            </div>
                            <Button type="submit" variant="primary" className="w-full">
                                {loading ? "Signing in…" : `Sign in as ${selectedRole}`}
                            </Button>
                        </>
                    ) : currentVariant === "forgot" ? (
                        <>
                            <Input
                                label="Email address"
                                type="email"
                                placeholder="you@example.com"
                                required
                                value={email}
                                onChange={(event) => setEmail(event.target.value)}
                            />
                            <Button type="submit" variant="primary" className="w-full">
                                Send reset link
                            </Button>
                        </>
                    ) : currentVariant === "reset" ? (
                        <>
                            <Input
                                label="New password"
                                type="password"
                                placeholder="Create a new password"
                                required
                                value={password}
                                onChange={(event) => setPassword(event.target.value)}
                            />
                            <Input
                                label="Confirm password"
                                type="password"
                                placeholder="Confirm your password"
                                required
                                value={confirmPassword}
                                onChange={(event) => setConfirmPassword(event.target.value)}
                            />
                            <Button type="submit" variant="primary" className="w-full">
                                Update password
                            </Button>
                        </>
                    ) : (
                        <>
                            <Input
                                label="Verification code"
                                type="text"
                                placeholder="Enter the 6-digit code"
                                required
                                value={otp}
                                onChange={(event) => setOtp(event.target.value)}
                            />
                            <Button type="submit" variant="primary" className="w-full">
                                {loading ? "Verifying…" : "Verify account"}
                            </Button>
                        </>
                    )}
                </form>

                {currentVariant === "login" || currentVariant === "register" ? (
                    <div className="mt-4 flex flex-wrap gap-3">
                        <SocialButton icon="G" label="Google" onClick={() => persistSocialAuth("google")} />
                        <SocialButton icon="" label="Apple" onClick={() => persistSocialAuth("apple")} />
                        <SocialButton icon="♪" label="TikTok" onClick={() => persistSocialAuth("tiktok")} />
                    </div>
                ) : null}

                <div className="mt-5 text-center text-sm text-zinc-600">
                    {currentVariant === "login" ? (
                        <>
                            New here? <button type="button" onClick={() => goToVariant("register")} className="font-semibold text-[#E91E63]">Create an account</button>
                        </>
                    ) : currentVariant === "register" ? (
                        <>
                            Already have an account? <button type="button" onClick={() => goToVariant("login")} className="font-semibold text-[#E91E63]">Sign in instead</button>
                        </>
                    ) : currentVariant === "forgot" ? (
                        <>
                            Remembered your password? <button type="button" onClick={() => goToVariant("login")} className="font-semibold text-[#E91E63]">Back to login</button>
                        </>
                    ) : currentVariant === "reset" ? (
                        <>
                            Need a new code? <Link href="/auth/forgot-password" className="font-semibold text-[#E91E63]">Request reset link</Link>
                        </>
                    ) : currentVariant === "mfa" ? (
                        <>
                            Need a new code? <button type="button" onClick={() => goToVariant("login")} className="font-semibold text-[#E91E63]">Back to sign in</button>
                        </>
                    ) : null}
                </div>

                {currentVariant === "register" ? (
                    <div className="mt-4 text-center">
                        <button
                            type="button"
                            onClick={() => {
                                resetRegisterFlow();
                                goToVariant("onboarding");
                            }}
                            className="text-sm font-semibold text-[#E91E63]"
                        >
                            Pick a different role
                        </button>
                    </div>
                ) : null}
            </Card>
        </div>
    );
}
