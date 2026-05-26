export interface AnalyticsMetric {
    label: string;
    value: string;
}

export interface AnalyticsInsight {
    title: string;
    description: string;
}

export interface TrendPoint {
    label: string;
    value: string;
    delta: string;
}

export const platformStats: AnalyticsMetric[] = [
    { label: "Live shoppers", value: "18.4k" },
    { label: "Active stores", value: "320" },
    { label: "Orders today", value: "$86k" },
];

export const platformRevenueSummary = [
    { label: "GMV", value: "$428k", delta: "+12%" },
    { label: "Conversion", value: "8.4%", delta: "+1.1%" },
    { label: "Repeat buyers", value: "24%", delta: "+3%" },
];

export const platformTrendMetrics: TrendPoint[] = [
    { label: "Live sessions", value: "144", delta: "+18%" },
    { label: "Campaign saves", value: "2.6k", delta: "+9%" },
    { label: "Checkout completion", value: "94%", delta: "+2%" },
];

export const creatorPerformancePulse: AnalyticsInsight[] = [
    { title: "Top live drop", description: "Creator Collaboration Hoodie" },
    { title: "Highest click-through", description: "Sunset bundle sets" },
    { title: "Fastest growth", description: "Livestream product pins" },
];

export const creatorRecommendedActions: AnalyticsInsight[] = [
    { title: "Share a recap", description: "Share a creator recap after the last livestream" },
    { title: "Promote the top bundle", description: "Promote the top-performing bundle in the next session" },
    { title: "Tune product pins", description: "Tune product pins for higher average order value" },
];

export const adminAnalyticsMetrics: AnalyticsMetric[] = [
    { label: "GMV", value: "$428k" },
    { label: "Orders", value: "1,284" },
    { label: "Live sessions", value: "144" },
];

export const adminOperationalSignals = [
    "Creator acquisition trending upward",
    "Live commerce conversion steady",
    "Merchant support queue stable",
];

export const adminWatchList = [
    "Promo spend trending upward",
    "All-time high for livestream shares",
    "Moderation queue requires review",
];

export const merchantAnalyticsMetrics: AnalyticsMetric[] = [
    { label: "Orders", value: "124" },
    { label: "Payouts", value: "3" },
    { label: "Products", value: "20" },
];

export const merchantOperationalFocus = [
    "Review live order health",
    "Confirm pending payouts",
    "Spot fast-moving products for promotion",
];

export const conversionHighlights = [
    "Live room CTA conversion is up 11%",
    "Bundle purchases lead the highest AOV segment",
    "Creator recommendations are driving repeat buyers",
];
