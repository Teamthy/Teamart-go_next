import PageHeader from "@/components/ui/PageHeader";
import { useRouter } from "next/navigation";

export default function CreatorOnboardingStart() {
    const router = useRouter();
    function handleApply() {
        // In production, trigger creator application logic here
        router.push("/creator");
    }
    return (
        <div className="min-h-screen flex flex-col items-center justify-center px-4 py-10">
            <PageHeader
                eyebrow="Creator onboarding"
                title="Apply as a Creator"
                description="Creators can host livestreams, launch drops, and grow their audience."
            />
            <button
                onClick={handleApply}
                className="mt-8 w-full max-w-sm rounded-full bg-pink-600 px-6 py-3 text-white font-semibold hover:bg-pink-700 transition"
            >
                Complete application
            </button>
        </div>
    );
}
