"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { AnimatePresence, motion } from "framer-motion";
import {
    ArrowLeft,
    ArrowRight,
    ArrowUp,
    Bell,
    Briefcase,
    Clock,
    CreditCard,
    Gift,
    Globe,
    Heart,
    LayoutGrid,
    Menu,
    Play,
    QrCode,
    RefreshCw,
    Search,
    Shield,
    ShoppingBag,
    ShoppingCart,
    Store,
    Wallet,
    X,
} from "lucide-react";

const navLinks = [
    { label: "Feed", href: "#feed" },
    { label: "Products", href: "#search" },
    { label: "Live", href: "#live" },
    { label: "Creator", href: "#creator" },
    { label: "Merchant", href: "#merchant" },
];

const featureCards = [
    {
        icon: <ShoppingBag className="w-6 h-6 text-fuchsia-500" />,
        title: "Creator storefronts",
        description: "Launch social-first storefronts with product pins, collections, and followable creator pages.",
    },
    {
        icon: <ShoppingCart className="w-6 h-6 text-cyan-500" />,
        title: "Premium commerce UX",
        description: "Amazon-level product pages, ratings, shipping clarity, and checkout trust baked in.",
    },
    {
        icon: <Store className="w-6 h-6 text-sky-500" />,
        title: "Seller command center",
        description: "Manage orders, inventory, pricing, promotions, and payouts from one unified workspace.",
    },
    {
        icon: <Globe className="w-6 h-6 text-emerald-500" />,
        title: "Multi-tenant platform",
        description: "Tenant-aware architecture built for creators, brands, merchants, and marketplace operators.",
    },
    {
        icon: <Wallet className="w-6 h-6 text-orange-500" />,
        title: "Live commerce flow",
        description: "Real-time video shopping with chat, reactions, buy buttons, and product pinning.",
    },
    {
        icon: <Shield className="w-6 h-6 text-indigo-500" />,
        title: "Trust and compliance",
        description: "Secure payments, buyer protections, and admin controls for safe marketplace growth.",
    },
];

const carouselCards = [
    {
        title: "Launch creator drops",
        subtitle: "with live social urgency.",
        description: "Blend short-form commerce, creator storytelling, and instant shopper action in a single feed.",
        bg: "bg-gradient-to-br from-fuchsia-600 to-orange-500",
        cardTitle: "Creator Live Drop",
        cardDesc: "Pin products, drive conversion, and surface creator collections to shoppers instantly.",
        icon: <Gift className="w-8 h-8" />,
    },
    {
        title: "Optimize conversion",
        subtitle: "with trusted ecommerce UX.",
        description: "Surface reviews, inventory, shipping, and recommendations tailored for purchase confidence.",
        bg: "bg-gradient-to-br from-slate-950 to-slate-700",
        cardTitle: "Product Detail Focus",
        cardDesc: "Build a premium commerce page with galleries, variants, and buy flows that convert.",
        icon: <ShoppingCart className="w-8 h-8" />,
    },
    {
        title: "Manage merchants",
        subtitle: "with TikTok Seller Center power.",
        description: "Get analytics, revenue tools, campaign controls, and creator collaboration signals in one place.",
        bg: "bg-gradient-to-br from-emerald-600 to-sky-500",
        cardTitle: "Seller Analytics",
        cardDesc: "Monitor top products, live sales metrics, and inventory health across brands and stores.",
        icon: <Briefcase className="w-8 h-8" />,
    },
];

const bnplTabs = [
    "Livestream drops",
    "Creator collections",
    "Store campaigns",
    "Product bundles",
];

const howItWorks = [
    {
        number: "01",
        title: "Launch your storefront",
        description: "Create a branded multi-tenant store and onboard your first creator.",
        color: "text-fuchsia-500",
    },
    {
        number: "02",
        title: "Publish social commerce content",
        description: "Build feed cards, live drop sessions, and creator-led product stories.",
        color: "text-cyan-500",
    },
    {
        number: "03",
        title: "Engage shoppers live",
        description: "Host livestream commerce events with pinned products, chat, and instant buy links.",
        color: "text-emerald-500",
    },
    {
        number: "04",
        title: "Convert and scale",
        description: "Track revenue, optimize campaigns, and automate product recommendations.",
        color: "text-orange-500",
    },
];

const businessTabs = [
    "Creator storefronts",
    "Seller dashboards",
    "AI recommendations",
    "Live commerce integrations",
];

export default function LandingPage() {
    const [mobileOpen, setMobileOpen] = useState(false);
    const [showDropdown, setShowDropdown] = useState(false);
    const [activeBnpl, setActiveBnpl] = useState(0);
    const [activeBusinessTab, setActiveBusinessTab] = useState(0);
    const [showTop, setShowTop] = useState(false);

    useEffect(() => {
        const onScroll = () => setShowTop(window.scrollY > 600);
        window.addEventListener("scroll", onScroll);
        return () => window.removeEventListener("scroll", onScroll);
    }, []);

    const heroActions = [
        { label: "Explore the feed", href: "/feed" },
        { label: "Launch your store", href: "/merchant" },
    ];

    return (
        <div className="min-h-screen bg-slate-950 text-white">
            <header className="sticky top-0 z-50 border-b border-white/10 backdrop-blur-xl bg-slate-950/95">
                <div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-4 sm:px-6 lg:px-8">
                    <Link href="/" className="flex items-center gap-3 text-lg font-semibold tracking-tight text-white">
                        <div className="grid h-11 w-11 place-items-center rounded-2xl bg-gradient-to-br from-fuchsia-500 to-orange-500 text-white shadow-lg shadow-fuchsia-500/20">
                            <ArrowUp className="h-5 w-5" />
                        </div>
                        Teamart
                    </Link>

                    <div className="hidden items-center gap-8 md:flex">
                        {navLinks.map((link) => (
                            <a key={link.href} href={link.href} className="text-sm text-white/70 transition hover:text-white">
                                {link.label}
                            </a>
                        ))}
                        <button
                            onMouseEnter={() => setShowDropdown(true)}
                            onMouseLeave={() => setShowDropdown(false)}
                            className="relative flex items-center gap-2 text-sm text-white/70 transition hover:text-white"
                        >
                            Solutions
                            <span className="text-white/50">▼</span>
                            <AnimatePresence>
                                {showDropdown && (
                                    <motion.div
                                        initial={{ opacity: 0, y: 12 }}
                                        animate={{ opacity: 1, y: 0 }}
                                        exit={{ opacity: 0, y: 12 }}
                                        className="absolute left-0 top-full mt-3 w-[320px] rounded-3xl border border-white/10 bg-slate-900/95 p-4 shadow-2xl"
                                    >
                                        <div className="grid gap-3">
                                            {[
                                                { title: "Live commerce", detail: "Host shoppable streams with chat and pinned offers." },
                                                { title: "Product marketplace", detail: "Premium pages with reviews, shipping, and buy flows." },
                                                { title: "Creator economy", detail: "Grow partnerships and amplify merchant reach." },
                                            ].map((item) => (
                                                <div key={item.title} className="rounded-3xl border border-white/10 bg-slate-950 p-4">
                                                    <h4 className="font-semibold text-white">{item.title}</h4>
                                                    <p className="mt-1 text-xs text-slate-400">{item.detail}</p>
                                                </div>
                                            ))}
                                        </div>
                                    </motion.div>
                                )}
                            </AnimatePresence>
                        </button>
                    </div>

                    <button className="md:hidden rounded-2xl border border-white/10 bg-white/5 p-3 text-white" onClick={() => setMobileOpen(!mobileOpen)}>
                        {mobileOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
                    </button>
                </div>

                <AnimatePresence>
                    {mobileOpen && (
                        <motion.div
                            initial={{ opacity: 0, height: 0 }}
                            animate={{ opacity: 1, height: "auto" }}
                            exit={{ opacity: 0, height: 0 }}
                            className="border-t border-white/10 bg-slate-950/95 md:hidden"
                        >
                            <div className="space-y-2 px-4 pb-4 pt-2">
                                {navLinks.map((link) => (
                                    <a key={link.href} href={link.href} className="block rounded-3xl px-4 py-3 text-sm text-white/80 hover:bg-white/5 hover:text-white">
                                        {link.label}
                                    </a>
                                ))}
                                <Link href="/auth/login" className="block rounded-3xl bg-fuchsia-500 px-4 py-3 text-center text-sm font-semibold text-white">
                                    Sign in
                                </Link>
                            </div>
                        </motion.div>
                    )}
                </AnimatePresence>
            </header>

            <main>
                <section className="relative overflow-hidden bg-[radial-gradient(circle_at_top,_rgba(255,255,255,0.12),_transparent_45%),linear-gradient(180deg,#050816_0%,#0e122a_100%)] px-4 pt-20 pb-24 sm:px-6 lg:px-8">
                    <div className="absolute inset-x-0 top-0 h-96 bg-[radial-gradient(circle_at_top,_rgba(192,132,252,0.24),_transparent_24%)]" />
                    <div className="relative mx-auto grid max-w-7xl gap-12 lg:grid-cols-[0.95fr_1.05fr] xl:gap-20">
                        <div className="max-w-2xl space-y-8 pt-10">
                            <div className="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-white/80 backdrop-blur">
                                <span className="inline-flex h-9 w-9 items-center justify-center rounded-full bg-fuchsia-500/20 text-fuchsia-300">
                                    <Play className="h-4 w-4" />
                                </span>
                                AI-native social commerce powering creators, merchants, and livestream marketplaces.
                            </div>

                            <div className="space-y-6">
                                <h1 className="text-5xl font-semibold tracking-tight text-white sm:text-6xl lg:text-7xl">
                                    The premium social commerce platform for creators, merchants, and connected shoppers.
                                </h1>
                                <p className="max-w-xl text-lg leading-8 text-slate-300 sm:text-xl">
                                    Teamart blends TikTok-style discovery, Amazon-grade checkout confidence, and TikTok Seller Center operations into one production-ready marketplace.
                                </p>
                                <div className="flex flex-col gap-4 sm:flex-row sm:items-center">
                                    {heroActions.map((action) => (
                                        <Link key={action.label} href={action.href} className="inline-flex items-center justify-center rounded-full bg-fuchsia-500 px-6 py-4 text-sm font-semibold text-white shadow-2xl shadow-fuchsia-500/20 transition hover:bg-fuchsia-400">
                                            {action.label}
                                            <ArrowRight className="ml-3 h-4 w-4" />
                                        </Link>
                                    ))}
                                </div>
                                <div className="grid gap-4 sm:grid-cols-2">
                                    <div className="rounded-3xl border border-white/10 bg-white/5 p-6">
                                        <p className="text-sm uppercase tracking-[0.3em] text-slate-400">Live creators</p>
                                        <p className="mt-3 text-3xl font-semibold text-white">5.2K+</p>
                                        <p className="mt-2 text-sm text-slate-400">Creators hosting commerce-driven livestreams.</p>
                                    </div>
                                    <div className="rounded-3xl border border-white/10 bg-white/5 p-6">
                                        <p className="text-sm uppercase tracking-[0.3em] text-slate-400">Monthly GMV</p>
                                        <p className="mt-3 text-3xl font-semibold text-white">$1.1M</p>
                                        <p className="mt-2 text-sm text-slate-400">High-volume marketplace checkout performance.</p>
                                    </div>
                                </div>
                            </div>

                            <div className="relative flex items-center justify-center">
                                <motion.div
                                    animate={{ y: [0, -18, 0] }}
                                    transition={{ duration: 4, repeat: Infinity, ease: "easeInOut" }}
                                    className="absolute -left-8 top-10 hidden h-24 w-24 rounded-3xl bg-fuchsia-500/20 blur-3xl sm:block"
                                />

                                <motion.div
                                    animate={{ y: [0, 16, 0] }}
                                    transition={{ duration: 4.2, repeat: Infinity, ease: "easeInOut", delay: 0.4 }}
                                    className="absolute right-6 top-24 hidden h-20 w-20 rounded-3xl bg-sky-500/20 blur-3xl sm:block"
                                />

                                <div className="relative w-full max-w-md overflow-hidden rounded-[3rem] border border-white/10 bg-slate-950/80 shadow-2xl shadow-slate-950/30">
                                    <div className="bg-slate-900/95 px-5 py-4 text-white">
                                        <div className="flex items-center justify-between">
                                            <div>
                                                <p className="text-xs uppercase tracking-[0.35em] text-slate-400">LIVE STORE</p>
                                                <p className="mt-1 text-sm text-white/80">Creator room • 18.3k viewers</p>
                                            </div>
                                            <button className="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/5 px-3 py-2 text-xs text-white/80">
                                                <Heart className="h-4 w-4 text-pink-400" />
                                                Save
                                            </button>
                                        </div>
                                    </div>

                                    <div className="space-y-4 p-5">
                                        <div className="aspect-[4/5] overflow-hidden rounded-[2rem] bg-slate-900">
                                            <div className="relative h-full w-full bg-[linear-gradient(135deg,#111827_0%,#0f172a_100%)]">
                                                <div className="absolute inset-x-0 bottom-0 h-24 bg-gradient-to-t from-slate-950/90 to-transparent" />
                                                <div className="absolute left-4 top-4 rounded-3xl bg-slate-950/80 px-3 py-2 text-xs font-semibold uppercase tracking-[0.3em] text-slate-100">
                                                    Live Now
                                                </div>
                                                <div className="absolute bottom-4 left-4 right-4 flex items-center justify-between text-white/90">
                                                    <div>
                                                        <p className="text-sm font-semibold">Trend Studio</p>
                                                        <p className="text-xs text-slate-300">Creator drop • 23 min left</p>
                                                    </div>
                                                    <div className="rounded-full bg-black/60 px-3 py-1 text-[11px] uppercase tracking-[0.24em] text-slate-200">
                                                        Shop
                                                    </div>
                                                </div>
                                            </div>
                                        </div>

                                        <div className="grid grid-cols-2 gap-4">
                                            <div className="rounded-3xl border border-white/10 bg-slate-900/90 p-4">
                                                <p className="text-xs uppercase tracking-[0.3em] text-slate-400">Featured</p>
                                                <p className="mt-3 text-lg font-semibold text-white">Creator Hoodie</p>
                                                <p className="mt-2 text-sm text-slate-400">$58 • Free shipping</p>
                                            </div>
                                            <div className="rounded-3xl border border-white/10 bg-slate-900/90 p-4">
                                                <p className="text-xs uppercase tracking-[0.3em] text-slate-400">Fast metrics</p>
                                                <p className="mt-3 text-lg font-semibold text-white">18.2% conversion</p>
                                                <p className="mt-2 text-sm text-slate-400">Live shoppers today</p>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </section>

                <section className="bg-slate-950/95 px-4 py-20 sm:px-6 lg:px-8">
                    <div className="mx-auto max-w-7xl">
                        <div className="grid gap-8 xl:grid-cols-[0.9fr_1.1fr]">
                            <div className="space-y-6">
                                <p className="text-sm uppercase tracking-[0.35em] text-fuchsia-500">Join Teamart</p>
                                <h2 className="text-4xl font-semibold text-white sm:text-5xl">Choose your role and start with customer-first onboarding.</h2>
                                <p className="max-w-xl text-base leading-8 text-slate-300">
                                    Every creator and merchant begins as a shopper, then unlocks creator or merchant capabilities with a secure, role-based onboarding flow.
                                </p>
                            </div>
                            <div className="grid gap-4 sm:grid-cols-3">
                                {[
                                    {
                                        title: "Shopper",
                                        description: "Browse feeds, watch livestreams, save products, and shop with confidence.",
                                        accent: "bg-fuchsia-500/10 text-fuchsia-300",
                                    },
                                    {
                                        title: "Creator",
                                        description: "Launch profile commerce, host live drops, pin products, and grow followers.",
                                        accent: "bg-sky-500/10 text-sky-300",
                                    },
                                    {
                                        title: "Merchant",
                                        description: "Manage stores, inventory, orders, promotions, and payouts at scale.",
                                        accent: "bg-emerald-500/10 text-emerald-300",
                                    },
                                ].map((card) => (
                                    <div key={card.title} className="rounded-3xl border border-white/10 bg-slate-900/95 p-6">
                                        <div className={`inline-flex rounded-full px-3 py-2 text-xs font-semibold ${card.accent}`}>
                                            {card.title}
                                        </div>
                                        <h3 className="mt-6 text-xl font-semibold text-white">{card.title}</h3>
                                        <p className="mt-3 text-sm leading-6 text-slate-400">{card.description}</p>
                                        <Link href="/auth/register" className="mt-6 inline-flex items-center gap-2 text-sm font-semibold text-white text-opacity-90 hover:text-white">
                                            Start with {card.title.toLowerCase()}
                                            <ArrowRight className="h-4 w-4" />
                                        </Link>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </div>
                </section>

                <section id="feed" className="bg-slate-950/95 px-4 py-20 sm:px-6 lg:px-8">
                    <div className="mx-auto max-w-7xl">
                        <div className="mb-12 flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
                            <div className="max-w-2xl">
                                <p className="text-sm uppercase tracking-[0.35em] text-fuchsia-500">TikTok-style product feed</p>
                                <h2 className="mt-4 text-4xl font-semibold text-white sm:text-5xl">Swipe, shop, and follow in one immersive feed.</h2>
                                <p className="mt-4 text-base leading-7 text-slate-300">
                                    Social commerce built for discovery, product obsession, and impulse purchase behavior.
                                </p>
                            </div>
                            <div className="flex flex-wrap gap-3">
                                <Link href="/feed" className="rounded-full border border-white/10 bg-white/5 px-6 py-3 text-sm font-semibold text-white shadow-white/5 transition hover:bg-white/10">
                                    Open feed
                                </Link>
                                <Link href="/search" className="rounded-full bg-fuchsia-500 px-6 py-3 text-sm font-semibold text-white shadow-xl shadow-fuchsia-500/25 transition hover:bg-fuchsia-400">
                                    Search products
                                </Link>
                            </div>
                        </div>

                        <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
                            {featureCards.map((feature) => (
                                <motion.div
                                    key={feature.title}
                                    initial={{ opacity: 0, y: 20 }}
                                    whileInView={{ opacity: 1, y: 0 }}
                                    viewport={{ once: true }}
                                    transition={{ duration: 0.4 }}
                                    className="rounded-3xl border border-white/10 bg-slate-900/80 p-6 shadow-xl shadow-slate-950/20"
                                >
                                    <div className="mb-5 inline-flex h-14 w-14 items-center justify-center rounded-3xl bg-white/5">
                                        {feature.icon}
                                    </div>
                                    <h3 className="text-xl font-semibold text-white">{feature.title}</h3>
                                    <p className="mt-3 text-sm leading-6 text-slate-400">{feature.description}</p>
                                </motion.div>
                            ))}
                        </div>
                    </div>
                </section>

                <section id="search" className="bg-white px-4 py-20 sm:px-6 lg:px-8">
                    <div className="mx-auto max-w-7xl">
                        <div className="grid gap-12 lg:grid-cols-[0.95fr_1.05fr] xl:gap-20">
                            <div className="space-y-6">
                                <p className="text-sm uppercase tracking-[0.35em] text-fuchsia-600">Discovery by design</p>
                                <h2 className="text-4xl font-semibold text-slate-950 sm:text-5xl">Search, suggestions, and live recommendations at once.</h2>
                                <p className="text-base leading-8 text-slate-600">
                                    Build a search experience that feels native to social shoppers and powerful enough for commerce discovery.
                                </p>
                                <div className="grid gap-4 sm:grid-cols-2">
                                    <div className="rounded-3xl border border-slate-200 bg-slate-50 p-6">
                                        <p className="text-sm uppercase tracking-[0.35em] text-slate-400">Trending</p>
                                        <p className="mt-4 text-2xl font-semibold text-slate-950">#CreatorDrop</p>
                                    </div>
                                    <div className="rounded-3xl border border-slate-200 bg-slate-50 p-6">
                                        <p className="text-sm uppercase tracking-[0.35em] text-slate-400">Live now</p>
                                        <p className="mt-4 text-2xl font-semibold text-slate-950">8 rooms</p>
                                    </div>
                                </div>
                            </div>

                            <div className="space-y-6 rounded-3xl border border-slate-200 bg-slate-50 p-6 shadow-sm">
                                <div className="rounded-3xl bg-white p-4 shadow-inner shadow-slate-200/20">
                                    <div className="flex items-center gap-3 rounded-3xl border border-slate-200 bg-slate-100 px-4 py-3">
                                        <Search className="h-5 w-5 text-slate-500" />
                                        <input
                                            type="search"
                                            placeholder="Search creators, products, livestreams..."
                                            className="w-full bg-transparent text-sm text-slate-900 outline-none placeholder:text-slate-400"
                                        />
                                    </div>
                                </div>
                                <div className="grid gap-4 sm:grid-cols-2">
                                    {[
                                        { title: "Trending products", value: "12.4K" },
                                        { title: "Top creators", value: "340" },
                                    ].map((item) => (
                                        <div key={item.title} className="rounded-3xl border border-slate-200 bg-white p-5">
                                            <p className="text-sm text-slate-500">{item.title}</p>
                                            <p className="mt-4 text-3xl font-semibold text-slate-950">{item.value}</p>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        </div>
                    </div>
                </section>

                <section id="live" className="bg-slate-950 px-4 py-20 sm:px-6 lg:px-8">
                    <div className="mx-auto max-w-7xl">
                        <div className="flex flex-col gap-8 lg:flex-row lg:items-end lg:justify-between">
                            <div className="max-w-2xl">
                                <p className="text-sm uppercase tracking-[0.35em] text-fuchsia-500">Livestream commerce</p>
                                <h2 className="mt-4 text-4xl font-semibold text-white sm:text-5xl">Turn every stream into a checkout-ready event.</h2>
                                <p className="mt-4 text-base leading-7 text-slate-300">
                                    Create immersive livestream rooms with pinned offers, reactions, and instant buyer paths.
                                </p>
                            </div>
                            <Link href="/live" className="self-start rounded-full bg-white px-6 py-3 text-sm font-semibold text-slate-950 transition hover:bg-slate-100">
                                View live rooms
                            </Link>
                        </div>

                        <div className="mt-12 grid gap-6 lg:grid-cols-3">
                            {[
                                { title: "Pinned offers", detail: "Highlight products directly in the stream." },
                                { title: "Viewer reactions", detail: "Real-time engagement for every room." },
                                { title: "Creator tools", detail: "Fast moderation, product pinning, and scheduling." },
                            ].map((item) => (
                                <motion.div
                                    key={item.title}
                                    initial={{ opacity: 0, y: 20 }}
                                    whileInView={{ opacity: 1, y: 0 }}
                                    viewport={{ once: true }}
                                    className="rounded-3xl border border-white/10 bg-white/5 p-6"
                                >
                                    <div className="inline-flex h-12 w-12 items-center justify-center rounded-3xl bg-fuchsia-500/10 text-fuchsia-300">
                                        <Store className="h-6 w-6" />
                                    </div>
                                    <h3 className="mt-5 text-xl font-semibold text-white">{item.title}</h3>
                                    <p className="mt-3 text-sm leading-6 text-slate-300">{item.detail}</p>
                                </motion.div>
                            ))}
                        </div>
                    </div>
                </section>

                <section id="creator" className="bg-white px-4 py-20 sm:px-6 lg:px-8">
                    <div className="mx-auto max-w-7xl">
                        <div className="grid gap-10 lg:grid-cols-[0.9fr_1.1fr] xl:gap-20">
                            <div className="space-y-6">
                                <p className="text-sm uppercase tracking-[0.35em] text-fuchsia-600">Creator toolkit</p>
                                <h2 className="text-4xl font-semibold text-slate-950 sm:text-5xl">Build creator-first commerce experiences.</h2>
                                <p className="text-base leading-8 text-slate-600">
                                    Help creators monetize with live rooms, social feeds, pinned products, and revenue analytics.
                                </p>
                            </div>

                            <div className="space-y-6 rounded-3xl border border-slate-200 bg-slate-50 p-8 shadow-sm">
                                <div className="flex items-center gap-4 rounded-3xl bg-white p-5 shadow-sm">
                                    <div className="rounded-3xl bg-fuchsia-500/10 p-3 text-fuchsia-500">
                                        <Gift className="h-5 w-5" />
                                    </div>
                                    <div>
                                        <p className="text-sm font-medium text-slate-950">Creator commissions</p>
                                        <p className="text-sm text-slate-500">Track earnings and performance across every campaign.</p>
                                    </div>
                                </div>
                                <div className="grid gap-4 sm:grid-cols-2">
                                    <div className="rounded-3xl border border-slate-200 bg-white p-6">
                                        <p className="text-sm text-slate-500">Followers</p>
                                        <p className="mt-3 text-3xl font-semibold text-slate-950">38K</p>
                                    </div>
                                    <div className="rounded-3xl border border-slate-200 bg-white p-6">
                                        <p className="text-sm text-slate-500">Stream engagement</p>
                                        <p className="mt-3 text-3xl font-semibold text-slate-950">92%</p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </section>

                <section id="merchant" className="bg-slate-950 px-4 py-20 sm:px-6 lg:px-8">
                    <div className="mx-auto max-w-7xl">
                        <div className="grid gap-10 lg:grid-cols-[0.9fr_1.1fr] xl:gap-20">
                            <div className="space-y-6">
                                <p className="text-sm uppercase tracking-[0.35em] text-fuchsia-500">Merchant command center</p>
                                <h2 className="text-4xl font-semibold text-white sm:text-5xl">Run your storefront with precision and speed.</h2>
                                <p className="text-base leading-7 text-slate-300">
                                    Detailed analytics, order flow, inventory controls, and campaign tools for high-volume merchants.
                                </p>
                            </div>

                            <div className="space-y-6 rounded-3xl border border-white/10 bg-white/5 p-8 shadow-xl">
                                <div className="flex flex-wrap gap-4">
                                    {bnplTabs.map((tab, index) => (
                                        <button
                                            key={tab}
                                            onClick={() => setActiveBnpl(index)}
                                            className={`rounded-full px-5 py-3 text-sm font-semibold transition ${activeBnpl === index ? "bg-white text-slate-950" : "bg-white/10 text-slate-300 hover:bg-white/20"}`}
                                        >
                                            {tab}
                                        </button>
                                    ))}
                                </div>
                                <div className="rounded-3xl bg-slate-900 p-6 text-slate-100">
                                    <h3 className="text-2xl font-semibold">{bnplTabs[activeBnpl]}</h3>
                                    <p className="mt-3 text-sm leading-7 text-slate-300">
                                        {[
                                            "Launch urgent livestream drops with product pins and social commerce boosts.",
                                            "Showcase curated creator collections with seamless buy journeys.",
                                            "Manage store campaigns, discounts, and livestream promos from one dashboard.",
                                            "Group complementary items into fast-add bundles that feel natural in feed.",
                                        ][activeBnpl]}
                                    </p>
                                </div>
                            </div>
                        </div>
                    </div>
                </section>

                <section className="bg-white px-4 py-20 sm:px-6 lg:px-8">
                    <div className="mx-auto max-w-7xl">
                        <div className="text-center">
                            <p className="text-sm uppercase tracking-[0.35em] text-fuchsia-600">How it works</p>
                            <h2 className="mt-4 text-4xl font-semibold text-slate-950 sm:text-5xl">From onboarding to live commerce in four steps.</h2>
                        </div>

                        <div className="mt-14 grid gap-6 lg:grid-cols-4">
                            {howItWorks.map((step) => (
                                <motion.div
                                    key={step.number}
                                    initial={{ opacity: 0, y: 20 }}
                                    whileInView={{ opacity: 1, y: 0 }}
                                    viewport={{ once: true }}
                                    transition={{ duration: 0.35 }}
                                    className="rounded-3xl border border-slate-200 bg-slate-50 p-7"
                                >
                                    <div className="flex items-center gap-4">
                                        <div className={`flex h-14 w-14 items-center justify-center rounded-3xl bg-slate-950 text-2xl text-white ${step.color}`}>
                                            {step.number}
                                        </div>
                                        <div>
                                            <h3 className="text-xl font-semibold text-slate-950">{step.title}</h3>
                                        </div>
                                    </div>
                                    <p className="mt-4 text-sm leading-7 text-slate-600">{step.description}</p>
                                </motion.div>
                            ))}
                        </div>
                    </div>
                </section>

                <section className="bg-slate-950 px-4 py-20 sm:px-6 lg:px-8">
                    <div className="mx-auto grid max-w-7xl gap-12 lg:grid-cols-[0.95fr_1.05fr] xl:gap-20">
                        <div className="space-y-6">
                            <p className="text-sm uppercase tracking-[0.35em] text-fuchsia-500">Business foundation</p>
                            <h2 className="text-4xl font-semibold text-white sm:text-5xl">A platform built to support creators, sellers, and admins.</h2>
                            <p className="text-base leading-7 text-slate-300">
                                Teamart delivers the workflows, dashboards, and APIs needed for a complete enterprise-ready social commerce system.
                            </p>
                        </div>
                        <div className="space-y-6 rounded-3xl border border-white/10 bg-white/5 p-8">
                            <div className="flex flex-wrap gap-3 bg-slate-900/80 rounded-3xl p-5">
                                {businessTabs.map((tab, index) => (
                                    <button
                                        key={tab}
                                        onClick={() => setActiveBusinessTab(index)}
                                        className={`rounded-full px-4 py-2 text-sm font-semibold transition ${activeBusinessTab === index ? "bg-fuchsia-500 text-white" : "bg-white/10 text-slate-200 hover:bg-white/20"}`}
                                    >
                                        {tab}
                                    </button>
                                ))}
                            </div>
                            <div className="rounded-3xl bg-slate-900 p-6 text-slate-100">
                                <h3 className="text-2xl font-semibold">{businessTabs[activeBusinessTab]}</h3>
                                <p className="mt-4 text-sm leading-7 text-slate-300">
                                    {[
                                        "Creator storefronts with social discovery, follow, and share paths.",
                                        "Seller dashboards with sales, inventory, and campaign intelligence.",
                                        "AI-recommended products, bundles, and livestream promos.",
                                        "Connected live video commerce with chat, pinning, and instant checkout.",
                                    ][activeBusinessTab]}
                                </p>
                            </div>
                        </div>
                    </div>
                </section>

                <section className="bg-gradient-to-br from-fuchsia-600 to-slate-950 px-4 py-20 sm:px-6 lg:px-8">
                    <div className="mx-auto max-w-5xl text-center text-white">
                        <div className="rounded-[3rem] border border-white/10 bg-white/5 p-12 shadow-2xl shadow-slate-950/30">
                            <span className="inline-flex rounded-full bg-white/10 px-4 py-2 text-xs uppercase tracking-[0.35em] text-slate-200">AI-native commerce</span>
                            <h2 className="mt-6 text-4xl font-semibold sm:text-5xl">Launch premium social commerce at scale.</h2>
                            <p className="mx-auto mt-6 max-w-2xl text-base leading-8 text-slate-200">
                                Teamart combines immersive discovery, trusted purchase flows, and commerce controls for marketplaces, creators, and brands.
                            </p>
                            <div className="mt-10 flex flex-col gap-4 sm:flex-row sm:justify-center">
                                <Link href="/auth/register" className="inline-flex items-center justify-center rounded-full bg-white px-8 py-4 text-sm font-semibold text-slate-950 transition hover:bg-slate-100">
                                    Start your marketplace
                                </Link>
                                <Link href="/dashboard" className="inline-flex items-center justify-center rounded-full border border-white/20 bg-white/5 px-8 py-4 text-sm font-semibold text-white transition hover:bg-white/10">
                                    Explore dashboard
                                </Link>
                            </div>
                        </div>
                    </div>
                </section>

                <footer className="bg-slate-950 px-4 py-16 text-slate-300 sm:px-6 lg:px-8">
                    <div className="mx-auto grid max-w-7xl gap-10 lg:grid-cols-4">
                        <div className="space-y-4">
                            <div className="flex items-center gap-3 text-white">
                                <div className="grid h-12 w-12 place-items-center rounded-3xl bg-gradient-to-br from-fuchsia-500 to-orange-500 text-white">
                                    <ArrowUp className="h-5 w-5" />
                                </div>
                                <div>
                                    <p className="text-xl font-semibold">Teamart</p>
                                    <p className="text-sm text-slate-400">AI-native social commerce platform</p>
                                </div>
                            </div>
                            <p className="text-sm text-slate-400">Premium live shopping, creator commerce, and seller center workflows built for modern marketplaces.</p>
                        </div>
                        <div>
                            <h3 className="mb-5 text-sm font-semibold uppercase tracking-[0.3em] text-slate-500">Product</h3>
                            <ul className="space-y-3 text-sm text-slate-400">
                                <li><a href="#feed" className="transition hover:text-white">Feed</a></li>
                                <li><a href="#search" className="transition hover:text-white">Search</a></li>
                                <li><a href="#live" className="transition hover:text-white">Live rooms</a></li>
                                <li><a href="#merchant" className="transition hover:text-white">Merchant tools</a></li>
                            </ul>
                        </div>
                        <div>
                            <h3 className="mb-5 text-sm font-semibold uppercase tracking-[0.3em] text-slate-500">Company</h3>
                            <ul className="space-y-3 text-sm text-slate-400">
                                <li><a href="#" className="transition hover:text-white">About</a></li>
                                <li><a href="#" className="transition hover:text-white">Blog</a></li>
                                <li><a href="#" className="transition hover:text-white">Careers</a></li>
                                <li><a href="#" className="transition hover:text-white">Contact</a></li>
                            </ul>
                        </div>
                        <div>
                            <h3 className="mb-5 text-sm font-semibold uppercase tracking-[0.3em] text-slate-500">Legal</h3>
                            <ul className="space-y-3 text-sm text-slate-400">
                                <li><a href="#" className="transition hover:text-white">Terms</a></li>
                                <li><a href="#" className="transition hover:text-white">Privacy</a></li>
                                <li><a href="#" className="transition hover:text-white">Security</a></li>
                            </ul>
                        </div>
                    </div>
                    <div className="mt-12 border-t border-white/10 pt-6 text-sm text-slate-500">© 2026 Teamart. Built for creators, merchants, and live shoppers.</div>
                </footer>
            </main>

            <AnimatePresence>
                {showTop && (
                    <motion.button
                        initial={{ opacity: 0, y: 24 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: 24 }}
                        onClick={() => window.scrollTo({ top: 0, behavior: "smooth" })}
                        className="fixed bottom-6 right-6 z-50 inline-flex h-14 w-14 items-center justify-center rounded-full border border-white/10 bg-white text-slate-950 shadow-2xl shadow-slate-950/20 transition hover:bg-slate-100"
                    >
                        <ArrowUp className="h-5 w-5" />
                    </motion.button>
                )}
            </AnimatePresence>
        </div>
    );
}
