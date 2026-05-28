import AuthForm from "@/components/auth/AuthForm";
import PageHeader from "@/components/ui/PageHeader";
import Link from "next/link";

export default function AuthForgotPage() {
    return (
        <div className="min-h-screen bg-[#FCE4EC] px-4 py-10 sm:px-6 lg:px-8">
            <div className="mx-auto max-w-3xl">
                <PageHeader
                    eyebrow="Recover access"
                    title="Forgot password"
                    description="Enter your email address and we’ll send a reset link so you can get back into your Teamart account."
                />

                <div className="mt-10 grid gap-8 lg:grid-cols-[0.75fr_0.25fr]">
                    <div className="rounded-[2rem] bg-white p-8 shadow-[0_25px_60px_rgba(233,30,99,0.12)]">
                        <AuthForm mode="forgot" />
                    </div>

                    <aside className="space-y-4 rounded-[2rem] bg-slate-950 p-6 text-slate-50 shadow-[0_25px_60px_rgba(0,0,0,0.12)]">
                        <p className="text-sm uppercase tracking-[0.32em] text-pink-300">Need help?</p>
                        <p className="text-sm leading-6 text-slate-300">If you don’t receive a reset email, check your spam folder or reach out to support.</p>
                        <Link href="/contact" className="inline-flex rounded-full bg-[#E91E63] px-5 py-3 text-sm font-semibold text-white">
                            Contact support
                        </Link>
                    </aside>
                </div>
            </div>
        </div>
    );
}
