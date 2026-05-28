<<<<<<< HEAD
import AuthTemplate from "@/components/auth/AuthTemplate";
=======
import AuthForm from "@/components/auth/AuthForm";
import PageHeader from "@/components/ui/PageHeader";
import Link from "next/link";
>>>>>>> 36e8d4c (feat(auth): production auth flows and onboarding UI)

export default function RegisterPage() {
    return (
        <div className="min-h-screen bg-[#FCE4EC] px-4 py-10 sm:px-6 lg:px-8">
            <div className="mx-auto max-w-3xl">
                <PageHeader
                    eyebrow="Create account"
                    title="Start your Teamart journey"
                    description="Join the creator commerce community and discover live shopping, storefront tools, and product drops."
                />

                <div className="mt-10 rounded-[2rem] bg-white p-8 shadow-[0_25px_60px_rgba(233,30,99,0.12)]">
                    <AuthForm mode="register" />
                    <div className="mt-6 flex flex-wrap items-center justify-between gap-3 text-sm text-slate-600">
                        <p>
                            Already have an account? <Link href="/auth/login" className="font-semibold text-[#E91E63]">Sign in</Link>
                        </p>
                        <Link href="/auth/forgot" className="font-semibold text-[#E91E63]">Forgot password?</Link>
                    </div>
                </div>
            </div>
        </div>
    );
}
