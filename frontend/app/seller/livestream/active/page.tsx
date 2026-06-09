"use client";

import SellerAppShell from "@/components/seller/SellerAppShell";
import PageHeader from "@/components/seller/PageHeader";
import SellerButton from "@/components/seller/Button";
import Badge from "@/components/seller/Badge";
import { Play, Eye, Calendar, Users } from "lucide-react";
import { fetchActiveLivestreams } from "@/lib/sellerApi";
import { useEffect, useState } from "react";

interface Livestream {
    id: string;
    title: string;
    creatorName: string;
    status: "live" | "scheduled" | "ended";
    viewers: number;
    revenue: number;
    startTime: string;
}

const initialStreams: Livestream[] = [];

const statusConfig = {
    live: { variant: "danger", label: "Live" },
    scheduled: { variant: "secondary", label: "Scheduled" },
    ended: { variant: "secondary", label: "Ended" },
} as const;

export default function ActiveLivestreamsPage() {
    const [streams, setStreams] = useState<Livestream[]>(initialStreams);

    useEffect(() => {
        let mounted = true;
        fetchActiveLivestreams().then((data) => {
            if (mounted) setStreams(data);
        });
        return () => {
            mounted = false;
        };
    }, []);

    return (
        <SellerAppShell>
            <PageHeader
                title="Active Livestreams"
                description="Monitor the streams that are currently live and manage the session in real-time."
                breadcrumbs={[
                    { label: "Dashboard", href: "/seller/dashboard" },
                    { label: "Livestream", href: "/seller/livestream" },
                    { label: "Active" },
                ]}
                actions={
                    <SellerButton variant="primary" size="md" icon={<Play className="w-4 h-4" />}>
                        Start New Stream
                    </SellerButton>
                }
            />

            <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
                {streams.map((stream) => (
                    <div key={stream.id} className="rounded-3xl border border-[var(--border-light)] bg-[var(--bg-white)] p-6 shadow-sm shadow-slate-200/10">
                        <div className="flex items-center justify-between mb-4">
                            <div>
                                <p className="text-sm text-[var(--text-tertiary)]">{stream.creatorName}</p>
                                <h3 className="text-lg font-semibold text-[var(--text-primary)]">{stream.title}</h3>
                            </div>
                            <Badge variant={statusConfig[stream.status].variant}>{statusConfig[stream.status].label}</Badge>
                        </div>

                        <div className="space-y-3 text-sm text-[var(--text-secondary)]">
                            <div className="flex items-center gap-2">
                                <Users className="w-4 h-4" />
                                <span>{stream.creatorName}</span>
                            </div>
                            <div className="flex items-center gap-2">
                                <Eye className="w-4 h-4" />
                                <span>{stream.viewers.toLocaleString()} viewers</span>
                            </div>
                            <div className="flex items-center gap-2">
                                <Calendar className="w-4 h-4" />
                                <span>{stream.startTime}</span>
                            </div>
                            <div className="text-sm text-[var(--text-primary)] font-semibold">
                                Revenue: ${stream.revenue.toFixed(2)}
                            </div>
                        </div>

                        <div className="mt-6 flex gap-2">
                            <SellerButton variant="primary" size="sm" className="flex-1">
                                View
                            </SellerButton>
                            <SellerButton variant="danger" size="sm" className="flex-1">
                                End Stream
                            </SellerButton>
                        </div>
                    </div>
                ))}
            </div>

            {streams.length === 0 && (
                <div className="rounded-3xl border border-[var(--border-light)] bg-[var(--bg-white)] p-12 text-center text-[var(--text-secondary)]">
                    <p className="text-lg font-semibold text-[var(--text-primary)] mb-2">No live streams right now</p>
                    <p className="mb-4">Start a livestream to engage viewers and drive sales instantly.</p>
                    <SellerButton variant="primary">Start a livestream</SellerButton>
                </div>
            )}
        </SellerAppShell>
    );
}
