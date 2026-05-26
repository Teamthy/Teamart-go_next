import Link from "next/link";
import Badge from "@/components/ui/badge";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";

interface LiveRoomCardProps {
    name: string;
    host: string;
    viewers: string;
    status: string;
    cta: string;
    href: string;
    summary?: string;
    pinnedProduct?: string;
    badge?: string;
}

export default function LiveRoomCard({ name, host, viewers, status, cta, href, summary, pinnedProduct, badge }: LiveRoomCardProps) {
    return (
        <Card className="p-5">
            <div className="flex items-start justify-between gap-3">
                <div>
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{host}</p>
                    <h3 className="mt-2 text-lg font-semibold text-zinc-900">{name}</h3>
                    {summary ? <p className="mt-3 text-sm leading-6 text-zinc-600">{summary}</p> : null}
                </div>
                <Badge tone={status === "Live now" ? "success" : "default"}>{status}</Badge>
            </div>
            <div className="mt-4 flex flex-wrap items-center justify-between gap-3">
                <p className="text-sm text-zinc-600">{viewers} watching right now</p>
                {badge ? <span className="rounded-full bg-[#FCE4EC] px-3 py-1 text-[11px] font-semibold text-[#E91E63]">{badge}</span> : null}
            </div>
            {pinnedProduct ? (
                <div className="mt-4 rounded-[20px] bg-[#FFF8FB] px-4 py-3">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Pinned product</p>
                    <p className="mt-2 text-sm font-semibold text-zinc-900">{pinnedProduct}</p>
                </div>
            ) : null}
            <Button asChild variant="primary" className="mt-4 w-full">
                <Link href={href}>{cta}</Link>
            </Button>
        </Card>
    );
}
