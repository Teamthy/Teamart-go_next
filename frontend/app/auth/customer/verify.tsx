import { useRouter, useSearchParams } from "next/navigation";
import PageHeader from "@/components/ui/PageHeader";

export default function CustomerVerify() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const firstName = searchParams.get("firstName") || "";
    const lastName = searchParams.get("lastName") || "";
    const email = searchParams.get("email") || "";

    function handleContinue() {
        // In production, trigger email verification logic here
        router.push(`/account?firstName=${encodeURIComponent(firstName)}&lastName=${encodeURIComponent(lastName)}&email=${encodeURIComponent(email)}`);
    }

    return (
        <div className="min-h-screen flex flex-col items-center justify-center px-4 py-10">
            <PageHeader
                eyebrow="Customer onboarding"
                title="Verify your email"
                description={`A verification link has been sent to ${email}. Please check your inbox.`}
            />
            <button
                onClick={handleContinue}
                className="mt-8 w-full max-w-sm rounded-full bg-pink-600 px-6 py-3 text-white font-semibold hover:bg-pink-700 transition"
            >
                Continue to dashboard
            </button>
        </div>
    );
}
