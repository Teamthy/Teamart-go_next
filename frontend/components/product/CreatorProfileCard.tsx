import Link from "next/link";
import type { Creator } from "@/lib/mock/creators";

interface CreatorProfileCardProps {
    creator: Creator;
}

export default function CreatorProfileCard({ creator }: CreatorProfileCardProps) {
    return (
        <section className="overflow-hidden rounded-[28px] border border-zinc-200 bg-white">
            <div className="flex flex-col gap-5 p-5 sm:flex-row sm:items-center sm:justify-between sm:p-6">
                <div className="flex items-center gap-4">
                    <img src={creator.avatar} alt={creator.name} className="h-16 w-16 rounded-[24px] object-cover" />
                    <div>
                        <p className="text-[11px] uppercase tracking-[0.24em] text-[#E91E63]">Creator spotlight</p>
                        <h3 className="mt-1 text-xl font-semibold text-zinc-900">{creator.name}</h3>
                        <p className="text-sm text-zinc-500">{creator.handle} • {creator.followers} followers</p>
                        <p className="mt-2 max-w-2xl text-sm leading-6 text-zinc-600">{creator.bio}</p>
                    </div>
                </div>
                <div className="flex flex-wrap gap-2">
                    <span className="rounded-full bg-[#FCE4EC] px-3 py-1 text-[11px] font-semibold text-[#E91E63]">
                        {creator.mutual ? "Mutual followers" : "New to you"}
                    </span>
                    <span className="rounded-full bg-zinc-100 px-3 py-1 text-[11px] font-semibold text-zinc-700">{creator.livestreamSchedule}</span>
                </div>
            </div>

            <div className="grid gap-3 border-t border-zinc-100 p-5 sm:grid-cols-3">
                <div className="rounded-[24px] bg-[#FFF8FB] p-4">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Top drop</p>
                    <p className="mt-2 text-lg font-semibold text-zinc-900">{creator.products[0] ?? "Featured bundle"}</p>
                </div>
                <div className="rounded-[24px] bg-[#FFF8FB] p-4">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Live shows</p>
                    <p className="mt-2 text-lg font-semibold text-zinc-900">{creator.engagement}</p>
                </div>
                <div className="rounded-[24px] bg-[#FFF8FB] p-4">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Best seller</p>
                    <p className="mt-2 text-lg font-semibold text-zinc-900">{creator.products[1] ?? creator.products[0] ?? "Creator favorite"}</p>
                </div>
            </div>

            <div className="flex flex-wrap gap-3 px-5 pb-5 sm:px-6">
                <Link href={`/creator/${creator.id}`} className="inline-flex rounded-[24px] bg-[#E91E63] px-4 py-3 text-sm font-semibold text-white">
                    View creator profile
                </Link>
                <Link href="/live" className="inline-flex rounded-[24px] border border-zinc-200 bg-white px-4 py-3 text-sm font-semibold text-zinc-900">
                    Join live stream
                </Link>
            </div>
        </section>
    );
}
