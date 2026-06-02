"use client";

import Link from "next/link";
import type { LiveRoomSummary } from "@/types/commerce";
import LiveBadge from "@/components/live/LiveBadge";
import ViewerCount from "@/components/live/ViewerCount";

interface LiveCardProps {
    room: LiveRoomSummary;
}

export default function LiveCard({ room }: LiveCardProps) {
    return (
        <article className="group overflow-hidden rounded-[32px] border border-slate-200/70 bg-white shadow-xl transition hover:-translate-y-1 hover:shadow-2xl">
            <div className="relative h-56 overflow-hidden bg-slate-900">
                <img src={room.thumbnail_url} alt={room.title} className="h-full w-full object-cover transition duration-500 group-hover:scale-105" />
                <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent" />
                <div className="absolute left-4 top-4">
                    <LiveBadge status={room.is_live ? "LIVE" : room.status.toUpperCase()} />
                </div>
                <div className="absolute left-4 bottom-4">
                    <ViewerCount viewers={room.viewers} />
                </div>
            </div>
            <div className="p-5">
                <div className="flex items-center gap-3">
                    <img src={room.host.avatar_url} alt={room.host.name} className="h-11 w-11 rounded-full object-cover" />
                    <div className="min-w-0 flex-1">
                        <p className="truncate text-sm font-semibold text-slate-900">{room.host.name}</p>
                        <p className="truncate text-xs text-slate-500">@{room.host.handle}</p>
                    </div>
                </div>
                <h3 className="mt-4 text-lg font-semibold text-slate-900">{room.title}</h3>
                <p className="mt-2 text-sm leading-6 text-slate-600">{room.short_description ?? "Join a new shopping stream with product drops, rapid reactions, and checkout-ready bundle moments."}</p>
                <div className="mt-4 rounded-3xl bg-slate-50 p-4 text-sm text-slate-600">
                    <p className="font-medium text-slate-900">Pinned product</p>
                    <p className="mt-2 text-sm">{room.pinned_product.name} · ${room.pinned_product.price.toFixed(2)}</p>
                </div>
                <Link
                    href={`/live/${room.id}`}
                    className="mt-5 inline-flex w-full items-center justify-center rounded-full bg-pink-500 px-4 py-3 text-sm font-semibold text-white transition hover:bg-pink-600"
                >
                    Watch Live
                </Link>
            </div>
        </article>
    );
}
