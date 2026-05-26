import Link from "next/link";
import { notFound } from "next/navigation";
import CreatorProfileCard from "@/components/product/CreatorProfileCard";
import ProductCard from "@/components/product/ProductCard";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import { creatorProfiles } from "@/lib/mock/creators";
import { products } from "@/lib/mock/products";

export default async function CreatorProfilePage({ params }: { params: Promise<{ slug: string }> }) {
    const { slug } = await params;
    const creator = creatorProfiles.find((item) => item.id === slug);

    if (!creator) {
        notFound();
    }

    const creatorProducts = products.filter((product) => product.creator === creator.name).slice(0, 4);

    return (
        <div className="space-y-8 pb-10">
            <CreatorProfileCard creator={creator} />

            <div className="grid gap-4 xl:grid-cols-[0.9fr_1.1fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Creator notes</p>
                    <h2 className="mt-3 text-xl font-semibold text-zinc-900">What this creator is spotlighting</h2>
                    <div className="mt-4 space-y-3 text-sm text-zinc-600">
                        <p>{creator.bio}</p>
                        <div className="rounded-[20px] bg-[#FFF8FB] px-4 py-3">
                            <p className="font-semibold text-zinc-900">Live schedule</p>
                            <p className="mt-1">{creator.livestreamSchedule}</p>
                        </div>
                    </div>
                    <div className="mt-5 flex flex-wrap gap-3">
                        <Button asChild variant="primary">
                            <Link href="/live">Join live</Link>
                        </Button>
                        <Button asChild variant="secondary">
                            <Link href="/products">Shop creator picks</Link>
                        </Button>
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Featured products</p>
                    <div className="mt-4 grid gap-4 md:grid-cols-2">
                        {creatorProducts.map((product) => (
                            <ProductCard key={product.id} product={product} />
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}
