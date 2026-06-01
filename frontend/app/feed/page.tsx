"use client";

import { useMemo, useState } from "react";
import Link from "next/link";
import PageHeader from "@/components/ui/PageHeader";
import FeedCard from "@/components/ui/FeedCard";
import StatCard from "@/components/ui/StatCard";
import Tabs from "@/components/ui/Tabs";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import Badge from "@/components/ui/badge";
import { getStoredCustomer } from "@/lib/auth-state";
import { feedItems } from "@/lib/mock/feed";
import { creators } from "@/lib/mock/creators";
import { products } from "@/lib/mock/products";
import { stores } from "@/lib/mock/stores";

const tabOptions = [
    { label: "For you", value: "all" },
    { label: "Live", value: "livestream now" },
    { label: "Creators", value: "creator post" },
    { label: "Merchants", value: "merchant spotlight" },
    { label: "Reviews", value: "customer review" },
    { label: "Promos", value: "product promo" },
];

const trendingTerms = [
    "#livebundle",
    "#creatordrop",
    "#beauty",
    "#giftable",
    "#limitedstock",
];

export default function FeedPage() {
    const [activeTab, setActiveTab] = useState("all");
    const customer = getStoredCustomer();

    const filteredItems = useMemo(() => {
        return activeTab === "all"
            ? feedItems
            : feedItems.filter((item) => item.kind === activeTab);
    }, [activeTab]);

    const liveCount = feedItems.filter(
        (item) => item.kind === "livestream now"
    ).length;

    return (
        <div className="space-y-8 pb-10">
            <PageHeader
                title="For you feed"
                description={
                    customer
                        ? `Welcome back, ${customer.firstName}. Your ${customer.favoriteCategory ?? "fashion"
                        } feed is tuned for live rooms, creator drops, and checkout-ready moments.`
                        : "A TikTok-native social shopping experience with live rooms, creator bundles, merchant spotlights, and instant checkout paths."
                }
                actions={
                    <>
                        <Button asChild variant="primary">
                            <Link href="/live">Jump into live</Link>
                        </Button>
                        <Button asChild variant="secondary">
                            <Link href="/search">Search now</Link>
                        </Button>
                    </>
                }
            />

            <div className="grid gap-4 md:grid-cols-3">
                <StatCard
                    label="Live rooms"
                    value={String(liveCount)}
                    helper="Creator and merchant rooms are on the move right now."
                />
                <StatCard
                    label="Storefront picks"
                    value="20"
                    helper="Fresh content is curated to keep discovery and checkout aligned."
                />
                <StatCard
                    label="Saved moments"
                    value="8"
                    helper="Your social commerce feed stays tuned for fast, useful shopping decisions."
                />
            </div>

            <Tabs tabs={tabOptions} active={activeTab} onChange={setActiveTab} />
        </div>
    );
}