import Link from "next/link";
import Badge from "@/components/ui/badge";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import ProductCard from "@/components/product/ProductCard";
import { products, recommendedProducts } from "@/lib/mock/products";
import {
    marketingCategories,
    marketingFaqItems,
    marketingFeatureCards,
    marketingHighlights,
    marketingLegalPoints,
    marketingPricingTiers,
    marketingSupportTopics,
} from "@/lib/mock/marketing";
import { renderHero, StatsGrid, titleCase } from "./common";

function renderMarketingPage(title: string, description: string, route: string) {
    const isContact = route === "contact";
    const isFaq = route === "faq";
    const isTerms = route === "terms";
    const isPrivacy = route === "privacy";

    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title, description, badge: "Pink commerce system" })}

            <div className="grid gap-4 xl:grid-cols-[1.1fr_0.9fr]">
                <Card className="p-5 sm:p-6">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-[#E91E63]">What you can expect</p>
                    <h2 className="mt-3 text-2xl font-semibold text-zinc-900">A consistent shopping experience across every screen</h2>
                    <div className="mt-4 grid gap-3 sm:grid-cols-2">
                        {marketingFeatureCards.map((feature) => (
                            <div key={feature.title} className="rounded-[24px] bg-[#FFF8FB] p-4">
                                <p className="text-base font-semibold text-zinc-900">{feature.title}</p>
                                <p className="mt-2 text-sm leading-6 text-zinc-600">{feature.description}</p>
                            </div>
                        ))}
                    </div>
                </Card>

                <Card className="p-5 sm:p-6">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Next steps</p>
                    <h2 className="mt-3 text-2xl font-semibold text-zinc-900">Browse the experience that matches your journey</h2>
                    <div className="mt-4 space-y-3">
                        <Link href="/feed" className="block rounded-[24px] bg-[#FCE4EC] px-4 py-3 text-sm font-semibold text-[#E91E63]">
                            Explore the feed
                        </Link>
                        <Link href="/creator" className="block rounded-[24px] bg-zinc-50 px-4 py-3 text-sm font-semibold text-zinc-800">
                            Open the creator desk
                        </Link>
                        <Link href="/cart" className="block rounded-[24px] bg-zinc-50 px-4 py-3 text-sm font-semibold text-zinc-800">
                            Review your cart
                        </Link>
                    </div>
                </Card>
            </div>

            {isContact ? (
                <div className="grid gap-4 lg:grid-cols-[0.9fr_1.1fr]">
                    <Card className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Contact</p>
                        <h2 className="mt-3 text-xl font-semibold text-zinc-900">Talk to the Teamart team</h2>
                        <div className="mt-4 space-y-3 text-sm text-zinc-600">
                            <p>Email: hello@teamart.co</p>
                            <p>Support: support@teamart.co</p>
                            <p>Hours: Mon–Fri, 9am–6pm PT</p>
                        </div>
                    </Card>
                    <Card className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">How can we help?</p>
                        <div className="mt-4 grid gap-3 sm:grid-cols-2">
                            {marketingSupportTopics.map((topic) => (
                                <div key={topic} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                    {topic}
                                </div>
                            ))}
                        </div>
                    </Card>
                </div>
            ) : null}

            {isFaq ? (
                <Card className="p-5 sm:p-6">
                    <div className="space-y-4">
                        {marketingFaqItems.map((item) => (
                            <details key={item.q} className="rounded-[24px] border border-zinc-100 bg-[#FFF8FB] px-4 py-3">
                                <summary className="cursor-pointer list-none text-sm font-semibold text-zinc-900">{item.q}</summary>
                                <p className="mt-3 text-sm leading-6 text-zinc-600">{item.a}</p>
                            </details>
                        ))}
                    </div>
                </Card>
            ) : null}

            {isTerms || isPrivacy ? (
                <Card className="p-5 sm:p-6">
                    <div className="space-y-4 text-sm leading-7 text-zinc-600">
                        <p>
                            This page mirrors the polished structure of the Teamart product experience while keeping the supporting copy concise and reader-friendly.
                        </p>
                        <p>
                            The content is organized into clear sections, supporting cards, and warm brand moments for a consistent mobile-first presentation.
                        </p>
                        <div className="grid gap-3 sm:grid-cols-3">
                            {marketingLegalPoints.map((point) => (
                                <div key={point} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-zinc-700">
                                    {point}
                                </div>
                            ))}
                        </div>
                    </div>
                </Card>
            ) : null}

            <StatsGrid />
        </div>
    );
}

function renderPricingPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({
                title: "Pricing",
                description: "Choose the growth plan that matches your creator, merchant, or operating model.",
                badge: "Flexible plans",
            })}
            <div className="grid gap-4 lg:grid-cols-3">
                {marketingPricingTiers.map((tier) => (
                    <Card key={tier.name} className="p-5">
                        <Badge tone="default">{tier.name}</Badge>
                        <p className="mt-4 text-3xl font-semibold text-zinc-900">{tier.price}</p>
                        <p className="mt-2 text-sm text-zinc-600">{tier.note}</p>
                        <div className="mt-4 space-y-2">
                            {tier.perks.map((perk) => (
                                <div key={perk} className="rounded-[20px] bg-[#FFF8FB] px-3 py-2 text-sm text-zinc-700">
                                    {perk}
                                </div>
                            ))}
                        </div>
                        <Button asChild variant="primary" className="mt-5 w-full">
                            <Link href="/auth/register">Start now</Link>
                        </Button>
                    </Card>
                ))}
            </div>
        </div>
    );
}

function renderExplorePage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({
                title: "Explore",
                description: "Browse the pink storefront experience with recommendations, live moments, and curated collections.",
                badge: "Discovery mode",
            })}
            <div className="grid gap-4 lg:grid-cols-[0.9fr_1.1fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Trending</p>
                    <h2 className="mt-3 text-xl font-semibold text-zinc-900">Everything your shoppers are following</h2>
                    <div className="mt-4 space-y-3">
                        {marketingHighlights.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                </Card>
                <div className="grid gap-4 md:grid-cols-2">
                    {products.map((product) => (
                        <ProductCard key={product.id} product={product} />
                    ))}
                </div>
            </div>
        </div>
    );
}

function renderCategoriesPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({
                title: "Categories",
                description: "Jump into discovery by inspiration, lifestyle, and shopping mood.",
                badge: "Browse by vibe",
            })}
            <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
                {marketingCategories.map((category) => (
                    <Card key={category.name} className="p-5">
                        <Badge tone="info">{category.name}</Badge>
                        <p className="mt-3 text-sm leading-6 text-zinc-600">{category.description}</p>
                        <Button asChild variant="secondary" className="mt-4">
                            <Link href={`/categories/${category.name.toLowerCase()}`}>Browse</Link>
                        </Button>
                    </Card>
                ))}
            </div>
        </div>
    );
}

function renderCategoryDetailPage(slug: string) {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({
                title: titleCase(slug),
                description: `A curated view of ${titleCase(slug)} products, creator picks, and live bundle suggestions.`,
                badge: "Category spotlight",
            })}
            <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
                {recommendedProducts.map((product) => (
                    <ProductCard key={product.id} product={product} />
                ))}
            </div>
        </div>
    );
}

function renderMaintenancePage() {
    return (
        <div className="space-y-8 pb-10">
            <Card className="p-8 text-center">
                <Badge tone="warning">Maintenance</Badge>
                <h1 className="mt-4 text-3xl font-semibold text-zinc-900">We’re refreshing the storefront</h1>
                <p className="mt-3 text-sm leading-7 text-zinc-600">
                    The Teamart experience is temporarily unavailable while we update the live commerce experience.
                    Check back soon for the next creator drop.
                </p>
                <div className="mt-6 flex justify-center gap-3">
                    <Button asChild variant="primary">
                        <Link href="/">Return home</Link>
                    </Button>
                    <Button asChild variant="secondary">
                        <Link href="/contact">Contact support</Link>
                    </Button>
                </div>
            </Card>
        </div>
    );
}

export function renderMarketing(slug: string[]) {
    const route = slug[0] ?? "home";
    const second = slug[1];

    switch (route) {
        case "about":
            return renderMarketingPage("About", "Meet the founding story behind a pink, live-first commerce platform.", route);
        case "pricing":
            return renderPricingPage();
        case "contact":
            return renderMarketingPage("Contact", "Reach the Teamart team for support, partnerships, and launch planning.", route);
        case "faq":
            return renderMarketingPage("FAQ", "Quick answers for creators, shoppers, and merchants using the Teamart experience.", route);
        case "terms":
            return renderMarketingPage("Terms", "A polished terms page that mirrors the customer-facing layout of the storefront.", route);
        case "privacy":
            return renderMarketingPage("Privacy", "A polished privacy section with clear, accessible design for your shopper support experience.", route);
        case "explore":
            return renderExplorePage();
        case "categories":
            if (second) {
                return renderCategoryDetailPage(second);
            }
            return renderCategoriesPage();
        case "maintenance":
            return renderMaintenancePage();
        case "gift-cards":
            return renderMarketingPage("Gift cards", "A polished gift card experience with easy discovery and purchase-ready actions.", route);
        case "coupon":
            return renderMarketingPage("Coupons", "A compact coupon experience layered into the storefront flow.", route);
        case "returns":
            return renderMarketingPage("Returns", "Support returns and exchanges with a thoughtful, shopper-friendly page.", route);
        case "search":
            return renderExplorePage();
        default:
            return renderMarketingPage(titleCase(route), `A pink, mobile-first version of the ${titleCase(route)} page, aligned with the Teamart storefront UI.`, route);
    }
}
