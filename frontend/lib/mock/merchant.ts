export interface MerchantMetric {
    label: string;
    value: string;
}

export interface MerchantPayout {
    id: string;
    amount: string;
    status: string;
}

export interface MerchantInventoryCount {
    name: string;
    count: string;
    note: string;
}

export const merchantMetrics: MerchantMetric[] = [
    { label: "Orders", value: "124" },
    { label: "Payouts", value: "3" },
    { label: "Products", value: "20" },
];

export const merchantActions = [
    { label: "Open orders", href: "/merchant/orders" },
    { label: "Manage products", href: "/merchant/products" },
    { label: "Inventory overview", href: "/merchant/inventory" },
];

export const merchantOperationalFocus = [
    "Review live order health",
    "Confirm pending payouts",
    "Spot fast-moving products for promotion",
];

export const merchantSettings = [
    "Default shipping profile enabled",
    "Auto-approval on low-risk orders",
    "Inventory alerts set to 10 units",
];

export const merchantQuickActions = [
    "Review orders",
    "Open payouts",
    "Manage inventory",
];

export const merchantPayouts: MerchantPayout[] = [
    { id: "P-308", amount: "$1,220", status: "Completed" },
    { id: "P-309", amount: "$1,750", status: "Pending" },
    { id: "P-310", amount: "$430", status: "Scheduled" },
];

export const merchantInventoryHighlights = [
    "Top performer is Signature Graphic Tee",
    "Low-stock items need replenishment planning",
    "Bundle-ready products are ready to push into live campaigns",
];

export const merchantInventoryCounts: MerchantInventoryCount[] = [
    { name: "In stock", count: "126", note: "Healthy availability for current drops" },
    { name: "Low stock", count: "8", note: "Needs replenishment planning" },
    { name: "Reserved", count: "14", note: "Held for pending orders" },
];

export const fulfillmentSummary = [
    "2 orders are packing now",
    "1 shipment is delayed at a carrier hub",
    "3 refunds are waiting on review",
];

export const merchantPayoutSummary = [
    "Next payout: May 31",
    "1 payout is pending approval",
    "Settlement notes are synced to merchant billing",
];

export const merchantProfiles = [
    { id: "m1", name: "Luma Home", category: "Home & wellness", store: "luma-home", live: "Live now", rating: "4.9", products: "34" },
    { id: "m2", name: "Glow Lab", category: "Beauty & wellness", store: "glow-lab", live: "Today 8 PM", rating: "4.8", products: "29" },
    { id: "m3", name: "Northstar Merch", category: "Fashion", store: "northstar-merch", live: "Tomorrow 6 PM", rating: "4.7", products: "41" },
    { id: "m4", name: "Studio Essentials", category: "Creator tools", store: "studio-essentials", live: "Live now", rating: "4.9", products: "38" },
    { id: "m5", name: "Verve Market", category: "Lifestyle", store: "verve-market", live: "Friday 7 PM", rating: "4.8", products: "26" },
    { id: "m6", name: "Pink Thread", category: "Fashion", store: "pink-thread", live: "Saturday 4 PM", rating: "4.7", products: "22" },
    { id: "m7", name: "Horizon Kit", category: "Travel & tech", store: "horizon-kit", live: "Live now", rating: "4.8", products: "31" },
    { id: "m8", name: "Cove House", category: "Home", store: "cove-house", live: "Sunday 5 PM", rating: "4.9", products: "28" },
    { id: "m9", name: "Atelier Join", category: "Stationery", store: "atelier-join", live: "Monday 6 PM", rating: "4.8", products: "33" },
    { id: "m10", name: "Craft Nest", category: "Accessories", store: "craft-nest", live: "Live now", rating: "4.7", products: "19" },
];
