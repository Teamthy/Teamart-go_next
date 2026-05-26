import Link from "next/link";
import { notFound } from "next/navigation";
import PageHeader from "@/components/ui/PageHeader";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import StatCard from "@/components/ui/StatCard";
import ProductCard from "@/components/product/ProductCard";
import { products } from "@/lib/mock/products";
import { stores } from "@/lib/mock/stores";

export default async function StoreDetailPage({ params }: { params: Promise<{ slug: string }> }) {
    const { slug } = await params;
    const store = stores.find((item) => item.slug === slug);

    if (!store) {
        notFound();
    }

    const featuredProducts = products.filter((item) => item.merchant === store.name).slice(0, 4);

    return (
        <div className="space-y-8 pb-10">
            <PageHeader
                title={store.name}
                description={store.tagline}
                actions={
                    <>
                        <Button asChild variant="primary">
                            <Link href="/live">Join live room</Link>
                        </Button>
                        <Button asChild variant="secondary">
                            <Link href="/stores">Back to stores</Link>
                        </Button>
                    </>
                }
            />

            <div className="grid gap-4 md:grid-cols-3">
                <StatCard label="Category" value={store.category} helper="A strong fit for shoppers browsing by merchant category." />
                <StatCard label="Followers" value={store.followers} helper="A growing audience for featured storefront content." />
                <StatCard label="Rating" value={store.rating} helper="Positive shopper sentiment and repeat-buyer alignment." />
            </div>

            <div className="grid gap-4 xl:grid-cols-[0.8fr_1.2fr]">
                <Card className="overflow-hidden p-0">
                    <img src={store.banner} alt={store.name} className="h-64 w-full object-cover" />
                    <div className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">About the store</p>
                        <h2 className="mt-3 text-xl font-semibold text-zinc-900">Live-first merchandising with premium storytelling</h2>
                        <p className="mt-3 text-sm leading-7 text-zinc-600">
                            This storefront is designed to keep product discovery fast, visually rich, and ready for creators or merchants who want a polished shopping moment.
                        </p>
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Store highlights</p>
                    <div className="mt-4 space-y-3">
                        {[
                            `Creator: ${store.creator}`,
                            `Live status: ${store.live}`,
                            `Products: ${store.products}`,
                            `Audience: ${store.followers}`,
                        ].map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">
                                {item}
                            </div>
                        ))}
                    </div>
                    <div className="mt-5 rounded-[24px] bg-zinc-950 p-4 text-white">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-pink-200">This week</p>
                        <p className="mt-3 text-sm leading-7 text-zinc-100">
                            The store is leaning into a mix of creator-led bundles, live room highlights, and fast-moving add-on items to keep browsers engaged longer.
                        </p>
                    </div>
                </Card>
            </div>

            <section className="space-y-4">
                <div className="flex flex-wrap items-end justify-between gap-4">
                    <div>
                        <p className="text-[11px] uppercase tracking-[0.2em] text-[#E91E63]">Featured products</p>
                        <h2 className="mt-2 text-2xl font-semibold text-zinc-900">Shop the best of {store.name}</h2>
                    </div>
                    <Button asChild variant="secondary">
                        <Link href="/products">Explore all products</Link>
                    </Button>
                </div>
                <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-4">
                    {featuredProducts.map((product) => (
                        <ProductCard key={product.id} product={product} />
                    ))}
                </div>
            </section>
        </div>
    );
}
