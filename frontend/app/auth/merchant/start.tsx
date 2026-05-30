import PageHeader from "@/components/ui/PageHeader";
import { useRouter } from "next/navigation";

export default function MerchantOnboardingStart() {
    const router = useRouter();
    function handleOpenStore() {
        // In production, trigger merchant onboarding logic here
        router.push("/merchant");
    }
    return (
        <div className="min-h-screen flex flex-col items-center justify-center px-4 py-10">
            <PageHeader
                eyebrow="Merchant onboarding"
                title="Open a Merchant Store"
                description="Merchants can open a store, manage products, and access analytics."
            />
            <button
                onClick={handleOpenStore}
                className="mt-8 w-full max-w-sm rounded-full bg-pink-600 px-6 py-3 text-white font-semibold hover:bg-pink-700 transition"
            >
                Complete onboarding
            </button>
        </div>
    );
}
