"use client";

import { useEffect } from "react";
import Link from "next/link";
import PageHeader from "@/components/ui/PageHeader";
import RoleCard from "@/components/ui/RoleCard";
import StatCard from "@/components/ui/StatCard";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import RouteGuard from "@/components/auth/RouteGuard";
import { getStoredCustomer, hasRole, setWorkspaceRole } from "@/lib/auth-state";
import { creatorProfiles } from "@/lib/mock/creators";

export default function CreatorPage() {
    const customer = getStoredCustomer();
    const isCreator = hasRole("creator");

    useEffect(() => {
        if (isCreator) {
            setWorkspaceRole("creator");
        }
    }, [isCreator]);

    return (
        <RouteGuard requiredRole="creator" redirectTo="/auth/creator">
            <div className="space-y-8 pb-10">
                <PageHeader
                    title="Creator hub"
                    description="Build your next campaign, manage your live room strategy, and keep your strongest product moments in front of your audience."
                    actions={
                        <>
                            <Button asChild variant="primary">
                                <Link href="/creator/studio">Open creator studio</Link>
                            </Button>
                            <Button asChild variant="secondary">
                                <Link href="/feed">See the feed</Link>
                            </Button>
                        </>
                    }
                />

                <div className="grid gap-4 md:grid-cols-3">
                    <StatCard label="Audience reach" value="84.2k" helper="Creator campaigns are performing strongly across live and social moments." />
                    <StatCard label="Avg. conversion" value="8.4%" helper="Bundle campaigns are converting faster than static showcases." />
                    <StatCard label="Creator mood" value="Ready" helper={customer ? `Welcome back, ${customer.firstName}.` : "Guest shoppers are exploring creator-led content."} />
                </div>

                <div className="grid gap-4 xl:grid-cols-[0.85fr_1.15fr]">
                    <Card className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Creator system</p>
                        <h2 className="mt-3 text-xl font-semibold text-zinc-900">Keep every drop polished, discoverable, and conversion-ready</h2>
                        <p className="mt-3 text-sm leading-7 text-zinc-600">
                            Curate categories, prepare live highlights, and keep your audience engaged with a studio experience that supports consistent product launches.
                        </p>
                        <div className="mt-5 space-y-3">
                            {[
                                "Plan drops around the strongest audience windows",
                                "Promote a featured product before each livestream",
                                "Use creator analytics to test what converts best",
                            ].map((item) => (
                                <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                    {item}
                                </div>
                            ))}
                        </div>
                    </Card>

                    <div className="space-y-4">
                        <RoleCard
                            title={isCreator ? "Creator access is active" : "Open creator access"}
                            description={isCreator ? "Your creator studio is ready for livestream planning, product pins, and audience-first merchandising." : "Create a creator account to unlock studio planning, live rooms, and revenue-focused content workflows."}
                            requirements={isCreator ? ["Studio insights are synced", "Creator analytics ready", "Product highlights enabled"] : ["Creator name and niche", "A social handle", "A business or brand email address"]}
                            ctaLabel={isCreator ? "Open creator studio" : "Create creator account"}
                            href={isCreator ? "/creator/studio" : "/auth/register?role=creator"}
                            tone={isCreator ? "success" : "warning"}
                        />
                    </div>
                </div>

                <div className="space-y-3">
                    <div>
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Creator network</p>
                        <h2 className="mt-2 text-xl font-semibold text-zinc-900">Fresh voices and market-ready content</h2>
                    </div>
                    <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
                        {creatorProfiles.slice(0, 3).map((creator) => (
                            <Card key={creator.id} className="p-5">
                                <div className="flex items-center gap-3">
                                    <img src={creator.avatar} alt={creator.name} className="h-12 w-12 rounded-full object-cover" />
                                    <div>
                                        <p className="font-semibold text-zinc-900">{creator.name}</p>
                                        <p className="text-sm text-zinc-500">{creator.handle}</p>
                                    </div>
                                </div>
                                <p className="mt-4 text-sm leading-6 text-zinc-600">{creator.bio}</p>
                                <div className="mt-4 flex flex-wrap gap-2 text-sm text-zinc-500">
                                    <span>{creator.followers} followers</span>
                                    <span>{creator.engagement}</span>
                                </div>
                            </Card>
                        ))}
                    </div>
                </div>
            </div>
        </RouteGuard>
    );
}
