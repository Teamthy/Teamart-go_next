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

export const sessionHistory = [
    { id: "s1", device: "Chrome on Windows", location: "New York, NY", active: true },
    { id: "s2", device: "iPhone 15", location: "San Francisco, CA", active: false },
];

export const cartItems = [
    { id: "1", name: "Sustainable Canvas Tote", price: "$32", qty: 1 },
    { id: "2", name: "Artist Collaboration Hoodie", price: "$74", qty: 2 },
];

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
