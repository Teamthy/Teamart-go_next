import Link from "next/link";
import PageHeader from "@/components/ui/PageHeader";
import Card from "@/components/ui/card";
import Button from "@/components/ui/button";

const providerBrand = {
    google: { label: "Google", color: "bg-white text-slate-950", border: "border border-slate-200" },
    facebook: { label: "Facebook", color: "bg-[#1877F2] text-white", border: "border-transparent" },
    apple: { label: "Apple", color: "bg-slate-950 text-white", border: "border-transparent" },
};

export default function AuthSocialPage({ params }: { params: { provider: string } }) {
    const provider = params.provider.toLowerCase();
    const brand = providerBrand[provider as keyof typeof providerBrand] ?? {
        label: provider.charAt(0).toUpperCase() + provider.slice(1),
        color: "bg-white text-slate-950",
        border: "border border-slate-200",
    };

    return (
        <div className="min-h-screen bg-[#FCE4EC] px-4 py-10 sm:px-6 lg:px-8">
            <div className="mx-auto max-w-4xl">
                <PageHeader
                    eyebrow="Social auth"
                    title={`Continue with ${brand.label}`}
                    description={`Use ${brand.label} to securely sign in to Teamart and access creator tools, storefronts, and live commerce.`}
                />

                <div className="mt-10 grid gap-8 lg:grid-cols-[0.7fr_0.3fr]">
                    <Card className="rounded-[2rem] bg-white p-8 shadow-[0_25px_60px_rgba(233,30,99,0.12)]">
                        <div className="space-y-6">
                            <div className="flex flex-wrap items-center justify-between gap-4">
                                <div>
                                    <p className="text-sm uppercase tracking-[0.32em] text-pink-500">Social login</p>
                                    <h2 className="mt-2 text-2xl font-semibold text-slate-950">Authenticate with {brand.label}</h2>
                                </div>
                                <span className={`rounded-full px-4 py-2 text-sm font-semibold ${brand.color} ${brand.border}`}>
                                    {brand.label}
                                </span>
                            </div>

                            <p className="text-sm leading-6 text-slate-600">Teamart will redirect you to {brand.label} to verify your account and then bring you back to the dashboard.</p>

                            <Button variant="primary" className="w-full">
                                Continue with {brand.label}
                            </Button>

                            <div className="rounded-3xl border border-pink-100 bg-pink-50 p-4 text-sm text-slate-700">
                                If you don’t currently have a Teamart account, one will be created automatically after sign in.
                            </div>
                        </div>
                    </Card>

                    <aside className="space-y-4 rounded-[2rem] bg-slate-950 p-6 text-slate-50 shadow-[0_25px_60px_rgba(0,0,0,0.12)]">
                        <p className="text-sm uppercase tracking-[0.32em] text-pink-300">Need another provider?</p>
                        <div className="grid gap-3">
                            <Link href="/auth/social/google" className="rounded-3xl bg-white px-4 py-3 text-center text-sm font-semibold text-slate-950">
                                Google
                            </Link>
                            <Link href="/auth/social/facebook" className="rounded-3xl bg-[#1877F2] px-4 py-3 text-center text-sm font-semibold text-white">
                                Facebook
                            </Link>
                            <Link href="/auth/social/apple" className="rounded-3xl bg-slate-950 px-4 py-3 text-center text-sm font-semibold text-white">
                                Apple
                            </Link>
                        </div>
                    </aside>
                </div>
            </div>
        </div>
    );
}
