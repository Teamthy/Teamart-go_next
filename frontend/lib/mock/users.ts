export interface UserSession {
    id: string;
    device: string;
    location: string;
    active: boolean;
    avatar: string;
}

export interface PlatformUser {
    id: string;
    name: string;
    role: string;
    status: string;
    region: string;
    avatar: string;
}

export interface AccountPreference {
    label: string;
    value: string;
}

export interface SavedProduct {
    id: string;
    name: string;
    price: string;
    note: string;
}

export interface MembershipTier {
    name: string;
    level: string;
    perks: string[];
}

export const sessionHistory: UserSession[] = [
    { id: "s1", device: "Chrome on Windows", location: "New York, NY", active: true, avatar: "👩‍💻" },
    { id: "s2", device: "iPhone 15", location: "San Francisco, CA", active: false, avatar: "📱" },
    { id: "s3", device: "Chrome on Mac", location: "Austin, TX", active: false, avatar: "💼" },
];

export const adminUsers: PlatformUser[] = [
    { id: "u1", name: "Avery Lane", role: "Creator operators", status: "Active", region: "North America", avatar: "🧑‍🎨" },
    { id: "u2", name: "Noah Patel", role: "Merchant admins", status: "Active", region: "Europe", avatar: "🛍️" },
    { id: "u3", name: "Sofia Kim", role: "Moderation support", status: "Review", region: "APAC", avatar: "🛡️" },
    { id: "u4", name: "Cameron Reyes", role: "Platform analysts", status: "Active", region: "North America", avatar: "📈" },
];

export const accountProfile = {
    name: "Mia Rivera",
    email: "mia@teamart.co",
    location: "San Francisco, CA",
    favoriteStorefront: "Creator drops",
    avatar: "https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=200&q=80",
    memberSince: "May 2024",
    plan: "Growth",
};

export const accountQuickActions = [
    "Personal details",
    "Security settings",
    "Notifications",
    "Orders and downloads",
];

export const accountPreferences: AccountPreference[] = [
    { label: "Discovery", value: "Personalized recommendations enabled" },
    { label: "Marketing", value: "Weekly creator drops and promotions" },
    { label: "Notifications", value: "Live stream reminders and order updates" },
    { label: "Language", value: "English (US)" },
];

export const savedProducts: SavedProduct[] = [
    { id: "saved-1", name: "Artist Collaboration Hoodie", price: "$74", note: "Saved for tonight’s live drop" },
    { id: "saved-2", name: "Curated Creator Notebook", price: "$18", note: "Saved for campaign planning" },
    { id: "saved-3", name: "Sustainable Canvas Tote", price: "$32", note: "Saved for gifting" },
];

export const accountMemberships: MembershipTier[] = [
    { name: "Creator drop pass", level: "Active", perks: ["Early access", "Live chat priority", "Seasonal perks"] },
    { name: "Live bundle club", level: "Active", perks: ["Bundle previews", "Member-only codes", "Priority support"] },
    { name: "VIP support add-on", level: "Paused", perks: ["Priority replies", "Dedicated checkout support"] },
];

export const accountSecurityChecks = [
    "Two-factor authentication enabled",
    "Recovery email synced",
    "Recent password rotation completed",
    "Trusted devices monitored",
];

export const accountPaymentMethods = [
    { label: "Visa ending in 2412", value: "Primary" },
    { label: "PayPal", value: "Secondary" },
    { label: "Apple Pay", value: "Fast checkout" },
];

export const accountWalletSummary = [
    "Available balance: $1,240.00",
    "Pending payouts: $260.00",
    "Card on file: Visa ending in 2412",
];

export const accountSupportTopics = [
    "Order status and delivery updates",
    "Returns and exchange requests",
    "Payment method changes",
    "Shipping and bundle questions",
];

export const accountOrderSnapshot = [
    "3 orders in progress",
    "1 order awaiting shipment",
    "2 completed deliveries",
    "Fast response time this week",
];

export const accountReturnRequests = [
    "Order O-6034 — exchange pending",
    "Order O-5981 — refund in review",
    "Order O-5920 — label ready to print",
];

export const accountBillingHistory = [
    "Subscription invoice paid",
    "Merchant fee summary available",
    "Customer support notes synced to billing history",
];

export const accountDownloadItems = [
    { label: "Invoices", value: "3 files ready to download" },
    { label: "Product guides", value: "2 creator bundles and setup docs" },
    { label: "Shipping labels", value: "1 label available for pickup" },
    { label: "Receipts", value: "5 recent purchases stored" },
];

export const accountAddressCards = [
    { label: "Default shipping", value: "42 Pink Lane, Apt 3, San Francisco, CA 94107" },
    { label: "Secondary billing", value: "88 Willow Street, Suite 2, Oakland, CA 94607" },
];

export const customers = [
    { id: "c1", name: "Mia Rivera", email: "mia@teamart.co", role: "customer", status: "Verified", region: "San Francisco" },
    { id: "c2", name: "Nadia Rose", email: "nadia@teamart.co", role: "customer", status: "Verified", region: "Seattle" },
    { id: "c3", name: "Theo James", email: "theo@teamart.co", role: "customer", status: "Verified", region: "Austin" },
    { id: "c4", name: "Amar Patel", email: "amar@teamart.co", role: "customer", status: "Verified", region: "Chicago" },
    { id: "c5", name: "Riley Chen", email: "riley@teamart.co", role: "customer", status: "Pending", region: "New York" },
    { id: "c6", name: "Asha Kim", email: "asha@teamart.co", role: "customer", status: "Verified", region: "Los Angeles" },
    { id: "c7", name: "Kira Stone", email: "kira@teamart.co", role: "customer", status: "Verified", region: "Atlanta" },
    { id: "c8", name: "Leo Ford", email: "leo@teamart.co", role: "customer", status: "Verified", region: "Denver" },
    { id: "c9", name: "Jade Brooks", email: "jade@teamart.co", role: "customer", status: "Pending", region: "Miami" },
    { id: "c10", name: "Noah Kim", email: "noah@teamart.co", role: "customer", status: "Verified", region: "Boston" },
];
