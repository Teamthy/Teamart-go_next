
import PageHeader from "@/components/ui/PageHeader";
import RoleCard from "@/components/ui/RoleCard";

export default function AuthLanding() {
    return (
        <div className="min-h-screen bg-[#F9F5F8] px-4 py-10 sm:px-6 lg:px-8">
            <div className="mx-auto max-w-4xl space-y-10">
                <PageHeader
                    eyebrow="Sign in or join"
                    title="Welcome to Teamart Social Commerce"
                    description="Choose your path: customer, creator, or merchant."
                />
                <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                    <RoleCard
                        title="Customer"
                        description="Shop, follow creators, join livestreams, and leave reviews."
                        requirements="Email verification required."
                        cta="Create account"
                        href="/auth/customer/first-name"
                    />
                    <RoleCard
                        title="Creator"
                        description="Apply to host livestreams, launch drops, and grow your audience."
                        requirements="Customer account required."
                        cta="Apply as creator"
                        href="/auth/creator/start"
                    />
                    <RoleCard
                        title="Merchant"
                        description="Open a store, manage products, and access merchant analytics."
                        requirements="Customer account required."
                        cta="Open merchant store"
                        href="/auth/merchant/start"
                    />
                </div>
            </div>
        </div>
    );
}
// End of new role-based cards section
