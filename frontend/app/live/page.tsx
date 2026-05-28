import Link from "next/link";
import PageHeader from "@/components/ui/PageHeader";
import StatCard from "@/components/ui/StatCard";
import LiveRoomCard from "@/components/ui/LiveRoomCard";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";

const liveRooms = [
    {
        name: "Sunset launch room",
        host: "Maya Chen",
        viewers: "2.4k",
        status: "Live now",
        cta: "Join live",
        href: "/live/sunset-launch",
        badge: "Popular",
        pinnedProduct: "Artist Collaboration Hoodie",
        summary: "Limited drop, bundle prompts, and chat reactions driving fast conversion.",
    },
    {
        name: "Mini haul with product pins",
        host: "Amina Blake",
        viewers: "1.8k",
        status: "Live now",
        cta: "Join live",
        href: "/live/mini-haul",
        badge: "Beauty",
        pinnedProduct: "Glow Cloud Serum",
        summary: "A short, high-energy room designed to guide shoppers from discovery to checkout.",
    },
    {
        name: "Creator Q&A and restock",
        host: "Jordan Park",
        viewers: "1.3k",
        status: "Live now",
        cta: "Join live",
        href: "/live/creator-qna",
        badge: "Tech",
        pinnedProduct: "Live Stream Ring Light",
        summary: "A focused room for product questions, live replies, and audience-led bundle discovery.",
    },
    {
        name: "Friday restock and merchant answers",
        host: "Luma Home",
        viewers: "1.1k",
        status: "Tonight",
        cta: "Preview room",
        href: "/live/friday-restock",
        badge: "Home",
        pinnedProduct: "Sustainable Canvas Tote",
        summary: "A calm merchant-led room that pairs restock news with easy, low-friction shipping clarity.",
    },
];

export default function LivePage() {
    return (
        <div className="space-y-8 pb-10">
            <PageHeader
                title="Live shopping rooms"
                description="Watch creator drops, join merchant Q&As, and move from discovery into checkout with a single room experience."
                actions={
                    <>
                        <Button asChild variant="primary">
                            <Link href="/feed">Back to feed</Link>
                        </Button>
                        <Button asChild variant="secondary">
                            <Link href="/products">Browse products</Link>
                        </Button>
                    </>
                }
            />

            <div className="grid gap-4 md:grid-cols-3">
                <StatCard label="Live rooms" value="4" helper="Creator and merchant rooms are currently active." />
                <StatCard label="Watchers" value="6.6k" helper="Audience size is trending upward this evening." />
                <StatCard label="Conversion" value="9.3%" helper="Live rooms outperform static product discovery." />
            </div>

            <div className="grid gap-4 xl:grid-cols-[0.9fr_1.1fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Room strategy</p>
                    <h2 className="mt-3 text-xl font-semibold text-zinc-900">Keep the shopping path fast and actionable</h2>
                    <p className="mt-3 text-sm leading-7 text-zinc-600">
                        Each room now highlights a hero product, surfaces a clear next step, and keeps the conversation focused on checkout-ready moments.
                    </p>
                    <div className="mt-5 space-y-3">
                        {[
                            "Pin top-performing products during room launch",
                            "Use chat prompts to drive quick bundle selections",
                            "Promote limited-time offers with a strong CTA",
                        ].map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                    <div className="mt-5 rounded-[24px] bg-zinc-950 p-4 text-white">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-pink-200">Live shopping pulse</p>
                        <p className="mt-3 text-sm leading-7 text-zinc-100">
                            The strongest rooms are pairing one hero product, one bundle suggestion, and one high-clarity CTA to reduce friction and improve checkout confidence.
                        </p>
                    </div>
                </Card>
                <div className="grid gap-4 md:grid-cols-2">
                    {liveRooms.map((room) => (
                        <LiveRoomCard key={room.name} {...room} />
                    ))}
                </div>
            </div>
        </div>
    );
}
