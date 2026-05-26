import Link from "next/link";
import { notFound } from "next/navigation";
import ProductCard from "@/components/product/ProductCard";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import { products } from "@/lib/mock/products";

const liveRooms = [
    {
        id: "sunset-launch",
        title: "Sunset launch room",
        host: "Maya Chen",
        viewers: "2.4k",
        status: "Live now",
        badge: "Popular",
        pinnedProductId: "5",
        summary: "A fast-moving drop where the hoodie, bundle styling, and limited-time savings are highlighted in real time.",
        description: "Shoppers are reacting to the latest drop, asking about fit and restock, and moving from discovery into checkout with a simple bundle path.",
        offer: "Bundle the hoodie with the tote and save 12% at checkout.",
        chat: [
            { author: "Lena", message: "The rose color looks amazing on camera." },
            { author: "Theo", message: "I’m adding the bundle now — the pricing feels really clear." },
            { author: "Amina", message: "The fit recommendation helped me decide in seconds." },
        ],
    },
    {
        id: "mini-haul",
        title: "Mini haul with product pins",
        host: "Amina Blake",
        viewers: "1.8k",
        status: "Live now",
        badge: "Beauty",
        pinnedProductId: "4",
        summary: "A beauty-first room highlighting wellness bundles, glow-focused products, and quick-shopping prompts.",
        description: "This room is tuned for a smooth, social-shopping experience with a strong product pin, audience questions, and repeat-purchase cues.",
        offer: "Try the serum and water bottle together for a curated self-care bundle.",
        chat: [
            { author: "Riley", message: "The packaging feels premium and giftable." },
            { author: "Kai", message: "I’m checking out the wellness set now." },
            { author: "Noah", message: "Love the quick bundle suggestion." },
        ],
    },
    {
        id: "creator-qna",
        title: "Creator Q&A and restock",
        host: "Jordan Park",
        viewers: "1.3k",
        status: "Live now",
        badge: "Tech",
        pinnedProductId: "6",
        summary: "A creator-led room where the livestream kit, setup advice, and restock timing are all in focus.",
        description: "The host is answering product questions, spotlighting the ring light, and helping shoppers decide which setup fits their workflow best.",
        offer: "Use the ring light in your next setup and save on the creator bundle.",
        chat: [
            { author: "Nadia", message: "The brightness controls sound exactly right for tutorials." },
            { author: "Mia", message: "I’m sold — the setup guide is so clear." },
            { author: "Zara", message: "This is the perfect add-on for my next stream." },
        ],
    },
    {
        id: "friday-restock",
        title: "Friday restock and merchant answers",
        host: "Luma Home",
        viewers: "1.1k",
        status: "Tonight",
        badge: "Home",
        pinnedProductId: "1",
        summary: "A merchant-led room centered on restock updates, shipping clarity, and the most giftable daily essentials.",
        description: "The room is built to answer the practical questions shoppers have around shipping, order timing, and giftable essentials.",
        offer: "Pick up the tote and companion bundle before the restock window closes.",
        chat: [
            { author: "Jules", message: "The shipping timing is super helpful." },
            { author: "Aria", message: "I’m getting the tote and the notebook together." },
            { author: "Talia", message: "This is the kind of room that makes shopping easy." },
        ],
    },
];

export default async function LiveRoomDetailPage({ params }: { params: Promise<{ id: string }> }) {
    const { id } = await params;
    const room = liveRooms.find((item) => item.id === id);

    if (!room) {
        notFound();
    }

    const pinnedProduct = products.find((product) => product.id === room.pinnedProductId) ?? products[0];

    return (
        <div className="space-y-8 pb-10">
            <div className="grid gap-6 xl:grid-cols-[1.1fr_0.9fr]">
                <Card className="p-5 sm:p-6">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-[#E91E63]">Live room</p>
                    <h1 className="mt-3 text-3xl font-semibold text-zinc-900">{room.title}</h1>
                    <p className="mt-3 text-sm leading-7 text-zinc-600">{room.description}</p>
                    <div className="mt-5 flex flex-wrap gap-3">
                        <span className="rounded-full bg-[#FCE4EC] px-3 py-1 text-[11px] font-semibold text-[#E91E63]">{room.badge}</span>
                        <span className="rounded-full bg-zinc-100 px-3 py-1 text-[11px] font-semibold text-zinc-700">{room.status}</span>
                        <span className="rounded-full bg-emerald-100 px-3 py-1 text-[11px] font-semibold text-emerald-800">{room.viewers} watching</span>
                    </div>
                    <div className="mt-5 rounded-[24px] bg-zinc-950 p-5 text-white">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-pink-200">Host spotlight</p>
                        <p className="mt-3 text-lg font-semibold">{room.host}</p>
                        <p className="mt-2 text-sm leading-7 text-zinc-100">{room.summary}</p>
                    </div>
                    <div className="mt-5 rounded-[24px] bg-[#FFF8FB] p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Offer</p>
                        <p className="mt-3 text-sm leading-7 text-zinc-700">{room.offer}</p>
                    </div>
                    <div className="mt-5 flex flex-wrap gap-3">
                        <Button asChild variant="primary">
                            <Link href="/cart">Add the pinned item</Link>
                        </Button>
                        <Button asChild variant="secondary">
                            <Link href="/live">Back to live rooms</Link>
                        </Button>
                    </div>
                </Card>

                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Pinned product</p>
                    <div className="mt-4">
                        <ProductCard product={pinnedProduct} />
                    </div>
                    <div className="mt-5 rounded-[24px] bg-[#FFF8FB] p-4">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Live chat</p>
                        <div className="mt-3 space-y-3">
                            {room.chat.map((entry) => (
                                <div key={entry.author} className="rounded-[20px] bg-white px-4 py-3 text-sm text-zinc-700">
                                    <span className="font-semibold text-zinc-900">{entry.author}</span>: {entry.message}
                                </div>
                            ))}
                        </div>
                    </div>
                </Card>
            </div>
        </div>
    );
}
