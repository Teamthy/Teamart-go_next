export interface StoreProfile {
    id: string;
    slug: string;
    name: string;
    category: string;
    rating: string;
    tagline: string;
    banner: string;
    followers: string;
    live: string;
    products: string;
    creator: string;
}

export const stores: StoreProfile[] = [
    { id: "s1", slug: "luma-home", name: "Luma Home", category: "Home & wellness", rating: "4.9", tagline: "Comfort-first essentials with giftable bundles and live shopper moments.", banner: "https://images.unsplash.com/photo-1505693416388-ac5ce068fe85?w=1200&q=80", followers: "28K", live: "Live now", products: "34", creator: "Sage Rivera" },
    { id: "s2", slug: "glow-lab", name: "Glow Lab", category: "Beauty & wellness", rating: "4.8", tagline: "Skincare and beauty bundles that convert well during creator-led product pins.", banner: "https://images.unsplash.com/photo-1524504388940-b1c1722653e1?w=1200&q=80", followers: "42K", live: "Today 8 PM", products: "29", creator: "Amina Blake" },
    { id: "s3", slug: "northstar-merch", name: "Northstar Merch", category: "Fashion", rating: "4.7", tagline: "Street-ready editorials and limited drops built for fast-moving audience interest.", banner: "https://images.unsplash.com/photo-1529139574466-a303027c1d8b?w=1200&q=80", followers: "51K", live: "Tomorrow 6 PM", products: "41", creator: "Maya Chen" },
    { id: "s4", slug: "studio-essentials", name: "Studio Essentials", category: "Creator tools", rating: "4.9", tagline: "A complete creator toolkit for livestream production, planning, and polished video setup.", banner: "https://images.unsplash.com/photo-1498050108023-c5249f4df085?w=1200&q=80", followers: "63K", live: "Live now", products: "38", creator: "Jordan Park" },
    { id: "s5", slug: "verve-market", name: "Verve Market", category: "Lifestyle", rating: "4.8", tagline: "Curated lifestyle goods focused on gifting, repeat buyers, and beautifully styled bundles.", banner: "https://images.unsplash.com/photo-1494438639946-1ebd1d20bf85?w=1200&q=80", followers: "34K", live: "Friday 7 PM", products: "26", creator: "Lena Grant" },
    { id: "s6", slug: "pink-thread", name: "Pink Thread", category: "Fashion", rating: "4.7", tagline: "Soft tailoring and elevated pieces for daily wear, creator collabs, and seasonal drops.", banner: "https://images.unsplash.com/photo-1529139574466-a303027c1d8b?w=1200&q=80", followers: "30K", live: "Saturday 4 PM", products: "22", creator: "Maya Chen" },
    { id: "s7", slug: "horizon-kit", name: "Horizon Kit", category: "Travel & tech", rating: "4.8", tagline: "Travel-first accessories and organizer solutions built for creators on the move.", banner: "https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?w=1200&q=80", followers: "24K", live: "Live now", products: "31", creator: "Jordan Park" },
    { id: "s8", slug: "cove-house", name: "Cove House", category: "Home", rating: "4.9", tagline: "Warm, tactile home goods ideal for seasonal bundles and highly shareable content.", banner: "https://images.unsplash.com/photo-1505693416388-ac5ce068fe85?w=1200&q=80", followers: "36K", live: "Sunday 5 PM", products: "28", creator: "Sage Rivera" },
    { id: "s9", slug: "atelier-join", name: "Atelier Join", category: "Stationery", rating: "4.8", tagline: "Premium planning and gifting staples made for creators and loyal repeat shoppers.", banner: "https://images.unsplash.com/photo-1517841905240-472988babdf9?w=1200&q=80", followers: "27K", live: "Monday 6 PM", products: "33", creator: "Lena Grant" },
    { id: "s10", slug: "craft-nest", name: "Craft Nest", category: "Accessories", rating: "4.7", tagline: "Accessory-driven merch and limited drops that feel polished, collectible, and giftable.", banner: "https://images.unsplash.com/photo-1494438639946-1ebd1d20bf85?w=1200&q=80", followers: "22K", live: "Live now", products: "19", creator: "Amina Blake" },
];

export const storeHighlights = stores.slice(0, 4);
