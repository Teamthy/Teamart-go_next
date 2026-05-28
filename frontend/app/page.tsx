<<<<<<< HEAD
import Link from "next/link";
import FeedCard from "@/components/ui/FeedCard";
import ProductCard from "@/components/product/ProductCard";
import SectionHeader from "@/components/ui/SectionHeader";
import Badge from "@/components/ui/badge";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import { products, recommendedProducts } from "@/lib/mock/products";
import { creators } from "@/lib/mock/creators";
import { feedItems } from "@/lib/mock/feed";
import { stores } from "@/lib/mock/stores";

const categories = ["Fashion", "Beauty", "Home", "Tech", "Fitness", "Accessories"];

export default function Home() {
  return (
    <div className="space-y-8 pb-12">
      <section className="rounded-[32px] bg-[linear-gradient(135deg,#fff8fb_0%,#ffffff_55%,#ecfdf5_100%)] p-5 sm:p-6">
        <div className="grid gap-6 xl:grid-cols-[1.1fr_0.9fr] xl:items-center">
          <div className="space-y-5">
            <Badge tone="default">TikTok + Amazon + Seller Central layer</Badge>
            <div className="space-y-3">
              <p className="text-[11px] uppercase tracking-[0.24em] text-[#E91E63]">Teamart-go_next</p>
              <h1 className="text-[30px] font-semibold tracking-tight text-zinc-900 sm:text-[34px]">
                Premium social commerce with scroll-first discovery and checkout-ready moments.
              </h1>
              <p className="max-w-2xl text-sm leading-7 text-zinc-600 sm:text-[15px]">
                Swipe through creator-led content, launch into live rooms, follow stores, and buy from a polished commerce surface designed to feel premium on mobile and desktop.
              </p>
            </div>
            <div className="flex flex-wrap gap-3">
              <Button asChild variant="primary">
                <Link href="/feed">Open the feed</Link>
              </Button>
              <Button asChild variant="secondary">
                <Link href="/search">Search products</Link>
              </Button>
              <Button asChild variant="secondary">
                <Link href="/live">Join live</Link>
              </Button>
            </div>
            <div className="grid gap-3 sm:grid-cols-3">
              {[
                { label: "Live shoppers", value: "18.4k" },
                { label: "Active stores", value: "320" },
                { label: "Orders today", value: "$86k" },
              ].map((stat) => (
                <div key={stat.label} className="rounded-[24px] bg-white px-4 py-4 shadow-[0_18px_45px_rgba(15,23,42,0.08)]">
                  <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{stat.label}</p>
                  <p className="mt-2 text-2xl font-semibold text-zinc-900">{stat.value}</p>
                </div>
              ))}
            </div>
          </div>

          <div className="space-y-4">
            <Card className="overflow-hidden p-0">
              <div className="bg-zinc-950 p-4 text-white sm:p-5">
                <div className="flex items-center justify-between gap-4">
                  <div>
                    <p className="text-[11px] uppercase tracking-[0.2em] text-pink-200">Today’s spotlight</p>
                    <p className="mt-2 text-lg font-semibold">Creator Collaboration Hoodie</p>
    <>
      <style>
        {`
          @import url("https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;500;600;700&display=swap");
=======
﻿import Link from "next/link";
import Navbar from "@/components/ui/Navbar";
import PageHeader from "@/components/ui/PageHeader";
import SearchBar from "@/components/ui/SearchBar";
import StatCard from "@/components/ui/StatCard";
import ProductGrid from "@/components/ui/ProductGrid";
import Card from "@/components/ui/card";
import EmptyState from "@/components/ui/EmptyState";
import { recommendedProducts, featuredProducts, categories } from "@/lib/mock-data";

export default function HomePage() {
  return (
    <div className="min-h-screen bg-slate-950 text-slate-50">
      <Navbar />
>>>>>>> 8018627 (feat(feed): tiktok-style feed, search and landing experience)

      <main className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        <div className="grid gap-10 lg:grid-cols-[1fr_0.44fr]">
          <section>
            <PageHeader
              eyebrow="Discover"
              title="Shop creator drops, livestreams, and curated collections"
              description="A discovery-led social commerce feed — live product drops, creator storefronts, and fast checkout built for modern shoppers."
            />

<<<<<<< HEAD
          body {
            background: #0b1120;
          }
        `}
      </style>

      <header className="flex min-h-screen flex-col items-center overflow-hidden bg-gradient-to-b from-slate-950 via-slate-900 to-slate-950 text-white">
        <nav className="sticky top-0 z-50 w-full border-b border-white/10 bg-white/5 backdrop-blur-md">
          <div className="flex items-center justify-between px-6 py-4 md:px-12 lg:px-20">
            <Link href="/" className="flex items-center gap-3">
              <div className="flex size-10 items-center justify-center rounded-xl bg-indigo-600 text-lg font-bold shadow-lg shadow-indigo-600/30">
                T
              </div>

              <div>
                <h1 className="text-lg font-semibold tracking-tight">Teamart</h1>
                <p className="-mt-1 text-xs text-slate-400">AI Commerce Platform</p>
              </div>
            </Link>

            <div
              className={`
                ${mobileOpen ? "max-md:w-full" : "max-md:w-0"}
                flex items-center gap-8 text-sm
                max-md:fixed
                max-md:left-0
                max-md:top-0
                max-md:h-screen
                max-md:flex-col
                max-md:justify-center
                max-md:overflow-hidden
                max-md:bg-slate-950/95
                max-md:backdrop-blur-xl
                max-md:transition-all
                max-md:duration-300
              `}
            >
              <Link href="#features" onClick={() => setMobileOpen(false)} className="text-slate-300 transition hover:text-white">
                Features
              </Link>
              <Link href="#creators" onClick={() => setMobileOpen(false)} className="text-slate-300 transition hover:text-white">
                Creators
              </Link>
              <Link href="#commerce" onClick={() => setMobileOpen(false)} className="text-slate-300 transition hover:text-white">
                Commerce
              </Link>
              <Link href="#pricing" onClick={() => setMobileOpen(false)} className="text-slate-300 transition hover:text-white">
                Pricing
              </Link>

              <button onClick={() => setMobileOpen(false)} className="rounded-lg bg-white p-2 text-black md:hidden">
                ✕
              </button>
            </div>

            <div className="hidden items-center gap-4 md:flex">
              <Link href="/dashboard" className="rounded-lg border border-white/10 px-5 py-2.5 transition hover:bg-white/5">
                Dashboard
              </Link>
              <Link href="/products" className="rounded-lg bg-indigo-600 px-5 py-2.5 font-medium shadow-lg shadow-indigo-600/20 transition hover:bg-indigo-500">
                Launch Store
              </Link>
            </div>

            <button onClick={() => setMobileOpen(true)} className="rounded-lg border border-white/10 bg-white/10 p-2 md:hidden">
              ☰
            </button>
          </div>
        </nav>

        <section className="mx-auto grid w-full max-w-7xl items-center gap-16 px-6 py-20 lg:grid-cols-2 lg:px-20 lg:py-28">
          <div>
            <div className="inline-flex items-center gap-2 rounded-full border border-indigo-500/20 bg-indigo-500/10 px-4 py-2 text-sm text-indigo-300">
              <span className="size-2 animate-pulse rounded-full bg-indigo-400"></span>
              AI Native Commerce Infrastructure
            </div>

            <h1 className="mt-6 text-5xl font-semibold leading-tight tracking-tight md:text-6xl">
              Build the next generation of <span className="text-indigo-400">creator commerce</span> with AI automation
            </h1>

            <p className="mt-6 max-w-2xl text-lg leading-8 text-slate-400">
              Teamart helps creators, brands, and digital businesses launch scalable AI-powered storefronts, automate workflows,
              manage users, process payments, and grow revenue from one unified platform.
            </p>

            <div className="mt-8 flex flex-wrap gap-4">
              <Link href="/auth/register" className="inline-block rounded-xl bg-indigo-600 px-7 py-4 font-medium shadow-xl shadow-indigo-600/20 transition hover:bg-indigo-500">
                Get Started
              </Link>
              <Link href="/docs" className="inline-block rounded-xl border border-white/10 px-7 py-4 transition hover:bg-white/5">
                View Documentation
              </Link>
            </div>

            <div className="mt-12 grid max-w-xl grid-cols-3 gap-6">
              <div>
                <h2 className="text-3xl font-bold">10K+</h2>
                <p className="mt-1 text-sm text-slate-400">Active Users</p>
              </div>
              <div>
                <h2 className="text-3xl font-bold">500+</h2>
                <p className="mt-1 text-sm text-slate-400">Creator Stores</p>
              </div>
              <div>
                <h2 className="text-3xl font-bold">$2M+</h2>
                <p className="mt-1 text-sm text-slate-400">Revenue Processed</p>
              </div>
            </div>
          </div>

          <div className="relative">
            <div className="absolute inset-0 rounded-full bg-indigo-600/20 blur-3xl"></div>
            <div className="relative rounded-3xl border border-white/10 bg-white/5 p-6 shadow-2xl backdrop-blur-xl">
              <div className="flex items-center justify-between border-b border-white/10 pb-4">
                <div>
                  <h3 className="text-lg font-semibold">Analytics Dashboard</h3>
                  <p className="text-sm text-slate-400">Real-time commerce insights</p>
                </div>
                <div className="rounded-full bg-emerald-500/20 px-3 py-1 text-sm text-emerald-400">Live</div>
              </div>

              <div className="mt-6 grid grid-cols-2 gap-4">
                <div className="rounded-2xl border border-white/5 bg-slate-900/70 p-5">
                  <p className="text-sm text-slate-400">Monthly Revenue</p>
                  <h2 className="mt-2 text-3xl font-bold">$48,320</h2>
                  <p className="mt-2 text-sm text-emerald-400">+18.4% growth</p>
                </div>

                <div className="rounded-2xl border border-white/5 bg-slate-900/70 p-5">
                  <p className="text-sm text-slate-400">AI Orders</p>
                  <h2 className="mt-2 text-3xl font-bold">12,847</h2>
                  <p className="mt-2 text-sm text-indigo-400">Automated workflows</p>
                </div>

                <div className="col-span-2 rounded-2xl border border-white/5 bg-slate-900/70 p-5">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm text-slate-400">User Growth</p>
                      <h2 className="mt-2 text-3xl font-bold">+324%</h2>
                    </div>
                    <div className="flex h-24 w-40 items-end gap-2 rounded-xl border border-white/10 bg-gradient-to-tr from-indigo-500/20 to-cyan-500/20 p-3">
                      <div className="h-10 w-4 rounded bg-indigo-400"></div>
                      <div className="h-16 w-4 rounded bg-indigo-400"></div>
                      <div className="h-12 w-4 rounded bg-indigo-400"></div>
                      <div className="h-20 w-4 rounded bg-indigo-400"></div>
                      <div className="h-14 w-4 rounded bg-indigo-400"></div>
                    </div>
                  </div>
                  <Badge tone="success">LIVE</Badge>
                </div>
                <p className="mt-3 text-sm leading-6 text-white/85">
                  A premium drop with bundle incentives, fast purchase paths, and creator-first storytelling built for social conversion.
                </p>
              </div>
              <div className="grid gap-3 p-4 sm:grid-cols-2 sm:p-5">
                <div className="rounded-[24px] bg-[#FFF8FB] p-4">
                  <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Revenue</p>
                  <p className="mt-2 text-2xl font-semibold text-zinc-900">$12.4k</p>
                </div>
                <div className="rounded-[24px] bg-[#FFF8FB] p-4">
                  <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Conversion</p>
                  <p className="mt-2 text-2xl font-semibold text-zinc-900">8.2%</p>

              <div className="mt-6 rounded-2xl border border-white/5 bg-slate-900/70 p-5">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="font-semibold">Registered Users</h3>
                    <p className="text-sm text-slate-400">Connected accounts across the platform</p>
                  </div>
                  <button className="text-sm text-indigo-400 hover:text-indigo-300">View All</button>
                </div>

                <div className="mt-4 space-y-3">
                  {[
                    "Creator Economy",
                    "Digital Products",
                    "Affiliate Commerce",
                    "AI Automation",
                  ].map((item, index) => (
                    <div key={index} className="flex items-center justify-between rounded-xl bg-white/5 px-4 py-3">
                      <div className="flex items-center gap-3">
                        <div className="flex size-10 items-center justify-center rounded-full bg-indigo-500/20">🚀</div>
                        <div>
                          <h4 className="font-medium">{item}</h4>
                          <p className="text-xs text-slate-400">Active Module</p>
                        </div>
                      </div>
                      <span className="text-sm text-emerald-400">Online</span>
                    </div>
=======
            <div className="mt-6">
              <SearchBar placeholder="Search products, creators, or live rooms..." />
            </div>

            <div className="mt-8 grid gap-6 lg:grid-cols-2">
              <Card>
                <h3 className="text-lg font-semibold text-slate-900">Featured categories</h3>
                <div className="mt-4 grid grid-cols-3 gap-3">
                  {categories.slice(0, 6).map((c) => (
                    <Link
                      key={c.slug}
                      href={`/categories/${c.slug}`}
                      className="rounded-2xl border border-white/6 bg-white/5 px-3 py-2 text-sm text-slate-100 hover:bg-white/6"
                    >
                      {c.name}
                    </Link>
>>>>>>> 8018627 (feat(feed): tiktok-style feed, search and landing experience)
                  ))}
                </div>
              </Card>

              <Card>
                <h3 className="text-lg font-semibold text-slate-900">Why Teamart</h3>
                <div className="mt-4 grid gap-3">
                  <StatCard label="Live rooms" value="24" helper="Active daily" />
                  <StatCard label="Creators" value="1.2k" helper="Top creators onboarded" />
                </div>
              </Card>
            </div>

            <div className="mt-8">
              <h2 className="text-xl font-semibold text-slate-50">Featured products</h2>
              <div className="mt-4">
                <ProductGrid products={featuredProducts.slice(0, 6)} />
              </div>
            </Card>

            <div className="grid gap-4 sm:grid-cols-2">
              <Card className="p-4">
                <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Social pulse</p>
                <p className="mt-3 text-lg font-semibold text-zinc-900">37 creator moments</p>
                <p className="mt-2 text-sm leading-6 text-zinc-600">Fresh social content is surfaced in one clean scrollable experience.</p>
              </Card>
              <Card className="p-4">
                <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Merchant tools</p>
                <p className="mt-3 text-lg font-semibold text-zinc-900">24 live offers</p>
                <p className="mt-2 text-sm leading-6 text-zinc-600">Promos, inventory, and product performance are all ready for professional merchandising.</p>
              </Card>
            </div>
<<<<<<< HEAD
          </div>
        </div>
      </section>

      <section className="space-y-3">
        <SectionHeader
          title="Browse by category"
          description="Jump into the most active shopping spaces without leaving a premium, social-first experience."
        />
        <div className="flex flex-wrap gap-2">
          {categories.map((category) => (
            <Link
              key={category}
              href="/search"
              className="rounded-full border border-zinc-200 bg-white px-4 py-2 text-sm font-semibold text-zinc-700"
            >
              {category}
            </Link>
          ))}
        </div>
      </section>

      <section className="grid gap-4 xl:grid-cols-[1fr_0.9fr]">
        <div className="space-y-4">
          <SectionHeader
            title="For your next scroll"
            description="A curated slice of creator drops, merchant promos, and commerce-first moments styled for a premium social feed."
          />
          <div className="grid gap-4 md:grid-cols-2">
            {feedItems.slice(0, 2).map((item) => (
              <FeedCard key={item.id} item={item} />
            ))}
          </div>
        </div>

        <div className="space-y-4">
          <Card className="p-5">
            <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Recommended stores</p>
            <div className="mt-3 space-y-3">
              {stores.slice(0, 4).map((store) => (
                <Link key={store.slug} href={`/stores/${store.slug}`} className="block rounded-[24px] bg-[#FFF8FB] p-3">
                  <p className="text-sm font-semibold text-zinc-900">{store.name}</p>
                  <p className="mt-1 text-xs text-zinc-600">{store.category} • {store.live}</p>
                </Link>
              ))}
            </div>
          </Card>

          <Card className="p-5">
            <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Creator spotlight</p>
            <div className="mt-3 grid gap-3 sm:grid-cols-2">
              {creators.map((creator) => (
                <Link key={creator.id} href={`/creator/${creator.id}`} className="rounded-[24px] bg-zinc-50 p-3">
                  <p className="text-sm font-semibold text-zinc-900">{creator.name}</p>
                  <p className="mt-1 text-xs text-zinc-500">{creator.followers}</p>
                </Link>
              ))}
            </div>
          </Card>
        </div>
      </section>

      <section className="space-y-4">
        <SectionHeader
          title="Trending products"
          description="A polished catalog slice with creator favorites and social-ready commerce picks."
        />
        <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
          {products.slice(0, 6).map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      </section>

      <section className="grid gap-4 lg:grid-cols-[0.9fr_1.1fr]">
        <Card className="p-5">
          <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Live commerce</p>
          <h2 className="mt-3 text-[22px] font-semibold text-zinc-900">Keep every drop and promotion conversion-ready</h2>
          <p className="mt-2 text-sm leading-7 text-zinc-600">
            Merge creator storytelling, merchant merchandising, and live shopping prompts into one premium commercial surface that helps shoppers move from discovery to checkout with less friction.
          </p>
          <div className="mt-4 flex flex-wrap gap-3">
            <Button asChild variant="primary">
              <Link href="/creator/studio">Open creator studio</Link>
            </Button>
            <Button asChild variant="secondary">
              <Link href="/merchant">Open merchant workspace</Link>
            </Button>
          </div>
        </Card>

        <div className="grid gap-4 sm:grid-cols-2">
          {recommendedProducts.slice(0, 4).map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      </section>
    </div>
        </section>
      </header>
    </>
  );
}
=======

            <div className="mt-10">
              <h2 className="text-xl font-semibold text-slate-50">Recommended for you</h2>
              <div className="mt-4">
                {recommendedProducts.length ? (
                  <ProductGrid products={recommendedProducts.slice(0, 6)} />
                ) : (
                  <EmptyState title="No recommendations yet" description="Complete onboarding to get personalized picks." />
                )}
              </div>
            </div>
          </section>

          <aside>
            <div className="sticky top-20 space-y-6">
              <Card>
                <h3 className="text-lg font-semibold text-slate-900">Live now</h3>
                <p className="mt-2 text-sm text-slate-400">Join live drops and score limited offers.</p>
              </Card>

              <Card>
                <h3 className="text-lg font-semibold text-slate-900">Quick actions</h3>
                <div className="mt-4 grid gap-3">
                  <Link href="/feed" className="rounded-full bg-fuchsia-500 px-4 py-2 text-sm font-semibold text-white text-center">Explore feed</Link>
                  <Link href="/auth/login" className="rounded-full border border-white/10 px-4 py-2 text-sm font-semibold text-white text-center">Sign in</Link>
                </div>
              </Card>
            </div>
          </aside>
        </div>
      </main>
    </div>
  );
}

>>>>>>> 8018627 (feat(feed): tiktok-style feed, search and landing experience)
