export interface MarketingFeatureCard {
    title: string;
    description: string;
}

export interface PricingTier {
    name: string;
    price: string;
    note: string;
    perks: string[];
}

export interface FaqItem {
    q: string;
    a: string;
}

export interface MarketingCategory {
    name: string;
    description: string;
}

export const marketingFeatureCards: MarketingFeatureCard[] = [
    {
        title: "Creator-first storefronts",
        description: "Launch shoppable moments with live bundles, creator highlights, and fast checkout flows.",
    },
    {
        title: "Live commerce insights",
        description: "Pulse-check orders, audience behavior, and product momentum from one beautifully organized view.",
    },
    {
        title: "Tailored customer journeys",
        description: "Guide shoppers from discovery and support to payment and post-purchase engagement.",
    },
];

export const marketingPricingTiers: PricingTier[] = [
    {
        name: "Starter",
        price: "$29",
        note: "For solo creators launching their first drop",
        perks: ["Unlimited product pins", "Basic analytics", "1 storefront theme"],
    },
    {
        name: "Growth",
        price: "$99",
        note: "For teams running live launches and bundles",
        perks: ["Live chat automations", "Advanced insights", "Team member access"],
    },
    {
        name: "Scale",
        price: "$249",
        note: "For high-volume merchants and full operations teams",
        perks: ["Custom branding", "Priority support", "Multi-tenant admin"],
    },
];

export const marketingFaqItems: FaqItem[] = [
    {
        q: "How quickly can I launch a storefront?",
        a: "Most creators connect their products, theme, and checkout flow in under 30 minutes.",
    },
    {
        q: "Can I run live shopping alongside my existing channels?",
        a: "Yes. Teamart is built to complement TikTok, Instagram, and your existing product catalog.",
    },
    {
        q: "Do you support creator payouts and merchant billing?",
        a: "Yes. The platform includes creator earnings, payout tracking, and merchant billing visibility.",
    },
];

export const marketingSupportTopics = [
    "Launching a new store",
    "Managing order fulfillment",
    "Live stream planning",
    "Creator payouts and earnings",
    "Discounts and coupon campaigns",
];

export const marketingCategories: MarketingCategory[] = [
    { name: "Beauty", description: "Glow drops, tools, and creator essentials" },
    { name: "Home", description: "Comfortable staples and daily rituals" },
    { name: "Fashion", description: "Fits, layering, and standout pieces" },
    { name: "Tech", description: "Smart devices and creator must-haves" },
    { name: "Wellness", description: "Calm, care, and self-care sets" },
    { name: "Live Drops", description: "Exclusive launches and creator bundles" },
];

export const marketingHighlights = [
    "Limited-edition bundles",
    "Creator-led gifting",
    "Live shopping moments",
];

export const marketingLegalPoints = [
    "Transparent policies",
    "Customer-first language",
    "Accessible, mobile-ready layout",
];
