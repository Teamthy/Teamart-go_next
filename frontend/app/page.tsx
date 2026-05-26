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
                </div>
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
  );
}
