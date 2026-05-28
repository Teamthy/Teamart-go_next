import { useState } from "react";
import { useRouter } from "next/navigation";
import PageHeader from "@/components/ui/PageHeader";

export default function CustomerFirstName() {
    const [firstName, setFirstName] = useState("");
    const router = useRouter();

    function handleNext(e: React.FormEvent) {
        e.preventDefault();
        if (firstName.trim()) {
            router.push("/auth/customer/last-name?firstName=" + encodeURIComponent(firstName));
        }
    }

    return (
        <div className="min-h-screen flex flex-col items-center justify-center px-4 py-10">
            <PageHeader
                eyebrow="Customer onboarding"
                title="What's your first name?"
                description="Let's get started with your customer account."
            />
            <form onSubmit={handleNext} className="mt-8 w-full max-w-sm space-y-6">
                <input
                    type="text"
                    className="w-full rounded-xl border border-slate-300 px-4 py-3 text-lg focus:border-pink-500 focus:ring-pink-500"
                    placeholder="First name"
                    value={firstName}
                    onChange={e => setFirstName(e.target.value)}
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
