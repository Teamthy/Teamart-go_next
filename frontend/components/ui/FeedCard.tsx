"use client";

import { useRef, useState } from "react";
import Link from "next/link";
import Badge from "@/components/ui/badge";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";

interface FeedCardProps {
    item: {
        kind: string;
        title: string;
        description: string;
        author: string;
        avatar: string;
        image?: string;
        likes: number;
        comments: number;
        shares: number;
        tag: string;
        cta: string;
        actionHref: string;
        actionLabel: string;
        music?: string;
        hashtags?: string[];
        merchantTags?: string[];
        productTags?: string[];
        creatorHref?: string;
        storeHref?: string;
        quickProducts?: { label: string; href: string }[];
        commentThreads?: string[];
        shareCopy?: string;
    };
}

const videoSrc = "https://interactive-examples.mdn.mozilla.net/media/cc0-videos/flower.mp4";

const creatorHrefMap: Record<string, string> = {
    "Maya Chen": "/creator/maya-chen",
    "Sage Rivera": "/creator/sage-rivera",
    "Jordan Park": "/creator/jordan-park",
    "Amina Blake": "/creator/amina-blake",
    "Lena Grant": "/creator/lena-grant",
};

const storeHrefMap: Record<string, string> = {
    "Pink Thread": "/stores/pink-thread",
    "Atelier Join": "/stores/atelier-join",
    "Luma Home": "/stores/luma-home",
    "Studio Essentials": "/stores/studio-essentials",
    "Glow Lab": "/stores/glow-lab",
    "Northstar Merch": "/stores/northstar-merch",
    "Verve Market": "/stores/verve-market",
    "Horizon Kit": "/stores/horizon-kit",
};

export default function FeedCard({ item }: FeedCardProps) {
    const [liked, setLiked] = useState(false);
    const [saved, setSaved] = useState(false);
    const [playing, setPlaying] = useState(true);
    const [activeDrawer, setActiveDrawer] = useState<null | "comments" | "share" | "shop">(null);
    const videoRef = useRef<HTMLVideoElement | null>(null);

    const creatorHref = item.creatorHref ?? creatorHrefMap[item.author] ?? "/creator/maya-chen";
    const storeHref = item.storeHref ?? storeHrefMap[item.author] ?? "/stores";
    const hashtags = item.hashtags ?? ["#socialcommerce", "#liveshopping", "#discover"];
    const merchantTags = item.merchantTags ?? [];
    const productTags = item.productTags ?? ["Featured", "Shop now"];
    const quickProducts = item.quickProducts ?? [{ label: item.cta, href: item.actionHref }];
    const shareCopy = item.shareCopy ?? "Check out this social commerce moment.";
    const commentThreads = item.commentThreads ?? ["Amazing drop.", "Absolutely saving this."];

    const handlePlayToggle = () => {
        if (!videoRef.current) {
            return;
        }

        if (playing) {
            videoRef.current.pause();
        } else {
            void videoRef.current.play();
        }

        setPlaying((current) => !current);
    };

    return (
        <article className="snap-start">
            <Card className="overflow-hidden border-white/70 bg-zinc-950 text-white shadow-[0_30px_80px_rgba(15,23,42,0.35)]">
                <div className="relative">
                    <video
                        ref={videoRef}
                        className="h-[420px] w-full object-cover md:h-[520px]"
                        src={videoSrc}
                        autoPlay
                        muted
                        loop
                        playsInline
                        poster={item.image}
                    />
                    <div className="absolute inset-0 bg-[linear-gradient(180deg,rgba(15,23,42,0.22),rgba(15,23,42,0.65))]" />
                    <div className="absolute inset-x-0 top-0 p-4">
                        <div className="flex items-start justify-between gap-3">
                            <div className="flex items-center gap-3">
                                <Link href={creatorHref} className="flex items-center gap-3">
                                    <img src={item.avatar} alt={item.author} className="h-11 w-11 rounded-full border-2 border-white/80 object-cover" />
                                    <div>
                                        <p className="text-sm font-semibold">{item.author}</p>
                                        <p className="text-[11px] text-white/80">{item.kind}</p>
                                    </div>
                                </Link>
                            </div>
                            <div className="flex flex-wrap justify-end gap-2">
                                <Badge tone="default">{item.tag}</Badge>
                                {item.kind === "livestream now" ? <Badge tone="success">LIVE</Badge> : null}
                            </div>
                        </div>
                    </div>

                    <div className="absolute inset-x-0 bottom-0 p-4 sm:p-5">
                        <div className="flex flex-wrap gap-2 text-[11px] text-white/90">
                            {hashtags.map((tag) => (
                                <span key={tag} className="rounded-full bg-black/35 px-2.5 py-1">
                                    {tag}
                                </span>
                            ))}
                        </div>
                        <h3 className="mt-3 text-[20px] font-semibold leading-tight sm:text-[22px]">{item.title}</h3>
                        <p className="mt-2 max-w-xl text-sm leading-6 text-white/85">{item.description}</p>
                        <div className="mt-3 flex flex-wrap gap-2 text-[12px] text-white/80">
                            <span>{item.music ?? "Audio by Teamart Studio"}</span>
                            {merchantTags.length > 0 ? <span>•</span> : null}
                            {merchantTags.map((tag) => (
                                <span key={tag}>{tag}</span>
                            ))}
                        </div>
                        <div className="mt-4 flex flex-wrap gap-2">
                            {productTags.map((tag) => (
                                <span key={tag} className="rounded-full bg-white/12 px-3 py-1 text-[11px] font-semibold text-white">
                                    {tag}
                                </span>
                            ))}
                        </div>
                    </div>
                </div>

                <div className="border-t border-white/10 bg-zinc-950/95 p-4 sm:p-5">
                    <div className="flex flex-wrap items-center justify-between gap-3">
                        <div className="flex flex-wrap gap-2">
                            <Button
                                type="button"
                                variant="ghost"
                                className="rounded-full border-white/15 bg-white/5 px-3 py-2 text-white"
                                onClick={() => setLiked((current) => !current)}
                            >
                                {liked ? "♥" : "♡"} {liked ? item.likes + 1 : item.likes}
                            </Button>
                            <Button
                                type="button"
                                variant="ghost"
                                className="rounded-full border-white/15 bg-white/5 px-3 py-2 text-white"
                                onClick={() => setActiveDrawer((current) => (current === "comments" ? null : "comments"))}
                            >
                                💬 {item.comments}
                            </Button>
                            <Button
                                type="button"
                                variant="ghost"
                                className="rounded-full border-white/15 bg-white/5 px-3 py-2 text-white"
                                onClick={() => setActiveDrawer((current) => (current === "share" ? null : "share"))}
                            >
                                ↗ {item.shares}
                            </Button>
                            <Button
                                type="button"
                                variant="ghost"
                                className="rounded-full border-white/15 bg-white/5 px-3 py-2 text-white"
                                onClick={() => setSaved((current) => !current)}
                            >
                                {saved ? "★" : "☆"} {saved ? "Saved" : "Save"}
                            </Button>
                            <Button
                                type="button"
                                variant="ghost"
                                className="rounded-full border-white/15 bg-white/5 px-3 py-2 text-white"
                                onClick={handlePlayToggle}
                            >
                                {playing ? "Pause" : "Play"}
                            </Button>
                        </div>
                        <div className="flex flex-wrap gap-2">
                            <Button
                                type="button"
                                variant="ghost"
                                className="rounded-full border-white/15 bg-white/5 px-3 py-2 text-white"
                                onClick={() => setActiveDrawer((current) => (current === "shop" ? null : "shop"))}
                            >
                                Shop
                            </Button>
                            <Button asChild variant="primary" className="rounded-full">
                                <Link href={storeHref}>Visit store</Link>
                            </Button>
                        </div>
                    </div>

                    {activeDrawer === "comments" ? (
                        <div className="mt-4 rounded-[24px] border border-white/10 bg-white/5 p-4">
                            <p className="text-[11px] uppercase tracking-[0.2em] text-white/60">Comments</p>
                            <div className="mt-3 space-y-2 text-sm text-white/85">
                                {commentThreads.map((comment) => (
                                    <div key={comment} className="rounded-[20px] bg-black/25 px-3 py-2">
                                        {comment}
                                    </div>
                                ))}
                            </div>
                        </div>
                    ) : null}

                    {activeDrawer === "share" ? (
                        <div className="mt-4 rounded-[24px] border border-white/10 bg-white/5 p-4">
                            <p className="text-[11px] uppercase tracking-[0.2em] text-white/60">Share</p>
                            <p className="mt-3 text-sm leading-6 text-white/85">{shareCopy}</p>
                            <div className="mt-4 flex flex-wrap gap-2">
                                <Button asChild variant="secondary" className="rounded-full bg-white text-zinc-950">
                                    <Link href={item.actionHref}>{item.actionLabel}</Link>
                                </Button>
                                <Button asChild variant="ghost" className="rounded-full border-white/15 bg-white/5 text-white">
                                    <Link href="/search">Search more</Link>
                                </Button>
                            </div>
                        </div>
                    ) : null}

                    {activeDrawer === "shop" ? (
                        <div className="mt-4 rounded-[24px] border border-white/10 bg-white/5 p-4">
                            <p className="text-[11px] uppercase tracking-[0.2em] text-white/60">Shop drawer</p>
                            <div className="mt-3 flex flex-wrap gap-2">
                                {quickProducts.map((product) => (
                                    <Button key={product.href} asChild variant="secondary" className="rounded-full bg-white text-zinc-950">
                                        <Link href={product.href}>{product.label}</Link>
                                    </Button>
                                ))}
                            </div>
                        </div>
                    ) : null}
                </div>
            </Card>
        </article>
    );
}
