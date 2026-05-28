export interface Creator {
    id: string;
    name: string;
    handle: string;
    avatar: string;
    followers: string;
    bio: string;
    mutual: boolean;
    engagement: string;
    livestreamSchedule: string;
    products: string[];
}

export interface CreatorMetric {
    label: string;
    value: string;
}

export interface CreatorAction {
    label: string;
    value: string;
}

export interface CreatorScheduleItem {
    title: string;
    time: string;
    status: string;
}

export const creators: Creator[] = [
    {
        id: "maya-chen",
        name: "Maya Chen",
        handle: "@mayastyles",
        avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=200&q=80",
        followers: "142K",
        bio: "A creator-led fashion studio mixing soft tailoring with live bundle drops.",
        mutual: false,
        engagement: "8.9% engagement",
        livestreamSchedule: "Today • 8 PM",
        products: ["Artist Collaboration Hoodie", "Sustainable Canvas Tote"],
    },
    {
        id: "jordan-park",
        name: "Jordan Park",
        handle: "@jordanpicks",
        avatar: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=200&q=80",
        followers: "89K",
        bio: "Tech and lifestyle picks that convert well in product pins and creator collabs.",
        mutual: true,
        engagement: "7.4% engagement",
        livestreamSchedule: "Tomorrow • 6 PM",
        products: ["Live Stream Ring Light", "Portable Charging Pad"],
    },
    {
        id: "sage-rivera",
        name: "Sage Rivera",
        handle: "@sageliving",
        avatar: "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=200&q=80",
        followers: "203K",
        bio: "Living-room favorites and premium home essentials for daily scroll-worthy bundles.",
        mutual: false,
        engagement: "9.3% engagement",
        livestreamSchedule: "Friday • 7 PM",
        products: ["Collector’s Edition Journal", "Artisan Candle Set"],
    },
];

export const creatorDashboardStats: CreatorMetric[] = [
    { label: "Live viewers", value: "1.8k" },
    { label: "Conversion", value: "9.2%" },
    { label: "Earnings", value: "$4.2k" },
];

export const creatorAnalyticsStats: CreatorMetric[] = [
    { label: "Reach", value: "84.2k" },
    { label: "Conversion", value: "8.4%" },
    { label: "Average order", value: "$46" },
];

export const creatorWorkspaceActions: CreatorAction[] = [
    { label: "Schedule the next livestream", value: "Plan the next creator drop" },
    { label: "Pin a featured product", value: "Highlight your top bundle" },
    { label: "Review creator analytics", value: "Stay on top of performance" },
];

export const creatorQuickLinks = [
    { label: "Products", href: "/creator/products" },
    { label: "Analytics", href: "/creator/analytics" },
    { label: "Livestream", href: "/creator/livestream" },
];

export const creatorStudioFocus = [
    "Review today’s performance and content priorities",
    "Pin the hero product for the next livestream",
    "Share the latest campaign recap with the team",
];

export const creatorStudioNextSteps = [
    "Open livestream planner",
    "Review catalog",
];

export const creatorSchedule: CreatorScheduleItem[] = [
    { title: "Sunset drop", time: "Today • 8 PM", status: "Live soon" },
    { title: "Creator Q&A", time: "Tomorrow • 6 PM", status: "Scheduled" },
    { title: "Weekend bundle reveal", time: "Friday • 7 PM", status: "Scheduled" },
];

export const creatorProfiles: Creator[] = [
    ...creators,
    {
        id: "amina-blake",
        name: "Amina Blake",
        handle: "@aminablend",
        avatar: "https://images.unsplash.com/photo-1488426862026-3ee34a7d66df?w=200&q=80",
        followers: "118K",
        bio: "Beauty and wellness creator focused on glow drops, mini routines, and fresh bundle storytelling.",
        mutual: true,
        engagement: "8.2% engagement",
        livestreamSchedule: "Saturday • 4 PM",
        products: ["Live Stream Ring Light", "Eco-Friendly Water Bottle"],
    },
    {
        id: "lena-grant",
        name: "Lena Grant",
        handle: "@lenaplan",
        avatar: "https://images.unsplash.com/photo-1544723795-3fb6469f5b39?w=200&q=80",
        followers: "164K",
        bio: "A planner and lifestyle creator who turns organization tools into daily commerce moments.",
        mutual: false,
        engagement: "7.8% engagement",
        livestreamSchedule: "Monday • 6 PM",
        products: ["Interactive Planner", "Collector’s Edition Journal"],
    },
    {
        id: "nate-owen",
        name: "Nate Owen",
        handle: "@natecurates",
        avatar: "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=200&q=80",
        followers: "96K",
        bio: "Curates premium home and travel essentials with a sharp eye for product storytelling and repeat buyers.",
        mutual: true,
        engagement: "8.5% engagement",
        livestreamSchedule: "Tuesday • 7 PM",
        products: ["Travel Tech Organizer", "Premium Leather Pouch"],
    },
    {
        id: "aria-sol",
        name: "Aria Sol",
        handle: "@ariasol",
        avatar: "https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=200&q=80",
        followers: "138K",
        bio: "Lifestyle creator blending fashion, home accents, and calm styling inspiration into memorable promotions.",
        mutual: false,
        engagement: "9.1% engagement",
        livestreamSchedule: "Wednesday • 5 PM",
        products: ["Signature Graphic Tee", "Artisan Candle Set"],
    },
    {
        id: "cora-grid",
        name: "Cora Grid",
        handle: "@coragrid",
        avatar: "https://images.unsplash.com/photo-1508214751196-bcfd4ca60f91?w=200&q=80",
        followers: "109K",
        bio: "A community-first creator known for converting reactions and audience questions into quick purchase decisions.",
        mutual: true,
        engagement: "8.4% engagement",
        livestreamSchedule: "Thursday • 8 PM",
        products: ["Collector’s Sticker Pack", "Branded Drawstring Bag"],
    },
    {
        id: "jules-fair",
        name: "Jules Fair",
        handle: "@julesfair",
        avatar: "https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=200&q=80",
        followers: "153K",
        bio: "A video-first creator building bundles that look polished on camera and convert on social proof.",
        mutual: false,
        engagement: "8.9% engagement",
        livestreamSchedule: "Sunday • 6 PM",
        products: ["Studio Notebook", "Portable Charging Pad"],
    },
];

export const creatorHighlights = [
    "Top product pin is the Creator Collaboration Hoodie",
    "Audience response is strongest on evening drops",
    "Bundle conversion grew 18% this week",
];
