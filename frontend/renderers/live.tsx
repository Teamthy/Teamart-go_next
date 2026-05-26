import Link from "next/link";
import Badge from "@/components/ui/badge";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import { liveActionPlan, liveRoomAudienceActions, liveRoomHighlights, liveSchedule, liveStats, liveStatsSummary, liveTopMoments, liveUpNext } from "@/lib/mock/live";
import { renderHero } from "./common";

function renderLivePage(title: string, description: string) {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title, description, badge: "Live" })}
            <div className="grid gap-4 md:grid-cols-3">
                {liveStats.map((item) => (
                    <Card key={item.label} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.label}</p>
                        <p className="mt-3 text-2xl font-semibold text-zinc-900">{item.value}</p>
                    </Card>
                ))}
            </div>
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Current stream</p>
                    <div className="mt-4 space-y-3">
                        {liveRoomHighlights.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">{item}</div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Up next</p>
                    <div className="mt-4 flex flex-wrap gap-3">
                        {liveUpNext.map((item) => (
                            <Button key={item.href} asChild variant={item.href === "/live/room" ? "primary" : "secondary"}>
                                <Link href={item.href}>{item.label}</Link>
                            </Button>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderLiveRoomPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Live room", description: "Join the shopping room and keep the chat, products, and next-step CTA aligned with the live moment.", badge: "Live room" })}
            <div className="grid gap-4 lg:grid-cols-[1.1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Live highlights</p>
                    <div className="mt-4 space-y-3">
                        {liveRoomHighlights.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">{item}</div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Audience actions</p>
                    <div className="mt-4 space-y-3">
                        {liveRoomAudienceActions.map((item) => (
                            <div key={item} className="rounded-[24px] bg-zinc-50 px-4 py-3 text-sm text-zinc-700">{item}</div>
                        ))}
                    </div>
                    <div className="mt-5 flex gap-3">
                        <Button asChild variant="primary"><Link href="/cart">Buy now</Link></Button>
                        <Button asChild variant="secondary"><Link href="/live">Browse rooms</Link></Button>
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderLiveStatsPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Live stats", description: "Review the audience surge, product momentum, and the top-performing moments that matter after the stream.", badge: "Live stats" })}
            <div className="grid gap-4 md:grid-cols-3">
                {liveStatsSummary.map((item) => (
                    <Card key={item.label} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.label}</p>
                        <p className="mt-3 text-2xl font-semibold text-zinc-900">{item.value}</p>
                    </Card>
                ))}
            </div>
            <div className="grid gap-4 lg:grid-cols-2">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Top moments</p>
                    <div className="mt-4 space-y-3">
                        {liveTopMoments.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">{item}</div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Action plan</p>
                    <div className="mt-4 space-y-3">
                        {liveActionPlan.map((item) => (
                            <div key={item} className="rounded-[24px] bg-zinc-50 px-4 py-3 text-sm text-zinc-700">{item}</div>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderLiveSchedulePage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Live schedule", description: "Keep the upcoming creator calendar and room plan clear, coordinated, and easy to revisit.", badge: "Live schedule" })}
            <div className="grid gap-4 lg:grid-cols-2">
                {liveSchedule.map((item) => (
                    <Card key={item.id} className="p-5">
                        <p className="text-lg font-semibold text-zinc-900">{item.name}</p>
                        <p className="mt-2 text-sm text-zinc-600">{item.time}</p>
                        <p className="mt-2 text-xs uppercase tracking-[0.15em] text-zinc-500">{item.status}</p>
                    </Card>
                ))}
            </div>
        </div>
    );
}

export function renderLive(slug: string[]) {
    const route = slug[0] ?? "live";
    const second = slug[1];

    if (route !== "live") {
        return renderLivePage("Live", "Explore live rooms, audience insights, and creator-led shopping moments in one polished surface.");
    }

    if (second === "room") return renderLiveRoomPage();
    if (second === "stats") return renderLiveStatsPage();
    if (second === "schedule") return renderLiveSchedulePage();

    return renderLivePage("Live", "Explore live rooms, audience insights, and creator-led shopping moments in one polished surface.");
}
