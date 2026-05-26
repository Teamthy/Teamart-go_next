import Link from "next/link";
import Badge from "@/components/ui/badge";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import SectionHeader from "@/components/ui/SectionHeader";

export function titleCase(value: string) {
    const cleaned = value.replace(/[-_]/g, " ");
    return cleaned
        .split(" ")
        .filter(Boolean)
        .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
        .join(" ");
}

export function renderHero({
    title,
    description,
    badge,
}: {
    title: string;
    description: string;
    badge?: string;
}) {
    return (
        <section className="rounded-[32px] bg-white p-5 sm:p-6">
            <div className="flex flex-col gap-5 lg:flex-row lg:items-end lg:justify-between">
                <div className="space-y-4">
                    {badge ? <Badge tone="default">{badge}</Badge> : null}
                    <div className="space-y-2">
                        <p className="text-[11px] uppercase tracking-[0.24em] text-[#E91E63]">Teamart-go_next</p>
                        <SectionHeader title={title} description={description} />
                    </div>
                </div>
                <div className="flex flex-wrap gap-3">
                    <Button asChild variant="primary">
                        <Link href="/feed">Explore marketplace</Link>
                    </Button>
                    <Button asChild variant="secondary">
                        <Link href="/search">Search products</Link>
                    </Button>
                </div>
            </div>
        </section>
    );
}

export function StatsGrid() {
    return (
        <div className="grid gap-3 sm:grid-cols-3">
            {[
                { label: "Live shoppers", value: "18.4k" },
                { label: "Creator campaigns", value: "214" },
                { label: "Orders today", value: "$86k" },
            ].map((item) => (
                <Card key={item.label} className="p-4">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.label}</p>
                    <p className="mt-2 text-2xl font-semibold text-zinc-900">{item.value}</p>
                </Card>
            ))}
        </div>
    );
}

export function SummaryList({ items }: { items: string[] }) {
    return (
        <div className="space-y-3">
            {items.map((item) => (
                <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                    {item}
                </div>
            ))}
        </div>
    );
}
