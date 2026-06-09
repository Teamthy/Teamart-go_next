import SellerAppShell from "@/components/seller/SellerAppShell";
import PageHeader from "@/components/seller/PageHeader";
import SectionPlaceholder from "@/components/seller/SectionPlaceholder";

const quickLinks = [
    { label: "Dashboard", href: "/seller/dashboard" },
    { label: "Products", href: "/seller/products" },
    { label: "Orders", href: "/seller/orders" },
    { label: "Analytics", href: "/seller/analytics" },
];

export default function SellerHomePage() {
    return (
        <SellerAppShell>
            <PageHeader
                title="Seller Center"
                description="Welcome back. Manage your storefront, orders, livestreams, creators, and growth tools from one unified hub."
                actions={
                    <div className="flex flex-wrap gap-3">
                        {quickLinks.map((link) => (
                            <a
                                key={link.href}
                                href={link.href}
                                className="inline-flex items-center justify-center rounded-lg border border-[var(--border-light)] bg-[var(--bg-white)] px-4 py-2 text-sm font-semibold text-[var(--text-primary)] transition hover:bg-[var(--bg-lighter)]"
                            >
                                {link.label}
                            </a>
                        ))}
                    </div>
                }
            />

            <SectionPlaceholder
                title="Your seller experience is ready"
                subtitle="Phase 7: Operational momentum"
                description="Phase 7 brings fuller navigation coverage and a unified hub for every seller workflow, including product launches, order triage, creator programs, livestream operations, and analytics."
                actions={quickLinks}
            />
        </SellerAppShell>
    );
}
