export const products = [
    {
        id: "1",
        name: "Sustainable Canvas Tote",
        price: "$32",
        description: "Everyday carry with recycled canvas and bold prints.",
        image: "/product-tote.jpg",
    },
    {
        id: "2",
        name: "Active Lifestyle Sneakers",
        price: "$98",
        description: "Comfort-first design for daily training and street style.",
        image: "/product-sneakers.jpg",
    },
    {
        id: "3",
        name: "Curated Creator Notebook",
        price: "$18",
        description: "A premium journal for planning creator launches and campaigns.",
        image: "/product-notebook.jpg",
    },
];

export const recommendedProducts = [
    {
        id: "5",
        name: "Artist Collaboration Hoodie",
        price: "$74",
        description: "Limited drop hoodie made with creator-first craftsmanship.",
        image: "/product-hoodie.jpg",
    },
    {
        id: "6",
        name: "Live Stream Ring Light",
        price: "$48",
        description: "Professional lighting for creator streams and product showcases.",
        image: "/product-light.jpg",
    },
];

export const featuredProducts = products;

export const categories = [
    { slug: "fashion", name: "Fashion" },
    { slug: "home", name: "Home" },
    { slug: "beauty", name: "Beauty" },
    { slug: "tech", name: "Tech" },
    { slug: "art", name: "Art" },
    { slug: "wellness", name: "Wellness" },
];

export const creators = [
    {
        id: "mia-rivera",
        name: "Mia Rivera",
        handle: "mia",
        bio: "Creator partner blending lifestyle drops, live commerce, and curated collections.",
        avatar: "/creator-mia.jpg",
        followers: "42.1k",
        liveStatus: "Live now",
        products: 28,
        rating: 4.9,
        category: "Lifestyle",
    },
    {
        id: "nina-park",
        name: "Nina Park",
        handle: "ninapark",
        bio: "Beauty and fashion creator serving limited edition drops in every stream.",
        avatar: "/creator-nina.jpg",
        followers: "31.4k",
        liveStatus: "Offline",
        products: 22,
        rating: 4.8,
        category: "Beauty",
    },
    {
        id: "leo-chan",
        name: "Leo Chan",
        handle: "leoch",
        bio: "Tech creator who launches gadgets, studio accessories, and livestream bundles.",
        avatar: "/creator-leo.jpg",
        followers: "18.7k",
        liveStatus: "Ending soon",
        products: 16,
        rating: 4.7,
        category: "Tech",
    },
    {
        id: "ava-simpson",
        name: "Ava Simpson",
        handle: "ava_sims",
        bio: "Home décor creator featuring pop-up drops and creator-curated service pieces.",
        avatar: "/creator-ava.jpg",
        followers: "26.9k",
        liveStatus: "Offline",
        products: 19,
        rating: 4.9,
        category: "Home",
    },
    {
        id: "kai-evans",
        name: "Kai Evans",
        handle: "kaievans",
        bio: "Creator focused on activewear, drops, and fast checkout bundles for followers.",
        avatar: "/creator-kai.jpg",
        followers: "12.2k",
        liveStatus: "Live now",
        products: 24,
        rating: 4.6,
        category: "Fashion",
    },
    {
        id: "zara-bell",
        name: "Zara Bell",
        handle: "zarabell",
        bio: "Studio entrepreneur and creator launching merch and livestream collaborations.",
        avatar: "/creator-zara.jpg",
        followers: "9.8k",
        liveStatus: "Offline",
        products: 12,
        rating: 4.7,
        category: "Creator",
    },
];

export const sessionHistory = [
    { id: "s1", device: "Chrome on Windows", location: "New York, NY", active: true },
    { id: "s2", device: "iPhone 15", location: "San Francisco, CA", active: false },
];

// cartItems and wishlistItems have moved to lib/mock/products

export const sellerProducts = [
    { id: "p1", sku: "TA-2024-01", name: "Artisanal Ceramic Mug", price: "$28", stock: 34, status: "Live", sales: 154 },
    { id: "p2", sku: "TA-2024-02", name: "Signature Graphic Tee", price: "$36", stock: 18, status: "Live", sales: 286 },
    { id: "p3", sku: "TA-2024-03", name: "Limited Edition Poster", price: "$22", stock: 12, status: "Live", sales: 92 },
    { id: "p4", sku: "TA-2024-04", name: "Premium Leather Pouch", price: "$45", stock: 9, status: "Low stock", sales: 78 },
    { id: "p5", sku: "TA-2024-05", name: "Eco-Friendly Water Bottle", price: "$26", stock: 46, status: "Live", sales: 130 },
    { id: "p6", sku: "TA-2024-06", name: "Designer Face Mask", price: "$18", stock: 72, status: "Live", sales: 214 },
    { id: "p7", sku: "TA-2024-07", name: "Custom Enamel Pin", price: "$12", stock: 82, status: "Live", sales: 430 },
    { id: "p8", sku: "TA-2024-08", name: "Branded Drawstring Bag", price: "$22", stock: 26, status: "Live", sales: 98 },
    { id: "p9", sku: "TA-2024-09", name: "Studio Notebook", price: "$16", stock: 55, status: "Live", sales: 168 },
    { id: "p10", sku: "TA-2024-10", name: "Collector’s Sticker Pack", price: "$10", stock: 140, status: "Live", sales: 600 },
    { id: "p11", sku: "TA-2024-11", name: "Signature Hoodie", price: "$72", stock: 22, status: "Live", sales: 58 },
    { id: "p12", sku: "TA-2024-12", name: "Creator Desk Lamp", price: "$54", stock: 13, status: "Low stock", sales: 46 },
    { id: "p13", sku: "TA-2024-13", name: "Artisan Candle Set", price: "$38", stock: 32, status: "Live", sales: 120 },
    { id: "p14", sku: "TA-2024-14", name: "Interactive Planner", price: "$29", stock: 48, status: "Live", sales: 77 },
    { id: "p15", sku: "TA-2024-15", name: "Travel Tech Organizer", price: "$44", stock: 19, status: "Low stock", sales: 53 },
    { id: "p16", sku: "TA-2024-16", name: "Live Stream Backdrop", price: "$88", stock: 8, status: "Low stock", sales: 37 },
    { id: "p17", sku: "TA-2024-17", name: "Collector’s Edition Journal", price: "$34", stock: 29, status: "Live", sales: 64 },
    { id: "p18", sku: "TA-2024-18", name: "Eco Canvas Sneakers", price: "$82", stock: 5, status: "Low stock", sales: 43 },
    { id: "p19", sku: "TA-2024-19", name: "Portable Charging Pad", price: "$46", stock: 21, status: "Live", sales: 89 },
    { id: "p20", sku: "TA-2024-20", name: "Premium Hoodie Box", price: "$98", stock: 15, status: "Live", sales: 72 },
];

export const sellerOrders = [
    { id: "O-5021", customer: "Mia R.", total: "$248", status: "Processing", items: 6, date: "May 20" },
    { id: "O-5018", customer: "Nate K.", total: "$124", status: "Shipped", items: 3, date: "May 19" },
    { id: "O-5015", customer: "Asha T.", total: "$398", status: "Pending", items: 9, date: "May 18" },
    { id: "O-5012", customer: "Lina W.", total: "$76", status: "Delivered", items: 2, date: "May 17" },
    { id: "O-5009", customer: "Theo P.", total: "$199", status: "Delivered", items: 5, date: "May 16" },
    { id: "O-5004", customer: "Jade B.", total: "$112", status: "Canceled", items: 1, date: "May 14" },
];

export const sellerPayouts = [
    { id: "P-308", amount: "$1,220", period: "May 11 - May 17", status: "Completed" },
    { id: "P-309", amount: "$1,750", period: "May 18 - May 24", status: "Pending" },
    { id: "P-310", amount: "$430", period: "May 25 - May 31", status: "Scheduled" },
];

export const accountProfile = {
    name: "Mia Rivera",
    email: "mia@teamart.co",
    location: "San Francisco, CA",
    favoriteStorefront: "Creator drops",
    shippingAddress: "142 Valencia St, San Francisco, CA 94103",
    memberSince: "May 2024",
    plan: "Growth",
};

export const accountOrderSnapshot = [
    "3 orders in progress",
    "1 order awaiting shipment",
    "2 completed deliveries",
    "Fast response time this week",
];

export const accountPreferences = [
    { label: "Discovery", value: "Personalized recommendations enabled" },
    { label: "Marketing", value: "Weekly creator drops and promotions" },
    { label: "Notifications", value: "Live stream reminders and order updates" },
    { label: "Language", value: "English (US)" },
];

export const accountPaymentMethods = [
    { label: "Visa ending in 2412", value: "Primary" },
    { label: "PayPal", value: "Secondary" },
    { label: "Apple Pay", value: "Fast checkout" },
];

export const accountBillingHistory = [
    "Subscription invoice paid",
    "Pending payment authorization",
    "Card update available for faster checkout",
];

export const customerOrders = [
    { id: "O-7124", date: "May 25", customer: "Mia Rivera", items: 4, total: "$122.00", status: "Shipped", delivery: "Arrives May 30" },
    { id: "O-7118", date: "May 22", customer: "Mia Rivera", items: 2, total: "$86.00", status: "Processing", delivery: "Preparing shipment" },
    { id: "O-7109", date: "May 18", customer: "Mia Rivera", items: 3, total: "$148.00", status: "Delivered", delivery: "Delivered May 23" },
    { id: "O-7084", date: "May 13", customer: "Mia Rivera", items: 1, total: "$34.00", status: "Delivered", delivery: "Delivered May 17" },
    { id: "O-7062", date: "May 10", customer: "Mia Rivera", items: 5, total: "$214.00", status: "Canceled", delivery: "Canceled" },
];
