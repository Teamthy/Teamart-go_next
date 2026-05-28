import Link from "next/link";
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

      <main className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        <div className="grid gap-10 lg:grid-cols-[1fr_0.44fr]">
          <section>
            <PageHeader
              eyebrow="Discover"
              title="Shop creator drops, livestreams, and curated collections"
              description="A discovery-led social commerce feed — live product drops, creator storefronts, and fast checkout built for modern shoppers."
            />

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
            </div>

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

