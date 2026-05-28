import AuthForm from "@/components/auth/AuthForm";
import PageHeader from "@/components/ui/PageHeader";
import Link from "next/link";

export default function AuthOtpPage() {
    return (
        <div className="min-h-screen bg-[#FCE4EC] px-4 py-10 sm:px-6 lg:px-8">
            <div className="mx-auto max-w-3xl">
                <PageHeader
                    eyebrow="Verification"
                    title="Enter your one-time code"
                    description="Confirm your identity with the secure verification code sent to your device."
                />

                <div className="mt-10 grid gap-8 lg:grid-cols-[0.75fr_0.25fr]">
                    <div className="rounded-[2rem] bg-white p-8 shadow-[0_25px_60px_rgba(233,30,99,0.12)]">
                        <AuthForm mode="otp" />
                    </div>

                    <aside className="space-y-4 rounded-[2rem] bg-slate-950 p-6 text-slate-50 shadow-[0_25px_60px_rgba(0,0,0,0.12)]">
                        <p className="text-sm uppercase tracking-[0.32em] text-pink-300">Quick tip</p>
                        <p className="text-sm leading-6 text-slate-300">Check the most recent message from Teamart if your code does not arrive immediately.</p>
                        <Link href="/auth/login" className="inline-flex rounded-full bg-[#E91E63] px-5 py-3 text-sm font-semibold text-white">
                            Back to login
                        </Link>
                    </aside>
                </div>
            </div>
        </div>
    );
}
