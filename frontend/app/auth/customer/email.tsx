import { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import PageHeader from "@/components/ui/PageHeader";

export default function CustomerEmail() {
    const [email, setEmail] = useState("");
    const router = useRouter();
    const searchParams = useSearchParams();
    const firstName = searchParams.get("firstName") || "";
    const lastName = searchParams.get("lastName") || "";

    function handleNext(e: React.FormEvent) {
        e.preventDefault();
        if (email.trim()) {
            router.push(`/auth/customer/verify?firstName=${encodeURIComponent(firstName)}&lastName=${encodeURIComponent(lastName)}&email=${encodeURIComponent(email)}`);
        }
    }

    return (
        <div className="min-h-screen flex flex-col items-center justify-center px-4 py-10">
            <PageHeader
                eyebrow="Customer onboarding"
                title="What's your email?"
                description="We'll send a verification link."
            />
            <form onSubmit={handleNext} className="mt-8 w-full max-w-sm space-y-6">
                <input
                    type="email"
                    className="w-full rounded-xl border border-slate-300 px-4 py-3 text-lg focus:border-pink-500 focus:ring-pink-500"
                    placeholder="Email address"
                    value={email}
                    onChange={e => setEmail(e.target.value)}
                    required
                />
                <button
                    type="submit"
                    className="w-full rounded-full bg-pink-600 px-6 py-3 text-white font-semibold hover:bg-pink-700 transition"
                >
                    Next
                </button>
            </form>
        </div>
    );
}
