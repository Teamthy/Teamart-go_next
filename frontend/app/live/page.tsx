"use client";

import Link from "next/link";
import { useLiveFeed } from "@/hooks/useLiveFeed";
import LiveCard from "@/components/live/LiveCard";
import Button from "@/components/ui/button";

export default function LivePage() {
    const { rooms, isLoading, error } = useLiveFeed();

    return (
        <div className="space-y-8 pb-10">
            <section className="rounded-[32px] border border-slate-200 bg-white p-8 shadow-xl shadow-slate-200/50">
                <div className="flex flex-col gap-6 xl:flex-row xl:items-center xl:justify-between">
                    <div className="max-w-3xl">
                        <p className="text-sm uppercase tracking-[0.3em] text-pink-600">Live rooms</p>
                        <h1 className="mt-4 text-4xl font-semibold tracking-tight text-slate-950">Live commerce rooms</h1>
                        <p className="mt-4 max-w-2xl text-base leading-8 text-slate-600">
                            Join trending creator and merchant streams, track viewer pulse, and shop pinned products directly from the room.
                        </p>
                    </div>
                    <div className="flex flex-wrap gap-3">
                        <Button variant="secondary" asChild>
                            <Link href="/feed">Back to feed</Link>
                        </Button>
                        <Button variant="primary" onClick={() => window.location.assign("/search")}>Browse products</Button>
                    </div>
                </div>
                <div className="mt-6 grid gap-4 sm:grid-cols-3">
                    <div className="rounded-[28px] bg-slate-950 p-5 text-white shadow-lg">
                        <p className="text-sm uppercase tracking-[0.2em] text-slate-400">Rooms live</p>
                        <p className="mt-3 text-3xl font-semibold">{rooms.length}</p>
                        <p className="mt-2 text-sm text-slate-300">Live hotel-style commerce streams available now.</p>
                    </div>
                    <div className="rounded-[28px] bg-slate-950 p-5 text-white shadow-lg">
                        <p className="text-sm uppercase tracking-[0.2em] text-slate-400">Engagement</p>
                        <p className="mt-3 text-3xl font-semibold">{rooms.reduce((sum, room) => sum + room.viewers, 0).toLocaleString()}</p>
                        <p className="mt-2 text-sm text-slate-300">Viewers across active rooms right now.</p>
                    </div>
                    <div className="rounded-[28px] bg-slate-950 p-5 text-white shadow-lg">
                        <p className="text-sm uppercase tracking-[0.2em] text-slate-400">Pinned products</p>
                        <p className="mt-3 text-3xl font-semibold">{rooms.length}</p>
                        <p className="mt-2 text-sm text-slate-300">Each room highlights one hero product for fast discovery.</p>
                    </div>
                </div>
            </section>

            <section className="space-y-4">
                <div className="grid gap-4 xl:grid-cols-3">
                    {rooms.map((room) => (
                        <LiveCard key={room.id} room={room} />
                    ))}
                </div>
                {isLoading ? (
                    <div className="rounded-[32px] border border-slate-200 bg-white p-8 text-center text-slate-600 shadow-lg">Loading live rooms...</div>
                ) : null}
                {error ? (
                    <div className="rounded-[32px] border border-rose-200 bg-rose-50 p-6 text-sm text-rose-700 shadow-lg">{error}</div>
                ) : null}
            </section>
        </div>
    );
}
