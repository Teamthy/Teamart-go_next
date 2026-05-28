import AuthForm from "@/components/auth/AuthForm";
import PageHeader from "@/components/ui/PageHeader";
import Link from "next/link";

export default function LoginPage() {
    return (
        <div className="min-h-screen bg-[#FCE4EC] px-4 py-10 sm:px-6 lg:px-8">
            <div className="mx-auto max-w-3xl">
                <PageHeader
                    eyebrow="Sign in"
                    title="Welcome back to Teamart"
                    description="Access your creator dashboard, orders, and livestream tools with a secure login."
                />

                <div className="mt-10 rounded-[2rem] bg-white p-8 shadow-[0_25px_60px_rgba(233,30,99,0.12)]">
                    <AuthForm mode="login" />
                    <div className="mt-6 flex flex-wrap items-center justify-between gap-3 text-sm text-slate-600">
                        <p>
                            New to Teamart? <Link href="/auth/register" className="font-semibold text-[#E91E63]">Create account</Link>
                        </p>
                        <Link href="/auth/forgot" className="font-semibold text-[#E91E63]">Forgot password?</Link>
                    </div>
                </div>
            </div>
        </div>
    );
}
